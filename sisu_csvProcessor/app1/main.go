package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

var maxProcesses = 100

type csvDataset struct {
	Email string `json:"Email"`
	Name  string `json:"Name"`
}

type counter struct {
	counter int
	mux     sync.Mutex
}

var pc counter
var rc counter

func main() {
	body, err := ioutil.ReadFile("data.csv")

	rowChan := make(chan []string, maxProcesses)

	if err != nil {
		log.Fatal("can not read csv file " + err.Error())
	}
	var wg sync.WaitGroup

	r := csv.NewReader(strings.NewReader(string(body)))
	go sendData(&wg, rowChan)
	for {
		r.Comma = ','
		row, err := r.Read()

		if err == io.EOF {
			log.Println("Finished processing file!")
			break
		}

		if err != nil {
			log.Fatal("error reading in csv:" + err.Error())
		}
		wg.Add(1)
		rowChan <- row
	}
	wg.Wait()
	log.Printf("Needed %d retries", rc.counter)
	log.Printf("Processed %d lines", pc.counter)

}

func sendData(wg *sync.WaitGroup, rowChan <-chan []string) {
	for {
		pc.mux.Lock()
		pc.counter++
		pc.mux.Unlock()

		row := <-rowChan
		ds := csvDataset{Email: row[1], Name: row[0]}
		json := encodeJSON(ds)
		makeRequest(json)
		wg.Done()
	}
}

func encodeJSON(ds csvDataset) []byte {
	b, err := json.Marshal(ds)
	if err != nil {
		log.Print("error reading csv:", err)
	}
	return b
}

func makeRequest(b []byte) {
	for {
		resp, errweb := http.Post("http://localhost:8080/givemedata/", "application/json", io.Reader(bytes.NewBuffer(b)))
		if errweb != nil {
			log.Fatal("error in the web -> " + errweb.Error())
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			break
		}

		if resp.StatusCode == 503 {
			// i, err := strconv.Atoi(resp.Header.Get("Retry-After"))
			// if err != nil {
			// 	log.Fatal("could not convert string to int -> " + err.Error())
			// }
			// fmt.Println("retry!")
			rc.mux.Lock()
			rc.counter++
			rc.mux.Unlock()
			// time.Sleep(time.Duration(i) * time.Second)
		}
	}
}
