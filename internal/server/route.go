package server

import "github.com/gofiber/swagger"

func (s *Server) initDormRoutes() {
	// dorm
	dormRoutes := s.app.Group("/dorms")
	dormRoutes.Post("/", s.authMiddleware.Auth, s.dormHandler.Create)
	dormRoutes.Get("/", s.dormHandler.GetAll)
	dormRoutes.Get("/:id", s.dormHandler.GetByID)
	dormRoutes.Patch("/:id", s.authMiddleware.Auth, s.dormHandler.Update)
	dormRoutes.Delete("/:id", s.authMiddleware.Auth, s.dormHandler.Delete)
}

func (s *Server) initAuthRoutes() {
	authRoutes := s.app.Group("/auth")
	authRoutes.Post("/register", s.userHandler.Register)
	authRoutes.Post("/login", s.userHandler.Login)
}

func (s *Server) initUserRoutes() {
	userRoutes := s.app.Group("/user")

	userRoutes.Get("/me", s.authMiddleware.Auth, s.userHandler.GetUserInfo)

	userRoutes.Post("/verify", s.userHandler.VerifyEmail)
	userRoutes.Post("/resetpassword", s.userHandler.ResetPasswordCreate)
	userRoutes.Post("/newpassword", s.authMiddleware.Auth, s.userHandler.ResetPassword)
	userRoutes.Patch("/", s.authMiddleware.Auth, s.userHandler.UpdateUserInformation)
	userRoutes.Delete("/", s.authMiddleware.Auth, s.userHandler.DeleteAccount)
}

func (s *Server) initExampleUploadRoutes() {
	s.app.Post("/upload/public", s.testUploadHandler.UploadToPublicBucketHandler)
	s.app.Post("/upload/private", s.testUploadHandler.UploadToPrivateBucketHandler)
	s.app.Get("/signedurl/*", s.testUploadHandler.GetSignedUrlHandler)
}

func (s *Server) initRoutes() {
	// greeting
	s.app.Get("/", s.greetingHandler.Greeting)

	// swagger
	s.app.Get("/swagger/*", swagger.HandlerDefault)

	s.initExampleUploadRoutes()
	s.initUserRoutes()
	s.initAuthRoutes()
	s.initDormRoutes()
}
