package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type request struct {
	Text      string `json:"text"`
	ContentId int    `json:"content_id"`
	ClientId  int    `json:"client_id"`
	Timestamp int64  `json:"timestamp"`
}

const requestAmount = 100000
const timeStep = 10000 // In ms
const routinesAmount = 10
const url = "http://localhost:8080"

func main() {
	tm := int64(time.Now().Unix()) * 1000

	ch := make(chan *request, routinesAmount*2)

	randArray := make([]byte, 64)

	for i := 0; i < routinesAmount; i++ {
		go func() {
			for {
				req := <-ch
				reqBody, _ := json.Marshal(req)
				fmt.Println("Sending " + strconv.Itoa(req.ContentId))
				http.Post(url, "application/json", bytes.NewReader(reqBody))
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
}
