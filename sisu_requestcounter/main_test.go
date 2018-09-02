package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestCountHandlerNumberOfCallsResets(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	//overwrite some vars from main.go
	timeframeInSec = 2
	fileToPersist = "tempTestFile"
	os.Remove(fileToPersist)
	initTSFile(fileToPersist) // we need a file

	calls := 0
	expectedBody := ""
	for calls <= 2 {
		expectedBody += fmt.Sprintf("Calls from the last %d sec -> %d\n", timeframeInSec, calls)
		calls++
		countHandler(w, req)
	}

	time.Sleep(3 * time.Second) // wait to prove that calls outside our time frame are not counted
	countHandler(w, req)

	expectedBody += fmt.Sprintf("Calls from the last %d sec -> %d\n", timeframeInSec, 0) // check if we get 0 back

	gotBodies, _ := ioutil.ReadAll(w.Result().Body)
	expectedBodies := []byte(expectedBody)

	if !bytes.Equal(gotBodies, expectedBodies) {
		t.Error("countHandler resuming after counting failed")
	}

	os.Remove(fileToPersist)

}
func TestCountHandlerNumberOfCalls(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	//overwrite some vars from main.go
	timeframeInSec = 10
	fileToPersist = "tempTestFile"
	os.Remove(fileToPersist)
	initTSFile(fileToPersist) // we need a file

	calls := 0
	expectedBody := ""
	for calls <= 10 {
		expectedBody += fmt.Sprintf("Calls from the last %d sec -> %d\n", timeframeInSec, calls)
		calls++
		countHandler(w, req)
	}

	gotBodies, _ := ioutil.ReadAll(w.Result().Body)
	expectedBodies := []byte(expectedBody)

	if !bytes.Equal(gotBodies, expectedBodies) {
		t.Error("countHandler failed")
	}

	os.Remove(fileToPersist)

}
func TestInitTSFile(t *testing.T) {
	initTSFile("test_temp_file")

	gotContent, _ := ioutil.ReadFile("test_temp_file")
	expectedContent := []byte("[]")

	if !bytes.Equal(gotContent, expectedContent) {
		t.Error("initTSFile failed")
	}
	os.Remove("test_temp_file")
}

func TestGetTSFromFile(t *testing.T) {
	tmpfile := getTempFileWithContents("[1,2]")
	gotInts := getTSFromFile(tmpfile.Name())
	expectedInts := []int{1, 2}

	if !reflect.DeepEqual(gotInts, expectedInts) {
		t.Error("getTSFromFile failed")
	}

	os.Remove(tmpfile.Name())
}

func TestWriteTSToFile(t *testing.T) {
	tmpfile := getTempFileWithContents("doesn't matter")
	expectedInts := []int{1, 2, 3}
	writeTSToFile(tmpfile.Name(), expectedInts)
	gotInts := getTSFromFile(tmpfile.Name())

	if !reflect.DeepEqual(gotInts, expectedInts) {
		t.Error("writeTSToFile failed")
	}

	os.Remove(tmpfile.Name())
}

func TestPurgeTSOlderThan(t *testing.T) {
	timestamps := []int{1, 2, 3}
	gotInts := purgeTSOlderThan(timestamps, 2)
	expectedInts := []int{2, 3}

	if !reflect.DeepEqual(gotInts, expectedInts) {
		t.Error("purgeTSOlderThan failed")
	}

	noTimestamps := []int{}
	gotNoInts := purgeTSOlderThan(noTimestamps, 2)
	if gotNoInts != nil {
		t.Error("purgeTSOlderThan with no timestamps failed")
	}

}

func getTempFileWithContents(provContent string) *os.File {
	content := []byte(provContent)
	tmpfile, err := ioutil.TempFile("", "")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
	return tmpfile
}
