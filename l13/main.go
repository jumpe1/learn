package main

import (
	"l13/pkg"
	"log"

	"github.com/go-ini/ini"
)

func main() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalf("Failed to read config.ini: %v", err)
	}

	channelToken := cfg.Section("LINE").Key("channel_token").String()
	userID := cfg.Section("LINE").Key("user_id").String()
	host := cfg.Section("ZTE").Key("host").String()
	password := cfg.Section("ZTE").Key("password").String()

	zte := pkg.NewZTEL13(host, password)

	loggedIn, err := zte.Login()
	if err != nil {
		log.Fatalf("Login error: %v", err)
		return
	}
	if !loggedIn {
		log.Fatal("Login failed")
		return
	}

	if err = zte.Reboot(); err != nil {
		log.Fatalf("Reboot error: %v", err)
		return
	}

	message := "Reboot 5G Rooter successfully."
	log.Println(message)
	lineMessageAPI := pkg.NewLineMessageAPI(channelToken, userID)
	if err = lineMessageAPI.Send(message); err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
}
