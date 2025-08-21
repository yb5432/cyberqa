// Package handlers implements the HTTP handlers for the Cyber Q&A API.
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"openai-api/pkg/database"
	"openai-api/pkg/models"
	"openai-api/pkg/openai"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// SubmitUserARequest represents the request body for submitting User A's answers.
type SubmitUserARequest struct {
	Answers      map[string]string `json:"answers"`
	ShareAnswers bool              `json:"shareAnswers"`
}

// SubmitUserAResponse represents the response body for submitting User A's answers.
type SubmitUserAResponse struct {
	Token string `json:"token"`
}

// SubmitUserBRequest represents the request body for submitting User B's answers.
type SubmitUserBRequest struct {
	Token        string            `json:"token"`
	Answers      map[string]string `json:"answers"`
	ShareAnswers bool              `json:"shareAnswers"`
}

// SubmitUserBResponse represents the response body for submitting User B's answers.
type SubmitUserBResponse struct {
	Success bool `json:"success"`
}

// ResultsResponse represents the response body for getting results.
type ResultsResponse struct {
	Compatibility int               `json:"compatibility"`
	Summary       string            `json:"summary"`
	UserAShared   bool              `json:"userAShared"`
	UserBShared   bool              `json:"userBShared"`
	UserAAnswers  map[string]string `json:"userAAnswers"`
	UserBAnswers  map[string]string `json:"userBAnswers"`
}

// OpenAIResponse represents the expected response structure from OpenAI.
type OpenAIResponse struct {
	Compatibility int    `json:"compatibility"`
	Summary       string `json:"summary"`
}

// QuestionUploadRequest represents the request body for uploading questions.
type QuestionUploadRequest []QuestionItem

// QuestionItem represents a single question item in the upload request.
type QuestionItem struct {
	ID               int      `json:"id"`
	QuestionText     string   `json:"question"`
	IsMultipleChoice bool     `json:"isMultipleChoice"`
	Options          []string `json:"options"`
}

// QuestionResponse represents a question in the response.
type QuestionResponse struct {
	ID               uint     `json:"id"`
	QuestionText     string   `json:"question"`
	IsMultipleChoice bool     `json:"isMultipleChoice"`
	Options          []string `json:"options"`
}

