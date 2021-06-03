package library

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jordanjohnston/desuplayer_v2/library"
	"github.com/jordanjohnston/desuplayer_v2/routes/util"
)

const subRoute = "library"

func Routes() map[string]util.RequestHandler {
	return map[string]util.RequestHandler{
		util.FormatRoute(util.BaseRoute, subRoute, "build"): buildLibrary,
	}
}

func buildLibrary(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	basePath := query.Get("musicDir")
	// withImages := query.Get("images")
	// getImages := withImages == "true"

	library.UnloadLibrary()
	err := library.BuildLibrary(basePath)
	if err != nil {
		log.Println("error getting music library ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error getting music library"))
		return
	}

	artists := library.GetAllArtists()
	jsonifiedResponse, err := json.Marshal(artists)
	if err != nil {
		log.Println("error converting music library to json ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error converting music library to json"))
		return
	}
	w.Write(jsonifiedResponse)
}
