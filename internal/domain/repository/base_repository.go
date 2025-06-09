package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
)

type BaseRepository[T any] interface {
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*T, error)
	FindAll(ctx context.Context, page, limit int) ([]*T, error)
	FindOne(ctx context.Context, query interface{}) (*T, error)
	FindMany(ctx context.Context, query interface{}, page, limit int) ([]*T, error)
	Count(ctx context.Context, query interface{}) (int64, error)
}

type baseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
	return &baseRepository[T]{db: db}
}

func (r *baseRepository[T]) Create(ctx context.Context, entity *T) error {
	err := r.db.WithContext(ctx).Create(entity).Error
	if err != nil {
		return fmt.Errorf("failed to create entity: %w", err)
	}
	return nil
}

func (r *baseRepository[T]) Update(ctx context.Context, entity *T) error {
	err := r.db.WithContext(ctx).Save(entity).Error
	if err != nil {
		return fmt.Errorf("failed to update entity: %w", err)
	}
	return nil
}

func (r *baseRepository[T]) Delete(ctx context.Context, id string) error {
	var entity T
	err := r.db.WithContext(ctx).Delete(&entity, id).Error
	if err != nil {
		return fmt.Errorf("failed to delete entity with ID %s: %w", id, err)
	}
	return nil
}

func (r *baseRepository[T]) FindByID(ctx context.Context, id string) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		return nil, fmt.Errorf("entity with ID %s not found: %w", id, err)
	}
	return &entity, nil
}

func (r *baseRepository[T]) FindAll(ctx context.Context, page, limit int) ([]*T, error) {
	offset := (page - 1) * limit
	var entities []*T
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&entities).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find all entities: %w", err)
	}
	return entities, nil
}

func (r *baseRepository[T]) FindOne(ctx context.Context, query interface{}) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).Where(query).First(&entity).Error
	if err != nil {
		return nil, fmt.Errorf("entity not found: %w", err)
	}
	return &entity, nil
}

func (r *baseRepository[T]) FindMany(ctx context.Context, query interface{}, page, limit int) ([]*T, error) {
	offset := (page - 1) * limit
	var entities []*T
	err := r.db.WithContext(ctx).Where(query).Offset(offset).Limit(limit).Find(&entities).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find many entities: %w", err)
	}
	return entities, nil
}

func (r *baseRepository[T]) Count(ctx context.Context, query interface{}) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Where(query).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count entities: %w", err)
	}
	return count, nil
}
