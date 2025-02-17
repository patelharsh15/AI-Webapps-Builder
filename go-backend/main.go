// main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"go-backend/constants"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

type TemplateResponse struct {
	Prompts   []string `json:"prompts"`
	UIPrompts []string `json:"uiPrompts"`
}

type ChatRequest struct {
	Messages []Message `json:"messages"`
}

type ChatResponse struct {
	Response string `json:"response"`
}

var (
	anthropicClient *AnthropicClient
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		log.Fatal("ANTHROPIC_API_KEY environment variable is required")
	}

	anthropicClient = NewAnthropicClient(apiKey)

	r := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST", "GET", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	})

	r.HandleFunc("/template", handleTemplate).Methods("POST")
	r.HandleFunc("/chat", handleChat).Methods("POST")

	handler := c.Handler(r)

	log.Println("Server starting on :3000")
	log.Fatal(http.ListenAndServe(":3000", handler))
}

func handleTemplate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Prompt string `json:"prompt"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Prompt == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	responseChan := make(chan string, 3)
	errorChan := make(chan error, 3)

	go func() {
		response, err := anthropicClient.GetTemplateType(req.Prompt)
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- response
	}()

	select {
	case err := <-errorChan:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	case answer := <-responseChan:
		var resp TemplateResponse

		switch answer {
		case "react":
			resp = TemplateResponse{
				Prompts:   []string{constants.BasePrompt, constants.ReactProjectPrompt},
				UIPrompts: []string{constants.ReactBasePrompt},
			}
		case "node":
			resp = TemplateResponse{
				Prompts:   []string{constants.NodeProjectPrompt},
				UIPrompts: []string{constants.NodeBasePrompt},
			}
		default:
			http.Error(w, "You can't access this", http.StatusForbidden)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	responseChan := make(chan string, 3)
	errorChan := make(chan error, 3)

	go func() {
		response, err := anthropicClient.Chat(req.Messages)
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- response
	}()

	select {
	case err := <-errorChan:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	case response := <-responseChan:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ChatResponse{Response: response})
	}
}