// SubmitUserA handles the POST /api/submit-user-a endpoint.
func SubmitUserA(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req SubmitUserARequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Generate a unique token (in a real application, you might want to use a more robust method)
	token := generateToken()

	// Convert answers to JSON string for storage
	answersJSON, err := json.Marshal(req.Answers)
	if err != nil {
		http.Error(w, "Failed to process answers", http.StatusInternalServerError)
		return
	}

	// Create UserA record
	userA := models.UserA{
		Token:        token,
		Answers:      string(answersJSON),
		ShareAnswers: req.ShareAnswers,
	}

	// Save to database
	if err := database.DB.Create(&userA).Error; err != nil {
		http.Error(w, "Failed to save user A data", http.StatusInternalServerError)
		return
	}

	// Create session record
	session := models.Session{
		Token:   token,
		UserAID: userA.ID,
	}
	if err := database.DB.Create(&session).Error; err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Return response
	response := SubmitUserAResponse{
		Token: token,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SubmitUserB handles the POST /api/submit-user-b endpoint.
func SubmitUserB(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req SubmitUserBRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Find session by token
	var session models.Session
	if err := database.DB.Preload("UserA").Where("token = ?", req.Token).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Invalid token", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to find session", http.StatusInternalServerError)
		return
	}

	// Convert answers to JSON string for storage
	answersJSON, err := json.Marshal(req.Answers)
	if err != nil {
		http.Error(w, "Failed to process answers", http.StatusInternalServerError)
		return
	}

	// Create UserB record
	userB := models.UserB{
		Token:        req.Token,
		Answers:      string(answersJSON),
		ShareAnswers: req.ShareAnswers,
	}

	// Save to database
	if err := database.DB.Create(&userB).Error; err != nil {
		http.Error(w, "Failed to save user B data", http.StatusInternalServerError)
		return
	}

	// Update session with UserB reference
	session.UserBID = &userB.ID
	if err := database.DB.Save(&session).Error; err != nil {
		http.Error(w, "Failed to update session", http.StatusInternalServerError)
		return
	}

	// Generate compatibility score and summary using OpenAI
	compatibility, summary, err := generateCompatibilityScore(session.UserA, userB)
	if err != nil {
		log.Printf("Failed to generate compatibility score: %v", err)
		// Continue without AI-generated content
		compatibility = 85           // Default value
		summary = "你们的答案很相似，有很好的默契！" // Default summary in Chinese
	}

	// Update session with compatibility score and summary
	session.Compatibility = compatibility
	session.Summary = summary
	if err := database.DB.Save(&session).Error; err != nil {
		log.Printf("Failed to save compatibility score: %v", err)
		// Continue without saving AI-generated content
	}

	// Return response
	response := SubmitUserBResponse{
		Success: true,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetResults handles the GET /api/results/{token} endpoint.
func GetResults(w http.ResponseWriter, r *http.Request) {
	// Get token from URL parameters
	vars := mux.Vars(r)
	token := vars["token"]

	// Find session by token
	var session models.Session
	if err := database.DB.Preload("UserA").Preload("UserB").Where("token = ?", token).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Invalid token", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to find session", http.StatusInternalServerError)
		return
	}

	// Check if UserB has submitted answers
	if session.UserB == nil {
		http.Error(w, "User B has not submitted answers yet", http.StatusNotFound)
		return
	}

	// Parse UserA answers
	var userAAnswers map[string]string
	if err := json.Unmarshal([]byte(session.UserA.Answers), &userAAnswers); err != nil {
		http.Error(w, "Failed to parse User A answers", http.StatusInternalServerError)
		return
	}

	// Parse UserB answers
	var userBAnswers map[string]string
	if err := json.Unmarshal([]byte(session.UserB.Answers), &userBAnswers); err != nil {
		http.Error(w, "Failed to parse User B answers", http.StatusInternalServerError)
		return
	}

	// Return response
	response := ResultsResponse{
		Compatibility: session.Compatibility,
		Summary:       session.Summary,
		UserAShared:   session.UserA.ShareAnswers,
		UserBShared:   session.UserB.ShareAnswers,
		UserAAnswers:  userAAnswers,
		UserBAnswers:  userBAnswers,
	}
	if !response.UserAShared {
		response.UserAAnswers = nil
	}
	if !response.UserBShared {
		response.UserBAnswers = nil
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// generateToken generates a simple token for identifying sessions.
func generateToken() string {
	// In a real application, you would use a more secure method
	// This is just a simple example
	return uuid.NewV4().String()
}

// generateCompatibilityScore generates a compatibility score and summary using OpenAI.
func generateCompatibilityScore(userA models.UserA, userB models.UserB) (int, string, error) {
	// Parse UserA answers
	var userAAnswers map[string]string
	if err := json.Unmarshal([]byte(userA.Answers), &userAAnswers); err != nil {
		return 0, "", err
	}

	// Parse UserB answers
	var userBAnswers map[string]string
	if err := json.Unmarshal([]byte(userB.Answers), &userBAnswers); err != nil {
		return 0, "", err
	}

	// Create a map with both users' answers
	answers := map[string]map[string]string{
		"userA": userAAnswers,
		"userB": userBAnswers,
	}

	// Convert answers to JSON string
	answersJSON, err := json.Marshal(answers)
	if err != nil {
		return 0, "", err
	}

	// Get system prompt from environment variable
	systemPrompt := os.Getenv("SYSTEM_PROMPT")
	if systemPrompt == "" {
		systemPrompt = openai.SystemPrompt
	}

	// Create OpenAI client
	config := openai.NewConfig().
		WithAPIKey(os.Getenv("OPENAI_API_KEY")).
		WithAPIBase(os.Getenv("OPENAI_API_BASE")).
		WithSystemPrompt(systemPrompt)

	client := openai.NewClient(config)

	llmModel := os.Getenv("MODELS")

	if llmModel == "" {
		llmModel = "gpt-3.5-turbo"
	}
	// Prepare the request
	request := &openai.ChatCompletionRequest{
		Model: llmModel, // Default model, can be overridden
		Messages: []openai.Message{
			{
				Role:    "user",
				Content: string(answersJSON),
			},
		},
		Temperature: 0.7,
		Stream:      false,
		MaxTokens:   4096,
	}

	// Send request to OpenAI
	ctx := context.Background()
	response, err := client.ChatCompletion(ctx, request)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get response from OpenAI: %w", err)
	}

	// fmt.Printf("OpenAI response: %v\n", response)

	// Check if we have a response
	if len(response.Choices) == 0 {
		return 0, "", fmt.Errorf("no choices in OpenAI response")
	}

	// Extract the content from the response
	content := response.Choices[0].Message.Content

	fmt.Printf("content: %v\n", content)

	// Try to parse the response as JSON
	var openAIResponse OpenAIResponse
	if err := json.Unmarshal([]byte(content), &openAIResponse); err != nil {
		fmt.Printf("Failed to parse OpenAI response as JSON: %v\n", err)
		// If JSON parsing fails, use the content as the summary and set a default compatibility
		return 85, content, nil
	}

	// Validate compatibility score
	if openAIResponse.Compatibility < 0 || openAIResponse.Compatibility > 100 {
		// If the compatibility score is out of range, use a default value
		openAIResponse.Compatibility = 85
	}

	return openAIResponse.Compatibility, openAIResponse.Summary, nil
}

// UploadQuestions handles the POST /api/questions/upload endpoint.
func UploadQuestions(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req QuestionUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("Error decoding request body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Delete all existing questions (delete-then-add approach)
	if err := database.DB.Where("1 = 1").Delete(&models.Question{}).Error; err != nil {
		http.Error(w, "Failed to delete existing questions", http.StatusInternalServerError)
		return
	}

	// Convert request items to database models
	var questions []models.Question
	for _, item := range req {
		// Convert options array to JSON string
		optionsJSON, err := json.Marshal(item.Options)
		if err != nil {
			http.Error(w, "Failed to process question options", http.StatusInternalServerError)
			return
		}

		question := models.Question{
			ID:               uint(item.ID),
			QuestionText:     item.QuestionText,
			IsMultipleChoice: item.IsMultipleChoice,
			Options:          string(optionsJSON),
		}
		questions = append(questions, question)
	}

	// Save new questions to database
	if len(questions) > 0 {
		if err := database.DB.Create(&questions).Error; err != nil {
			http.Error(w, "Failed to save questions", http.StatusInternalServerError)
			return
		}
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Questions uploaded successfully"})
}

// GetQuestions handles the GET /api/questions endpoint.
func GetQuestions(w http.ResponseWriter, r *http.Request) {
	// Get all questions from database
	var dbQuestions []models.Question
	if err := database.DB.Find(&dbQuestions).Error; err != nil {
		http.Error(w, "Failed to retrieve questions", http.StatusInternalServerError)
		return
	}

	// Convert database models to response format
	var questions []QuestionResponse
	for _, dbQuestion := range dbQuestions {
		// Parse options JSON string
		var options []string
		if err := json.Unmarshal([]byte(dbQuestion.Options), &options); err != nil {
			// If parsing fails, use empty array
			options = []string{}
		}

		question := QuestionResponse{
			ID:               dbQuestion.ID,
			QuestionText:     dbQuestion.QuestionText,
			IsMultipleChoice: dbQuestion.IsMultipleChoice,
			Options:          options,
		}
		questions = append(questions, question)
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}
