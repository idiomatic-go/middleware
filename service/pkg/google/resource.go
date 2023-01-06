package google

import (
	"io"
	"log"
	"net/http"
)

const (
	uri = "https://www.google.com/search?q=test"
)

func Search() []byte {
	newReq, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Printf("error: %v\n", err)
		return nil
	}
	resp, err2 := http.DefaultClient.Do(newReq)
	if err2 != nil {
		log.Printf("error: %v\n", err2)
		return nil
	}
	return ReadBody(resp)
}

func ReadBody(resp *http.Response) []byte {
	defer resp.Body.Close()
	var bytes []byte
	bytes, _ = io.ReadAll(resp.Body)
	return bytes
}
