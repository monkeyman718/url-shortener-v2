package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "math/rand"

    _ "github.com/lib/pq"
    "github.com/gorilla/mux"
)

type urlMap struct {
    shortCode   string
    longURL     string
}

var sqlconn *sql.DB
var err error

const shortURLLength = 6
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"


func main(){
    initDB()
    defer sqlconn.Close()

    router := mux.NewRouter()
    router.HandleFunc("/shorten", shortenerHandler )

    fmt.Println("Listening on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", router))
}

func initDB(){
    sqlStr := "host=localhost port=5432 user=postgres password=postgres sslmode=disable"    
    // Open psql connection
    sqlconn, err = sql.Open("postgres", sqlStr)

    if err != nil {
        panic(err)
        return
    }

    err = sqlconn.Ping()
    
    // Check for errors connecting
    if err != nil {
        fmt.Printf("Cannot connect to database: %v", err)
        return 
    }
    
    fmt.Println("Connection to database successful!")
}

func insertShortCode(urlItem *urlMap){
    sqlStr := `INSERT INTO "urls" (short_code, long_url) VALUES($1, $2)` 
    _, err := sqlconn.Exec(sqlStr, urlItem.shortCode,urlItem.longURL)

    if err != nil {
        panic(err)
    }
}

func generateShortURL() string {
    shortCode := make([]byte, shortURLLength)
    for i := 0; i < shortURLLength; i++ {
        shortCode[i] = charset[rand.Intn(shortURLLength)]
    }

    return string(shortCode)
}

func shortenerHandler(w http.ResponseWriter, r *http.Request){
    shortCode := generateShortURL() 
    longURL := r.URL.Query().Get("url")

    var urlItem urlMap 
    urlItem.shortCode = shortCode
    urlItem.longURL = longURL

    insertShortCode(&urlItem)

    fmt.Fprintf(w, "http://localhost:8080/%v\n", shortCode)
}
//func redirectHandler(w http.ResponseWriter, r *http.Response){}
