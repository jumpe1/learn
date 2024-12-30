package main

import (
	"fmt"
	"l13/pkg"
	"log"

	"github.com/go-ini/ini"
)

func main() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalf("Failed to read config.ini: %v", err)
	}

	host := cfg.Section("ZTE").Key("host").String()
	password := cfg.Section("ZTE").Key("password").String()

	zte := pkg.NewZTEL13(host, password)

	loggedIn, err := zte.Login()
	if err != nil {
		log.Fatalf("Login error: %v", err)
	}
	if !loggedIn {
		fmt.Println("Login failed")
		return
	}

	err = zte.Reboot()
	if err != nil {
		log.Fatalf("Reboot error: %v", err)
	}

	fmt.Println("Reboot command sent successfully.")
}
