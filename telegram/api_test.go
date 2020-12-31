package telegram

import (
	"log"
	"testing"
)

func TestNewTelegramBot(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	apiToken := "paste your own token"
	NewTelegramBot(apiToken)
}
