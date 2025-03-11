package server

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/handlers"
)

type handler struct {
	greeting       *handlers.GreetingHandler
	user           ports.UserHandler
	exampleUpload  *handlers.TestUploadHandler
	dorm           ports.DormHandler
	leasingHistory ports.LeasingHistoryHandler
	ownershipProof ports.OwnershipProofHandler
}

func (s *Server) initHandler() {
	greeting := handlers.NewGreetingHandler()
	user := handlers.NewUserHandler(s.service.user)
	exampleUpload := handlers.NewTestUploadHandler(s.storage)
	dorm := handlers.NewDormHandler(s.service.dorm)
	leasingHistory := handlers.NewLeasingHistoryHandler(s.service.leasingHistory)
	ownershipProof := handlers.NewOwnershipProofHandler(s.service.ownershipProof, s.storage)

	s.handler = &handler{
		greeting:       greeting,
		user:           user,
		exampleUpload:  exampleUpload,
		dorm:           dorm,
		leasingHistory: leasingHistory,
		ownershipProof: ownershipProof,
	}
}
