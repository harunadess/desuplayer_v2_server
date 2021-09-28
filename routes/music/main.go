package music

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jordanjohnston/desuplayer_v2/library"
	"github.com/jordanjohnston/desuplayer_v2/routes/util"
)

const subRoute = "music"

// TODO: other routes
func Routes() map[string]util.RequestHandler {
	return map[string]util.RequestHandler{
		util.FormatRoute(util.BaseRoute, subRoute, "getSong"):       getSong,
		util.FormatRoute(util.BaseRoute, subRoute, "getSongMeta"):   getSongMeta,
		util.FormatRoute(util.BaseRoute, subRoute, "getAllArtists"): getAllArtists,
	}
}

func writeResponse(a interface{}, w http.ResponseWriter) {
	bytes, err := json.Marshal(a)
	if err != nil {
		log.Println("error writing response ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
		return
	}
	w.Write(bytes)
	bytes = nil
}

func getSong(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	song := query.Get("path")
	bs, ok := library.GetSong(song)
	if !ok {
		log.Printf("song does not exist %v", song)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("song not found"))
		return
	}
	w.Write(bs)
	bs = nil
}

func getSongMeta(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	path := query.Get("path")
	albumTitle := query.Get("album")
	albumArtist := query.Get("artist")

	meta, ok := library.GetSongMeta(path, albumTitle, albumArtist)
	if !ok {
		log.Printf("failed to get song with parameters: %v %v %v", albumTitle, albumArtist, path)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("song not found"))
		return
	}
	writeResponse(meta, w)
}

func getAllArtists(w http.ResponseWriter, r *http.Request) {
	bs := library.GetAllAlbums()
	writeResponse(bs, w)
}
