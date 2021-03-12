package main

import (
	"html/template"
	"net/http"
)

func main() {
	//the index page runs the index html file in port 8080
	//javascript websockets is used to connect to the broker and the data recieved is used to plot the real time data using chart.js

	http.HandleFunc("/", indexPage)
	http.ListenAndServe(":8081", nil) //starts a server at 8000 port
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	page:="hello"

	//passing our data to index html
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, page)
}
