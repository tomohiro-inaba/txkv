package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/tinaba/txkv"
)

const (
	DB_PATH  = "./var/txkv.db"
	GC_CYCLE = 30 * time.Second
)

var db *txkv.DB = nil

func beginHandler(w http.ResponseWriter, r *http.Request) {
	writableParam := r.URL.Query().Get("writable")
	writable := writableParam == "true"
	tx := db.Begin(writable)
	fmt.Fprintf(w, "txid=%s\n", tx.ID())
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")
	if len(paths) != 3 {
		fmt.Fprintf(w, "key is required\n")
		return
	}
	key := paths[2]

	txid := r.URL.Query().Get("txid")
	if txid == "" {
		fmt.Fprintf(w, "txid is required\n")
		return
	}
	tx, err := db.GetTx(txid)
	if err != nil {
		log.Fatal(err)
	}

	value, err := tx.Read(txkv.Key(key))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "value=%s\n", value)
}

func writeHandler(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")
	if len(paths) != 3 {
		fmt.Fprintf(w, "key is required\n")
		return
	}
	key := paths[2]

	value := r.URL.Query().Get("value")
	if value == "" {
		fmt.Fprintf(w, "value is required\n")
		return
	}

	txid := r.URL.Query().Get("txid")
	if txid == "" {
		fmt.Fprintf(w, "txid is required\n")
		return
	}
	tx, err := db.GetTx(txid)
	if err != nil {
		log.Fatal(err)
	}

	if err := tx.Write(txkv.Key(key), txkv.Value(value)); err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "ok\n")
}

func commitHandler(w http.ResponseWriter, r *http.Request) {
	txid := r.URL.Query().Get("txid")
	if txid == "" {
		fmt.Fprintf(w, "txid is required\n")
		return
	}
	tx, err := db.GetTx(txid)
	if err != nil {
		log.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "ok\n")
}

func rollbackHandler(w http.ResponseWriter, r *http.Request) {
	txid := r.URL.Query().Get("txid")
	if txid == "" {
		fmt.Fprintf(w, "txid is required\n")
		return
	}
	tx, err := db.GetTx(txid)
	if err != nil {
		log.Fatal(err)
	}

	if err := tx.Rollback(); err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "ok\n")
}

func gcHandler(w http.ResponseWriter, r *http.Request) {
	gc()
	fmt.Fprintf(w, "gc completed\n")
}

func gc() {
	tx := db.Begin(true)
	tx.GC()
	tx.Commit()
}

func main() {
	db, _ = txkv.Open(DB_PATH)
	defer db.Close()

	http.HandleFunc("/begin", beginHandler)
	http.HandleFunc("/read/", readHandler)
	http.HandleFunc("/write/", writeHandler)
	http.HandleFunc("/commit", commitHandler)
	http.HandleFunc("/rollback", rollbackHandler)
	http.HandleFunc("/gc", gcHandler)

	go func() {
		for {
			time.Sleep(GC_CYCLE)
			log.Printf("gc start\n")
			gc()
			log.Printf("gc completed\n")
		}
	}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
