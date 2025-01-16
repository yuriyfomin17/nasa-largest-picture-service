package httpserver

import "net/http"

func (h *HttpServer) GetLargestPicture(w http.ResponseWriter, r *http.Request) {
	sol := r.URL.Query().Get("sol")

	picture, err := h.largestPictureService.FindLargestPicture(r.Context(), sol)
	if err != nil {
		BadRequest("invalid-picture", err, &w, r)
	}
	if err != nil {
		BadRequest("invalid-picture", err, &w, r)
		return
	}
	RespondOk(picture, &w)
}
