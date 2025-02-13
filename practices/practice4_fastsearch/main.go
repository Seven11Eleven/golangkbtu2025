package main

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Row struct {
	ID        int    `xml:"id" json:"id"`
	Age       int    `xml:"age" json:"age"`
	FirstName string `xml:"first_name" json:"first_name"`
	LastName  string `xml:"last_name" json:"last_name"`
	Gender    string `xml:"gender" json:"gender"`
	About     string `xml:"about" json:"about"`
}

type Root struct {
	Rows []Row `xml:"row"`
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	file, err := os.Open("./dataset.xml")
	if err != nil {
		log.Print(err)
		http.Error(w, "error opening dataset", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	user := &Root{}
	xmlDec := xml.NewDecoder(file)
	if err := xmlDec.Decode(&user); err != nil {
		log.Print(err)
		http.Error(w, "error decoding dataset", http.StatusInternalServerError)
		return
	}

	var res []Row
	for _, v := range user.Rows {
		if query == "" || contains(query, v.LastName) || contains(query, v.FirstName) || contains(query, v.About) {

			//log.Print(v)
			//return
			res = append(res, v)
		}
		if limit > 0 && len(res) <= limit {
			break
		}

	}

	if len(res) == 0 {
		http.Error(w, "user not found", http.StatusNotFound)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)

}

func contains(q, v string) bool {
	return strings.Contains(strings.ToUpper(q), strings.ToUpper(v))
}

func main() {

	http.HandleFunc("/", SearchServer)

	//s := &http.Server{
	//	Addr:         "127.0.0.1:1337",
	//	Handler:      ,
	//	ReadTimeout:  10,
	//	WriteTimeout: 10,
	//}
	//if err != nil {
	//	fmt.Print(err)
	//	return
	//}
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
