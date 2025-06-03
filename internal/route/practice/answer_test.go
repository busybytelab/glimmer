package practice

import (
	"bytes"
	"encoding/json"
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

	// Run migrations in the temporary directory
	cmd := exec.Command("go", "run", "../../../cmd/glimmer/main.go", "migrate", "--dir", pbDataDir)
	if err := cmd.Run(); err != nil {
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
