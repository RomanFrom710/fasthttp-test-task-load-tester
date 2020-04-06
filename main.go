package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/paulbellamy/ratecounter"
)

type request struct {
	Text      string `json:"text"`
	ContentId int    `json:"content_id"`
	ClientId  int    `json:"client_id"`
	Timestamp int64  `json:"timestamp"`
}

const requestAmount = 1000000
const timeStep = 1 // In ms
const url = "http://localhost:8080"

func main() {
	tm := int64(time.Now().Unix()) * 1000
	routinesAmount := runtime.NumCPU()
	counter := ratecounter.NewRateCounter(1 * time.Second)
	ch := make(chan *request, routinesAmount*2)

	randArray := make([]byte, 64)

	for i := 0; i < routinesAmount; i++ {
		go func() {
			for {
				req := <-ch
				reqBody, _ := json.Marshal(req)
				http.Post(url, "application/json", bytes.NewReader(reqBody))
				counter.Incr(1)
			}
		}()
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	for i := 0; i < requestAmount; i++ {
		rand.Read(randArray)
		req := request{
			Text:      hex.EncodeToString(randArray),
			ContentId: i,
			ClientId:  r1.Intn(10) + 1,
			Timestamp: tm,
		}
		ch <- &req
		tm += timeStep
	}

	fmt.Println(counter.Rate())
}
