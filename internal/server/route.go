package server

import (
	"github.com/yokeTH/go-pkg/scalar"
)

func (s *Server) initRoutes() {
	// greeting
	s.app.Get("/", s.handler.greeting.Greeting)

	// swagger
	s.app.Get("/swagger/*", scalar.DefaultHandler)

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
	userRoutes.Get("/:id", s.handler.user.GetUserByID)

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
	historyRoutes := s.app.Group("/history", s.authMiddleware.Auth)
	historyRoutes.Post("/:id/review", s.handler.leasingHistory.CreateReview)
	historyRoutes.Patch("/:id/review", s.handler.leasingHistory.UpdateReview)
	historyRoutes.Delete("/:id/review", s.handler.leasingHistory.DeleteReview)
	historyRoutes.Post("/:id/review/image", s.handler.leasingHistory.UploadReviewImage)
	historyRoutes.Delete("/review/image/:url", s.handler.leasingHistory.DeleteReviewImageByURL)
	historyRoutes.Get("/me", s.handler.leasingHistory.GetByUserID)
	historyRoutes.Get("/bydorm/:id", s.handler.leasingHistory.GetByDormID)
	historyRoutes.Get("/:id", s.handler.leasingHistory.GetByID)
	historyRoutes.Patch("/:id", s.handler.leasingHistory.SetEndTimestamp)
	historyRoutes.Delete("/:id", s.handler.leasingHistory.Delete)
	historyRoutes.Post("/:id/review/report", s.handler.leasingHistory.ReportReview)
}

func (s *Server) initLeasingRequestRoutes() {
	requestRoutes := s.app.Group("/request", s.authMiddleware.Auth)
	requestRoutes.Post("/:id", s.handler.leasingRequest.Create)
	requestRoutes.Get("/me", s.handler.leasingRequest.GetByUserID)
	requestRoutes.Patch("/:id/approve", s.handler.leasingRequest.Approve)
	requestRoutes.Patch("/:id/reject", s.handler.leasingRequest.Reject)
	requestRoutes.Patch("/:id/cancel", s.handler.leasingRequest.Cancel)
	requestRoutes.Delete("/:id", s.authMiddleware.RequireAdmin, s.handler.leasingRequest.Delete)
	requestRoutes.Get("/bydorm/:id", s.handler.leasingRequest.GetByDormID)
}

func (s *Server) initOrderRoutes() {
	orderRoutes := s.app.Group("/order", s.authMiddleware.Auth)
	orderRoutes.Post("/", s.handler.order.CreateOrder)
	orderRoutes.Get("/:id", s.handler.order.GetOrderByID)
	orderRoutes.Get("/unpaid/me", s.handler.order.GetMyUnpaidOrder)
	orderRoutes.Get("/unpaid/:id", s.handler.order.GetUnpaidOrderByUserID)
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
	ownershipRoutes.Post("/:id/approve", s.authMiddleware.Auth, s.authMiddleware.RequireAdmin, s.handler.ownershipProof.Approve)
	ownershipRoutes.Post("/:id/reject", s.authMiddleware.Auth, s.authMiddleware.RequireAdmin, s.handler.ownershipProof.Reject)

}

func (s *Server) initReceiptRoutes() {
	receiptRoutes := s.app.Group("/receipt", s.authMiddleware.Auth)
	receiptRoutes.Get("/", s.handler.receipt.GetByUserID)
}

func (s *Server) initContractRoutes() {
	contractRoutes := s.app.Group("/contract", s.authMiddleware.Auth)
	contractRoutes.Patch("/:contractID/sign", s.handler.contract.SignContract)
	contractRoutes.Patch("/:contractID/cancel", s.handler.contract.CancelContract)
	contractRoutes.Get("/:contractID", s.handler.contract.GetContractByContractID)
	contractRoutes.Get("/", s.handler.contract.GetContractByUserID)
	contractRoutes.Get("/:dormID", s.handler.contract.GetContractByDormID)
	contractRoutes.Delete("/:contractID", s.handler.contract.Delete)

}

func (s *Server) initSupportRoutes() {
	supportRoutes := s.app.Group("/support", s.authMiddleware.Auth)
	supportRoutes.Post("/", s.handler.support.Create)
	supportRoutes.Get("/", s.handler.support.GetAll)
	supportRoutes.Patch("/:id", s.authMiddleware.RequireAdmin, s.handler.support.UpdateStatus)
}

func (s *Server) initAdminRoutes() {
	adminRoutes := s.app.Group("/admin", s.authMiddleware.Auth, s.authMiddleware.RequireAdmin)
	adminRoutes.Patch("/user/:id/ban", s.handler.user.BanUser)
	adminRoutes.Patch("/user/:id/unban", s.handler.user.UnbanUser)
	adminRoutes.Get("/lessee/pending", s.handler.user.GetPending)
	adminRoutes.Patch("/lessee/:id/verify", s.handler.user.VerifyStudentVerification)
	adminRoutes.Patch("/lessee/:id/reject", s.handler.user.RejectStudentVerification)
	adminRoutes.Get("/reviews/reported", s.handler.leasingHistory.GetReportedReviews)
	adminRoutes.Delete("/reviews/:id", s.handler.leasingHistory.DeleteReview)
}
