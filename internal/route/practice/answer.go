package practice

import (
	"encoding/json"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/busybytelab.com/glimmer/internal/domain"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/rs/zerolog/log"
)

type (
	AnswerRoute interface {
		HandleEvaluateAnswer(e *core.RequestEvent) error
		HandleProcessAnswer(e *core.RequestEvent) error
	}

	answerRoute struct{}

	// AnswerEvaluationRequest defines the request body for answer evaluation
	AnswerEvaluationRequest struct {
		PracticeItemId string `json:"practiceItemId"`
		UserAnswer     string `json:"userAnswer"`
	}

	// AnswerEvaluationResponse defines the response for answer evaluation
	AnswerEvaluationResponse struct {
		IsCorrect bool `json:"isCorrect"`
	}

	// ProcessAnswerRequest defines the request body for processing answer
	ProcessAnswerRequest struct {
		PracticeItemId   string `json:"practiceItemId"`
		UserAnswer       string `json:"userAnswer"`
		PracticeSession  string `json:"practiceSession"`
		LearnerId        string `json:"learnerId"`
		HintLevelReached int    `json:"hintLevelReached"`
	}

	// ProcessAnswerResponse defines the response for processing answer
	ProcessAnswerResponse struct {
		IsCorrect        bool    `json:"isCorrect"`
		Score            float64 `json:"score"`
		Feedback         string  `json:"feedback"`
		HintLevelReached int     `json:"hintLevelReached"`
		AttemptNumber    int     `json:"attemptNumber"`
	}
)

func NewAnswerRoute() AnswerRoute {
	return &answerRoute{}
}

func (r *answerRoute) HandleEvaluateAnswer(e *core.RequestEvent) error {
	// Auth check
	if e.Auth == nil {
		return apis.NewUnauthorizedError("You must be logged in", nil)
	}

	// 1. Parse request JSON body
	var req AnswerEvaluationRequest
	if err := e.BindBody(&req); err != nil {
		return e.BadRequestError("Invalid request body", err)
	}

	// Validate request
	if req.PracticeItemId == "" {
		return e.BadRequestError("PracticeItemId is required", nil)
	}

	// 2. Load practice item from DB
	practiceItem, err := e.App.FindRecordById(domain.CollectionPracticeItems, req.PracticeItemId)
	if err != nil {
		log.Error().Err(err).Str("practiceItemId", req.PracticeItemId).Msg("Failed to find practice item")
		return e.NotFoundError("Practice item not found", err)
	}

	// 3. Get correct answer from practice item
	correctAnswer := practiceItem.GetString("correct_answer")
	if correctAnswer == "" || correctAnswer == "null" {
		return e.BadRequestError("Practice item has no correct answer", nil)
	}

	// 4. Evaluate answer
	isCorrect := false
	if req.UserAnswer != "" {
		correctAnswerStr := getCleanCorrectAnswer(correctAnswer)
		isCorrect = checkAnswer(practiceItem, req.UserAnswer, correctAnswerStr)
	}

	// 5. Return response
	return e.JSON(http.StatusOK, AnswerEvaluationResponse{
		IsCorrect: isCorrect,
	})
}

