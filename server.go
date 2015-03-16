package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func handlVote(w http.ResponseWriter, r *http.Request) {
	queryString := r.URL.Query()
	ids, ok := queryString["id"]
	callbacks, exist := queryString["callback"]
	if !ok {
		http.Error(w, "missing voteID", http.StatusBadRequest)
	}
	re, ok := queryString["result"]
	if !ok {
		http.Error(w, "missing resultString", http.StatusBadRequest)
	}
	v := Vote{}
	err := json.Unmarshal([]byte(re[0]), &v.Detail)
	if err != nil {
		http.Error(w, "error while voting", http.StatusNotImplemented)
	}
	err = vote(ids[0], v)
	if err != nil {
		http.Error(w, "oops", http.StatusNotImplemented)
	}
	result := "success"
	if exist {
		result = callbacks[0] + "(" + result + ")"
	}
	w.Write([]byte(result))
}

func handlNewVote(w http.ResponseWriter, r *http.Request) {
	vote := Vote{}
	var err error
	vote.Title = r.PostFormValue("title")
	vote.Description = r.PostFormValue("description")
	vote.Vtype, err = strconv.Atoi(r.PostFormValue("vtype"))
	callbacks, ok := r.URL.Query()["callback"]
	if err != nil || vote.Vtype < 0 || vote.Vtype > 1 {
		http.Error(w, "wrong voteID", http.StatusBadRequest)
	}
	err = json.Unmarshal([]byte(r.PostFormValue("detail")), &vote.Detail)
	if err != nil {
		http.Error(w, "wrong detail", http.StatusBadRequest)
	}
	var result string
	result, err = newVote(vote)
	if err != nil {
		http.Error(w, "wrong detail", http.StatusNotImplemented)
	}
	if ok {
		result = callbacks[0] + "(" + result + ")"
	}
	w.Write([]byte(result))
}

func handlGetResult(w http.ResponseWriter, r *http.Request) {
	id, ok := r.URL.Query()["id"]
	callbacks, exist := r.URL.Query()["callback"]
	if !ok {
		http.Error(w, "wrong id", http.StatusBadRequest)
	}
	ret, err := getResult(id[0])
	result := string(ret)
	if err != nil {
		http.Error(w, "error while geting result", http.StatusNotImplemented)
	}
	if exist {
		result = callbacks[0] + "(" + result + ")"
	}
	w.Write([]byte(result))
}

func main() {
	http.HandleFunc("/vote", handlVote)
	http.HandleFunc("/result", handlGetResult)
	http.HandleFunc("/newvote", handlNewVote)
	if err := http.ListenAndServe("localhost:2000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
