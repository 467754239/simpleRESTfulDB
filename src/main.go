package main

import (
    "fmt"
    "log"
    "net/http"
    "strings"
    "sync"
    "html"
)

type db struct {
    db_map map[string]string
    l      sync.Mutex
}

var DB_g = db{db_map:make(map[string]string)}

type Hello struct{}

func (h Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    // uri拆解
    url_list := strings.Split(html.EscapeString(r.URL.Path), "/")

    // GET
    if url_list[1] == "get" && len(url_list) >= 3 {
	DB_g.l.Lock()
	defer DB_g.l.Unlock()
	key := url_list[2]
	value, state := DB_g.db_map[key]

	if state {
	    w.Write([]byte(fmt.Sprint("\n")))
	    w.Write([]byte(value))
	} else {
	    w.WriteHeader(http.StatusNotFound)
	    w.Write([]byte(fmt.Sprint("\n")))
	    w.Write([]byte(fmt.Sprintf("key:%s not found.", key)))
	}

	// SET
    } else if url_list[1] == "set" && len(url_list) == 4 {
	DB_g.l.Lock()
	defer DB_g.l.Unlock()
	key, value := url_list[2], url_list[3]
	DB_g.db_map[key] = value
    } else {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(fmt.Sprint("error request length.")))
    }
    return
}

func main() {
    var h Hello
    // 每接收一个请求自动开一个协程 涉及到协程就会有并发 那么就要加锁.
    err := http.ListenAndServe(":8002", h)
    if err != nil {
	log.Fatal(err)
    }
}