package egoxml

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
)

func Serve() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/log", logHandler)
	router.HandleFunc("/api/parse", parseHandler)
	router.HandleFunc("/api/excludes", excludeHandler)
	router.HandleFunc("/api/combines", combineHandler)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	fmt.Println("Starting Server on http://localhost:8100")
	errs := make(chan error, 2)

	// Catch Interrupt for graceful shutdown
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		errs <- http.ListenAndServe("localhost:8100", router)
	}()

	<-errs
	defer db.Close()
}
func structExpose(data interface{}, w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func excludeHandler(w http.ResponseWriter, r *http.Request) {
	structExpose(GetConfig().Exclude, w, r)
}

func combineHandler(w http.ResponseWriter, r *http.Request) {
	structExpose(GetConfig().Combine, w, r)
}

func parseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	create()
	cfg := GetConfig()
	entries := Parse(cfg.Savegame)

	lastEntry := fetchLastEntry()
	err := bulkInsert(lastEntry, entries)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	entries := fetchLogs()
	structExpose(entries, w, r)
}
