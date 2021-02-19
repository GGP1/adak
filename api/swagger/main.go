package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	//go:embed index.html
	indexHTML string
	port      = flag.Int64("port", 8080, "server port")
)

func main() {
	flag.Parse()

	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Serving API docs on http://localhost%s\n", addr)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.Write([]byte(indexHTML))
	})

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed listening on %s: %v", addr, err)
	}
}
