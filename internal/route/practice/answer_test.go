package practice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/busybytelab.com/glimmer/internal/domain"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/router"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestApp(t *testing.T) *tests.TestApp {
	// Create temporary directory for PocketBase data
	tmpDir, err := os.MkdirTemp("", "pb_test_*")
	require.NoError(t, err)
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	// Create pb_data directory structure
	pbDataDir := filepath.Join(tmpDir, "pb_data")
	err = os.MkdirAll(pbDataDir, 0755)
	require.NoError(t, err)

	// Get the project root directory by looking for go.mod
	projectRoot, err := findProjectRoot()
	require.NoError(t, err)

	// Set environment variables for the test
	os.Setenv("DB_DISABLE_AUTO_MIGRATE", "false")
	os.Setenv("POCKETBASE_DATA_DIR", pbDataDir)
	os.Setenv("APP_NAME", "TestApp")
	os.Setenv("APP_URL", "http://localhost")
	os.Setenv("SENDER_ADDRESS", "test@example.com")
	os.Setenv("SENDER_NAME", "Test Sender")
	defer os.Unsetenv("DB_DISABLE_AUTO_MIGRATE")
	defer os.Unsetenv("POCKETBASE_DATA_DIR")
	defer os.Unsetenv("APP_NAME")
	defer os.Unsetenv("APP_URL")
	defer os.Unsetenv("SENDER_ADDRESS")
	defer os.Unsetenv("SENDER_NAME")

	// Run migrations in the temporary directory
	cmd := exec.Command("go", "run", filepath.Join(projectRoot, "cmd/glimmer/main.go"), "migrate", "--dir", pbDataDir)
	cmd.Env = append(os.Environ(), "DB_DISABLE_AUTO_MIGRATE=false")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()

	if err != nil {
		t.Logf("Migration stdout: %s", stdout.String())
		t.Logf("Migration stderr: %s", stderr.String())
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Create test app with custom data directory
	app, err := tests.NewTestApp(pbDataDir)
	require.NoError(t, err)
	t.Cleanup(func() {
		app.Cleanup()
	})

	return app
}

// findProjectRoot finds the project root directory by looking for go.mod
func findProjectRoot() (string, error) {
	// Start from the current directory
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Keep going up until we find go.mod
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		// Go up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// We've reached the root of the filesystem
			return "", fmt.Errorf("could not find project root (go.mod)")
		}
		dir = parent
	}
}

