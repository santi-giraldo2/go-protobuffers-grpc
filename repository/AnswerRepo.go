package repository

import (
	"context"

	"enerbit.com/go/grpc/models"
)

func SetAnswer(ctx context.Context, answer *models.Answer) error {
	return implementation.SetAnswer(ctx, answer)
}
