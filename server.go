package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"os"
)

var db = MemoryDB{}
var port = os.Getenv("PORT")

func init() {
	db.Init()
	if port == "" {
		port = "4444"
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

func main() {
	service := new(Service)
	err := rpc.Register(service)
	if err != nil {
		log.Fatal(":)")
	}
	http.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		enableCors(&w)
		w.Header().Set("Content-Type", "application/json")
		res := NewRPCRequest(r.Body).Call()
		io.Copy(w, res)
	})
	log.Println("Started rpc server")
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
