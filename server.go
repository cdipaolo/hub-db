package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"time"

	"github.com/cdipaolo/hub-db/crawl"
)

var (
	pageCount uint64
	errors    []error
	running   bool
	startTime time.Time
)

func init() {
	errors = []error{}

	http.Handle("/", Get(HandleStatus))
}

// CrawlPornHub crawls pages from pornhub
// for images, starting at 'start' (inclusive)
// and ending at 'stop' (exclusive). Pages
// are dumped to the filepath described in
// the config file
func CrawlPornHub(start, stop uint64) {
	running = true
	for i := start; i < stop; i++ {
		p, err := crawl.GetBasePageByPageNumber(i)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		log.Printf("CRAWL [page = %v, albums = %v]", i, len(p.Albums))

		// save to json and dump to file
		bytes, err := json.Marshal(p)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		err = ioutil.WriteFile(path.Join(Config.DumpPath, strconv.FormatUint(i, 10)+"_page.json"), bytes, os.ModePerm)
		if err != nil {
			errors = append(errors, err)
		}
		pageCount++

		time.Sleep(Config.TimeDelaySeconds * time.Second)
	}
	running = false
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()
	err := ParseConfigFromFile()
	if err != nil {
		panic(fmt.Sprintf("ERROR: error parsing the config from file. %v", err))
	}

	startTime = time.Now()

	/*delta := (Config.EndPage - Config.StartPage) / uint64(runtime.NumCPU())
	for i := Config.StartPage; i < Config.EndPage-1; i += delta {
		log.Printf("INIT : 1 Crawler crawling [%v, %v)", i, i+delta)
		go CrawlPornHub(i, i+delta)
	}*/
	go CrawlPornHub(Config.StartPage, Config.EndPage)

	log.Printf("Listening at http://127.0.0.1%v ...\n", Config.portString)
	log.Fatal(http.ListenAndServe(Config.portString, nil))
}
