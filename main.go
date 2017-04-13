package main

import (
	"encoding/json"
	"fmt"
	"log"
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
	r := mux.NewRouter()
	r.HandleFunc("/mars/manifest/{rover}", getManifest)
	r.HandleFunc("/mars/photos/{rover}/sol/{sol}", getImagesBySol)
	r.HandleFunc("/mars/photos/{rover}/earthdate/{date}", getImagesByEarthDate)
	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("serving on port 8080")
	log.Fatal(server.ListenAndServe())
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
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func getImagesBySol(w http.ResponseWriter, r *http.Request) {
	var data *marsrover.PhotoResponse
	rover := mux.Vars(r)["rover"]
	sol, err := strconv.Atoi(mux.Vars(r)["sol"])
	if err != nil {
		log.Fatal(err)
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
		}
		c.Set(key, data, cache.DefaultExpiration)
	}
	json, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
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
		}
		c.Set(key, data, cache.DefaultExpiration)
	}
	json, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
