package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type MessengerEvent struct {
	Object string       `json:"object"`
	Entry  []EventEntry `json:"entry"`
}

type EventEntry struct {
	PageID   string     `json:"id"`
	Time     int        `json:"time"`
	Envelope []Envelope `json:"messaging"`
}

type Envelope struct {
	Sender    Sender    `json:"sender"`
	Recipient Recipient `json:"recipient"`
	Time      int       `json:"timestamp`
	Message   Message   `json:"message"`
}

type Sender struct {
	PSID string `json:"id"`
}

type Recipient struct {
	PageId string `json:"id"`
}

type Message struct {
	ID   string `json:"mid"`
	Text string `json:"text"`
}

func main() {
	http.HandleFunc("/webhook", handler)
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		handleGet(w, r)
	}

	if r.Method == "POST" {
		handlePost(r)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	verifyToken := "mooz_is_cool"
	query := r.URL.Query()
	token := query.Get("hub.verify_token")
	mode := query.Get("hub.mode")
	challenge := query.Get("hub.challenge")

	if mode == "subscribe" && token == verifyToken {
		w.Write([]byte(challenge))
	} else {
		http.Error(w, "wicked beasty, thou shalt not pass", http.StatusForbidden)
	}
}

func handlePost(r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	event := MessengerEvent{}

	err = json.Unmarshal(body, &event)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(event.Entry[0].Envelope[0].Message.Text)
}
