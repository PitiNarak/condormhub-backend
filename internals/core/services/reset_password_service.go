package services

import "github.com/PitiNarak/condormhub-backend/internals/core/domain"

func (s *UserService) ResetPasswordCreate(email string) (domain.User, error) {
	user, err := s.UserRepo.GetUserViaEmail(email)
	if err != nil {
		return user, err
	}
	return user, nil
}
