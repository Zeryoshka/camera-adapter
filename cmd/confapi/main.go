package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	ServerPort uint
}

func main() {
	uprConfPath := flag.String("upr-conf", "", "path to local yaml config")
	serverPort := flag.Int("port", 80, "port for conf-server")
	flag.Parse()

	staticServer := http.FileServer(http.Dir("./confapi-static/"))
	http.Handle("/", staticServer)
	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			f, err := os.Open(*uprConfPath)
			if err != nil {
				log.Println("can't read config with path:", uprConfPath, "; error:", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			configData, _ := io.ReadAll(f)
			w.Write(configData)
		} else if r.Method == "POST" {
			w.WriteHeader(http.StatusCreated)
			w.Header().Add("Content-Type", "text/plain")
			f, err := os.Create(*uprConfPath)
			if err != nil {
				log.Println("can't read config with path:", uprConfPath, "; error:", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			configData, _ := io.ReadAll(r.Body)
			_, err = f.Write(configData)
			if err != nil {
				log.Println("can't write to file cause:", err)
			}
			err = f.Close()
			if err != nil {
				log.Println("can't close file cause:", err)
			}

		} else {
			log.Println("Method:", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	err := http.ListenAndServe("0.0.0.0:"+strconv.Itoa(*serverPort), nil)
	log.Fatal(err)
}
