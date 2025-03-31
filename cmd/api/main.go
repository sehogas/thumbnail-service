package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/nfnt/resize"

	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
)

var (
	Version string = "development"
)

func main() {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3010"
	}

	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(10000000),
	)
	if err != nil {
		log.Fatal(err)
	}

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(10*time.Minute),
		cache.ClientWithRefreshKey("opn"),
	)
	if err != nil {
		log.Fatal(err)
	}

	handler := http.HandlerFunc(thumbnailHandler)

	mux := http.NewServeMux()

	mux.Handle("/", cacheClient.Middleware(handler))
	mux.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		response := struct {
			Version string `json:"version"`
		}{
			Version: Version,
		}
		json, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	})
	log.Printf("Starting server on port %v", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PORT), mux))
}

func thumbnailHandler(w http.ResponseWriter, r *http.Request) {
	width := r.URL.Query().Get("width")
	if width == "" {
		width = "80"
	}
	intWidth, err := strconv.ParseUint(width, 10, 32)
	if err != nil {
		intWidth = 80
	}

	height := r.URL.Query().Get("height")
	if height == "" {
		height = "0"
	}
	intHeight, err := strconv.ParseUint(height, 10, 32)
	if err != nil {
		intHeight = 0
	}

	filename := r.URL.Query().Get("filename")
	if filename == "" {
		filename = "filename_thumbnail.jpg"
	}

	urlImage := r.URL.Query().Get("url")
	if urlImage == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	urlImageUnescape, err := url.QueryUnescape(urlImage)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	res, err := http.Get(urlImageUnescape)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Println("status code: ", res.StatusCode)
		http.Error(w, "", res.StatusCode)
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if len(body) < 50 {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	img, err := jpeg.Decode(bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusNotFound)
		return
	}

	m := resize.Resize(uint(intWidth), uint(intHeight), img, resize.Lanczos3)

	var buffer bytes.Buffer
	if err := jpeg.Encode(&buffer, m, nil); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", filename))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

// func create_thumbnail(filename string, width uint, height uint, newFilename string) error {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return err
// 	}

// 	img, err := jpeg.Decode(file)
// 	if err != nil {
// 		return err
// 	}
// 	file.Close()

// 	m := resize.Resize(width, height, img, resize.Lanczos3)

// 	out, err := os.Create(newFilename)
// 	if err != nil {
// 		return err
// 	}
// 	defer out.Close()

// 	// write new image to file
// 	return jpeg.Encode(out, m, nil)
// }
