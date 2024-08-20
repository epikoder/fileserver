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
	_port := "8100"
	for i, v := range os.Args {
		if v == "-p" {
			_port = os.Args[i+1]
		}
	}

	port := flag.String("p", _port, "port to serve on")
	directory := flag.String("d", DIRECTORY, "the directory of static file to host")
	flag.Parse()

	os.Mkdir(DIRECTORY, os.ModePerm)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			http.FileServer(http.Dir(*directory)).ServeHTTP(w, r)
			return
		}

		// Ensure the request is a multipart form-data request
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			http.Error(w, "Expected multipart/form-data content type", http.StatusBadRequest)
			return
		}

		log.Println(r.Header)
		f, h, err := r.FormFile("file")
		if err != nil {
			log.Println("Error retrieving the file:", err)
			http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		b, err := io.ReadAll(f)
		if err != nil {
			log.Println("Error reading the file:", err)
			http.Error(w, "Error reading the file", http.StatusInternalServerError)
			return
		}

		arr := strings.Split(h.Filename, ".")
		name := uuid.NewString() + "." + arr[len(arr)-1]
		err = os.WriteFile(DIRECTORY+name, b, os.ModePerm)
		if err != nil {
			log.Println("Error saving the file:", err)
			http.Error(w, "Error saving the file", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(name))
	})
	http.Handle("/", mux)
	log.Println("Starting FileServer...")
	address := "0.0.0.0:" + *port
	log.Printf("Serving %s on HTTP %s\n", *directory, address)
	log.Fatal(http.ListenAndServe(address, nil))
}
