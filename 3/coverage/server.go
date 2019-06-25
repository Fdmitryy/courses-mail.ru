package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"golang-2019-1"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// тут писать SearchServer

var (
	GoodToken = "123"
	Filename  = "dataset.xml"
)

type row struct {
	Id        int    `xml:"id"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Name      string
	Age       int    `xml:"age"`
	About     string `xml:"about"`
	Gender    string `xml:"gender"`
}

type root struct {
	List []row `xml:"row"`
}

type rows []row

func (s rows) Len() int      { return len(s) }
func (s rows) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ById struct{ rows }
type ByName struct{ rows }
type ByAge struct{ rows }

func (s ByName) Less(i, j int) bool { return s.rows[i].Name < s.rows[j].Name }
func (s ById) Less(i, j int) bool   { return s.rows[i].Id < s.rows[j].Id }
func (s ByAge) Less(i, j int) bool  { return s.rows[i].Age < s.rows[j].Age }

type ByIdRev struct{ rows }
type ByNameRev struct{ rows }
type ByAgeRev struct{ rows }

func (s ByNameRev) Less(i, j int) bool { return s.rows[i].Name > s.rows[j].Name }
func (s ByIdRev) Less(i, j int) bool   { return s.rows[i].Id > s.rows[j].Id }
func (s ByAgeRev) Less(i, j int) bool  { return s.rows[i].Age > s.rows[j].Age }

var orders = map[int]func(orderField string, resp rows) error{
	golang_2019_1.OrderByAsIs: func(orderField string, resp rows) error { return nil },
	golang_2019_1.OrderByAsc:  sorting,
	golang_2019_1.OrderByDesc: sortingReverse,
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	query := r.URL.Query().Get("query")
	orderField := r.URL.Query().Get("order_field")
	orderBy, _ := strconv.Atoi(r.URL.Query().Get("order_by"))
	AccessToken := r.Header.Get("AccessToken")

	if AccessToken != GoodToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fileBuffer, err := ioutil.ReadFile(Filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	inputData := []byte(fileBuffer)
	users := new(root)
	xml.Unmarshal(inputData, &users)
	resp := rows{}

	if query == "" {
		resp = users.List
	} else {
		for _, v := range users.List {
			v.Name = v.FirstName + " " + v.LastName
			if strings.Contains(v.Name, query) || strings.Contains(v.About, query) {
				resp = append(resp, v)
			}
		}
	}

	if offset >= len(resp) {
		marshErr := golang_2019_1.SearchErrorResponse{Error: "Empty answer"}
		res, _ := json.Marshal(marshErr)
		fmt.Fprintf(w, string(res))
		return
	}

	if limit == 1 {
		w.WriteHeader(http.StatusBadRequest)
		marshErr := golang_2019_1.SearchErrorResponse{Error: "Bad limit, empty answer"}
		res, _ := json.Marshal(marshErr)
		fmt.Fprintf(w, string(res))
		return
	}

	if offset+limit >= len(resp) {
		limit = len(resp) - offset
	}

	function, exist := orders[orderBy]
	if !exist {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	errField := function(orderField, resp)
	if errField != nil {
		w.WriteHeader(http.StatusBadRequest)
		marshErr := golang_2019_1.SearchErrorResponse{Error: "ErrorBadOrderField"}
		res, _ := json.Marshal(marshErr)
		fmt.Fprintf(w, string(res))
		return
	}

	resp = resp[offset : offset+limit]
	res, _ := json.Marshal(&resp)
	fmt.Fprintf(w, string(res))
}

func sorting(orderField string, resp rows) error {
	switch orderField {
	case "Id":
		sort.Sort(ById{resp})
		return nil
	case "Age":
		sort.Sort(ByAge{resp})
		return nil
	case "Name":
		sort.Sort(ByName{resp})
		return nil
	case "":
		sort.Sort(ByName{resp})
		return nil
	default:
		return errors.New(golang_2019_1.ErrorBadOrderField)
	}
}

func sortingReverse(orderField string, resp rows) error {
	switch orderField {
	case "Id":
		sort.Sort(ByIdRev{resp})
		return nil
	case "Age":
		sort.Sort(ByAgeRev{resp})
		return nil
	case "Name":
		sort.Sort(ByNameRev{resp})
		return nil
	case "":
		sort.Sort(ByNameRev{resp})
		return nil
	default:
		return errors.New(golang_2019_1.ErrorBadOrderField)
	}
}
