package routes

import (
	"net/http"

	"github.com/jordanjohnston/desuplayer_v2/routes/library"
	"github.com/jordanjohnston/desuplayer_v2/routes/music"
)

func SetUpRequestHandlers(m *http.ServeMux) {
	for route, handler := range music.Routes() {
		m.HandleFunc(route, handler)
	}
	for route, handler := range library.Routes() {
		m.HandleFunc(route, handler)
	}
}
