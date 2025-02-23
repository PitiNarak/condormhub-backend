package server

import "github.com/gofiber/swagger"

func (s *Server) initRoutes() {
	// greeting
	s.app.Get("/", s.handler.greeting.Greeting)

	// swagger
	s.app.Get("/swagger/*", swagger.HandlerDefault)

	s.initExampleUploadRoutes()
	s.initUserRoutes()
	s.initAuthRoutes()
	s.initDormRoutes()
	s.initLeasingHistoryRoutes()
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
	dormRoutes := s.app.Group("/history")
	dormRoutes.Post("/", s.authMiddleware.Auth, s.handler.leasingHistory.Create)
	dormRoutes.Get("/me", s.authMiddleware.Auth, s.handler.leasingHistory.GetByUserID)
	dormRoutes.Get("/:id", s.authMiddleware.Auth, s.handler.leasingHistory.GetByDormID)
}
