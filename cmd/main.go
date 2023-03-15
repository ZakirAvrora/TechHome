package main

import (
	"context"
	"github.com/ZakirAvrora/TechHome/config"
	"github.com/ZakirAvrora/TechHome/internals/db"
	"github.com/ZakirAvrora/TechHome/internals/service"
	"github.com/ZakirAvrora/TechHome/pkg/cache"
	"github.com/ZakirAvrora/TechHome/pkg/server"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

const DefaultCacheSize = 1000
const DefaultCacheTTl = 30 * time.Second

func main() {

	conf := config.NewConfig(".env")

	conn, err := db.DatabaseInit(conf.Database)
	if err != nil {
		log.Fatalln("error in database initialization:", err)
	}

	if err := db.JsonUploadToDb(conn, "links.json"); err != nil {
		log.Fatalln("error in json file loading to database:", err)
	}

	memCache := cache.NewMemoryCache(DefaultCacheSize, DefaultCacheTTl)
	go memCache.RunCleaner()
	defer memCache.StopCleaner()

	svr := &http.Server{Addr: "0.0.0.0:8081", Handler: service.NewService(conn, memCache)}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	//Run server with gracefully shutdown capability
	server.Run(serverCtx, serverStopCtx, svr)
	<-serverCtx.Done()
}
