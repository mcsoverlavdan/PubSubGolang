package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"time"
)

//the addr variable has the address to the broker in this case localhost port 8000
var addr = flag.String("addr", "localhost:8000", "http service address")

func main(){
	//the url uses the addr variable we created in line 22 and point to the path /broker

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/broker"}
	log.Printf("connecting to %s", u.String())

	//the dial function connects the publisher to the broker and if returned error dial err is printed

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	//c holds the *websocket.Conn websocket connection variable with which we can read incoming data

	_, message, err := c.ReadMessage()
	if err != nil {
				log.Println("read:", err)
				return
			}
	fmt.Println(message)

	//simulating iot type data , random values withint range 10 to 500 is passed to the broker
	//each news value is sent to the broker to publish under the topic "iot" ;   msg vairable is written in Json format as
	//the broker takes in messages and  JSON is unmarshled to a Message struct
	//action parameter is given as publish

	for i:=0;i<10000;i++{
		m:=rand.Intn(500 - 10) + 10
		//newsnotification:="{\"value\":\""+ strconv.Itoa(i)+"\""
		newsnotification:=strconv.Itoa(m)
		msg:="{\"action\":\"publish\",\"topic\":\"iot\",\"message\":\""+newsnotification+"\"}"
		fmt.Println(msg)
		if err = c.WriteMessage(1, []byte(msg)); err != nil {
			fmt.Println(err)
		}
		time.Sleep(5000*time.Millisecond)

	}



}
