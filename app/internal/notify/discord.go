package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func Notify(message string) {
	host, _ := os.Hostname()
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")

	log.Printf("🔔 [%s] %s\n", host, message)

	if webhookURL == "" {
		log.Println("DISCORD_WEBHOOK_URL not defined")
		return
	}

	SendWebhook(webhookURL, map[string]interface{}{
		"content": fmt.Sprintf("[%s] %s", host, message),
	})
}

func SendWebhook(url string, payload interface{}) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("erro ao converter payload: %v", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("falha ao disparar webhook: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook retornou status inesperado: %d", resp.StatusCode)
	}

	return nil
}