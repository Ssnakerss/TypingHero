package web

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
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

// Text samples organized by difficulty (1-10)
var textSamples = map[int][]string{
	1:  {"cat", "dog", "sun", "run", "the", "and", "big", "red"},
	2:  {"The cat runs.", "A dog barks.", "The sun shines.", "I like apples.", "She reads books."},
	3:  {"The quick brown fox jumps over the lazy dog.", "She sells seashells by the seashore.", "How much wood would a woodchuck chuck?"},
	4:  {"The weather today is pleasant and sunny.", "Walking in the park is very relaxing.", "Learning a new language takes time and practice."},
	5:  {"Programming is both an art and a science that requires logical thinking.", "The ancient castle stood on a hill overlooking the peaceful village below."},
	6:  {"Software development involves writing code, testing applications, and fixing bugs.", "The scientist carefully conducted the experiment to test her groundbreaking hypothesis."},
	7:  {"Machine learning algorithms can analyze large amounts of data to identify patterns and make predictions.", "The entrepreneur founded a successful startup that revolutionized the technology industry."},
	8:  {"Cryptocurrency transactions are recorded on a decentralized blockchain ledger that ensures transparency and security.", "The philosophical debate about artificial intelligence consciousness continues to divide experts in the field."},
	9:  {"Quantum computing leverages the principles of quantum mechanics to perform calculations at unprecedented speeds.", "The interdisciplinary research combines advances in neuroscience, computer science, and cognitive psychology."},
	10: {"The implementation of distributed systems requires careful consideration of consistency, availability, and partition tolerance according to CAP theorem.", "Advanced natural language processing models utilize transformer architectures with attention mechanisms to achieve state-of-the-art results."},
}

func generateText(difficulty int) string {
	if difficulty < 1 {
		difficulty = 1
	}
	if difficulty > 10 {
		difficulty = 10
	}

	samples := textSamples[difficulty]
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

func StartWeb() {
	rand.Seed(time.Now().UnixNano())

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
}
