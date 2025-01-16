package services

import (
	"biggest-mars-pictures/internal/app/domain"
	"context"
)

type PictureServiceInterface interface {
	FindLargestPicture(ctx context.Context, sol string) (*domain.NasaPicture, error)
}
