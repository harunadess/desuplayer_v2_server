package routes

import (
	"net/http"

	"github.com/jordanjohnston/desuplayer_v2/routes/library"
	"github.com/jordanjohnston/desuplayer_v2/routes/music"
	"github.com/jordanjohnston/desuplayer_v2/routes/util"
)

var routes = []map[string]util.RequestHandler{music.Routes(), library.Routes()}

func SetUpRequestHandlers(m *http.ServeMux) {
	// for route, handler := range music.Routes() {
	// 	m.HandleFunc(route, handler)
	// }
	// for route, handler := range library.Routes() {
	// 	m.HandleFunc(route, handler)
	// }
	for _, route := range routes {
		for path, handler := range route {
			m.HandleFunc(path, handler)
		}
	}
}
