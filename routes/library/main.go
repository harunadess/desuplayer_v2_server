package library

import (
	"encoding/json"
	"log"
	"net/http"

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
	withImages := query.Get("images")
	getImages := withImages == "true"

	library.UnloadLibrary()
	musicLibrary, err := directoryscaper.GetAllInDirectory(baseDir, getImages)
	if err != nil {
		log.Println("error getting music library ", err)
		w.WriteHeader(http.StatusInternalServerError)
		// todo: better handling
		w.Write([]byte("error getting music library"))
		return
	}
	library.SetLibrary(musicLibrary)

	jsonifiedLib, err := json.Marshal(musicLibrary)
	if err != nil {
		log.Println("error converting music library to json ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error converting music library to json"))
		return
	}
	w.Write(jsonifiedLib)
	writeLibraryToJSON(musicLibrary)
}

func writeLibraryToJSON(library directoryscaper.MusicLibrary) {
	fileio.WriteToJSON(library, fileio.AbsPath("/library.json"))
}

func getLibrary(w http.ResponseWriter, r *http.Request) {
	lib := library.GetLibrary()
	if lib == nil {
		log.Println("library does not exist (need build)")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("library does not exist. please build library first."))
		return
	}
	jsonifiedLib, err := json.Marshal(lib)
	if err != nil {
		log.Println("error converting music library to json ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error converting music library to json"))
		return
	}
	log.Println("jsonifiedLib len ", len(jsonifiedLib))
	w.Write(jsonifiedLib)
}
