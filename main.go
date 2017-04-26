package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
	"github.com/r4d1n/marsrover"
)

var mars *marsrover.Client
var c *cache.Cache

func init() {
	mars = marsrover.NewClient(os.Getenv("NASA_API_KEY"))
	c = cache.New(60*time.Minute, 60*time.Minute)
}

func main() {
	port := flag.Int("port", 3333, "the port that the service should listen on")
	r := mux.NewRouter()
	r.HandleFunc("/mars/manifest/{rover}", getManifest)
	r.HandleFunc("/mars/photos/{rover}/sol/{sol}", getImagesBySol)
	r.HandleFunc("/mars/photos/{rover}/earthdate/{date}", getImagesByEarthDate)
	flag.Parse()
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Printf("serving on port %d \n", *port)
	fmt.Println(server.ListenAndServe())
}

func getManifest(w http.ResponseWriter, r *http.Request) {
	var data *marsrover.Manifest
	rover := mux.Vars(r)["rover"]
	key := fmt.Sprintf("manifest:%s", rover)
	if x, found := c.Get(key); found {
		fmt.Printf("manifest:%s is in the cache \n", rover)
		data = x.(*marsrover.Manifest)
	} else {
		fmt.Printf("manifest:%s is NOT in the cache \n", rover)
		var err error
		data, err = mars.GetManifest(rover)
		if err != nil {
			fmt.Println(err)
		}
		c.Set(key, data, cache.DefaultExpiration)
	}
	json, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func getImagesBySol(w http.ResponseWriter, r *http.Request) {
	var data *marsrover.PhotoResponse
	rover := mux.Vars(r)["rover"]
	sol, err := strconv.Atoi(mux.Vars(r)["sol"])
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
	}
	key := fmt.Sprintf("sol:%s:%d", rover, sol)
	if x, found := c.Get(key); found {
		fmt.Printf("%s is in the cache \n", key)
		data = x.(*marsrover.PhotoResponse)
	} else {
		fmt.Printf("%s is NOT in the cache \n", key)
		var err error
		data, err = mars.GetImagesBySol(rover, sol)
		if err != nil {
			fmt.Println(err)
			http.NotFound(w, r)
		}
		c.Set(key, data, cache.DefaultExpiration)
	}
	json, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func getImagesByEarthDate(w http.ResponseWriter, r *http.Request) {
	var data *marsrover.PhotoResponse
	rover := mux.Vars(r)["rover"]
	date := mux.Vars(r)["date"]
	key := fmt.Sprintf("date:%s:%s", rover, date)
	if x, found := c.Get(key); found {
		fmt.Printf("%s is in the cache \n", key)
		data = x.(*marsrover.PhotoResponse)
	} else {
		fmt.Printf("%s is NOT in the cache \n", key)
		var err error
		data, err = mars.GetImagesByEarthDate(rover, date)
		if err != nil {
			fmt.Println(err)
			http.NotFound(w, r)
		}
		c.Set(key, data, cache.DefaultExpiration)
	}
	json, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
