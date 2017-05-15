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
	"github.com/r4d1n/marsrover"
)

var mars *marsrover.Client
var pool = newPool()

func init() {
	mars = marsrover.NewClient(os.Getenv("NASA_API_KEY"))
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
	var j []byte
	rover := mux.Vars(r)["rover"]
	key := fmt.Sprintf("manifest:%s", rover)
	conn := pool.Get()
	defer conn.Close()
	if reply, _ := conn.Do("GET", key); reply != nil {
		fmt.Printf("manifest:%s is in the cache \n", rover)
		j = reply.([]byte)
	} else {
		fmt.Printf("manifest:%s is NOT in the cache \n", rover)
		var data *marsrover.Manifest
		var err error
		data, err = mars.GetManifest(rover)
		j, err = json.Marshal(data)
		_, err = conn.Do("SET", key, j)
		if err != nil {
			fmt.Println(err)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func getImagesBySol(w http.ResponseWriter, r *http.Request) {
	var j []byte
	conn := pool.Get()
	defer conn.Close()
	rover := mux.Vars(r)["rover"]
	sol, err := strconv.Atoi(mux.Vars(r)["sol"])
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
	}
	key := fmt.Sprintf("sol:%s:%d", rover, sol)
	if reply, _ := conn.Do("GET", key); reply != nil {
		fmt.Printf("%s is in the cache \n", key)
		j = reply.([]byte)
	} else {
		fmt.Printf("%s is NOT in the cache \n", key)
		var data *marsrover.PhotoResponse
		var err error
		data, err = mars.GetImagesBySol(rover, sol)
		j, err = json.Marshal(data)
		_, err = conn.Do("SET", key, j)
		if err != nil {
			fmt.Println(err)
			http.NotFound(w, r)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func getImagesByEarthDate(w http.ResponseWriter, r *http.Request) {
	var j []byte
	conn := pool.Get()
	defer conn.Close()
	rover := mux.Vars(r)["rover"]
	date := mux.Vars(r)["date"]
	key := fmt.Sprintf("date:%s:%s", rover, date)
	if reply, _ := conn.Do("GET", key); reply != nil {
		fmt.Printf("%s is in the cache \n", key)
		j = reply.([]byte)
	} else {
		fmt.Printf("%s is NOT in the cache \n", key)
		var data *marsrover.PhotoResponse
		var err error
		data, err = mars.GetImagesByEarthDate(rover, date)
		j, err = json.Marshal(data)
		_, err = conn.Do("SET", key, j)
		if err != nil {
			fmt.Println(err)
			http.NotFound(w, r)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}
