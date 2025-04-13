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
	s.initLeasingRequestRoutes()
	s.initOrderRoutes()
	s.initTransactionRoutes()
	s.initOwnershipProofRoutes()
	s.initReceiptRoutes()
	s.initContractRoutes()
	s.initSupportRoutes()
	s.initAdminRoutes()
}

func (s *Server) initExampleUploadRoutes() {
	s.app.Post("/upload/public", s.handler.exampleUpload.UploadToPublicBucketHandler)
	s.app.Post("/upload/private", s.handler.exampleUpload.UploadToPrivateBucketHandler)
	s.app.Get("/signedurl/*", s.handler.exampleUpload.GetSignedUrlHandler)
}

func (s *Server) initUserRoutes() {
	userRoutes := s.app.Group("/user")

	userRoutes.Get("/me", s.authMiddleware.Auth, s.handler.user.GetUserInfo)
	userRoutes.Get("/income", s.authMiddleware.Auth, s.handler.user.GetLessorIncome)
	userRoutes.Get("/:id", s.authMiddleware.Auth, s.handler.user.GetUserByID)

	userRoutes.Post("/verify", s.handler.user.VerifyEmail)
	userRoutes.Post("/resetpassword", s.handler.user.ResetPasswordCreate)
	userRoutes.Post("/newpassword", s.handler.user.ResetPassword)
	userRoutes.Post("/resend", s.handler.user.ResendVerificationEmailHandler)

	userRoutes.Patch("/firstfill", s.authMiddleware.Auth, s.handler.user.FirstFillInformation)
	userRoutes.Patch("/", s.authMiddleware.Auth, s.handler.user.UpdateUserInformation)
	userRoutes.Delete("/", s.authMiddleware.Auth, s.handler.user.DeleteAccount)

	userRoutes.Post("/studentEvidence", s.authMiddleware.Auth, s.handler.user.UploadStudentEvidence)
	userRoutes.Get("/:id/studentEvidence", s.authMiddleware.Auth, s.handler.user.GetStudentEvidenceByID)

	userRoutes.Post("/profilePic", s.authMiddleware.Auth, s.handler.user.UploadProfilePicture)
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
	dormRoutes.Delete("/images/:url", s.authMiddleware.Auth, s.handler.dorm.DeleteDormImageByURL)
	dormRoutes.Delete("/:id", s.authMiddleware.Auth, s.handler.dorm.Delete)
	dormRoutes.Post("/:id/images", s.authMiddleware.Auth, s.handler.dorm.UploadDormImage)
	dormRoutes.Get("/owner/:id", s.handler.dorm.GetByOwnerID)
}

func (s *Server) initLeasingHistoryRoutes() {
	historyRoutes := s.app.Group("/history")
	historyRoutes.Post("/review/:id", s.authMiddleware.Auth, s.handler.leasingHistory.CreateReview)
	historyRoutes.Patch("/review/:id", s.authMiddleware.Auth, s.handler.leasingHistory.UpdateReview)
	historyRoutes.Delete("/review/:id", s.authMiddleware.Auth, s.handler.leasingHistory.DeleteReview)
	historyRoutes.Get("/me", s.authMiddleware.Auth, s.handler.leasingHistory.GetByUserID)
	historyRoutes.Get("/bydorm/:id", s.authMiddleware.Auth, s.handler.leasingHistory.GetByDormID)
	historyRoutes.Patch("/:id", s.authMiddleware.Auth, s.handler.leasingHistory.SetEndTimestamp)
	historyRoutes.Delete("/:id", s.authMiddleware.Auth, s.handler.leasingHistory.Delete)
}

func (s *Server) initLeasingRequestRoutes() {
	historyRoutes := s.app.Group("/request")
	historyRoutes.Post("/:id", s.authMiddleware.Auth, s.handler.leasingRequest.Create)
	historyRoutes.Get("/me", s.authMiddleware.Auth, s.handler.leasingRequest.GetByUserID)
	historyRoutes.Patch("/:id/approve", s.authMiddleware.Auth, s.handler.leasingRequest.Approve)
	historyRoutes.Patch("/:id/reject", s.authMiddleware.Auth, s.handler.leasingRequest.Reject)
	historyRoutes.Patch("/:id/cancel", s.authMiddleware.Auth, s.handler.leasingRequest.Cancel)
	historyRoutes.Delete("/:id", s.authMiddleware.Auth, s.handler.leasingRequest.Delete)
	historyRoutes.Get("/bydorm/:id", s.authMiddleware.Auth, s.handler.leasingRequest.GetByDormID)
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
	ownershipRoutes.Post("/:id/upload", s.authMiddleware.Auth, s.handler.ownershipProof.UploadFile)
	ownershipRoutes.Delete("/:id", s.authMiddleware.Auth, s.handler.ownershipProof.Delete)
	ownershipRoutes.Get("/:id", s.handler.ownershipProof.GetByDormID)
	ownershipRoutes.Post("/:id/approve", s.authMiddleware.Auth, s.handler.ownershipProof.Approve)
	ownershipRoutes.Post("/:id/reject", s.authMiddleware.Auth, s.handler.ownershipProof.Reject)

}

func (s *Server) initReceiptRoutes() {
	ownershipRoutes := s.app.Group("/receipt")
	ownershipRoutes.Get("/", s.authMiddleware.Auth, s.handler.receipt.GetByUserID)
}

func (s *Server) initContractRoutes() {
	contractRoutes := s.app.Group("/contract")
	contractRoutes.Patch("/:contractID/sign", s.authMiddleware.Auth, s.handler.contract.SignContract)
	contractRoutes.Patch("/:contractID/cancel", s.authMiddleware.Auth, s.handler.contract.CancelContract)
	contractRoutes.Get("/:contractID", s.authMiddleware.Auth, s.handler.contract.GetContractByContractID)
	contractRoutes.Get("/", s.authMiddleware.Auth, s.handler.contract.GetContractByUserID)
	contractRoutes.Get("/:dormID", s.authMiddleware.Auth, s.handler.contract.GetContractByDormID)
	contractRoutes.Delete("/:contractID", s.authMiddleware.Auth, s.handler.contract.Delete)

}

func (s *Server) initSupportRoutes() {
	supportRoutes := s.app.Group("/support")
	supportRoutes.Post("/", s.authMiddleware.Auth, s.handler.support.Create)
	supportRoutes.Get("/", s.authMiddleware.Auth, s.handler.support.GetAll)
	supportRoutes.Patch("/:id", s.authMiddleware.Auth, s.handler.support.UpdateStatus)
}

func (s *Server) initAdminRoutes() {
	adminRoutes := s.app.Group("/admin", s.authMiddleware.Auth, s.authMiddleware.RequireAdmin)
	adminRoutes.Patch("/user/:id/ban", s.handler.user.BanUser)
	adminRoutes.Patch("/user/:id/unban", s.handler.user.UnbanUser)
}
