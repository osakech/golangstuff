package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"
)

const maxParallel = 50 // set the number of parallel processes
const slowMode = true  // slow mode, to see the buffer filling up

type recvUserData struct {
	Email string
	Name  string
}

type userData struct {
	Email string
	Name  string
	Count int
}

var inMemoryUserdata map[string]userData

func initMap() {
	inMemoryUserdata = make(map[string]userData)
}

var mapMutex sync.Mutex

var dataChan = make(chan []byte, maxParallel)
var wg sync.WaitGroup

func givemedata(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Panicln("Body Read Error -> " + err.Error())
	}

	select {
	case dataChan <- b:
	default:
		// w.Header().Set("Retry-After", "1")
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("503 - Server busy!"))
	}
}

func processData(dataChan <-chan []byte, wg *sync.WaitGroup) {

	for {
		data := <-dataChan
		var ds recvUserData
		errjson := json.Unmarshal(data, &ds)
		if errjson != nil {
			log.Panic("JSON Unmarshal Error -> " + errjson.Error())
		}
		saveData(&ds)
	}
}

func saveData(ds *recvUserData) {
	mapMutex.Lock() // mutex our map
	if slowMode {
		time.Sleep(100 * time.Microsecond)
	}
	localuserData, ok := inMemoryUserdata[ds.Email]
	if !ok {
		data := userData{Email: ds.Email, Name: ds.Name, Count: 1}
		inMemoryUserdata[ds.Email] = data
	} else {
		localuserData.Count++
		inMemoryUserdata[ds.Email] = localuserData
	}
	mapMutex.Unlock()
}

func monitoring(w http.ResponseWriter, r *http.Request) {
	type wrapedForSortingUserdata struct {
		Key   string
		Value userData
	}
	var sortedUserData []wrapedForSortingUserdata
	for k, v := range inMemoryUserdata {
		sortedUserData = append(sortedUserData, wrapedForSortingUserdata{k, v})
	}

	sort.Slice(sortedUserData, func(i, j int) bool {
		if sortedUserData[i].Value.Count > sortedUserData[j].Value.Count { // First sort by count descending 10,9,8...
			return true
		}
		if sortedUserData[i].Value.Count < sortedUserData[j].Value.Count {
			return false
		}
		return sortedUserData[i].Value.Email < sortedUserData[j].Value.Email // Then sort by Email ascending a,b,c ...
	})

	for _, v := range sortedUserData {
		fmt.Fprintf(w, "%s  %d \n", v.Value.Email, v.Value.Count)
	}
}

func main() {
	initMap()
	go processData(dataChan, &wg)
	http.HandleFunc("/givemedata/", givemedata)
	http.HandleFunc("/monitoring/", monitoring)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
