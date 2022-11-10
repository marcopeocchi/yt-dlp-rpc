package pkg

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"

	"golang.org/x/net/websocket"
)

// Package available variables
var (
	db     = MemoryDB{}
	port   = os.Getenv("PORT")
	driver = os.Getenv("YT_DLP_PATH")
)

func init() {
	db.Init()
	if port == "" {
		port = "4444"
	}
}

// Enable WebSockets as transport protocol
func serveWS(ws *websocket.Conn) {
	jsonrpc.ServeConn(ws)
}

func serveHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")
	res := NewRPCRequest(r.Body).Call()
	io.Copy(w, res)
}

// Enable CORS for every origin
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

// Run blocking HTTP Server with WS upgrader
func RunBlocking() {
	service := new(Service)
	err := rpc.Register(service)
	if err != nil {
		log.Fatal("Something has gone terribly wrong :)")
	}
	http.HandleFunc("/rpc", serveHTTP)
	http.Handle("/rpc-ws", websocket.Handler(serveWS))

	log.Println("Started RPC server")
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
