package domain

import "biggest-mars-pictures/internal/app/repository"

type NasaPicture struct {
	Url  string
	Size int64
}

func ToNasaPicture(repoImage repository.Image) *NasaPicture {
	return &NasaPicture{
		Url:  repoImage.Url,
		Size: repoImage.Size,
	}
}

func ToRepoImage(nasaPicture *NasaPicture, sol int) *repository.Image {
	return &repository.Image{
		Sol:  sol,
		Url:  nasaPicture.Url,
		Size: nasaPicture.Size,
	}
}
