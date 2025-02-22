package server

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/handlers"
)

type Handler struct {
	greeting      *handlers.GreetingHandler
	user          ports.UserHandler
	exampleUpload *handlers.TestUploadHandler
	dorm          ports.DormHandler
}
