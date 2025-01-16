package repository

type ImageRepoInterface interface {
	GetImageById(sol int) (*Image, error)
	SaveImage(image *Image) error
}
