package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	GoogleOIDCConfig  = &oauth2.Config{}
	GoogleOAuthConfig = &oauth2.Config{}
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	}

	GoogleOAuthConfig = getConfig(
		"http://localhost:8080/callback",
		[]string{"https://www.googleapis.com/auth/userinfo.email"},
	)
	GoogleOIDCConfig = getConfig(
		"http://localhost:8080/callback/oidc",
		[]string{"openid", "https://www.googleapis.com/auth/userinfo.email"},
	)

	http.HandleFunc("/", handleMain)
	http.HandleFunc("/oauth", handleOAuth)
	http.HandleFunc("/oidc", handleOIDC)
	http.HandleFunc("/callback", handleCallback)
	http.HandleFunc("/callback/oidc", handleCallbackOIDC)
	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getConfig(redirectURL string, scopes []string) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  redirectURL,
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
		<a href="/oauth">OAuth 2.0</a><br/>
		<a href="/oidc">OpenID Connect</a>
	`)
}

func handleOAuth(w http.ResponseWriter, r *http.Request) {
	url := GoogleOAuthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleOIDC(w http.ResponseWriter, r *http.Request) {
	url := GoogleOIDCConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	// oauthは、access tokenを使って情報を取得する

	code := r.URL.Query().Get("code")

	token, err := GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	userInfo, err := getUserInfo(token.AccessToken)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	fmt.Fprintf(w, "Access Token: %s\nUser Email: %s", token.AccessToken, userInfo.Email)
}

func handleCallbackOIDC(w http.ResponseWriter, r *http.Request) {
	// oidcは、id tokenを使って情報を取得する
	// oauthは、発行されたaccess tokenが流出した場合、第三者がユーザーになりすまし、リソースにアクセスできてしまう可能性がある。
	// oidcは、id tokenとaccess tokenの両方を必要とすることで、ユーザーが認証されたことを保証してリソースアクセスする。

	code := r.URL.Query().Get("code")

	oauth2Token, err := GoogleOIDCConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		fmt.Fprintf(w, "Error: No ID token found")
		return
	}
	idToken, err := oidcVerifyAndParseIDToken(oauth2Token, rawIDToken)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	fmt.Fprintf(w, "Access Token: %s\nUser Email: %s", oauth2Token.AccessToken, idToken.Email)
}

func oidcVerifyAndParseIDToken(oauth2Token *oauth2.Token, rawIDToken string) (*tokenPayload, error) {
	provider, err := oidc.NewProvider(context.Background(), "https://accounts.google.com")
	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %v", err)
	}

	// IDトークン検証用の検証器を作成
	verifier := provider.Verifier(&oidc.Config{ClientID: GoogleOIDCConfig.ClientID})

	// IDトークンを検証
	idToken, err := verifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID Token: %v", err)
	}

	// トークンペイロードを取得
	payload := &tokenPayload{}
	if err := idToken.Claims(&payload); err != nil {
		return nil, fmt.Errorf("failed to parse ID Token claims: %v", err)
	}

	return payload, nil
}

func getUserInfo(accessToken string) (*userInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo userInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

type userInfo struct {
	Email string `json:"email"`
}

type tokenPayload struct {
	Issuer   string `json:"iss"`
	Subject  string `json:"sub"`
	Audience string `json:"aud"`
	Expiry   uint64 `json:"exp"`
	IssuedAt uint64 `json:"iat"`
	Email    string `json:"email"`
}
