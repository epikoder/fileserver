package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
)

const DIRECTORY = "public/"

func main() {
	port := flag.String("p", "8100", "port to serve on")
	directory := flag.String("d", DIRECTORY, "the directory of static file to host")
	flag.Parse()

	os.Mkdir(DIRECTORY, os.ModePerm)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			http.FileServer(http.Dir(*directory)).ServeHTTP(w, r)
			return
		}
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(500)
			return
		}
		f, h, err := r.FormFile("file")
		if err != nil {
			w.WriteHeader(500)
			return
		}
		b, err := io.ReadAll(f)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		arr := strings.Split(h.Filename, ".")
		name := uuid.NewString() + "." + arr[len(arr)-1]
		os.WriteFile(DIRECTORY+name, b, os.ModePerm)
		w.Write([]byte(name))
	})
	http.Handle("/", mux)
	log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
