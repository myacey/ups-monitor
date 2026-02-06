package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"sync/atomic"
)

var (
	portToConn = flag.String("port", "COM7", "serial port")
	frequency  = flag.Int("freq", 5000, "frequency to get port data in ms")
)

var lastStatus atomic.Value

func main() {
	flag.Parse()

	go readUPSStatus()

	http.HandleFunc("/api/v1/stats", func(w http.ResponseWriter, r *http.Request) {
		status := lastStatus.Load()
		if status == nil {
			http.Error(w, "no data yet", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	})

	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Println("HTTP server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("http server failed:", err)
	}
}
