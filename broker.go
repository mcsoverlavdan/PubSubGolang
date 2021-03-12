package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
)


//A server application calls the Upgrader.Upgrade method from an HTTP request handler to get a *Conn:

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
//the clients connected will be appended to the clients slice and the subscribers to the subscription slice
type PubSub struct {
	Clients       []Client
	Subscriptions []Subscription
}
//the client id which is auto generated id using autoid function and the connection variable stores the connection value
type Client struct {
	Id         string
	Connection *websocket.Conn
}
//the subscription type has topic he has subscribed to  and client value
type Subscription struct {
	Topic  string
	Client Client
}
//autoid generates a unique id to differentiate clients
func autoId() string {

	return uuid.Must(uuid.NewV4()).String()
}
//the sent method of client type is to sent a message to the websoccect connection
func (client *Client) Send(message [] byte) error {

	return client.Connection.WriteMessage(1, message)

}

//addclient method is used to append new c;ient details to the client slice
func (ps *PubSub) AddClient(client Client) *PubSub {

	ps.Clients = append(ps.Clients, client)

	fmt.Println("adding new client to the list", client.Id, len(ps.Clients))

	payload := []byte("Hello Client ID:" +
		client.Id)
	//sent the client id to the connection
	err := client.Connection.WriteMessage(1, payload)
	if err != nil {
		log.Println(err)
	}
	return ps

}
//when connection is cloed remove client from the subscriptions and the client list
func (ps *PubSub) RemoveClient(client Client) (*PubSub) {

	// first remove all subscriptions by this client

	for index, sub := range ps.Subscriptions {

		if client.Id == sub.Client.Id {
			//slice data is deleted using append data till that value and after that value
			ps.Subscriptions = append(ps.Subscriptions[:index], ps.Subscriptions[index+1:]...)
		}
	}

	// remove client from the list

	for index, c := range ps.Clients {

		if c.Id == client.Id {
			ps.Clients = append(ps.Clients[:index], ps.Clients[index+1:]...)
		}

	}

	return ps
}
//a variable ps of tpe pub sub is created
var ps PubSub

//----------------------------------------------------------------------connection handler function........

//----------------------for passing data to channels creating struct (start).....
type QueueData struct{
	clinetdetails Client
	payloadmessage []byte
}
type DataPasser struct {
	queue chan QueueData
}

//this funtion is used for websocket connections
func (p *DataPasser) ConnectNodes(w http.ResponseWriter, r *http.Request) {
	//upgrader checkorgin is returned true so that any websocket connections arent refsed
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true

	}

	fmt.Println("starting connec.......")
	//The Conn type represents a WebSocket connection.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	//the client details are added to the client Val with struct Client
	ClientVal := Client{
		Id:         autoId(),
		Connection: conn,
	}
	//client is then added to the client list using Addclinet function
	ps.AddClient(ClientVal)

	//reading incoming data from the client in and infinte loop until the connection from the client is closed
	for {
		_, payload, err := conn.ReadMessage()
		fmt.Println("payload msg",payload)
		//the data where payload is the messafe and client val is the client details are passed to the channel queue
		p.queue<-QueueData{ClientVal,payload}


		if err != nil {
			log.Println(err)
			ps.RemoveClient(ClientVal)
			fmt.Println("after removing client list", len(ps.Clients))
			return
		}
	}


}
//......................... struct  to unmarshal JSON tupe----------------
type Message struct {
	Action  string          `json:"action"`
	Topic   string          `json:"topic"`
	Message json.RawMessage `json:"message"`
}

//the log method is run as a go routine and check for data in queue channel and does the publish subscribe methods accordingly
//checking the message sent by the client
func (p *DataPasser) log() {
	//for range in the channel queue
	for item := range p.queue {
		//m of tpe Message used for unmarshalling the JSON
		m := Message{}
		//the message is unmarshaled to m
		err := json.Unmarshal(item.payloadmessage, &m)
		if err != nil {
			fmt.Println("This is not correct message payload",err,"payload : ",item.payloadmessage)
		}
		fmt.Println("Client Id", item.clinetdetails.Id)
		fmt.Println("Messsage : ",m.Topic,m.Message,m.Action)

		//for action mentioned in the action parameter publish and subscribe methods are done
		if m.Action =="publish"{
			var subscriptions  []Subscription
			//checking for all the sibscriptions under the topic mentioned
			for _, subscription := range ps.Subscriptions {
				if subscription.Topic == m.Topic {
					subscriptions = append(subscriptions, subscription)
				}
			}
			//for the subscriptions send the message
			for _, sub := range subscriptions {

				fmt.Printf("Sending to client id %s message is %s \n", sub.Client.Id, m.Message)
				//sub.Client.Connection.WriteMessage(1, message)

				err := sub.Client.Send(m.Message)
				if err !=nil{
					fmt.Println("not able to send ",err)
				}

			}


		}else if m.Action =="subscribe"{
			fmt.Println("inside subscribe")
			var subscriptions  []Subscription
			//check if subscription already exixsts for the same client for that topic
			for _, subscription := range ps.Subscriptions {
				if subscription.Client.Id == item.clinetdetails.Id && subscription.Topic == m.Topic {
					subscriptions = append(subscriptions, subscription)

				}
			}
			if len(subscriptions) > 0 {

				// client is subscribed this topic before

				continue
			}

			newSubscription := Subscription{
				Topic:  m.Topic,
				Client: item.clinetdetails,
			}
			//append the client to the subscritption list with the topic mentioned
			ps.Subscriptions = append(ps.Subscriptions, newSubscription)
		}


	}

}




//.........................................main page for broker to know status. of clients sonnected and subscriptions..............
func index(w http.ResponseWriter,r *http.Request){
	fmt.Println("the index function is starting ")

	msg:="<h1>Clients:</h1><br>"
	for _,value:= range ps.Clients{
		msg+="Client Id"+value.Id +"<br>"
	}
	fmt.Println(msg)
	msg+="<h1>Subscriptions:</h1><br>"
	for _,value:= range ps.Subscriptions{
		msg+="Client Id   :"+value.Client.Id +" Topic its subscribed to is "+value.Topic+"<br>"
	}
	w.Header().Set("Content-Type", "text/html")
	_, err := fmt.Fprint(w, msg)
	if err !=nil{
		fmt.Println("error in Fprint",err)
	}

}

func main(){
	//the channel queue is used to pass data to the log funtion when new messages are recieved from the client
	queue:= make(chan QueueData)

	passer := &DataPasser{queue: queue}
	//log funtion sorts the request to publish and subscribe
	go passer.log()


	//the inedx funtion is used to view the current clients whic are connected
	http.HandleFunc("/",index)
	//connect nodes established the websocket connection between clients 
	http.HandleFunc("/broker",passer.ConnectNodes)
	err:=http.ListenAndServe(":8000",nil)

	if err !=nil{
		fmt.Println("error in listen and serve",err)
	}
}
