package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Point struct {
	UID      string
	Interest string
}

var (
	err error
	c   *mgo.Collection
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	dbURL := os.Getenv("DBURL")
	if dbURL == "" {
		dbURL = "127.0.0.1"
	}

	session, err := mgo.Dial(dbURL)
	check(err)
	defer session.Close()
	c = session.DB("db").C("points")

	http.HandleFunc("/", root)

	http.HandleFunc("/insert", insert)
	http.HandleFunc("/find", find)
	http.HandleFunc("/remove", remove)
	fmt.Println("Serving on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func root(w http.ResponseWriter, r *http.Request) {

	readCookie, _ := r.Cookie("csID")
	fmt.Println("readCookie:", readCookie)

	http.SetCookie(w, &http.Cookie{Name: "csID", Value: "def456"})
	r.ParseForm()
	fmt.Println(r.Form)

	io.WriteString(w, " ")
}

func insert(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	err := c.Insert(Point{r.FormValue("uid"), r.FormValue("interest")})
	check(err)
	io.WriteString(w, "Successfully inserted "+r.FormValue("uid")+"\n")
}

func find(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	res := Point{}
	err = c.Find(bson.M{"uid": r.FormValue("uid")}).One(&res)
	check(err)
	resJSON, err := json.Marshal(res)
	check(err)
	io.WriteString(w, string(resJSON)+"\n")
}

func remove(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	err = c.Remove(bson.M{"uid": r.FormValue("uid")})
	check(err)
	io.WriteString(w, "Successfully removed "+r.FormValue("uid")+"\n")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
