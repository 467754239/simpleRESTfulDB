package store

import (
    simplejson "github.com/bitly/go-simplejson"
    "io/ioutil"
    "log"
    "encoding/json"
    "os"
)

func Dump(data map[string]interface{}) {
    j, err := json.Marshal(data)
    if err != nil {
	log.Fatal("json.Marshal error", err)
	return
    }
    ioutil.WriteFile("C:/dump.db", j, 0x644)

}

func Load() map[string]interface{} {
    _, err := os.Stat("C:/dump.db")
    if err == nil {
	j, err := ioutil.ReadFile("C:/dump.db")
	if err != nil {
	    log.Fatal("read file error from disk", err)
	    return nil
	}
	json_j, err := simplejson.NewJson(j)
	if err != nil {
	    log.Fatal("json file error", err)
	    return nil
	}
	m, _ := json_j.Map()
	return m
    } else {
        log.Println("file not exists.")
	return nil
    }

}
