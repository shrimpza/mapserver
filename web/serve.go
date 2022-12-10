package web

import (
	"embed"
	"mapserver/app"
	"mapserver/public"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func Serve(ctx *app.App) {
	fields := logrus.Fields{
		"bindhost": ctx.Config.BindHost,
		"port":     ctx.Config.Port,
		"webroot":  ctx.Config.WebRoot,
		"webdev":   ctx.Config.Webdev,
	}
	logrus.WithFields(fields).Info("Starting http server")

	root := ctx.Config.WebRoot
	apiroot := root + "api/"

	api := NewApi(ctx)
	router := mux.NewRouter()
	apirouter := mux.NewRouter()

	tiles := &Tiles{ctx: ctx}
	tiles.Init()
	apirouter.Handle(apiroot+"tile/{layerId}/{x}/{y}/{zoom}", tiles)
	apirouter.HandleFunc(apiroot+"config", api.GetConfig)
	apirouter.HandleFunc(apiroot+"stats", api.GetStats)
	apirouter.HandleFunc(apiroot+"media/{filename}", api.GetMedia)
	apirouter.HandleFunc(apiroot+"minetest", api.PostMinetestData)
	apirouter.HandleFunc(apiroot+"mapobjects/", api.QueryMapobjects)
	apirouter.HandleFunc(apiroot+"colormapping", api.GetColorMapping)
	apirouter.HandleFunc(apiroot+"viewblock/{x}/{y}/{z}", api.GetBlockData)

	if ctx.Config.MapObjects.Areas {
		apirouter.Handle(apiroot+"areas", &AreasHandler{ctx: ctx})
	}

	ws := NewWS(ctx)
	apirouter.Handle(apiroot+"ws", ws)

	ctx.Tilerenderer.Eventbus.AddListener(ws)
	ctx.WebEventbus.AddListener(ws)

	if ctx.Config.WebApi.EnableMapblock {
		//mapblock endpoint
		apirouter.HandleFunc(apiroot+"mapblock/", api.GetMapBlockData)
	}

	router.PathPrefix(apiroot).HandlerFunc(apirouter.ServeHTTP)

	if ctx.Config.EnablePrometheus {
		router.Handle(root+"metrics", promhttp.Handler())
	}

	// static files
	if ctx.Config.Webdev {
		logrus.Print("using live mode")
		fs := http.FileServer(http.FS(os.DirFS("public")))
		router.Handle(root, http.StripPrefix(root, fs))

	} else {
		logrus.Print("using embed mode")
		fs := http.FileServer(http.FS(public.Files))
		router.PathPrefix(root).Handler(http.StripPrefix(root, fs))
		//mux.HandleFunc("/", CachedServeFunc(fs.ServeHTTP))
	}

	// main entry point
	http.HandleFunc("/", router.ServeHTTP)
	err := http.ListenAndServe(ctx.Config.BindHost+":"+strconv.Itoa(ctx.Config.Port), nil)
	if err != nil {
		panic(err)
	}
}

func getFileSystem(useLocalfs bool, content embed.FS) http.FileSystem {
	if useLocalfs {
		log.Print("using live mode")
		return http.FS(os.DirFS("public"))
	}

	log.Print("using embed mode")
	return http.FS(content)
}
