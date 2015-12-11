package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

func CrawlPornHub(start, stop uint64) {
	for i := start; i <= stop; i++ {
		p, err := crawl.GetBasePageByPageNumber(i)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		// save to json and dump to file
		bytes, err := json.Marshal(p)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		err = ioutil.WriteFile(Config.DumpPath+strconv.FormatUint(i, 10)+"_page.json", bytes, os.ModePerm)
		if err != nil {
			errors = append(errors, err)
		}
	}
}

func main() {
	flag.Parse()
	err := ParseConfigFromFile()
	if err != nil {
		panic(fmt.Sprintf("ERROR: error parsing the config from file. %v", err))
	}

	startTime = time.Now()

	log.Printf("Listening at http://127.0.0.1%v ...\n", Config.portString)
	log.Fatal(http.ListenAndServe(Config.portString, nil))
}
