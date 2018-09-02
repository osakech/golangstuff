// Package main Uses the standard library and creates a HTTP server
// that on each request responds with a counter of the total number
// of requests that it has received during the previous 60 seconds
// (moving window). The server should continue to return the
// correct numbers after restarting it, by persisting data to a file.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// fileToPersist is the name of the file we persist our data
var fileToPersist = "ts_persistence.json"

// timeframeInSec is the time in seconds of our timeframe
var timeframeInSec = 60

// fileMutex is used ensure mutually exclusive access to our file
var fileMutex sync.Mutex

// countHandler tells the user how many times the service has been called
// in our specified time frame
func countHandler(w http.ResponseWriter, r *http.Request) {
	// In case more than one user decides to use our service we make sure only one at a time
	// is able to read and write the file and the timestamps are sorted
	fileMutex.Lock()
	defer fileMutex.Unlock()

	// Get previously persisted timestamps
	timestamps := getTSFromFile(fileToPersist)

	// Get rid of timestamps older than our needed timeframe
	timestamps = purgeTSOlderThan(timestamps, int(time.Now().Unix())-timeframeInSec)

	// Feedback to the user
	fmt.Fprintf(w, "Calls from the last %d sec -> %d\n", timeframeInSec, len(timestamps))

	// add current timestamp to slice
	timestamps = append(timestamps, int(time.Now().Unix()))

	// persist our data
	writeTSToFile(fileToPersist, timestamps)

	return
}

// getTSFromFile gets the users persisted access times from a file
// It returns a slice of ints with unix timestamps
func getTSFromFile(filename string) []int {
	content, err := ioutil.ReadFile(filename)
	handleError(err)

	var timestamps []int
	errJSON := json.Unmarshal(content, &timestamps)
	handleError(errJSON)

	return timestamps
}

// purgeTSOlderThan removes all timestamps older than the specified timeframe
// It retuns a slice of ints with unix timestamps if timestamps were found inside our
// desired timeframe and nil if not
func purgeTSOlderThan(timestamps []int, timeFrame int) []int {
	for i, v := range timestamps {
		if v >= timeFrame {
			return timestamps[i:] // all TS from now on are inside our timeframe
		}
	}
	return nil
}

// writeTSToFile persists our sclice of ints as JSON in a file
func writeTSToFile(filepath string, timestamps []int) {
	jsonBytes, errjson := json.Marshal(timestamps)
	handleError(errjson)
	err := ioutil.WriteFile(filepath, jsonBytes, 0777)
	handleError(err)
}

// initTSFile initializes a new file with an empty json array inside in case there
// is no file
func initTSFile(filepath string) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		timestamps := []int{}
		writeTSToFile(filepath, timestamps)
	}
}

// handleError does generic error handlig
func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initTSFile(fileToPersist)
	http.HandleFunc("/counter/", countHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
