package repository

import (
	"context"

	"gorm.io/gorm"
)

// Repository is a generic interface for basic CRUD operations.
type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	FindByID(ctx context.Context, id any) (*T, error)
	FindOne(ctx context.Context, where map[string]any) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, where map[string]any) error
}

// GormRepository is a generic GORM-backed implementation.
type GormRepository[T any] struct {
	db *gorm.DB
}

func NewGormRepository[T any](db *gorm.DB) *GormRepository[T] {
	return &GormRepository[T]{db: db}
}

func (r *GormRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *GormRepository[T]) FindByID(ctx context.Context, id any) (*T, error) {
	var t T
	if err := r.db.WithContext(ctx).First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *GormRepository[T]) FindOne(ctx context.Context, where map[string]any) (*T, error) {
	var t T
	if err := r.db.WithContext(ctx).Where(where).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *GormRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *GormRepository[T]) Delete(ctx context.Context, where map[string]any) error {
	var t T
	return r.db.WithContext(ctx).Where(where).Delete(&t).Error
}
