package middleware

import (
	"fmt"
	"image_server/server"
	"image_server/utils"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}
type Routes []Route

// assign handlers to all routes
func NewRouter() *mux.Router {
	router := mux.NewRouter()

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = middleware(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

//middleware: dump incoming request
func middleware(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && strings.Contains(r.URL.Path, "/images") {
			utils.Log(fmt.Sprintf("Request URL: %v", r.URL.Path))
		} else {
			req_dump, _ := httputil.DumpRequest(r, true)
			utils.Log(fmt.Sprintf("Request Dump: %v", string(req_dump)))
		}
		inner.ServeHTTP(w, r)
	})
}

// routes
var routes = Routes{
	Route{
		"CreateAlbum",
		"POST",
		"/albums",
		server.CreateAlbum,
	},
	Route{
		"DeleteAlbum",
		"DELETE",
		"/albums/{album_name}",
		server.DeleteAlbum,
	},
	Route{
		"UploadImage",
		"POST",
		"/albums/{album_name}/images",
		server.UploadImage,
	},
	Route{
		"GetAllImages",
		"GET",
		"/albums/{album_name}/images",
		server.GetAllImages,
	},
	Route{
		"GetImage",
		"GET",
		"/albums/{album_name}/images/{image_name}",
		server.DownloadImage,
	},
	Route{
		"DeleteImage",
		"Delete",
		"/albums/{album_name}/images/{image_name}",
		server.DeleteImage,
	},
}
