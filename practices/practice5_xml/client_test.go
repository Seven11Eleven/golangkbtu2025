package main

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

// код писать тут

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

//searcherParams.Add("limit", strconv.Itoa(req.Limit))
//searcherParams.Add("offset", strconv.Itoa(req.Offset))
//searcherParams.Add("query", req.Query)
//searcherParams.Add("order_field", req.OrderField)
//searcherParams.Add("order_by", strconv.Itoa(req.OrderBy))

//	searcherReq.Header.Add("AccessToken", srv.AccessToken)

func SearchServer(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	orderField := r.URL.Query().Get("order_field")
	accessToken := r.Header.Get("AccessToken")

	if accessToken != "1337" {
		http.Error(w, `{"Error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	validFields := map[string]bool{"query": true, "": false}
	if !validFields[orderField] {
		http.Error(w, `{"Error": "ErrorBadOrderField"}`, http.StatusBadRequest)
		return
	}

	file, err := os.Open("./dataset.xml")
	if err != nil {
		log.Print(err)
		http.Error(w, `{"Error": "error opening dataset"}`, http.StatusInternalServerError)
		return
	}
	defer file.Close()

	user := &Root{}
	xmlDec := xml.NewDecoder(file)
	if err := xmlDec.Decode(&user); err != nil {
		log.Print(err)
		http.Error(w, `{"Error": "error decoding dataset"}`, http.StatusInternalServerError)
		return
	}

	var res []Row
	for _, v := range user.Rows {
		if query == "" || contains(query, v.LastName) || contains(query, v.FirstName) || contains(query, v.About) {
			res = append(res, v)
		}
		if limit > 0 && len(res) <= limit {
			break
		}
	}

	if len(res) == 0 {
		http.Error(w, `{"Error": "user not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, `{"Error": "failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}

func contains(q, v string) bool {
	return strings.Contains(strings.ToUpper(q), strings.ToUpper(v))
}

type TestCase struct {
	Request *SearchRequest
	IsError bool
	LenRes  int
}

var TestCases = []TestCase{{
	Request: &SearchRequest{
		Limit:      1,
		Offset:     1,
		Query:      "",
		OrderField: "query",
		OrderBy:    1,
	},
	IsError: false,
	LenRes:  1,
}, {
	Request: &SearchRequest{
		Limit:      -10,
		OrderField: "query",
	},
	IsError: true,
	LenRes:  0,
},
	{
		Request: &SearchRequest{
			Limit:      26,
			OrderField: "query",
		},
		IsError: false,
		LenRes:  0,
	},
	{
		Request: &SearchRequest{
			Limit:      26,
			Offset:     -1,
			OrderField: "query",
		},
		IsError: true,
		LenRes:  0,
	},
}

func TestCase1(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	testSearch := &SearchClient{
		AccessToken: "1337",
		URL:         ts.URL,
	}

	for numCase, testCase := range TestCases {
		res, err := testSearch.FindUsers(*testCase.Request)
		if err != nil && !testCase.IsError {
			t.Errorf("case [%d] unexpected error: %#v", numCase, err)
		}
		if err == nil && testCase.IsError {
			t.Errorf("case [%d] expected error, got nil", numCase)
		}
		if testCase.LenRes != 0 && testCase.LenRes != len(res.Users) {
			t.Errorf("results not matched, wait %d, take %d, case num: %d", testCase.LenRes, len(res.Users), numCase)
		}
	}
}

func TestTimeoutError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	defer ts.Close()

	testSearch := &SearchClient{
		AccessToken: "1337",
		URL:         ts.URL,
	}

	_, err := testSearch.FindUsers(SearchRequest{})
	if err == nil {
		t.Errorf("expected timeout error, got nil")
	}

	if !strings.Contains(err.Error(), "timeout") {
		t.Errorf("expected timeout error, got: %v", err)
	}
}

func TestUnknownError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	testSearch := &SearchClient{
		AccessToken: "1337",
		URL:         ts.URL,
	}

	_, err := testSearch.FindUsers(SearchRequest{})
	if err == nil {
		t.Errorf("expected unknown error, got nil")
	}
}

func TestAccessToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	defer ts.Close()

	testSearch := &SearchClient{
		AccessToken: "133",
		URL:         ts.URL,
	}

	_, err := testSearch.FindUsers(SearchRequest{})
	if err == nil {
		t.Errorf("expected access token error, got nil")
	}

	if !strings.Contains(err.Error(), "Bad AccessToken") {
		t.Errorf("expected access token error, got: %v", err)
	}
}

func TestBadRequestHandling(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	testSearch := &SearchClient{
		AccessToken: "1337",
		URL:         ts.URL,
	}

	_, err := testSearch.FindUsers(SearchRequest{OrderField: "invalid"})
	if err == nil || !strings.Contains(err.Error(), "OrderFeld invalid") {
		t.Errorf("expected OrderField error, got: %v", err)
	}
}

//func TestCase2(t *testing.T) {
//	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
//	defer ts.Close()
//	tc := TestCase{
//		Request: &SearchRequest{
//			Limit:      25,
//			Offset:     1,
//			Query:      "",
//			OrderField: "query",
//			OrderBy:    1,
//		},
//		IsError: false,
//		LenRes:  1,
//	}
//	testSearch := &SearchClient{
//		AccessToken: "1337",
//		URL:         ts.URL,
//	}
//
//	res, err := testSearch.FindUsers(tc)
//}
