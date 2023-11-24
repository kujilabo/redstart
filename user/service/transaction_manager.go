package service

import "context"

type TransactionManager interface {
	Do(ctx context.Context, fn func(rf RepositoryFactory) error) error
}