func TestHandleEvaluateAnswer(t *testing.T) {
	// Setup PocketBase test instance with temporary data directory
	app := setupTestApp(t)

	// Create test practice item
	collection, err := app.FindCollectionByNameOrId(domain.CollectionPracticeItems)
	require.NoError(t, err)

	practiceItem := core.NewRecord(collection)
	// Store correct answer as JSON string
	correctAnswerJSON, err := json.Marshal("test answer")
	require.NoError(t, err)
	practiceItem.Set("correct_answer", string(correctAnswerJSON))
	err = app.SaveNoValidate(practiceItem)
	require.NoError(t, err)

	// Create a practice item with no correct answer
	practiceItemNoAnswer := core.NewRecord(collection)
	err = app.SaveNoValidate(practiceItemNoAnswer)
	require.NoError(t, err)

	// Create a practice item with a raw number as correct answer
	practiceItemRawNumber := core.NewRecord(collection)
	practiceItemRawNumber.Set("correct_answer", "12345")
	err = app.SaveNoValidate(practiceItemRawNumber)
	require.NoError(t, err)

	// Create a practice item with backticks in correct answer
	practiceItemBackticks := core.NewRecord(collection)
	practiceItemBackticks.Set("correct_answer", "`code`")
	err = app.SaveNoValidate(practiceItemBackticks)
	require.NoError(t, err)

	// Create a practice item with triple backticks in correct answer
	practiceItemTripleBackticks := core.NewRecord(collection)
	practiceItemTripleBackticks.Set("correct_answer", "```code block```")
	err = app.SaveNoValidate(practiceItemTripleBackticks)
	require.NoError(t, err)

	// Create a test user once for all tests
	userCollection, err := app.FindCollectionByNameOrId("users")
	require.NoError(t, err)
	user := core.NewRecord(userCollection)
	user.Set("email", "test@example.com")
	user.Set("password", "test123")
	err = app.SaveNoValidate(user)
	require.NoError(t, err)

	tests := []struct {
		name           string
		request        AnswerEvaluationRequest
		expectedStatus int
		expectedBody   AnswerEvaluationResponse
		setupAuth      bool
		errorMessage   string
	}{
		{
			name: "correct answer",
			request: AnswerEvaluationRequest{
				PracticeItemId: practiceItem.Id,
				UserAnswer:     "test answer",
			},
			expectedStatus: http.StatusOK,
			expectedBody: AnswerEvaluationResponse{
				IsCorrect: true,
			},
			setupAuth: true,
		},
		{
			name: "incorrect answer",
			request: AnswerEvaluationRequest{
				PracticeItemId: practiceItem.Id,
				UserAnswer:     "wrong answer",
			},
			expectedStatus: http.StatusOK,
			expectedBody: AnswerEvaluationResponse{
				IsCorrect: false,
			},
			setupAuth: true,
		},
		{
			name: "case insensitive match",
			request: AnswerEvaluationRequest{
				PracticeItemId: practiceItem.Id,
				UserAnswer:     "TEST ANSWER",
			},
			expectedStatus: http.StatusOK,
			expectedBody: AnswerEvaluationResponse{
				IsCorrect: true,
			},
			setupAuth: true,
		},
		{
			name: "whitespace insensitive match",
			request: AnswerEvaluationRequest{
				PracticeItemId: practiceItem.Id,
				UserAnswer:     "  test answer  ",
			},
			expectedStatus: http.StatusOK,
			expectedBody: AnswerEvaluationResponse{
				IsCorrect: true,
			},
			setupAuth: true,
		},
		{
			name: "missing practice item id",
			request: AnswerEvaluationRequest{
				UserAnswer: "test answer",
			},
			expectedStatus: http.StatusBadRequest,
			setupAuth:      true,
			errorMessage:   "PracticeItemId is required",
		},
		{
			name: "non-existent practice item",
			request: AnswerEvaluationRequest{
				PracticeItemId: "non-existent-id",
				UserAnswer:     "test answer",
			},
			expectedStatus: http.StatusNotFound,
			setupAuth:      true,
			errorMessage:   "Practice item not found",
		},
		{
			name: "empty user answer",
			request: AnswerEvaluationRequest{
				PracticeItemId: practiceItem.Id,
				UserAnswer:     "",
			},
			expectedStatus: http.StatusOK,
			expectedBody: AnswerEvaluationResponse{
				IsCorrect: false,
			},
			setupAuth: true,
		},
		{
			name: "practice item with no correct answer",
			request: AnswerEvaluationRequest{
				PracticeItemId: practiceItemNoAnswer.Id,
				UserAnswer:     "test answer",
			},
			expectedStatus: http.StatusBadRequest,
			setupAuth:      true,
			errorMessage:   "Practice item has no correct answer",
		},
		{
			name: "correct answer is raw number",
			request: AnswerEvaluationRequest{
				PracticeItemId: practiceItemRawNumber.Id,
				UserAnswer:     "12345",
			},
			expectedStatus: http.StatusOK,
			expectedBody: AnswerEvaluationResponse{
				IsCorrect: true,
			},
			setupAuth: true,
		},
		{
			name: "incorrect answer for raw number",
			request: AnswerEvaluationRequest{
				PracticeItemId: practiceItemRawNumber.Id,
				UserAnswer:     "54321",
			},
			expectedStatus: http.StatusOK,
			expectedBody: AnswerEvaluationResponse{
				IsCorrect: false,
			},
			setupAuth: true,
		},
		{
			name: "correct answer with backticks",
			request: AnswerEvaluationRequest{
				PracticeItemId: practiceItemBackticks.Id,
				UserAnswer:     "`code`",
			},
			expectedStatus: http.StatusOK,
			expectedBody: AnswerEvaluationResponse{
				IsCorrect: true,
			},
			setupAuth: true,
		},
		{
			name: "correct answer with triple backticks",
			request: AnswerEvaluationRequest{
				PracticeItemId: practiceItemTripleBackticks.Id,
				UserAnswer:     "```code block```",
			},
			expectedStatus: http.StatusOK,
			expectedBody: AnswerEvaluationResponse{
				IsCorrect: true,
			},
			setupAuth: true,
		},
		{
			name: "correct answer with backticks, user answer without",
			request: AnswerEvaluationRequest{
				PracticeItemId: practiceItemBackticks.Id,
				UserAnswer:     "code",
			},
			expectedStatus: http.StatusOK,
			expectedBody: AnswerEvaluationResponse{
				IsCorrect: true,
			},
			setupAuth: true,
		},
		{
			name: "correct answer without backticks, user answer with",
			request: AnswerEvaluationRequest{
				PracticeItemId: practiceItem.Id,
				UserAnswer:     "`test answer`",
			},
			expectedStatus: http.StatusOK,
			expectedBody: AnswerEvaluationResponse{
				IsCorrect: true,
			},
			setupAuth: true,
		},
		{
			name: "correct answer with triple backticks, user answer without",
			request: AnswerEvaluationRequest{
				PracticeItemId: practiceItemTripleBackticks.Id,
				UserAnswer:     "code block",
			},
			expectedStatus: http.StatusOK,
			expectedBody: AnswerEvaluationResponse{
				IsCorrect: true,
			},
			setupAuth: true,
		},
		{
			name: "unauthorized request",
			request: AnswerEvaluationRequest{
				PracticeItemId: practiceItem.Id,
				UserAnswer:     "test answer",
			},
			expectedStatus: http.StatusUnauthorized,
			setupAuth:      false,
		},
		{
			name: "invalid json body",
			request: AnswerEvaluationRequest{
				PracticeItemId: practiceItem.Id,
				UserAnswer:     "test answer",
			},
			expectedStatus: http.StatusBadRequest,
			setupAuth:      true,
			errorMessage:   "Invalid request body",
		},
		{
			name: "correct answer without backticks, user answer with triple backticks",
			request: AnswerEvaluationRequest{
				PracticeItemId: practiceItem.Id,
				UserAnswer:     "```test answer```",
			},
			expectedStatus: http.StatusOK,
			expectedBody: AnswerEvaluationResponse{
				IsCorrect: true,
			},
			setupAuth: true,
		},
		{
			name: "correct answer with single backtick, user answer with triple backticks",
			request: AnswerEvaluationRequest{
				PracticeItemId: practiceItemBackticks.Id,
				UserAnswer:     "```code```",
			},
			expectedStatus: http.StatusOK,
			expectedBody: AnswerEvaluationResponse{
				IsCorrect: true,
			},
			setupAuth: true,
		},
		{
			name: "correct answer with triple backticks, user answer with single backticks",
			request: AnswerEvaluationRequest{
				PracticeItemId: practiceItemTripleBackticks.Id,
				UserAnswer:     "`code block`",
			},
			expectedStatus: http.StatusOK,
			expectedBody: AnswerEvaluationResponse{
				IsCorrect: true,
			},
			setupAuth: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			var err error

			if tt.name == "invalid json body" {
				body = []byte(`{invalid json}`)
			} else {
				body, err = json.Marshal(tt.request)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/glimmer/v1/practice/evaluate-answer", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rec := httptest.NewRecorder()

			// Create request event
			var authUser *core.Record
			if tt.setupAuth {
				authUser = user
			}
			e := &core.RequestEvent{
				App:  app,
				Auth: authUser,
				Event: router.Event{
					Response: rec,
					Request:  req,
				},
			}

			// Create route handler
			route := NewAnswerRoute()

			// Handle request
			err = route.HandleEvaluateAnswer(e)

			// If error, simulate framework error handling
			if err != nil {
				t.Logf("Test case '%s' received error: %v", tt.name, err)

				// Determine status code based on error message
				statusCode := http.StatusInternalServerError
				errMsg := err.Error()
				switch {
				case errMsg == "You must be logged in.":
					statusCode = http.StatusUnauthorized
				case errMsg == "PracticeItemId is required." ||
					errMsg == "Invalid request body." ||
					errMsg == "Practice item has no correct answer.":
					statusCode = http.StatusBadRequest
				case errMsg == "Practice item not found.":
					statusCode = http.StatusNotFound
				}
				rec.WriteHeader(statusCode)

				// Create error response
				errorResponse := map[string]string{
					"message": errMsg,
				}
				json.NewEncoder(rec.Body).Encode(errorResponse)
			}

			// Check response
			if tt.expectedStatus != http.StatusOK {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
				if tt.errorMessage != "" {
					var errorResponse map[string]string
					err := json.Unmarshal(rec.Body.Bytes(), &errorResponse)
					require.NoError(t, err)
					assert.Contains(t, errorResponse["message"], tt.errorMessage)
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			// Parse response body
			var response AnswerEvaluationResponse
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			require.NoError(t, err)

			// Check response content
			assert.Equal(t, tt.expectedBody.IsCorrect, response.IsCorrect)
		})
	}
}

func TestHandleProcessAnswer(t *testing.T) {
	// Setup PocketBase test instance with temporary data directory
	app := setupTestApp(t)

	// Create test user once
	userCollection, err := app.FindCollectionByNameOrId("users")
	require.NoError(t, err)
	user := core.NewRecord(userCollection)
	user.Set("email", "test@example.com")
	user.Set("password", "test123")
	err = app.SaveNoValidate(user)
	require.NoError(t, err)

	// Create test learner once
	learnerCollection, err := app.FindCollectionByNameOrId(domain.CollectionLearners)
	require.NoError(t, err)
	learner := core.NewRecord(learnerCollection)
	learner.Set("nickname", "Test Learner")
	learner.Set("user", user.Id)
	err = app.SaveNoValidate(learner)
	require.NoError(t, err)

	t.Run("correct answer with no hints", func(t *testing.T) {
		// Create practice item with hints
		collection, err := app.FindCollectionByNameOrId(domain.CollectionPracticeItems)
		require.NoError(t, err)
		practiceItem := core.NewRecord(collection)
		correctAnswerJSON, err := json.Marshal("test answer")
		require.NoError(t, err)
		practiceItem.Set("correct_answer", string(correctAnswerJSON))
		hintsJSON, err := json.Marshal([]string{"First hint", "Second hint", "Third hint"})
		require.NoError(t, err)
		practiceItem.Set("hints", string(hintsJSON))
		err = app.SaveNoValidate(practiceItem)
		require.NoError(t, err)

		// Create practice session
		sessionCollection, err := app.FindCollectionByNameOrId(domain.CollectionPracticeSessions)
		require.NoError(t, err)
		session := core.NewRecord(sessionCollection)
		session.Set("learner", learner.Id)
		session.Set("status", "active")
		err = app.SaveNoValidate(session)
		require.NoError(t, err)

		request := ProcessAnswerRequest{
			PracticeItemId:   practiceItem.Id,
			UserAnswer:       "test answer",
			PracticeSession:  session.Id,
			LearnerId:        learner.Id,
			HintLevelReached: 0,
		}

		body, err := json.Marshal(request)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/glimmer/v1/practice/process-answer", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		e := &core.RequestEvent{
			App:  app,
			Auth: user,
			Event: router.Event{
				Response: rec,
				Request:  req,
			},
		}

		route := NewAnswerRoute()
		err = route.HandleProcessAnswer(e)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response ProcessAnswerResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response.IsCorrect)
		assert.Equal(t, 1.0, response.Score)
		assert.Equal(t, "Excellent! You got it right on your own!", response.Feedback)
		assert.Equal(t, 0, response.HintLevelReached)
		assert.Equal(t, 1, response.AttemptNumber)
	})

	t.Run("correct answer with hints used (score 0.5)", func(t *testing.T) {
		// Create practice item with hints
		collection, err := app.FindCollectionByNameOrId(domain.CollectionPracticeItems)
		require.NoError(t, err)
		practiceItem := core.NewRecord(collection)
		correctAnswerJSON, err := json.Marshal("test answer 2")
		require.NoError(t, err)
		practiceItem.Set("correct_answer", string(correctAnswerJSON))
		hintsJSON, err := json.Marshal([]string{"First hint", "Second hint", "Third hint"})
		require.NoError(t, err)
		practiceItem.Set("hints", string(hintsJSON))
		err = app.SaveNoValidate(practiceItem)
		require.NoError(t, err)

		// Create practice session
		sessionCollection, err := app.FindCollectionByNameOrId(domain.CollectionPracticeSessions)
		require.NoError(t, err)
		session := core.NewRecord(sessionCollection)
		session.Set("learner", learner.Id)
		session.Set("status", "active")
		err = app.SaveNoValidate(session)
		require.NoError(t, err)

		request := ProcessAnswerRequest{
			PracticeItemId:   practiceItem.Id,
			UserAnswer:       "test answer 2",
			PracticeSession:  session.Id,
			LearnerId:        learner.Id,
			HintLevelReached: 2, // Using 2 hints to get score 0.5
		}

		body, err := json.Marshal(request)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/glimmer/v1/practice/process-answer", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		e := &core.RequestEvent{
			App:  app,
			Auth: user,
			Event: router.Event{
				Response: rec,
				Request:  req,
			},
		}

		route := NewAnswerRoute()
		err = route.HandleProcessAnswer(e)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response ProcessAnswerResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response.IsCorrect)
		assert.InDelta(t, 0.5, response.Score, 0.01) // 1.0 - (2/(3+1)) = 0.5
		assert.Equal(t, "Good job! You found the right answer with some help from the hints.", response.Feedback)
		assert.Equal(t, 2, response.HintLevelReached)
		assert.Equal(t, 1, response.AttemptNumber)
	})

	t.Run("incorrect answer", func(t *testing.T) {
		// Create practice item
		collection, err := app.FindCollectionByNameOrId(domain.CollectionPracticeItems)
		require.NoError(t, err)
		practiceItem := core.NewRecord(collection)
		correctAnswerJSON, err := json.Marshal("correct answer")
		require.NoError(t, err)
		practiceItem.Set("correct_answer", string(correctAnswerJSON))
		err = app.SaveNoValidate(practiceItem)
		require.NoError(t, err)

		// Create practice session
		sessionCollection, err := app.FindCollectionByNameOrId(domain.CollectionPracticeSessions)
		require.NoError(t, err)
		session := core.NewRecord(sessionCollection)
		session.Set("learner", learner.Id)
		session.Set("status", "active")
		err = app.SaveNoValidate(session)
		require.NoError(t, err)

		request := ProcessAnswerRequest{
			PracticeItemId:   practiceItem.Id,
			UserAnswer:       "wrong answer",
			PracticeSession:  session.Id,
			LearnerId:        learner.Id,
			HintLevelReached: 0,
		}

		body, err := json.Marshal(request)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/glimmer/v1/practice/process-answer", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		e := &core.RequestEvent{
			App:  app,
			Auth: user,
			Event: router.Event{
				Response: rec,
				Request:  req,
			},
		}

		route := NewAnswerRoute()
		err = route.HandleProcessAnswer(e)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response ProcessAnswerResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.False(t, response.IsCorrect)
		assert.Equal(t, 0.0, response.Score)
		assert.Equal(t, "That's not correct. Consider using the hints for guidance.", response.Feedback)
		assert.Equal(t, 0, response.HintLevelReached)
		assert.Equal(t, 1, response.AttemptNumber)
	})

	t.Run("case insensitive match", func(t *testing.T) {
		// Create practice item
		collection, err := app.FindCollectionByNameOrId(domain.CollectionPracticeItems)
		require.NoError(t, err)
		practiceItem := core.NewRecord(collection)
		correctAnswerJSON, err := json.Marshal("case test")
		require.NoError(t, err)
		practiceItem.Set("correct_answer", string(correctAnswerJSON))
		err = app.SaveNoValidate(practiceItem)
		require.NoError(t, err)

		// Create practice session
		sessionCollection, err := app.FindCollectionByNameOrId(domain.CollectionPracticeSessions)
		require.NoError(t, err)
		session := core.NewRecord(sessionCollection)
		session.Set("learner", learner.Id)
		session.Set("status", "active")
		err = app.SaveNoValidate(session)
		require.NoError(t, err)

		request := ProcessAnswerRequest{
			PracticeItemId:   practiceItem.Id,
			UserAnswer:       "CASE TEST",
			PracticeSession:  session.Id,
			LearnerId:        learner.Id,
			HintLevelReached: 0,
		}

		body, err := json.Marshal(request)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/glimmer/v1/practice/process-answer", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		e := &core.RequestEvent{
			App:  app,
			Auth: user,
			Event: router.Event{
				Response: rec,
				Request:  req,
			},
		}

		route := NewAnswerRoute()
		err = route.HandleProcessAnswer(e)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response ProcessAnswerResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response.IsCorrect)
		assert.Equal(t, 1.0, response.Score)
		assert.Equal(t, "Excellent! You got it right on your own!", response.Feedback)
		assert.Equal(t, 0, response.HintLevelReached)
		assert.Equal(t, 1, response.AttemptNumber)
	})

	t.Run("missing practice item id", func(t *testing.T) {
		request := ProcessAnswerRequest{
			UserAnswer:       "test answer",
			PracticeSession:  "session-id",
			LearnerId:        learner.Id,
			HintLevelReached: 0,
		}

		body, err := json.Marshal(request)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/glimmer/v1/practice/process-answer", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		e := &core.RequestEvent{
			App:  app,
			Auth: user,
			Event: router.Event{
				Response: rec,
				Request:  req,
			},
		}

		route := NewAnswerRoute()
		err = route.HandleProcessAnswer(e)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "PracticeItemId is required")
	})

	t.Run("unauthorized request", func(t *testing.T) {
		request := ProcessAnswerRequest{
			PracticeItemId:   "item-id",
			UserAnswer:       "test answer",
			PracticeSession:  "session-id",
			LearnerId:        learner.Id,
			HintLevelReached: 0,
		}

		body, err := json.Marshal(request)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/glimmer/v1/practice/process-answer", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		e := &core.RequestEvent{
			App:  app,
			Auth: nil, // No auth
			Event: router.Event{
				Response: rec,
				Request:  req,
			},
		}

		route := NewAnswerRoute()
		err = route.HandleProcessAnswer(e)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "You must be logged in")
	})

	t.Run("practice item with no correct answer", func(t *testing.T) {
		// Create practice item with no correct answer
		collection, err := app.FindCollectionByNameOrId(domain.CollectionPracticeItems)
		require.NoError(t, err)
		practiceItem := core.NewRecord(collection)
		err = app.SaveNoValidate(practiceItem)
		require.NoError(t, err)

		// Create practice session
		sessionCollection, err := app.FindCollectionByNameOrId(domain.CollectionPracticeSessions)
		require.NoError(t, err)
		session := core.NewRecord(sessionCollection)
		session.Set("learner", learner.Id)
		session.Set("status", "active")
		err = app.SaveNoValidate(session)
		require.NoError(t, err)

		request := ProcessAnswerRequest{
			PracticeItemId:   practiceItem.Id,
			UserAnswer:       "test answer",
			PracticeSession:  session.Id,
			LearnerId:        learner.Id,
			HintLevelReached: 0,
		}

		body, err := json.Marshal(request)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/glimmer/v1/practice/process-answer", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		e := &core.RequestEvent{
			App:  app,
			Auth: user,
			Event: router.Event{
				Response: rec,
				Request:  req,
			},
		}

		route := NewAnswerRoute()
		err = route.HandleProcessAnswer(e)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Practice item has no correct answer")
	})

	t.Run("correct answer is a raw number", func(t *testing.T) {
		// Create practice item with a raw number as correct answer
		collection, err := app.FindCollectionByNameOrId(domain.CollectionPracticeItems)
		require.NoError(t, err)
		practiceItem := core.NewRecord(collection)
		practiceItem.Set("correct_answer", "12345")
		err = app.SaveNoValidate(practiceItem)
		require.NoError(t, err)

		// Create practice session
		sessionCollection, err := app.FindCollectionByNameOrId(domain.CollectionPracticeSessions)
		require.NoError(t, err)
		session := core.NewRecord(sessionCollection)
		session.Set("learner", learner.Id)
		session.Set("status", "active")
		err = app.SaveNoValidate(session)
		require.NoError(t, err)

		request := ProcessAnswerRequest{
			PracticeItemId:   practiceItem.Id,
			UserAnswer:       "12345",
			PracticeSession:  session.Id,
			LearnerId:        learner.Id,
			HintLevelReached: 0,
		}

		body, err := json.Marshal(request)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/glimmer/v1/practice/process-answer", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		e := &core.RequestEvent{
			App:  app,
			Auth: user,
			Event: router.Event{
				Response: rec,
				Request:  req,
			},
		}

		route := NewAnswerRoute()
		err = route.HandleProcessAnswer(e)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response ProcessAnswerResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response.IsCorrect)
		assert.Equal(t, 1.0, response.Score)
	})

	t.Run("correct answer contains backticks", func(t *testing.T) {
		// Create practice item with backticks
		collection, err := app.FindCollectionByNameOrId(domain.CollectionPracticeItems)
		require.NoError(t, err)
		practiceItem := core.NewRecord(collection)
		practiceItem.Set("correct_answer", "`code`")
		err = app.SaveNoValidate(practiceItem)
		require.NoError(t, err)

		// Create practice session
		sessionCollection, err := app.FindCollectionByNameOrId(domain.CollectionPracticeSessions)
		require.NoError(t, err)
		session := core.NewRecord(sessionCollection)
		session.Set("learner", learner.Id)
		session.Set("status", "active")
		err = app.SaveNoValidate(session)
		require.NoError(t, err)

		request := ProcessAnswerRequest{
			PracticeItemId:   practiceItem.Id,
			UserAnswer:       "`code`",
			PracticeSession:  session.Id,
			LearnerId:        learner.Id,
			HintLevelReached: 0,
		}

		body, err := json.Marshal(request)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/glimmer/v1/practice/process-answer", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		e := &core.RequestEvent{
			App:  app,
			Auth: user,
			Event: router.Event{
				Response: rec,
				Request:  req,
			},
		}

		route := NewAnswerRoute()
		err = route.HandleProcessAnswer(e)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response ProcessAnswerResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response.IsCorrect)
		assert.Equal(t, 1.0, response.Score)
	})

	t.Run("correct answer with backticks and user answer without", func(t *testing.T) {
		// Create practice item with backticks
		collection, err := app.FindCollectionByNameOrId(domain.CollectionPracticeItems)
		require.NoError(t, err)
		practiceItem := core.NewRecord(collection)
		practiceItem.Set("correct_answer", "```go\nfunc main() {}\n```")
		err = app.SaveNoValidate(practiceItem)
		require.NoError(t, err)

		// Create practice session
		sessionCollection, err := app.FindCollectionByNameOrId(domain.CollectionPracticeSessions)
		require.NoError(t, err)
		session := core.NewRecord(sessionCollection)
		session.Set("learner", learner.Id)
		session.Set("status", "active")
		err = app.SaveNoValidate(session)
		require.NoError(t, err)

		request := ProcessAnswerRequest{
			PracticeItemId:   practiceItem.Id,
			UserAnswer:       "go\nfunc main() {}\n",
			PracticeSession:  session.Id,
			LearnerId:        learner.Id,
			HintLevelReached: 0,
		}

		body, err := json.Marshal(request)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/glimmer/v1/practice/process-answer", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		e := &core.RequestEvent{
			App:  app,
			Auth: user,
			Event: router.Event{
				Response: rec,
				Request:  req,
			},
		}

		route := NewAnswerRoute()
		err = route.HandleProcessAnswer(e)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response ProcessAnswerResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response.IsCorrect)
		assert.Equal(t, 1.0, response.Score)
	})

	t.Run("correct answer without backticks and user with triple", func(t *testing.T) {
		// Create practice item
		collection, err := app.FindCollectionByNameOrId(domain.CollectionPracticeItems)
		require.NoError(t, err)
		practiceItem := core.NewRecord(collection)
		correctAnswerJSON, err := json.Marshal("plain text answer")
		require.NoError(t, err)
		practiceItem.Set("correct_answer", string(correctAnswerJSON))
		err = app.SaveNoValidate(practiceItem)
		require.NoError(t, err)

		// Create practice session
		sessionCollection, err := app.FindCollectionByNameOrId(domain.CollectionPracticeSessions)
		require.NoError(t, err)
		session := core.NewRecord(sessionCollection)
		session.Set("learner", learner.Id)
		session.Set("status", "active")
		err = app.SaveNoValidate(session)
		require.NoError(t, err)

		request := ProcessAnswerRequest{
			PracticeItemId:   practiceItem.Id,
			UserAnswer:       "```plain text answer```",
			PracticeSession:  session.Id,
			LearnerId:        learner.Id,
			HintLevelReached: 0,
		}

		body, err := json.Marshal(request)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/glimmer/v1/practice/process-answer", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		e := &core.RequestEvent{
			App:  app,
			Auth: user,
			Event: router.Event{
				Response: rec,
				Request:  req,
			},
		}

		route := NewAnswerRoute()
		err = route.HandleProcessAnswer(e)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response ProcessAnswerResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response.IsCorrect)
		assert.Equal(t, 1.0, response.Score)
	})
}

