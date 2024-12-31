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
		return
	}

	host := cfg.Section("ZTE").Key("host").String()
	password := cfg.Section("ZTE").Key("password").String()
	statusFilename := cfg.Section("ETC").Key("reboot_status_file_path").String()

	zte := pkg.NewZTEL13(host, password)
	logger := pkg.NewLogger()

	loggedIn, err := zte.Login()
	if err != nil {
		logger.LogAndWriteStatus(statusFilename, fmt.Sprintf("Login error: %v", err))
		return
	}
	if !loggedIn {
		logger.LogAndWriteStatus(statusFilename, fmt.Sprintf("Login failed: %v", err))
		return
	}

	if err = zte.Reboot(); err != nil {
		logger.LogAndWriteStatus(statusFilename, fmt.Sprintf("Reboot error: %v", err))
		return
	}

	message := "Reboot 5G Router successfully."
	logger.LogAndWriteStatus(statusFilename, message)
}
