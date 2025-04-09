package server

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	handler1 "github.com/PitiNarak/condormhub-backend/internal/handler"
)

type handler struct {
	greeting       *handler1.GreetingHandler
	user           ports.UserHandler
	exampleUpload  *handler1.TestUploadHandler
	dorm           ports.DormHandler
	leasingHistory ports.LeasingHistoryHandler
	order          ports.OrderHandler
	tsx            ports.TransactionHandler
	ownershipProof ports.OwnershipProofHandler
	contract       ports.ContractHandler
	leasingRequest ports.LeasingRequestHandler
	receipt        ports.ReceiptHandler
	support        ports.SupportHandler
}

func (s *Server) initHandler() {
	greeting := handler1.NewGreetingHandler()
	user := handler1.NewUserHandler(s.service.user)
	exampleUpload := handler1.NewTestUploadHandler(s.storage)
	dorm := handler1.NewDormHandler(s.service.dorm)
	leasingHistory := handler1.NewLeasingHistoryHandler(s.service.leasingHistory)
	order := handler1.NewOrderHandler(s.service.order)
	tsx := handler1.NewTransactionHandler(s.service.tsx, s.stripeConfig)
	ownershipProof := handler1.NewOwnershipProofHandler(s.service.ownershipProof, s.storage)
	contract := handler1.NewContractHandler(s.service.contract)
	leasingRequest := handler1.NewLeasingRequestHandler(s.service.leasingRequest)
	receipt := handler1.NewReceiptHandler(s.service.receipt)
	support := handler1.NewSupportHandler(s.service.support)

	s.handler = &handler{
		greeting:       greeting,
		user:           user,
		exampleUpload:  exampleUpload,
		dorm:           dorm,
		leasingHistory: leasingHistory,
		order:          order,
		tsx:            tsx,
		ownershipProof: ownershipProof,
		contract:       contract,
		leasingRequest: leasingRequest,
		receipt:        receipt,
		support:        support,
	}
}
