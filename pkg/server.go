package pkg

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"

	jsoniter "github.com/json-iterator/go"
	"golang.org/x/net/websocket"
	"goytdlp.rpc/m/pkg/cli"
)

// Package available variables
var (
	port   string
	driver string

	db   = MemoryDB{}
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

func init() {
	db.New()
	flag.StringVar(&port, "port", "4444", "port where RPC server will listen")
	flag.StringVar(&driver, "driver", "yt-dlp", "yt-dlp executable path")
	flag.Parse()
}

// Enable WebSockets as transport protocol
func serveWS(ws *websocket.Conn) {
	log.Println(ws.Request().RemoteAddr, "connected")
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
		log.Println(cli.Format(
			"You're running this program as root (UID 0)", cli.Yellow,
		))
		log.Println(cli.Format(
			"This isn't reccomended unless you're using Docker", cli.Yellow,
		))
	}

	service := new(Service)
	err := rpc.Register(service)
	if err != nil {
		log.Fatal("Something has gone terribly wrong :)")
	}

	http.HandleFunc("/rpc", serveHTTP)
	http.Handle("/rpc-ws", websocket.Handler(serveWS))

	fmt.Println(`        __         ____         ___  ___  _____`)
	fmt.Println(`  __ __/ /________/ / /__  ____/ _ \/ _ \/ ___/`)
	fmt.Println(` / // / __/___/ _  / / _ \/___/ , _/ ___/ /__`)
	fmt.Println(` \_, /\__/    \_._/_/ .__/   /_/|_/_/   \___/`)
	fmt.Println(`/___/              /_/`)

	fmt.Println(
		"\n"+cli.Format("/rpc", cli.BgBlue)+"\t HTTP POST Handler",
		"\n"+cli.Format("/rpc-ws", cli.BgBlue)+"\t WebSocket Handler\n",
	)
	log.Println("Driver in use:", driver)
	log.Println("Started RPC server on port", port)

	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), nil)
	if err != nil {
		log.Fatalln(cli.Format(
			fmt.Sprintf("Failed to bind port %s: already in use", port), cli.BgRed,
		))
	}
}
