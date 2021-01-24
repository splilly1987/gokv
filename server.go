package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"net/http"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
	kv "github.com/splilly1987/gokv/db"
)

var (
	dbfile = "test.db"
)

//InfoServer - For our HTTP Server
type InfoServer struct {
	router     *mux.Router
	httpServer *http.Server
}

func main() {
	var srv InfoServer
	srv.router = mux.NewRouter()

	srv.router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s: Health Requested\n", time.Now().Format(time.RFC3339))
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})
	srv.router.HandleFunc("/kv/{key}", GetKeyValueHandler).Methods("GET")
	srv.router.HandleFunc("/kv/{key}", PutKeyValueHandler).Methods("POST", "PUT", "UPDATE")
	srv.router.HandleFunc("/kv/{key}", DeleteKeyValueHandler).Methods("DELETE")
	srv.router.HandleFunc("/kv", ListKeysHandler).Methods("GET")

	srv.httpServer = &http.Server{
		Handler:      srv.router,
		Addr:         "127.0.0.1:9000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.httpServer.ListenAndServe())
}

//GetKeyValueHandler - Get Person
func GetKeyValueHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["key"]

	store, _ := kv.Open(dbfile)
	defer store.Close()
	var val map[string]string
	store.Get(name, &val)
	json.NewEncoder(w).Encode(val)
}

//PutKeyValueHandler - Add and Update key with value
func PutKeyValueHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	decoder := json.NewDecoder(r.Body)
	store, _ := kv.Open(dbfile)
	defer store.Close()
	var val map[string]string
	if err := decoder.Decode(&val); err != nil {
		log.Println(err)
	}

	store.Put(vars["key"], val)
}

//ListKeysHandler - Get all Keys
func ListKeysHandler(w http.ResponseWriter, r *http.Request) {
	store, _ := kv.Open(dbfile)
	defer store.Close()
	list, _ := store.ListKeys()

	njson := simplejson.New()
	njson.Set("keys", list)

	payload, err := njson.MarshalJSON()
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

//DeleteKeyValueHandler - Delete key
func DeleteKeyValueHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	store, _ := kv.Open(dbfile)
	defer store.Close()
	log.Printf("Deleting: %s\n", vars["key"])
	store.Delete(vars["key"])
}
