package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "os"
    "context"
    "flag"
    "os/signal"
    "time"
)

var port = "8090"

type MiddlewareFunc func(http.Handler) http.Handler

func postsHandler(w http.ResponseWriter, r *http.Request) {
    response := fmt.Sprintf("Posts page")
    fmt.Fprint(w, response)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    response := fmt.Sprintf("<h1>Post %s</h1>",id)
    fmt.Fprint(w, response)
}

func main() {

    var wait time.Duration
    flag.DurationVar(&wait, "graceful-timeout", time.Second * 15, "duration before exit connection to finish")
    flag.Parse()

    router := mux.NewRouter()
    posts := router.PathPrefix("/posts").Subrouter()
    posts.HandleFunc("/",postsHandler)
    posts.HandleFunc("/{id:[0-9]+}", postHandler)
    router.Use(loggingMiddleware)
    http.Handle("/", router)

//    log.Println("Starting on ", port)
//    http.ListenAndServe(":"+port, nil)

    srv := &http.Server{
        Addr:         ":"+port,
        WriteTimeout: time.Second * 15,
        ReadTimeout:  time.Second * 15,
        IdleTimeout:  time.Second * 60,
        Handler:      router,
    }

    go func() {
        if err := srv.ListenAndServe(); err != nil {
            log.Fatalln(err)
        }
    }()

    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)

    <-c

    ctx, cancel := context.WithTimeout(context.Background(), wait)
    defer cancel()

    srv.Shutdown(ctx)

    log.Println("Shutdown...")
    os.Exit(0)
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println(r.RequestURI)
        next.ServeHTTP(w, r)
    })
}
