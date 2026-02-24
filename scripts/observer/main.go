package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gorilla/websocket"
)

const (
	baseURL     = "http://localhost:8080"
	wsURL       = "ws://localhost:8080/ws"
	outFile     = "chat_sample.md"
	maxMessages = 1000
)

func main() {
	log.Println("Starting Chat Observer...")

	// 1. Create observer user
	username := fmt.Sprintf("observer_%d", time.Now().Unix())
	password := "observer123"

	log.Printf("Registering observer user: %s", username)
	registerPayload := map[string]string{
		"user_id":      username,
		"password":     password,
		"display_name": "K6 Observer",
	}
	body, _ := json.Marshal(registerPayload)
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Registration failed: %v", err)
	}
	resp.Body.Close()

	// 2. Login to get token
	log.Println("Logging in to retrieve JWT...")
	loginPayload := map[string]string{
		"user_id":  username,
		"password": password,
	}
	body, _ = json.Marshal(loginPayload)
	req, _ = http.NewRequest(http.MethodPost, baseURL+"/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatalf("Login failed: %v", err)
	}

	// Extract cookies
	var token string
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "access_token" {
			token = cookie.Value
			break
		}
	}
	resp.Body.Close()
	
	if token == "" {
		log.Fatal("No access_token cookie found")
	}

	// 3. Connect to WebSocket
	wsURI := fmt.Sprintf("%s?token=%s", wsURL, token)
	log.Printf("Connecting to WebSocket: %s", wsURL)
	c, _, err := websocket.DefaultDialer.Dial(wsURI, nil)
	if err != nil {
		log.Fatalf("WebSocket connection failed: %v", err)
	}
	defer c.Close()

	// Prepare output file
	f, err := os.Create(outFile)
	if err != nil {
		log.Fatalf("Could not create output file: %v", err)
	}
	defer f.Close()

	header := fmt.Sprintf("# STORM Day - Chat Interception Sample\n\n**Date:** %s\n**Description:** Cet échantillon capture les %d premiers messages échangés sur le système générés pendant le tir de charge k6.\n\n---\n\n", time.Now().Format(time.RFC1123), maxMessages)
	f.WriteString(header)

	msgCount := 0
	log.Printf("Listening for messages (will capture up to %d messages)...", maxMessages)

	// Clean input to detect valid chat messages
	contentRe := regexp.MustCompile(`"content":"([^"]+)"`)
	
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		// Look for content matching our k6 virtual users output
		matches := contentRe.FindSubmatch(message)
		if len(matches) > 1 {
			msgCount++
			content := string(matches[1])
			
			// Extract User ID from the parsed message (format: "Hello from user_xxx_VU at TIME")
			logLine := fmt.Sprintf("- `%s` - **Message:** %s\n", time.Now().Format("15:04:05.000"), content)
			f.WriteString(logLine)
			
			if msgCount%100 == 0 {
				log.Printf("Captured %d messages...", msgCount)
			}
			
			if msgCount >= maxMessages {
				log.Printf("Reached %d messages limit. Saving report.", maxMessages)
				break
			}
		}
	}

	f.WriteString("\n---\n*Observation terminée.*")
	log.Printf("Report saved to %s", outFile)
}
