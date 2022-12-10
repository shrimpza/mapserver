package web

import (
	"github.com/gorilla/mux"
	"mapserver/public"
	"net/http"
)

func (api *Api) GetMedia(resp http.ResponseWriter, req *http.Request) {
	// /api/media/{filename}/{x}/{y}/{zoom}
	vars := mux.Vars(req)

	if len(vars) != 1 {
		resp.WriteHeader(500)
		resp.Write([]byte("wrong number of arguments"))
		return
	}

	filename := vars["filename"]
	fallback, hasfallback := req.URL.Query()["fallback"]

	content := api.Context.MediaRepo[filename]

	if content == nil && hasfallback && len(fallback) > 0 {
		var err error
		content, err = public.Files.ReadFile("pics/" + fallback[0])
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if content != nil {
		resp.Write(content)
		resp.Header().Add("content-type", "image/png")
		return
	}

	resp.WriteHeader(404)
	resp.Write([]byte(filename))
}
