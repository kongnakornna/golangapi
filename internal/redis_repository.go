package internal

import (
	"context"
)

type RedisRepository[M any] interface {
	Create(ctx context.Context, key string, exp *M, seconds int) error
	Get(ctx context.Context, key string) (*M, error)
	Delete(ctx context.Context, key string) error
	Sadd(ctx context.Context, key string, value string) error
	Sadds(ctx context.Context, key string, values []string) error
	Srem(ctx context.Context, key string, value string) error
	SIsMember(ctx context.Context, key string, value string) (bool, error)
	Expire(ctx context.Context, key string, seconds int) error
	// Simple key-value operations for revocation checking
	Set(ctx context.Context, key string, value string, seconds int) error
	GetString(ctx context.Context, key string) (string, error)
	// SMembers(ctx context.Context, key string) ([]string, error)
}
