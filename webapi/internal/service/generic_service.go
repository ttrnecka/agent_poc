package service

import (
	"context"
	"fmt"

	"github.com/ttrnecka/agent_poc/webapi/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DependencyDeleteFunc func(ctx context.Context, parentID primitive.ObjectID) error

type GenericService[T any] interface {
	All(context.Context) ([]T, error)
	Get(context.Context, string) (*T, error)
	Delete(context.Context, string) error
	Update(context.Context, primitive.ObjectID, *T) (primitive.ObjectID, error)
	RegisterDependencies(...DependencyDeleteFunc)
}

type genericService[T any] struct {
	MainRepo    repository.GenericRepository[T]
	dependecies []DependencyDeleteFunc
}

func NewGenericService[T any](r repository.GenericRepository[T]) GenericService[T] {
	return &genericService[T]{
		MainRepo: r,
	}
}

func (s *genericService[T]) All(ctx context.Context) ([]T, error) {
	return s.MainRepo.All(ctx)
}

func (s *genericService[T]) Get(ctx context.Context, id string) (*T, error) {
	idp, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return s.MainRepo.GetByID(ctx, idp)
}

func (s *genericService[T]) Delete(ctx context.Context, id string) error {
	idp, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	err = s.DeleteDependencies(ctx, idp)
	if err != nil {
		return err
	}
	return s.MainRepo.HardDeleteByID(ctx, idp)
}

func (s *genericService[T]) Update(ctx context.Context, id primitive.ObjectID, item *T) (primitive.ObjectID, error) {
	if id.IsZero() {
		return s.MainRepo.Create(ctx, item)
	}
	return id, s.MainRepo.UpdateByID(ctx, id, item)
}

func (s *genericService[T]) DeleteDependencies(ctx context.Context, parentID primitive.ObjectID) error {
	for _, fn := range s.dependecies {
		if err := fn(ctx, parentID); err != nil {
			return fmt.Errorf("failed to delete dependency: %w", err)
		}
	}
	return nil
}

func (s *genericService[T]) RegisterDependencies(fn ...DependencyDeleteFunc) {
	s.dependecies = append(s.dependecies, fn...)
}
