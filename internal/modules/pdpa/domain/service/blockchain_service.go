package service

import "context"

type BlockchainService interface {
	RecordHash(ctx context.Context, hash string) (string, error)
}