func (r *answerRoute) HandleProcessAnswer(e *core.RequestEvent) error {
	// Auth check
	if e.Auth == nil {
		return apis.NewUnauthorizedError("You must be logged in", nil)
	}

	// 1. Parse request JSON body
	var req ProcessAnswerRequest
	if err := e.BindBody(&req); err != nil {
		return e.BadRequestError("Invalid request body", err)
	}

	// Validate request
	if req.PracticeItemId == "" {
		return e.BadRequestError("PracticeItemId is required", nil)
	}
	if req.PracticeSession == "" {
		return e.BadRequestError("PracticeSession is required", nil)
	}
	if req.LearnerId == "" {
		return e.BadRequestError("LearnerId is required", nil)
	}

	// 2. Load practice item from DB
	practiceItem, err := e.App.FindRecordById(domain.CollectionPracticeItems, req.PracticeItemId)
	if err != nil {
		log.Error().Err(err).Str("practiceItemId", req.PracticeItemId).Msg("Failed to find practice item")
		return e.NotFoundError("Practice item not found", err)
	}

	// 3. Get correct answer from practice item
	correctAnswer := practiceItem.GetString("correct_answer")
	if correctAnswer == "" || correctAnswer == "null" {
		return e.BadRequestError("Practice item has no correct answer", nil)
	}

	// 4. Evaluate answer correctness
	isCorrect := false
	if req.UserAnswer != "" {
		correctAnswerStr := getCleanCorrectAnswer(correctAnswer)
		isCorrect = checkAnswer(practiceItem, req.UserAnswer, correctAnswerStr)
	}

	// 5. Calculate score based on correctness and hint usage
	score := calculateScore(isCorrect, req.HintLevelReached, practiceItem)

	// 6. Generate feedback
	feedback := generateFeedback(isCorrect, req.HintLevelReached, score)

	// 7. Get existing result or create new one
	existingResult, err := getLatestResult(e.App, req.PracticeItemId, req.PracticeSession)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get latest result")
		return e.InternalServerError("Failed to get existing result", err)
	}

	now := time.Now().Format(time.RFC3339)
	var attemptNumber int

	if existingResult != nil {
		// Update existing result
		attemptNumber = existingResult.GetInt("attempt_number") + 1

		existingResult.Set("answer", req.UserAnswer)
		existingResult.Set("is_correct", isCorrect)
		existingResult.Set("score", score)
		existingResult.Set("feedback", feedback)
		existingResult.Set("hint_level_reached", req.HintLevelReached)
		existingResult.Set("attempt_number", attemptNumber)
		existingResult.Set("submitted_at", now)

		if err := e.App.Save(existingResult); err != nil {
			log.Error().Err(err).Msg("Failed to update practice result")
			return e.InternalServerError("Failed to update result", err)
		}
	} else {
		// Create new result
		attemptNumber = 1

		collection, err := e.App.FindCollectionByNameOrId(domain.CollectionPracticeResults)
		if err != nil {
			log.Error().Err(err).Msg("Failed to find practice_results collection")
			return e.InternalServerError("Failed to find results collection", err)
		}

		newResult := core.NewRecord(collection)
		newResult.Set("practice_item", req.PracticeItemId)
		newResult.Set("practice_session", req.PracticeSession)
		newResult.Set("learner", req.LearnerId)
		newResult.Set("answer", req.UserAnswer)
		newResult.Set("is_correct", isCorrect)
		newResult.Set("score", score)
		newResult.Set("feedback", feedback)
		newResult.Set("hint_level_reached", req.HintLevelReached)
		newResult.Set("attempt_number", attemptNumber)
		newResult.Set("started_at", now)
		newResult.Set("submitted_at", now)

		if err := e.App.Save(newResult); err != nil {
			log.Error().Err(err).Msg("Failed to create practice result")
			return e.InternalServerError("Failed to create result", err)
		}
	}

	log.Info().
		Str("practiceItemId", req.PracticeItemId).
		Str("userAnswer", req.UserAnswer).
		Bool("isCorrect", isCorrect).
		Float64("score", score).
		Int("hintLevel", req.HintLevelReached).
		Int("attempt", attemptNumber).
		Msg("Answer processed")

	// 8. Return response
	return e.JSON(http.StatusOK, ProcessAnswerResponse{
		IsCorrect:        isCorrect,
		Score:            score,
		Feedback:         feedback,
		HintLevelReached: req.HintLevelReached,
		AttemptNumber:    attemptNumber,
	})
}

// normalizeAnswerString trims leading/trailing whitespace and backticks (single or triple).
func normalizeAnswerString(s string) string {
	s = strings.TrimSpace(s)
	// Handle triple backticks first as they are more specific
	if len(s) >= 6 && strings.HasPrefix(s, "```") && strings.HasSuffix(s, "```") {
		s = s[3 : len(s)-3]
	} else if len(s) >= 2 && strings.HasPrefix(s, "`") && strings.HasSuffix(s, "`") {
		s = s[1 : len(s)-1]
	}
	return strings.TrimSpace(s)
}