// TestHandleProcessAnswerSecondAttempt tests that second attempts increment attempt number
func TestHandleProcessAnswerSecondAttempt(t *testing.T) {
	app := setupTestApp(t)

	// Create test data
	collection, err := app.FindCollectionByNameOrId(domain.CollectionPracticeItems)
	require.NoError(t, err)

	practiceItem := core.NewRecord(collection)
	correctAnswerJSON, err := json.Marshal("correct answer")
	require.NoError(t, err)
	practiceItem.Set("correct_answer", string(correctAnswerJSON))
	err = app.SaveNoValidate(practiceItem)
	require.NoError(t, err)

	userCollection, err := app.FindCollectionByNameOrId("users")
	require.NoError(t, err)
	user := core.NewRecord(userCollection)
	user.Set("email", "test@example.com")
	user.Set("password", "test123")
	err = app.SaveNoValidate(user)
	require.NoError(t, err)

	learnerCollection, err := app.FindCollectionByNameOrId(domain.CollectionLearners)
	require.NoError(t, err)
	learner := core.NewRecord(learnerCollection)
	learner.Set("nickname", "Test Learner")
	learner.Set("user", user.Id)
	err = app.SaveNoValidate(learner)
	require.NoError(t, err)

	sessionCollection, err := app.FindCollectionByNameOrId(domain.CollectionPracticeSessions)
	require.NoError(t, err)
	session := core.NewRecord(sessionCollection)
	session.Set("learner", learner.Id)
	session.Set("status", "active")
	err = app.SaveNoValidate(session)
	require.NoError(t, err)

	route := NewAnswerRoute()

	// First attempt - incorrect answer
	req1 := ProcessAnswerRequest{
		PracticeItemId:   practiceItem.Id,
		UserAnswer:       "wrong answer",
		PracticeSession:  session.Id,
		LearnerId:        learner.Id,
		HintLevelReached: 0,
	}

	body1, err := json.Marshal(req1)
	require.NoError(t, err)

	httpReq1 := httptest.NewRequest(http.MethodPost, "/api/glimmer/v1/practice/process-answer", bytes.NewReader(body1))
	httpReq1.Header.Set("Content-Type", "application/json")
	rec1 := httptest.NewRecorder()

	e1 := &core.RequestEvent{
		App:  app,
		Auth: user,
		Event: router.Event{
			Response: rec1,
			Request:  httpReq1,
		},
	}

	err = route.HandleProcessAnswer(e1)
	require.NoError(t, err)

	var response1 ProcessAnswerResponse
	err = json.Unmarshal(rec1.Body.Bytes(), &response1)
	require.NoError(t, err)
	assert.Equal(t, 1, response1.AttemptNumber)
	assert.False(t, response1.IsCorrect)

	// Second attempt - correct answer
	req2 := ProcessAnswerRequest{
		PracticeItemId:   practiceItem.Id,
		UserAnswer:       "correct answer",
		PracticeSession:  session.Id,
		LearnerId:        learner.Id,
		HintLevelReached: 1,
	}

	body2, err := json.Marshal(req2)
	require.NoError(t, err)

	httpReq2 := httptest.NewRequest(http.MethodPost, "/api/glimmer/v1/practice/process-answer", bytes.NewReader(body2))
	httpReq2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()

	e2 := &core.RequestEvent{
		App:  app,
		Auth: user,
		Event: router.Event{
			Response: rec2,
			Request:  httpReq2,
		},
	}

	err = route.HandleProcessAnswer(e2)
	require.NoError(t, err)

	var response2 ProcessAnswerResponse
	err = json.Unmarshal(rec2.Body.Bytes(), &response2)
	require.NoError(t, err)
	assert.Equal(t, 2, response2.AttemptNumber)
	assert.True(t, response2.IsCorrect)

	// Verify database record was updated
	results, err := app.FindRecordsByFilter(
		domain.CollectionPracticeResults,
		"practice_item = {:practiceItem} && practice_session = {:session}",
		"-created",
		1,
		0,
		map[string]any{
			"practiceItem": practiceItem.Id,
			"session":      session.Id,
		},
	)
	require.NoError(t, err)
	require.Len(t, results, 1)

	result := results[0]
	// The answer is stored as JSON string, so we need to unmarshal it for comparison
	var storedAnswer string
	err = json.Unmarshal([]byte(result.GetString("answer")), &storedAnswer)
	require.NoError(t, err)
	assert.Equal(t, "correct answer", storedAnswer)
	assert.True(t, result.GetBool("is_correct"))
	assert.Equal(t, 1, result.GetInt("hint_level_reached"))
	assert.Equal(t, 2, result.GetInt("attempt_number"))
	assert.Greater(t, result.GetFloat("score"), 0.0)
}

