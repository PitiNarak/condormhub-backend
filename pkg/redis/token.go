package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *Redis) SetAccessToken(ctx context.Context, userId uuid.UUID, accessToken string) error {
	accessTokenKey := fmt.Sprintf("access_token:%s", userId)

	err := r.client.Set(ctx, accessTokenKey, accessToken, time.Hour*time.Duration(r.config.AccessTokenExpireHrs)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) SetRefreshToken(ctx context.Context, userId uuid.UUID, refreshToken string) error {
	refreshTokenKey := fmt.Sprintf("refresh_token:%s", userId)

	err := r.client.Set(ctx, refreshTokenKey, refreshToken, time.Hour*time.Duration(r.config.RefreshTokenExpireHrs)).Err()
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

func (r *Redis) DeleteToken(ctx context.Context, userId uuid.UUID) error {
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