func checkAnswer(item *core.Record, userAnswer string, correctAnswer string) bool {
	normalizedUserAnswer := strings.ToLower(normalizeAnswerString(userAnswer))
	normalizedCorrectAnswer := strings.ToLower(normalizeAnswerString(correctAnswer))
	correct := normalizedUserAnswer == normalizedCorrectAnswer

	log.Info().
		Str("practiceItemId", item.Id).
		Str("userAnswer", userAnswer).
		Str("correctAnswer", correctAnswer).
		Str("normalizedUserAnswer", normalizedUserAnswer).
		Str("normalizedCorrectAnswer", normalizedCorrectAnswer).
		Bool("isCorrect", correct).
		Msg("Answer evaluation")

	return correct
}

// getCleanCorrectAnswer unmarshals the correct answer, falling back to a raw string if needed.
// This handles cases where the answer is a JSON-encoded string (e.g., "\"some text\"") or a
// raw value (e.g., "1", "some text").
func getCleanCorrectAnswer(rawCorrectAnswer string) string {
	var correctAnswerStr string
	if err := json.Unmarshal([]byte(rawCorrectAnswer), &correctAnswerStr); err != nil {
		// If unmarshalling fails, it's likely not a JSON-encoded string.
		// Log a warning and use the raw value.
		log.Warn().Err(err).Str("correctAnswer", rawCorrectAnswer).Msg("Could not unmarshal correct answer as JSON, using raw value")
		return rawCorrectAnswer
	}
	return correctAnswerStr
}

// calculateScore computes the score based on correctness and hint usage
// Max score is 1.0, min score is 0.0
func calculateScore(isCorrect bool, hintLevelReached int, practiceItem *core.Record) float64 {
	if !isCorrect {
		return 0.0
	}

	// Get total number of hints available
	hintsJson := practiceItem.GetString("hints")
	var hints []string
	totalHints := 0

	if hintsJson != "" && hintsJson != "null" {
		if err := json.Unmarshal([]byte(hintsJson), &hints); err == nil {
			totalHints = len(hints)
		}
	}

	// If no hints are available or none were used, full score
	if totalHints == 0 || hintLevelReached == 0 {
		return 1.0
	}

	// Calculate score reduction based on hint usage
	// Each hint used reduces the score by (1.0 / (totalHints + 1))
	// This ensures that using all hints still gives some score (> 0)
	scoreReduction := float64(hintLevelReached) / float64(totalHints+1)
	score := 1.0 - scoreReduction

	// Ensure score doesn't go below a minimum threshold (0.1)
	if score < 0.1 {
		score = 0.1
	}

	// Round to 0 decimal places
	return math.Round(score)
}

// generateFeedback creates appropriate feedback based on performance
func generateFeedback(isCorrect bool, hintLevelReached int, score float64) string {
	if !isCorrect {
		if hintLevelReached > 0 {
			return "Not quite right. Try reviewing the hints and give it another shot!"
		}
		return "That's not correct. Consider using the hints for guidance."
	}

	// Correct answer feedback
	if hintLevelReached == 0 {
		return "Excellent! You got it right on your own!"
	} else if score >= 0.8 {
		return "Great work! You used hints wisely and got the correct answer."
	} else if score >= 0.5 {
		return "Good job! You found the right answer with some help from the hints."
	} else {
		return "Correct! Keep practicing to build your confidence without using as many hints."
	}
}

// getLatestResult retrieves the most recent practice result for a given item and session
func getLatestResult(app core.App, practiceItemId, sessionId string) (*core.Record, error) {
	results, err := app.FindRecordsByFilter(
		domain.CollectionPracticeResults,
		"practice_item = {:practiceItem} && practice_session = {:session}",
		"-created",
		1,
		0,
		map[string]any{
			"practiceItem": practiceItemId,
			"session":      sessionId,
		},
	)

	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	return results[0], nil
}
