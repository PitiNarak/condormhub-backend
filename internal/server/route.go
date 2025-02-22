package server

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
