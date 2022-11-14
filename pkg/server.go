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
	"goytdlp.rpc/m/pkg/cli"
)

// Package available variables
var (
	db     = MemoryDB{}
	port   = os.Getenv("PORT")
	driver = os.Getenv("YT_DLP_PATH")
)

func init() {
	db.New()
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
	uid := os.Getuid()
	if uid == 0 {
		log.Println(cli.Yellow, "You're running this program as root (UID 0)", cli.Reset)
		log.Println(cli.Yellow, "This isn't reccomended unless you're using Docker", cli.Reset)
	}

	service := new(Service)
	err := rpc.Register(service)
	if err != nil {
		log.Fatal("Something has gone terribly wrong :)")
	}

	http.HandleFunc("/rpc", serveHTTP)
	http.Handle("/rpc-ws", websocket.Handler(serveWS))

	log.Printf("Started RPC server on port %s\n/rpc\t-> HTTP POST Handler\n/rpc-ws\t-> WebSocket Handler\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
