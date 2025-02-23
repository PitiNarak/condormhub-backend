package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *Redis) SetAccessToken(ctx context.Context, userId uuid.UUID, accessToken string, ttl time.Duration) error {
	accessTokenKey := fmt.Sprintf("access_token:%s", userId)

	err := r.client.Set(ctx, accessTokenKey, accessToken, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) SetRefreshToken(ctx context.Context, userId uuid.UUID, refreshToken string, ttl time.Duration) error {
	refreshTokenKey := fmt.Sprintf("refresh_token:%s", userId)

	err := r.client.Set(ctx, refreshTokenKey, refreshToken, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) GetAccessToken(ctx context.Context, userId uuid.UUID) (string, error) {
	accessTokenKey := fmt.Sprintf("access_token:%s", userId)

	accessToken, err := r.client.Get(ctx, accessTokenKey).Result()
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (r *Redis) GetRefreshToken(ctx context.Context, userId uuid.UUID) (string, error) {
	refreshTokenKey := fmt.Sprintf("refresh_token:%s", userId)

	refreshToken, err := r.client.Get(ctx, refreshTokenKey).Result()
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (r *Redis) DeleteAccessTokenAndRefreshToken(ctx context.Context, userId uuid.UUID) error {
	accessTokenKey := fmt.Sprintf("access_token:%s", userId)

	err := r.client.Del(ctx, accessTokenKey).Err()
	if err != nil {
		return err
	}

	refreshTokenKey := fmt.Sprintf("refresh_token:%s", userId)
	err = r.client.Del(ctx, refreshTokenKey).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) SetVerificationToken(ctx context.Context, userID uuid.UUID, token string, ttl time.Duration) error {
	verificationTokenKey := fmt.Sprintf("verification_token:%s", userID)

	err := r.client.Set(ctx, verificationTokenKey, token, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) GetVerificationToken(ctx context.Context, userID uuid.UUID) (string, error) {
	verificationTokenKey := fmt.Sprintf("verification_token:%s", userID)

	token, err := r.client.Get(ctx, verificationTokenKey).Result()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (r *Redis) DeleteVerificationToken(ctx context.Context, userID uuid.UUID) error {
	verificationTokenKey := fmt.Sprintf("verification_token:%s", userID)

	err := r.client.Del(ctx, verificationTokenKey).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) SetResetToken(ctx context.Context, userID uuid.UUID, token string, ttl time.Duration) error {
	resetTokenKey := fmt.Sprintf("reset_token:%s", userID)

	err := r.client.Set(ctx, resetTokenKey, token, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) GetResetToken(ctx context.Context, userID uuid.UUID) (string, error) {
	resetTokenKey := fmt.Sprintf("reset_token:%s", userID)

	token, err := r.client.Get(ctx, resetTokenKey).Result()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (r *Redis) DeleteResetToken(ctx context.Context, userID uuid.UUID) error {
	resetTokenKey := fmt.Sprintf("reset_token:%s", userID)

	err := r.client.Del(ctx, resetTokenKey).Err()
	if err != nil {
		return err
	}

	return nil
}
