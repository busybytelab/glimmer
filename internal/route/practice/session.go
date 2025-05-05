package practice

import (
	"github.com/busybytelab.com/glimmer/internal/llm"
	"github.com/pocketbase/pocketbase/core"
)

type (
	SessionRoute interface {
		HandleCreatePracticeSession(e *core.RequestEvent) error
	}

	sessionRoute struct {
		llmService *llm.Service
	}
)

func NewPracticeSessionRoute(llmService *llm.Service) SessionRoute {
	return &sessionRoute{
		llmService: llmService,
	}
}

func (r *sessionRoute) HandleCreatePracticeSession(e *core.RequestEvent) error {

	return e.Next()
}
