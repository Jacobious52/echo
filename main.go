package main

import (
	"bufio"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"
)

var addr = flag.String("addr", ":8080", "http service address")
var jsTemplate = template.Must(template.ParseFiles("./template/chat.js"))

func hotReload() {
	ticker := time.NewTicker(time.Second * 2)
	info, err := os.Stat("./template/chat.js")
	if err != nil {
		log.Println("warning chat.js:", err)
	}
	lastMod := info.ModTime()
	go func() {
		for t := range ticker.C {
			info, err := os.Stat("./template/chat.js")
			if err != nil {
				log.Println("warning chat.js:", err)
			}
			if info.ModTime() != lastMod {
				jsTemplate = template.Must(template.ParseFiles("./template/chat.js"))
				lastMod = info.ModTime()
				log.Println("reload chat.js at", t)
			}
		}
	}()
}

func serveJS(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	w.Header().Set("Content-Type", "text/javascript")
	jsTemplate.Execute(w, r.Host)
}

func logMux(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func readInput(l chan string) {
	defer close(l)
	for {
		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		l <- line
	}
}

func main() {
	flag.Parse()

	hub := newHub()
	go hub.run()

	hotReload()

	// read input async
	line := make(chan string)
	go readInput(line)

	// keep waiting for stdin on a separate thread so we don't block
	go func() {
		for l := range line {
			switch l {
			case "quit", "exit", "close", "kill", "die":
				log.Println("exiting..")
				os.Exit(0)
			case "clear", "clean":
				hub.savedMessages = nil
				hub.savedMessages = make([][]byte, 0)
				log.Println("saved messages cleared")
			default:
				log.Println("unknown command:", l)
			}
		}
	}()

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/template/chat.js", serveJS)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	log.Println("server started.. awaiting requests..")

	err := http.ListenAndServe(*addr, logMux(http.DefaultServeMux))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
