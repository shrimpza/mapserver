package web

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"image/color"
	"mapserver/app"
	"mapserver/coords"
	"mapserver/tilerenderer"
	"net/http"
	"strconv"
)

type Tiles struct {
	ctx   *app.App
	blank []byte
}

func (t *Tiles) Init() {
	t.blank = tilerenderer.CreateBlankTile(color.RGBA{255, 255, 255, 255})
}

func (t *Tiles) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// /api/tile/{layerId}/{x}/{y}/{zoom}
	vars := mux.Vars(req)

	if len(vars) != 4 {
		resp.WriteHeader(500)
		resp.Write([]byte("wrong number of arguments"))
		return
	}

	timer := prometheus.NewTimer(tileServeDuration)
	defer timer.ObserveDuration()

	layerid, _ := strconv.Atoi(vars["layerId"])
	x, _ := strconv.Atoi(vars["x"])
	y, _ := strconv.Atoi(vars["y"])
	zoom, _ := strconv.Atoi(vars["zoom"])

	c := coords.NewTileCoords(x, y, zoom, layerid)
	tile, err := t.ctx.TileDB.GetTile(c)

	if err != nil {
		resp.WriteHeader(500)
		resp.Write([]byte(err.Error()))

	} else {
		resp.Header().Add("content-type", "image/png")

		if tile == nil {
			resp.Write(t.blank)
			//TODO: cache/layer color

		} else {
			tilesCumulativeSize.Add(float64(len(tile)))
			resp.Write(tile)

		}
	}
}
