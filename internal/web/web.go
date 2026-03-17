package web

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/Ssnakerss/TypingHero/internal/models"
	"github.com/Ssnakerss/TypingHero/internal/web/router"
)

type TextRequest struct {
	Difficulty int `json:"difficulty"`
}

type ResultRequest struct {
	OriginalText string `json:"originalText"`
	UserInput    string `json:"userInput"`
	TimeTakenSec int    `json:"timeTakenSec"`
}

type TextResponse struct {
	Text string `json:"text"`
}

type ResultResponse struct {
	Errors    int     `json:"errors"`
	WPM       float64 `json:"wpm"`
	Accuracy  float64 `json:"accuracy"`
	CharCount int     `json:"charCount"`
	TimeTaken int     `json:"timeTaken"`
}

func generateText(difficulty int) string {
	if difficulty < 1 {
		difficulty = 1
	}
	if difficulty > 10 {
		difficulty = 10
	}

	samples := models.TextPools[difficulty]
	selected := samples[rand.Intn(len(samples))]

	// For higher difficulties, combine multiple sentences
	if difficulty >= 5 {
		count := difficulty / 2
		var result strings.Builder
		result.WriteString(selected)
		for i := 1; i < count; i++ {
			result.WriteString(" ")
			result.WriteString(samples[rand.Intn(len(samples))])
		}
		return result.String()
	}

	return selected
}

func getTextHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	text := generateText(req.Difficulty)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TextResponse{Text: text})
}

func calculateResultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Calculate errors (character-by-character comparison)
	errors := 0
	original := strings.TrimSpace(req.OriginalText)
	userInput := strings.TrimSpace(req.UserInput)

	minLen := len(original)
	if len(userInput) < minLen {
		minLen = len(userInput)
	}

	for i := 0; i < minLen; i++ {
		if original[i] != userInput[i] {
			errors++
		}
	}

	// Count missing or extra characters
	errors += abs(len(original) - len(userInput))

	// Calculate accuracy
	charCount := len(original)
	var accuracy float64 = 100
	if charCount > 0 {
		accuracy = float64(charCount-errors) / float64(charCount) * 100
		if accuracy < 0 {
			accuracy = 0
		}
	}

	// Calculate WPM (words per minute)
	// Standard: 5 characters = 1 word
	words := float64(charCount) / 5.0
	var wpm float64 = 0
	if req.TimeTakenSec > 0 {
		wpm = words / float64(req.TimeTakenSec) * 60
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ResultResponse{
		Errors:    errors,
		WPM:       wpm,
		Accuracy:  accuracy,
		CharCount: charCount,
		TimeTaken: req.TimeTakenSec,
	})
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func StartWeb(ctx context.Context, db models.Storage) {

	r := router.New()

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// API routes
	http.HandleFunc("/api/text", getTextHandler)
	http.HandleFunc("/api/result", calculateResultHandler)

	fmt.Println("Strating File Server...")
	// Serve static files
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	fmt.Println("Typing Trainer server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
	<-ctx.Done()

}
