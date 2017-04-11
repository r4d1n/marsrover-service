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
	"github.com/r4d1n/marsrover"
)

var c *marsrover.Client

func init() {
	c = marsrover.NewClient(os.Getenv("NASA_API_KEY"))
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
	rover := mux.Vars(r)["rover"]
	m, err := c.GetManifest(rover)
	if err != nil {
		log.Fatal(err)
	}
	data, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func getImagesBySol(w http.ResponseWriter, r *http.Request) {
	rover := mux.Vars(r)["rover"]
	sol, err := strconv.Atoi(mux.Vars(r)["sol"])
	if err != nil {
		log.Fatal(err)
	}
	p, err := c.GetImagesBySol(rover, sol)
	if err != nil {
		log.Fatal(err)
	}
	data, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func getImagesByEarthDate(w http.ResponseWriter, r *http.Request) {
	rover := mux.Vars(r)["rover"]
	date := mux.Vars(r)["date"]
	p, err := c.GetImagesByEarthDate(rover, date)
	if err != nil {
		log.Fatal(err)
	}
	data, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
