package main

import (
    "fmt"
    "log"
    "net/http"
    "strings"
    "sync"
    "html"
    "github.com/467754239/simpleRESTfulDB/src/store"
)

type db struct {
    db_map map[string]interface{}
    l      sync.Mutex
}

var DB_g = db{db_map:make(map[string]interface{})}

type Hello struct{}

func (h Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    // uri拆解
    url_list := strings.Split(html.EscapeString(r.URL.Path), "/")

    // GET
    if url_list[1] == "get" && len(url_list) == 3 {
	DB_g.l.Lock()
	defer DB_g.l.Unlock()
	key := url_list[2]
	value, state := DB_g.db_map[key]
	if state {
	    w.Write([]byte(fmt.Sprint("\n")))
	    value, _ := value.(string)
	    w.Write([]byte(value))
	    log.Printf("ACTION: GET, %s, sucessful.", key)
	} else {
	    w.WriteHeader(http.StatusNotFound)
	    w.Write([]byte(fmt.Sprint("\n")))
	    w.Write([]byte(fmt.Sprintf("key:%s not found.", key)))
	    log.Printf("ACTION: GET, %s, failed.", key)
	}

	// SET
    } else if url_list[1] == "set" && len(url_list) == 4 {
	DB_g.l.Lock()
	defer DB_g.l.Unlock()
	key, value := url_list[2], url_list[3]
	log.Printf("ACTION: SET, %s=%s", key, value)
	DB_g.db_map[key] = value
	log.Println(DB_g.db_map)
        // 持久存储到磁盘
	store.Dump(DB_g.db_map)
        log.Println("dump file finish.")
    } else {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(fmt.Sprint("error request length.")))
    }
    return
}

func main() {
    var h Hello

    // load json data
    if store.Load() != nil {
        DB_g.db_map = store.Load()
    }

    // 每接收一个请求自动开一个协程 涉及到协程就会有并发 那么就要加锁.
    err := http.ListenAndServe(":8013", h)
    if err != nil {
	log.Fatal(err)
    }
}