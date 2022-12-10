package web

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"mapserver/coords"
	"net/http"
	"strconv"
)

type ViewBlock struct {
	BlockMapping map[int]string `json:"blockmapping"`
	ContentId    []int          `json:"contentid"`
}

func (api *Api) GetBlockData(resp http.ResponseWriter, req *http.Request) {
	// /api/viewblock/{x}/{y}/{z}
	vars := mux.Vars(req)

	if len(vars) != 3 {
		resp.WriteHeader(500)
		resp.Write([]byte("wrong number of arguments"))
		return
	}

	x, _ := strconv.Atoi(vars["x"])
	y, _ := strconv.Atoi(vars["y"])
	z, _ := strconv.Atoi(vars["z"])

	c := coords.NewMapBlockCoords(x, y, z)
	mb, err := api.Context.MapBlockAccessor.GetMapBlock(c)

	if err != nil {
		resp.WriteHeader(500)
		resp.Write([]byte(err.Error()))

	} else {

		var vb *ViewBlock
		if mb != nil {
			vb = &ViewBlock{}
			vb.BlockMapping = mb.BlockMapping
			vb.ContentId = mb.Mapdata.ContentId
		}

		resp.Header().Add("content-type", "application/json")
		json.NewEncoder(resp).Encode(vb)

	}
}
