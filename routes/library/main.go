package library

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/jordanjohnston/desuplayer_v2/directoryscaper"
	"github.com/jordanjohnston/desuplayer_v2/fileio"
	"github.com/jordanjohnston/desuplayer_v2/library"
	"github.com/jordanjohnston/desuplayer_v2/routes/util"
)

const subRoute = "library"

func Routes() map[string]util.RequestHandler {
	return map[string]util.RequestHandler{
		util.FormatRoute(util.BaseRoute, subRoute, "build"): buildLibrary,
		util.FormatRoute(util.BaseRoute, subRoute, "get"):   getLibrary,
	}
}

func buildLibrary(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	baseDir := query.Get("musicDir")

	musicLibrary, err := directoryscaper.GetAllInDirectory(baseDir, true)
	if err != nil {
		log.Println("error getting music library ", err)
		w.WriteHeader(http.StatusInternalServerError)
		// todo: better handling
		w.Write([]byte("error getting music library"))
		return
	}
	library.LoadLibrary()

	jsonifiedLib, err := json.Marshal(musicLibrary)
	if err != nil {
		log.Println("error converting music library to json ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error converting music library to json"))
		return
	}
	w.Write(jsonifiedLib)
}

func getLibrary(w http.ResponseWriter, r *http.Request) {
	file, err := fileio.ReadSingleFile("/library.json")
	if err != nil {
		log.Printf("error getting library: %v", err)
		if strings.Contains(err.Error(), "cannot find") {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("could not find library file - run 'build library' again." + err.Error()))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error getting track " + err.Error()))
		}
		return
	}
	w.Write(file)
}
