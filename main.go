package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/3d0c/imagio/config"
	"github.com/3d0c/imagio/imgproc"
	"github.com/3d0c/imagio/query"
	. "github.com/3d0c/imagio/utils"
	"github.com/golang/groupcache"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

var cacheGroup *groupcache.Group

func initCacheGroup() {

	pool := groupcache.NewHTTPPool(config.Get().Listen())
	pool.Set(config.Get().CachePeers()...)

	cacheGroup = groupcache.NewGroup("imagio-storage", config.Get().CacheSize(), groupcache.GetterFunc(
		func(ctx context.Context, key string, dest groupcache.Sink) error {
			log.Println("cache miss")
			dest.SetBytes(imgproc.Do(
				Construct(new(query.Options), key).(*query.Options),
			))
			return nil
		}),
	)
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}

func main() {
	dumpcfg := flag.Bool("dumpcfg", false, "Dump config.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage: %s [OPTIONS]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *dumpcfg {
		config.Get().DumpCfg()
		os.Exit(0)
	}

	initCacheGroup()

	log.Printf("Service listen on %v\n", config.Get().Listen())

	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			var data []byte
			var ctx context.Context

			cacheGroup.Get(ctx, r.URL.String(), groupcache.AllocatingByteSliceSink(&data))

			http.ServeContent(w, r, r.URL.String(), time.Now(), bytes.NewReader(data))
		},
	)

	http.HandleFunc("/nocache",
		func(w http.ResponseWriter, r *http.Request) {
			var result []byte
			result = imgproc.Do(Construct(new(query.Options), r.URL).(*query.Options))
			w.Write(result)
		},
	)

	http.HandleFunc("/stat",
		func(w http.ResponseWriter, r *http.Request) {
			// awesome stat. not implemented yet.
		},
	)

	log.Fatal(http.ListenAndServe(config.Get().Listen(), nil))
}
