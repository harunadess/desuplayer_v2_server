package music

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jordanjohnston/desuplayer_v2/library"
	"github.com/jordanjohnston/desuplayer_v2/routes/util"
)

const subRoute = "music"

func Routes() map[string]util.RequestHandler {
	return map[string]util.RequestHandler{
		// util.FormatRoute(util.BaseRoute, subRoute, "getTrack"):      getTrack,
		// util.FormatRoute(util.BaseRoute, subRoute, "getArtist"):     getArtist,
		// util.FormatRoute(util.BaseRoute, subRoute, "getAlbum"):      getAlbum,
		// util.FormatRoute(util.BaseRoute, subRoute, "getGenre"):      getGenre,
		util.FormatRoute(util.BaseRoute, subRoute, "getAllArtists"): getAllArtists,
		// util.FormatRoute(util.BaseRoute, subRoute, "getAllAlbums"):  getAllAlbums,
		// util.FormatRoute(util.BaseRoute, subRoute, "getAllGenres"):  getAllGenres,
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
}

// func getTrack(w http.ResponseWriter, r *http.Request) {
// 	query := r.URL.Query()
// 	track := query.Get("track")
// 	bs, ok := library.GetSong(track)
// 	if !ok {
// 		log.Printf("track does not exist %v", track)
// 		w.WriteHeader(http.StatusNotFound)
// 		w.Write([]byte("track not found"))
// 		return
// 	}
// 	writeResponse(bs, w)
// }

// func getArtist(w http.ResponseWriter, r *http.Request) {
// 	query := r.URL.Query()
// 	artist := query.Get("artist")
// 	bs, ok := library.GetArtist(artist)
// 	if !ok {
// 		log.Printf("artist does not exist %v", artist)
// 		w.WriteHeader(http.StatusNotFound)
// 		w.Write([]byte("artist not found"))
// 		return
// 	}
// 	writeResponse(bs, w)
// }

// func getAlbum(w http.ResponseWriter, r *http.Request) {
// 	query := r.URL.Query()
// 	album := query.Get("album")
// 	bs, ok := library.GetAlbum(album)
// 	if !ok {
// 		log.Printf("album does not exist %v", album)
// 		w.WriteHeader(http.StatusNotFound)
// 		w.Write([]byte("album not found"))
// 		return
// 	}
// 	writeResponse(bs, w)
// }

// func getGenre(w http.ResponseWriter, r *http.Request) {
// 	query := r.URL.Query()
// 	genre := query.Get("genre")
// 	bs, ok := library.GetGenre(genre)
// 	if !ok {
// 		log.Printf("genre does not exist %v", genre)
// 		w.WriteHeader(http.StatusNotFound)
// 		w.Write([]byte("genre not found"))
// 		return
// 	}
// 	writeResponse(bs, w)
// }

func getAllArtists(w http.ResponseWriter, r *http.Request) {
	bs := library.GetAllArtists()
	writeResponse(bs, w)
}

// func getAllAlbums(w http.ResponseWriter, r *http.Request) {
// 	bs, _ := library.GetAllAlbums()
// 	writeResponse(bs, w)
// }

// func getAllGenres(w http.ResponseWriter, r *http.Request) {
// 	bs, _ := library.GetAllGenres()
// 	writeResponse(bs, w)
// }
