package services

import (
	"biggest-mars-pictures/internal/app/clients/nasa"
	"biggest-mars-pictures/internal/app/domain"
	"biggest-mars-pictures/internal/app/repository"
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"sort"
	"strconv"
)

type LargestPictureService struct {
	client          *nasa.Client
	imageRepository repository.ImageRepoInterface
}

func NewLargestPictureService(client *nasa.Client, imageRepository repository.ImageRepoInterface) *LargestPictureService {
	return &LargestPictureService{
		client:          client,
		imageRepository: imageRepository,
	}
}

func (lps *LargestPictureService) FindLargestPicture(ctx context.Context, sol string) (*domain.NasaPicture, error) {
	solInt, err2 := strconv.Atoi(sol)
	if err2 != nil {
		return nil, err2
	}
	imageFromRepository, err2 := lps.imageRepository.GetImageById(solInt)
	if err2 == nil {
		fmt.Printf("found image from repository")
		return domain.ToNasaPicture(*imageFromRepository), nil
	}
	photos, err := lps.client.FindNasaPhotos(ctx, sol)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	nasaPhotoChannels := make(chan domain.NasaPicture, len(photos.Photos))

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(len(photos.Photos))
	for _, photo := range photos.Photos {
		g.Go(func() error {
			size, currError := lps.client.FindPhotoSize(ctx, photo.ImageSrc)
			if currError != nil {
				return fmt.Errorf("findLargestPicture %w", currError)
			} else {
				currNasaPicture := domain.NasaPicture{
					Url: photo.ImageSrc, Size: size,
				}
				select {
				case <-ctx.Done():
					fmt.Println("error was called so context is done")
					return ctx.Err()
				case nasaPhotoChannels <- currNasaPicture:
					return nil
				}
			}
		})
	}
	if errorGroup := g.Wait(); errorGroup != nil {
		return nil, fmt.Errorf("%w", errorGroup)
	}
	close(nasaPhotoChannels)
	nasaPhotos := make([]domain.NasaPicture, len(photos.Photos))
	for nasaPhoto := range nasaPhotoChannels {
		nasaPhotos = append(nasaPhotos, nasaPhoto)
	}
	sort.Slice(nasaPhotos, func(i, j int) bool {
		return nasaPhotos[i].Size >= nasaPhotos[j].Size
	})
	err2 = lps.imageRepository.SaveImage(domain.ToRepoImage(&nasaPhotos[0], solInt))
	if err2 != nil {
		return nil, fmt.Errorf("%w", err2)
	}
	return &nasaPhotos[0], nil
}
