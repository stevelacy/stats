package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	// Client is the global redis pool
	Client = newClient()
	port   = ":4004"
)

var transPixel = "\x47\x49\x46\x38\x39\x61\x01\x00\x01\x00\x80\x00\x00\x00\x00\x00\x00\x00\x00\x21\xF9\x04\x01\x00\x00\x00\x00\x2C\x00\x00\x00\x00\x01\x00\x01\x00\x00\x02\x02\x44\x01\x00\x3B"

// StatsResponse is the total count of each statistic
type StatsResponse struct {
	Count int `json:"count"`
}

func main() {
	router := httprouter.New()
	router.GET("/:domain", Record)
	router.GET("/:domain/stats", Stats)
	router.GET("/:domain/stat.gif", RecordGif)
	router.POST("/:domain", Record)
	fmt.Println("starting on", port)
	log.Fatal(http.ListenAndServe(port, router))
}

// Record gets/sets the count to the url
func Record(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	url := ps.ByName("domain")
	if url == "/favicon.ico" {
		NotFound(w, r, ps)
		return
	}

	val, err := Client.Get(url).Result()
	if err == redis.Nil {
		set(url, 1)
		return
	}
	if err != nil {
		return
	}
	converted, err := strconv.Atoi(val)
	if err != nil {
		return
	}
	set(url, converted+1)
}

// RecordGif returns a gif image after incrementing the count
func RecordGif(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	Record(w, r, ps)
	w.Header().Set("Content-Type", "image/gif")
	w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	w.Header().Set("Expires", "Wed, 11 Nov 1998 11:11:11 GMT")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Pragma", "no-cache")
	fmt.Fprintf(w, transPixel)
}

// Stats returns the statistics for that domain/route
func Stats(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	url := ps.ByName("domain")
	val, err := Client.Get(url).Result()
	if err == redis.Nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		return
	}
	query := r.URL.Query()
	if query.Get("output") == "text" {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, val)
		return
	}

	converted, err := strconv.Atoi(val)
	if err != nil {
		return
	}
	stats := StatsResponse{
		Count: converted,
	}
	res, err := json.Marshal(stats)
	w.Write([]byte(res))
}

// NotFound returns when route is not found
func NotFound(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)
}

func set(url string, val int) error {
	return Client.Set(url, val, 0).Err()
}

func newClient() *redis.Client {
	connURL := "localhost:6379"
	if os.Getenv("REDIS_URL") != "" {
		connURL = os.Getenv("REDIS_URL")
	}
	return redis.NewClient(&redis.Options{
		Addr: connURL,
		DB:   10,
	})
}
