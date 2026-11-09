package main

import (
 "database/sql"
 "encoding/json"
 "fmt"
 "log"
 "net/http"
 "os"

 _ "github.com/lib/pq"
)

type Item struct {
 ID int `json:"id"`
 Name string `json:"name"`
}

func main() {
 dsn := "host=" + env("DB_HOST", "db") +
  " port=" + env("DB_PORT", "5432") +
  " user=" + env("DB_USER", "app_user") +
  " password=" + env("DB_PASSWORD", "secure_password") +
  " dbname=" + env("DB_NAME", "myapp") +
  " sslmode=disable"

 db, err := sql.Open("postgres", dsn)
 if err != nil {
  log.Fatal(err)
 }
 defer db.Close()

 staticDir := env("STATIC_DIR", "/app/static")
 fs := http.FileServer(http.Dir(staticDir))

 http.HandleFunc("/api/items", func(w http.ResponseWriter, r *http.Request) {
  rows, err := db.Query(`SELECT id, name FROM items ORDER BY id`)
  if err != nil {
   http.Error(w, err.Error(), 500)
   return
  }
  defer rows.Close()

  items := []Item{}
  for rows.Next() {
   var it Item
   if err := rows.Scan(&it.ID, &it.Name); err != nil {
    http.Error(w, err.Error(), 500)
    return
   }
   items = append(items, it)
  }
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(items)
 })

 http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  fmt.Fprint(w, "OK")
 })

 http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
  if err := db.Ping(); err != nil {
   http.Error(w, "DB Error", 500)
   return
  }
  w.WriteHeader(http.StatusOK)
  fmt.Fprint(w, "Data Connection OK")
 })

 http.HandleFunc("/api/items/count", func(w http.ResponseWriter, r *http.Request) {
  var count int
  err := db.QueryRow(`SELECT COUNT(*) FROM items`).Scan(&count)
  if err != nil {
   http.Error(w, err.Error(), 500)
   return
  }
  w.Header().Set("Content-Type", "text/plain")
  fmt.Fprint(w, count)
 })

 http.Handle("/", fs)

 addr := ":" + env("APP_PORT", "3000")
 log.Printf("listening on %s", addr)
 if err := http.ListenAndServe(addr, nil); err != nil {
  log.Fatal(err)
 }
}

func env(key, def string) string {
 if v := os.Getenv(key); v != "" {
  return v
 }
 return def
}