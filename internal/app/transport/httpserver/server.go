package httpserver

import "biggest-mars-pictures/internal/app/services"

type HttpServer struct {
	largestPictureService services.PictureServiceInterface
}

func NewHttpServer(largestPictureService services.PictureServiceInterface) *HttpServer {
	return &HttpServer{
		largestPictureService: largestPictureService,
	}
}
