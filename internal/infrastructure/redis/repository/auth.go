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
		prefix_key: "seller_session",
	}
}

func (r *SessionRepository) CreateSession(ctx context.Context, session *domain.RefreshSession) error {
	data, err := json.Marshal(session)

	if err != nil {
		return err
	}

	ttl := time.Until(session.ExpiresAt)
	if ttl <= 0 {
		return errors.New("Expiration time must be in the future")
	}

	key := fmt.Sprintf("%s:%s", r.prefix_key, session.RefreshToken)
	return r.client.Set(ctx, key, data, ttl).Err()
}

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
