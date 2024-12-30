package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

type DictToClass map[string]interface{}

type ZTEL13 struct {
	host     string
	password string
	client   *http.Client
}

func NewZTEL13(host, password string) *ZTEL13 {
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		panic(err)
	}

	return &ZTEL13{
		host:     host,
		password: password,
		client: &http.Client{
			Jar: jar,
		},
	}
}

func (z *ZTEL13) getCmdProcess(cmd string, withTs bool) (DictToClass, error) {
	baseURL := fmt.Sprintf("http://%s/goform/goform_get_cmd_process", z.host)
	q := url.Values{}
	q.Set("isTest", "false")
	q.Set("cmd", cmd)
	if withTs {
		q.Set("_", fmt.Sprintf("%d", jqueryNow()))
	}
	if strings.Contains(cmd, ",") {
		q.Set("multi_data", "1")
	}

	reqURL := baseURL + "?" + q.Encode()
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", fmt.Sprintf("http://%s/", z.host))
	req.Header.Set("Sec-GPC", "1")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := z.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result DictToClass
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (z *ZTEL13) setCmdProcess(cmd string, params map[string]string) (DictToClass, error) {
	baseURL := fmt.Sprintf("http://%s/goform/goform_set_cmd_process", z.host)

	form := url.Values{}
	form.Set("isTest", "false")
	form.Set("goformId", cmd)
	for k, v := range params {
		form.Set(k, v)
	}

	req, err := http.NewRequest("POST", baseURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", fmt.Sprintf("http://%s", z.host))
	req.Header.Set("Referer", fmt.Sprintf("http://%s/", z.host))
	req.Header.Set("Sec-GPC", "1")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := z.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result DictToClass
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (z *ZTEL13) ld() (DictToClass, error) {
	return z.getCmdProcess("LD", true)
}

func (z *ZTEL13) rd() (DictToClass, error) {
	return z.getCmdProcess("wa_inner_version,cr_version,RD", false)
}

func (z *ZTEL13) Login() (bool, error) {
	ldResp, err := z.ld()
	if err != nil {
		return false, err
	}

	ldVal, ok := ldResp["LD"].(string)
	if !ok {
		return false, fmt.Errorf("LD field not found or invalid in response")
	}

	p1 := passwordAlgorithmsCookie(z.password)
	pw := passwordAlgorithmsCookie(p1 + ldVal)

	res, err := z.setCmdProcess("LOGIN", map[string]string{
		"password": pw,
	})
	if err != nil {
		return false, err
	}

	if val, exists := res["result"].(string); exists && val == "0" {
		return true, nil
	}
	return false, nil
}

func (z *ZTEL13) Reboot() error {
	rdResp, err := z.rd()
	if err != nil {
		return err
	}

	waInner, ok1 := rdResp["wa_inner_version"].(string)
	crVersion, ok2 := rdResp["cr_version"].(string)
	rdVal, ok3 := rdResp["RD"].(string)

	if !ok1 || !ok2 || !ok3 {
		return fmt.Errorf("response fields not found or invalid")
	}

	a1 := hexSha256(waInner + crVersion)
	ad := hexSha256(a1 + rdVal)

	_, err = z.setCmdProcess("REBOOT_DEVICE", map[string]string{
		"AD": ad,
	})
	if err != nil {
		return err
	}

	return nil
}
