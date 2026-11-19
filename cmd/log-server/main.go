// Package main serves sample NGINX logs for local development.
package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"time"
)

// main запускает локальный HTTP-сервер с примерами логов.
func main() {
	port := flag.String("port", "8000", "server port")
	flag.Parse()

	address := net.JoinHostPort("127.0.0.1", *port)
	server := &http.Server{
		Addr:              address,
		Handler:           http.FileServer(http.Dir("sample_logs")),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("sample logs: http://%s/access.txt", address)
	log.Fatal(server.ListenAndServe())
}
