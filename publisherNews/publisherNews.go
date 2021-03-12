package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

//the news fecthing is taken from our previous Assingment's initial news fetch function
//the function main connects to the broker and sends the data we get to the broker in json format

//the addr variable has the address to the broker in this case localhost port 8000
var addr = flag.String("addr", "localhost:8000", "http service address")

//for unmarshaling the xml code we need to have a do structure to which we can unmarshal its easy to use the site taht
//automatically converts it xml to struct https://www.onlinetool.io/xmltogo/
//here we are taking india today rss
//https://www.indiatoday.in/rss/1206584

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Media   string   `xml:"media,attr"`
	Atom    string   `xml:"atom,attr"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Text          string `xml:",chardata"`
		Title         string `xml:"title"`
		Description   string `xml:"description"`
		Link          string `xml:"link"`
		LastBuildDate string `xml:"lastBuildDate"`
		Generator     string `xml:"generator"`
		Image         struct {
			Text        string `xml:",chardata"`
			URL         string `xml:"url"`
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			Description string `xml:"description"`
		} `xml:"image"`
		Item []struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			//Description testtype `xml:"description"`
			Description string `xml:"description"`
			PubDate     string `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
}


//value we need are the only values we need from so we create a structure
//here the problem we faced was the link when passed as string ,the html template when used with href considers it as
//unsecure link , it got fixed when used with template url

type valueWeNeed struct{
	Title       string
	//Link        string
	Link template.URL
	//Description testtype
	Description string
	PubDate     string
}
//the final data we get from the main page is stored as a global variable so other functions can access it
var finalFinalData []valueWeNeed

//for using go routines using wait groups we decare a variable wg
var wg sync.WaitGroup



func extractData(queue chan<- valueWeNeed,title,link,description,pubdate string) {
	//defer func(){
	//	fmt.Println("Done!")
	//}()
	defer wg.Done()
	e := strings.Index(description, "</a>")
	description=description[e+4:]
	link = strings.Replace(link, " ", "", -1)//Remove all the trailing spaces in the link
	p := valueWeNeed{Title: title,
		Link        :template.URL(link),
		Description :description,
		PubDate     :pubdate,
	}
	//fmt.Println(p)
	//runtime.Gosched()
	queue<-p

}
func dataToWebsite(s Rss) []valueWeNeed{
	//finalData := make(map[int]valueWeNeed)
	var finalData []valueWeNeed
	queue := make(chan valueWeNeed,50)
	//this for loop created 50 new go rounites here we are waiting for the routines to finish
	for i:=0;i<50;i++ {

		value:=s.Channel.Item[i] //each items we got from the xml

		//fmt.Println("entering....")
		wg.Add(1)
		go extractData(queue,value.Title,value.Link,value.Description,value.PubDate)
		//fmt.Println("i value ",i,runtime.NumGoroutine()) //to check the go routines

	}
	wg.Wait()
	close(queue)



	//looping over the channel buffer
	for elem := range queue {
		finalData=append(finalData,elem)

	}

	//returning the final data
	return finalData


}


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

	//-------------------------------------------------------first assignment part for feting news---------

	resp, err := http.Get("https://www.indiatoday.in/rss/1206584")
	if err != nil {
		log.Println("error in http get:", err)
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.Body)

	var s Rss
	err = xml.Unmarshal(bytes, &s)
	if err!=nil{
		fmt.Println("xml error ",err)
	}
	resp.Body.Close()

	finalFinalData=dataToWebsite(s)
	//-------------------------------------------------------end of first assignment part for feting news--------------

	//each news title is sent to the broker to publish under the topic "news" ;   msg vairable is written in Json format as
	//the broker takes in messages and  JSON is unmarshled to a Message struct
	//action parameter is given as publish 
	for _,news:=range finalFinalData{
		newsnotification:=news.Title
		fmt.Println(newsnotification)
		msg:="{\"action\":\"publish\",\"topic\":\"news\",\"message\":\""+newsnotification+"\"}"
		if err = c.WriteMessage(1, []byte(msg)); err != nil {
			fmt.Println(err)
		}
		time.Sleep(1000*time.Millisecond)
	}



}
