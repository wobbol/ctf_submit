package main

import (
	"fmt"
	"html"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

type rootHandler struct {
	ctf CTF
}

func (h *rootHandler) submitHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/submit")
	status := "None"
	if r.Method == http.MethodPost {
		flag := html.EscapeString(r.PostFormValue("flag"))
		name := html.EscapeString(r.PostFormValue("name"))
		err := h.ctf.SubmitFlag(name, flag)
		fmt.Println(h.ctf)
		if err != nil {
			fmt.Println(err)
			status = "Failure"
		} else {
			status = "Success"
		}
		h.ctf.SaveCTF("./whatctf")
	}
	submit, _ := template.ParseFiles("./submit.html")
	submit.Execute(w, nil)
	io.WriteString(w, "Last submit: "+status)
	h.scoreboardHandler(w, r)
}
func (h *rootHandler) scoreboardHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/scoreboard")
	t, _ := template.ParseFiles("./scoreboard.html")
	t.Execute(w, h.ctf)
}
func (h *rootHandler) scorebreakdownHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/scorebreakdown")
	t, _ := template.ParseFiles("./scorebreakdown.html")
	t.Execute(w, h.ctf)
}
func (h *rootHandler) flagboardHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/flagboard")
	t, _ := template.ParseFiles("./flagboard.html")
	t.Execute(w, h.ctf)
}

func (h *rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.ToLower(r.URL.Path)
	style, _ := template.ParseFiles("./pageheader.html")
	style.Execute(w, nil)
	if r.URL.Path == "/" {
		h.scoreboardHandler(w, r)
		h.flagboardHandler(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/submit") {
		h.submitHandler(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/scoreboard") {
		h.scoreboardHandler(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/flagboard") {
		h.flagboardHandler(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/scorebreakdown") {
		h.scorebreakdownHandler(w, r)
		return
	}

}
func httpsredirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://localhost:8081"+r.RequestURI, http.StatusMovedPermanently)
}

func WriteExample() {
	ctf := CTF{}
	ctf.Flagboard.AddFlag(Flag{Data: "aflag", Desc: "example Flag1", Pts: 10})
	ctf.Flagboard.AddFlag(Flag{Data: "bflag", Desc: "example Flag2", Pts: 10})
	ctf.Flagboard.AddFlag(Flag{Data: "cflag", Desc: "example Flag3", Pts: 10})
	ctf.Flagboard.AddFlag(Flag{Data: "dflag", Pts: 10})

	ctf.Scoreboard.AddUser(User{Name: "aardvark"})
	ctf.Scoreboard.AddUser(User{Name: "apple"})
	ctf.SubmitFlag("aardvark", "aflag")
	ctf.SaveCTF("./example")
}

func main() {
	// Done have a file with all the flag values and their worth
	// Done have a file with all the players and their collected flags
	// Done display the player list and their points via http
	//      collect the flags and player names only via TLS/https POST request

	WriteExample()
	//ctf2 := NewCTF("./example")
	//fmt.Printf("%+v\n", ctf2)
	fmt.Println("ctf-scoreboard")
	example := rootHandler{ctf: *NewCTF("./example")}
	go http.ListenAndServe(":8080", http.HandlerFunc(httpsredirectHandler))

	//http.Handle("/submit", new(submitHandler))
	http.Handle("/", &example)
	log.Fatal(http.ListenAndServeTLS(":8081", "./cert.pem", "./key.pem", nil))
}
