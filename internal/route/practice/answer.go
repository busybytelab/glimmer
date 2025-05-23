package practice

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/busybytelab.com/glimmer/internal/domain"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/rs/zerolog/log"
)

type (
	AnswerRoute interface {
		HandleEvaluateAnswer(e *core.RequestEvent) error
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
		// Unmarshal the correct answer from JSON string
		var correctAnswerStr string
		if err := json.Unmarshal([]byte(correctAnswer), &correctAnswerStr); err != nil {
			log.Error().Err(err).Str("correctAnswer", correctAnswer).Msg("Failed to unmarshal correct answer")
			return e.InternalServerError("Invalid correct answer format", err)
		}

		normalizedUserAnswer := strings.TrimSpace(strings.ToLower(req.UserAnswer))
		normalizedCorrectAnswer := strings.TrimSpace(strings.ToLower(correctAnswerStr))
		isCorrect = normalizedUserAnswer == normalizedCorrectAnswer
	}

	log.Info().
		Str("practiceItemId", req.PracticeItemId).
		Str("userAnswer", req.UserAnswer).
		Str("correctAnswer", correctAnswer).
		Bool("isCorrect", isCorrect).
		Msg("Evaluated answer")

	// 5. Return response
	return e.JSON(http.StatusOK, AnswerEvaluationResponse{
		IsCorrect: isCorrect,
	})
}
