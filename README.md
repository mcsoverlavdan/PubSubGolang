# PubSubGolang
# Distributed Systems
A simple golang program on Publish Subscribe


# Requirements:
```golang
go get "github.com/gorilla/websocket"
go get "github.com/satori/go.uuid"
```


#### -run broker.go (open http://localhost:8000/  to view at the clients and subscribers need to refersh to 				get the latest count)
#### -run subscriberNews.go (open http://localhost:8080/   to view the output {topic news})
#### -run publisherNews.go (the news fetching part is taken from our previous assignment )

#### -run publisherIot.go
#### -run subscriberIot.go (open http://localhost:8081/  this will create a realtime graph of incoming 				iot or integer input {topic iot} )

The broker.go takes in clients request,appends them to the clients slice and based on the “action “and “topic” parameters it is then appended to subscriptions slice or is used to publish data based on the topic. 
Clients send publish and subscribe request using

 '{"action":"  ","topic":"  ",”message”:”message data”}'

Communication is established using websockets importing 
"github.com/gorilla/websocket"

Unique id to the clients are given using uuid importing
"github.com/satori/go.uuid"

Producer uses Dial function to connect to broker 
websocket.DefaultDialer.Dial(u.String(), nil)

Subscriber uses javascript 
conn = new WebSocket("ws://localhost:8000/broker"); 
And chart js to plot data 

![Image of ARCH](https://github.com/mcsoverlavdan/PubSubGolang/blob/main/pubsubassgn.png)


# Simple publisher program:

```golang
//url address to connect to
var addr = flag.String("addr", "localhost:8000", "http service address"


func main(){
	//url with path
   u := url.URL{Scheme: "ws", Host: *addr, Path: "/broker"}
   log.Printf("connecting to %s", u.String())

//Dial with url to connect to the broker

   c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)


   if err != nil {
      log.Fatal("dial:", err)
   }

//to read messages from broker 

_, message, err := c.ReadMessage()

if err != nil {
         log.Println("read:", err)
         return
      }
fmt.Println(message)

newsnotification:="hello this is a new news "
msg:="{\"action\":\"publish\",\"topic\":\"abc\",\"message\":\""+newsnotification+"\"}"

//msg is sent using write message and is passed as bytes

if err = c.WriteMessage(1, []byte(msg)); err != nil {
   fmt.Println(err)
}
```







# Output
Broker uses local host 8000 to display clients connected and the subscriptions and the topic


![Image of output1](https://github.com/mcsoverlavdan/PubSubGolang/blob/main/Picture1.png)


News subscriber displays the currently received  news from the topic “news”.

![Image of output2](https://github.com/mcsoverlavdan/PubSubGolang/blob/main/Picture2.png)





The integer values recieved are used to plot realtime data using chart.js cdn
Here in this case two clients (browsers) are subscribed to the topic iot,where they receive data of daily corona positive cases

![Image of graph](https://github.com/mcsoverlavdan/PubSubGolang/blob/main/Picture3.png)

# References :
https://www.youtube.com/watch?v=yyREnTgRTQ0