func TestGetCleanCorrectAnswer(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expected   string
		expectWarn bool
	}{
		{
			name:       "JSON-encoded string",
			input:      `"hello world"`,
			expected:   "hello world",
			expectWarn: false,
		},
		{
			name:       "raw number string",
			input:      "123",
			expected:   "123",
			expectWarn: true,
		},
		{
			name:       "raw string with spaces",
			input:      "hello world",
			expected:   "hello world",
			expectWarn: true,
		},
		{
			name:       "JSON-encoded number as string",
			input:      `"123"`,
			expected:   "123",
			expectWarn: false,
		},
		{
			name:       "string with single backticks",
			input:      "`code`",
			expected:   "`code`",
			expectWarn: true,
		},
		{
			name:       "string with triple backticks",
			input:      "```code block```",
			expected:   "```code block```",
			expectWarn: true,
		},
		{
			name:       "JSON-encoded string with backticks",
			input:      "\"`code`\"",
			expected:   "`code`",
			expectWarn: false,
		},
		{
			name:       "empty string",
			input:      "",
			expected:   "",
			expectWarn: true,
		},
		{
			name:       "JSON-encoded empty string",
			input:      `""`,
			expected:   "",
			expectWarn: false,
		},
		{
			name:       "null value",
			input:      "null",
			expected:   "",
			expectWarn: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			var logBuf bytes.Buffer
			originalLogger := log.Logger
			log.Logger = zerolog.New(&logBuf)
			defer func() { log.Logger = originalLogger }()

			actual := getCleanCorrectAnswer(tt.input)
			assert.Equal(t, tt.expected, actual)

			logOutput := logBuf.String()
			if tt.expectWarn {
				assert.Contains(t, logOutput, "Could not unmarshal correct answer as JSON")
				assert.Contains(t, logOutput, tt.input)
			} else {
				assert.Empty(t, logOutput)
			}
		})
	}
}

func TestNormalizeAnswerString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "no change", input: "hello", expected: "hello"},
		{name: "simple trim space", input: "  hello  ", expected: "hello"},
		{name: "single backticks", input: "`hello`", expected: "hello"},
		{name: "single backticks with space", input: "  `hello`  ", expected: "hello"},
		{name: "single backticks with inner space", input: "`  hello  `", expected: "hello"},
		{name: "triple backticks", input: "```hello```", expected: "hello"},
		{name: "triple backticks with space", input: "  ```hello```  ", expected: "hello"},
		{name: "triple backticks with inner space", input: "```  hello  ```", expected: "hello"},
		{name: "mismatched backticks start", input: "`hello``", expected: "hello"},
		{name: "mismatched backticks end", input: "``hello`", expected: "hello"},
		{name: "backticks in middle", input: "he`llo", expected: "he`llo"},
		{name: "empty string", input: "", expected: ""},
		{name: "only backticks", input: "``", expected: ""},
		{name: "only triple backticks", input: "``````", expected: ""},
		{name: "only spaces", input: "   ", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := normalizeAnswerString(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
