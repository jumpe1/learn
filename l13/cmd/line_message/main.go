package main

import (
	"l13/pkg"
	"log"
	"os"
	"strings"

	"github.com/go-ini/ini"
)

func main() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalf("Failed to read config.ini: %v", err)
	}

	channelToken := cfg.Section("LINE").Key("channel_token").String()
	userID := cfg.Section("LINE").Key("user_id").String()
	statusFilename := cfg.Section("ETC").Key("reboot_status_file_path").String()

	content, err := os.ReadFile(statusFilename)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
		return
	}

	message := strings.TrimSpace(string(content))
	lineMessageAPI := pkg.NewLineMessageAPI(channelToken, userID)
	if err = lineMessageAPI.Send(message); err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
}
