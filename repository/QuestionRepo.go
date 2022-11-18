package repository

import (
	"context"

	"enerbit.com/go/grpc/models"
)

func SetQuestion(ctx context.Context, question *models.Question) error {
	return implementation.SetQuestion(ctx, question)
}
