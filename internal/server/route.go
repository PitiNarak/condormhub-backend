package server

import (
	"github.com/PitiNarak/condormhub-backend/docs"
	"github.com/gofiber/swagger"
	"github.com/swaggo/swag"
)

func (s *Server) initRoutes() {
	// greeting
	s.app.Get("/", s.handler.greeting.Greeting)

	// swagger
	swag.Register(docs.SwaggerInfo.InfoInstanceName, docs.SwaggerInfo)
	s.app.Get("/swagger/*", swagger.HandlerDefault)

	s.initExampleUploadRoutes()
	s.initUserRoutes()
	s.initAuthRoutes()
	s.initDormRoutes()
	s.initLeasingHistoryRoutes()
	s.initOrderRoutes()
	s.initTransactionRoutes()
	s.initOwnershipProofRoutes()
}

func (s *Server) initExampleUploadRoutes() {
	s.app.Post("/upload/public", s.handler.exampleUpload.UploadToPublicBucketHandler)
	s.app.Post("/upload/private", s.handler.exampleUpload.UploadToPrivateBucketHandler)
	s.app.Get("/signedurl/*", s.handler.exampleUpload.GetSignedUrlHandler)
}

func (s *Server) initUserRoutes() {
	userRoutes := s.app.Group("/user")

	userRoutes.Get("/me", s.authMiddleware.Auth, s.handler.user.GetUserInfo)

	userRoutes.Post("/verify", s.handler.user.VerifyEmail)
	userRoutes.Post("/resetpassword", s.handler.user.ResetPasswordCreate)
	userRoutes.Post("/newpassword", s.handler.user.ResetPassword)
	userRoutes.Patch("/", s.authMiddleware.Auth, s.handler.user.UpdateUserInformation)
	userRoutes.Delete("/", s.authMiddleware.Auth, s.handler.user.DeleteAccount)
}

func (s *Server) initAuthRoutes() {
	authRoutes := s.app.Group("/auth")
	authRoutes.Post("/register", s.handler.user.Register)
	authRoutes.Post("/login", s.handler.user.Login)
	authRoutes.Post("/refresh", s.handler.user.RefreshToken)
}

func (s *Server) initDormRoutes() {
	// dorm
	dormRoutes := s.app.Group("/dorms")
	dormRoutes.Post("/", s.authMiddleware.Auth, s.handler.dorm.Create)
	dormRoutes.Get("/", s.handler.dorm.GetAll)
	dormRoutes.Get("/:id", s.handler.dorm.GetByID)
	dormRoutes.Patch("/:id", s.authMiddleware.Auth, s.handler.dorm.Update)
	dormRoutes.Delete("/:id", s.authMiddleware.Auth, s.handler.dorm.Delete)
}

func (s *Server) initLeasingHistoryRoutes() {
	historyRoutes := s.app.Group("/history")
	historyRoutes.Post("/:id", s.authMiddleware.Auth, s.handler.leasingHistory.Create)
	historyRoutes.Get("/me", s.authMiddleware.Auth, s.handler.leasingHistory.GetByUserID)
	historyRoutes.Get("/bydorm/:id", s.authMiddleware.Auth, s.handler.leasingHistory.GetByDormID)
	historyRoutes.Patch("/:id", s.authMiddleware.Auth, s.handler.leasingHistory.SetEndTimestamp)
	historyRoutes.Delete("/:id", s.authMiddleware.Auth, s.handler.leasingHistory.Delete)
}

func (s *Server) initOrderRoutes() {
	orderRoutes := s.app.Group("/order")
	orderRoutes.Post("/", s.authMiddleware.Auth, s.handler.order.CreateOrder)
	orderRoutes.Get("/:id", s.authMiddleware.Auth, s.handler.order.GetOrderByID)
	orderRoutes.Get("/unpaid/me", s.authMiddleware.Auth, s.handler.order.GetMyUnpaidOrder)
	orderRoutes.Get("/unpaid/:id", s.authMiddleware.Auth, s.handler.order.GetUnpaidOrderByUserID)
}

func (s *Server) initTransactionRoutes() {
	tsxRoutes := s.app.Group("/transaction")
	tsxRoutes.Post("/", s.authMiddleware.Auth, s.handler.tsx.CreateTransaction)
	tsxRoutes.Post("/webhook", s.handler.tsx.Webhook)
}

func (s *Server) initOwnershipProofRoutes() {
	ownershipRoutes := s.app.Group("/ownership")
	ownershipRoutes.Post("/create", s.authMiddleware.Auth, s.handler.ownershipProof.Create)
	ownershipRoutes.Delete("/:id", s.authMiddleware.Auth, s.handler.ownershipProof.Delete)
	ownershipRoutes.Get("/:id", s.handler.ownershipProof.GetByDormID)
	ownershipRoutes.Post("/update", s.authMiddleware.Auth, s.handler.ownershipProof.Update)
	ownershipRoutes.Post("/:id/approve", s.authMiddleware.Auth, s.handler.ownershipProof.Approve)
	ownershipRoutes.Post("/:id/reject", s.authMiddleware.Auth, s.handler.ownershipProof.Reject)

}
