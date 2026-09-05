package redis_repository

import (
	"HailowSellerService/internal/domain"
	cache "HailowSellerService/pkg/redis"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionRepository struct {
	client     *cache.RedisClient
	prefix_key string
}

func NewSessionRepository(client *cache.RedisClient) *SessionRepository {
	return &SessionRepository{
		client:     client,
		prefix_key: "session",
	}
}

// CreateSession saves a new session, deleting the old token for the same device
func (r *SessionRepository) CreateSession(ctx context.Context, session *domain.RefreshSession) error {
	ttl := time.Until(session.ExpiresAt)
	if ttl <= 0 {
		return errors.New("Expiration time must be in the future")
	}

	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	tokenKey := fmt.Sprintf("%s:%s", r.prefix_key, session.RefreshToken)
	userDevicesKey := fmt.Sprintf("user_devices:%s", session.SellerID)

	// Using a pipeline for atomic operations
	pipe := r.client.Pipeline()

	// Checking if this device had an old token and remove it
	oldToken, err := r.client.HGet(ctx, userDevicesKey, session.DeviceID).Result()
	if err == nil && oldToken != "" {
		oldTokenKey := fmt.Sprintf("%s:%s", r.prefix_key, oldToken)
		pipe.Del(ctx, oldTokenKey)
	}

	// Saving a new session using the token key
	pipe.Set(ctx, tokenKey, data, ttl)

	// Update refresh_token
	pipe.HSet(ctx, userDevicesKey, session.DeviceID, session.RefreshToken)

	_, err = pipe.Exec(ctx)
	return err
}

// GetSessionByToken finds a session by refresh token
func (r *SessionRepository) GetSessionByToken(ctx context.Context, token string) (*domain.RefreshSession, error) {
	key := fmt.Sprintf("%s:%s", r.prefix_key, token)
	data, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, domain.ErrSessionNotFound
	} else if err != nil {
		return nil, err
	}

	var session domain.RefreshSession
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// DeleteSession deletes the current device token
func (r *SessionRepository) DeleteSession(ctx context.Context, token string) error {
	session, err := r.GetSessionByToken(ctx, token)
	if err != nil {
		return err
	}

	tokenKey := fmt.Sprintf("%s:%s", r.prefix_key, token)
	userDevicesKey := fmt.Sprintf("user_devices:%s", session.SellerID)

	pipe := r.client.Pipeline()
	pipe.Del(ctx, tokenKey)
	pipe.HDel(ctx, userDevicesKey, session.DeviceID)

	_, err = pipe.Exec(ctx)
	return err
}

// DeleteAllUserSessions logs out a user from all devices
func (r *SessionRepository) DeleteAllUserSessions(ctx context.Context, userID string) error {
	userDevicesKey := fmt.Sprintf("user_devices:%s", userID)

	devicesMap, err := r.client.HGetAll(ctx, userDevicesKey).Result()
	if err != nil {
		return err
	}

	pipe := r.client.Pipeline()

	// Deleting all session keys
	for _, token := range devicesMap {
		tokenKey := fmt.Sprintf("%s:%s", r.prefix_key, token)
		pipe.Del(ctx, tokenKey)
	}

	// Removing Device Hash
	pipe.Del(ctx, userDevicesKey)

	_, err = pipe.Exec(ctx)
	return err
}
