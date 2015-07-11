package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app  = kingpin.New("github-assigner", "")
	addr = app.Flag("addr", "Address to listen on for HTTP.").Default(":3000").OverrideDefaultFromEnvar("ADDR").String()
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	m := http.NewServeMux()

	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		before := time.Now()
		fmt.Printf("time=%q message=\"request\" request=%q\n", time.Now().Format(time.RFC3339Nano), r.URL)
		defer func() {
			after := time.Now()
			dur := after.Sub(before)

			fmt.Printf("time=%q message\"response\" request=%q took=%s took_ms=%d\n", after.Format(time.RFC3339Nano), r.URL, dur, dur)
		}()

		var v interface{}
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		spew.Dump(v)
	})

	if err := http.ListenAndServe(*addr, m); err != nil {
		panic(err)
	}
}
