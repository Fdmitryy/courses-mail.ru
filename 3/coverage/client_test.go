package main

import (
	"golang-2019-1"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

// тут писать код тестов

type TestCase struct {
	Name    string
	Req     golang_2019_1.SearchRequest
	Ans     *golang_2019_1.SearchResponse
	isError bool
}

func TestSearchClient_FindUsers(t *testing.T) {
	cases := []TestCase{
		TestCase{
			Name: "check order id, desc",
			Req: golang_2019_1.SearchRequest{
				Limit:      26,
				Offset:     0,
				Query:      "we",
				OrderField: "Id",
				OrderBy:    golang_2019_1.OrderByDesc,
			},
			Ans: &golang_2019_1.SearchResponse{
				Users: []golang_2019_1.User{
					{
						Id:     20,
						Name:   "Lowery York",
						Age:    27,
						About:  "Dolor enim sit id dolore enim sint nostrud deserunt. Occaecat minim enim veniam proident mollit Lorem irure ex. Adipisicing pariatur adipisicing aliqua amet proident velit. Magna commodo culpa sit id.\n",
						Gender: "male",
					},
					{
						Id:     10,
						Name:   "Henderson Maxwell",
						Age:    30,
						About:  "Ex et excepteur anim in eiusmod. Cupidatat sunt aliquip exercitation velit minim aliqua ad ipsum cillum dolor do sit dolore cillum. Exercitation eu in ex qui voluptate fugiat amet.\n",
						Gender: "male",
					},
					{
						Id:     4,
						Name:   "Owen Lynn",
						Age:    30,
						About:  "Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n",
						Gender: "male",
					},
				},
				NextPage: false,
			},
			isError: false,
		},
		TestCase{
			Name: "check order name, desc, offset",
			Req: golang_2019_1.SearchRequest{
				Limit:      20,
				Offset:     1,
				Query:      "we",
				OrderField: "Name",
				OrderBy:    golang_2019_1.OrderByDesc,
			},
			Ans: &golang_2019_1.SearchResponse{
				Users: []golang_2019_1.User{
					{
						Id:     20,
						Name:   "Lowery York",
						Age:    27,
						About:  "Dolor enim sit id dolore enim sint nostrud deserunt. Occaecat minim enim veniam proident mollit Lorem irure ex. Adipisicing pariatur adipisicing aliqua amet proident velit. Magna commodo culpa sit id.\n",
						Gender: "male",
					},
					{
						Id:     10,
						Name:   "Henderson Maxwell",
						Age:    30,
						About:  "Ex et excepteur anim in eiusmod. Cupidatat sunt aliquip exercitation velit minim aliqua ad ipsum cillum dolor do sit dolore cillum. Exercitation eu in ex qui voluptate fugiat amet.\n",
						Gender: "male",
					},
				},
				NextPage: false,
			},
			isError: false,
		},
		TestCase{
			Name: "check order age, desc, limit < data",
			Req: golang_2019_1.SearchRequest{
				Limit:      2,
				Offset:     0,
				Query:      "we",
				OrderField: "Age",
				OrderBy:    golang_2019_1.OrderByDesc,
			},
			Ans: &golang_2019_1.SearchResponse{
				Users: []golang_2019_1.User{
					{
						Id:     4,
						Name:   "Owen Lynn",
						Age:    30,
						About:  "Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n",
						Gender: "male",
					},
					{
						Id:     10,
						Name:   "Henderson Maxwell",
						Age:    30,
						About:  "Ex et excepteur anim in eiusmod. Cupidatat sunt aliquip exercitation velit minim aliqua ad ipsum cillum dolor do sit dolore cillum. Exercitation eu in ex qui voluptate fugiat amet.\n",
						Gender: "male",
					},
				},
				NextPage: true,
			},
			isError: false,
		},
		TestCase{
			Name: "check order empty string(Name), desc",
			Req: golang_2019_1.SearchRequest{
				Limit:      20,
				Offset:     0,
				Query:      "we",
				OrderField: "",
				OrderBy:    golang_2019_1.OrderByDesc,
			},
			Ans: &golang_2019_1.SearchResponse{
				Users: []golang_2019_1.User{
					{
						Id:     4,
						Name:   "Owen Lynn",
						Age:    30,
						About:  "Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n",
						Gender: "male",
					},
					{
						Id:     20,
						Name:   "Lowery York",
						Age:    27,
						About:  "Dolor enim sit id dolore enim sint nostrud deserunt. Occaecat minim enim veniam proident mollit Lorem irure ex. Adipisicing pariatur adipisicing aliqua amet proident velit. Magna commodo culpa sit id.\n",
						Gender: "male",
					},
					{
						Id:     10,
						Name:   "Henderson Maxwell",
						Age:    30,
						About:  "Ex et excepteur anim in eiusmod. Cupidatat sunt aliquip exercitation velit minim aliqua ad ipsum cillum dolor do sit dolore cillum. Exercitation eu in ex qui voluptate fugiat amet.\n",
						Gender: "male",
					},
				},
				NextPage: false,
			},
			isError: false,
		},
		TestCase{
			Name: "check bad order field",
			Req: golang_2019_1.SearchRequest{
				Limit:      15,
				Offset:     0,
				Query:      "Boyd",
				OrderField: "Bad",
				OrderBy:    golang_2019_1.OrderByDesc,
			},
			isError: true,
		},
		TestCase{
			Name: "check order id, asc",
			Req: golang_2019_1.SearchRequest{
				Limit:      20,
				Offset:     0,
				Query:      "we",
				OrderField: "Id",
				OrderBy:    golang_2019_1.OrderByAsc,
			},
			Ans: &golang_2019_1.SearchResponse{
				Users: []golang_2019_1.User{
					{
						Id:     4,
						Name:   "Owen Lynn",
						Age:    30,
						About:  "Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n",
						Gender: "male",
					},
					{
						Id:     10,
						Name:   "Henderson Maxwell",
						Age:    30,
						About:  "Ex et excepteur anim in eiusmod. Cupidatat sunt aliquip exercitation velit minim aliqua ad ipsum cillum dolor do sit dolore cillum. Exercitation eu in ex qui voluptate fugiat amet.\n",
						Gender: "male",
					},
					{
						Id:     20,
						Name:   "Lowery York",
						Age:    27,
						About:  "Dolor enim sit id dolore enim sint nostrud deserunt. Occaecat minim enim veniam proident mollit Lorem irure ex. Adipisicing pariatur adipisicing aliqua amet proident velit. Magna commodo culpa sit id.\n",
						Gender: "male",
					},
				},
				NextPage: false,
			},
			isError: false,
		},
		TestCase{
			Name: "check order name, asc, limit + offset >= data",
			Req: golang_2019_1.SearchRequest{
				Limit:      2,
				Offset:     1,
				Query:      "we",
				OrderField: "Name",
				OrderBy:    golang_2019_1.OrderByAsc,
			},
			Ans: &golang_2019_1.SearchResponse{
				Users: []golang_2019_1.User{
					{
						Id:     20,
						Name:   "Lowery York",
						Age:    27,
						About:  "Dolor enim sit id dolore enim sint nostrud deserunt. Occaecat minim enim veniam proident mollit Lorem irure ex. Adipisicing pariatur adipisicing aliqua amet proident velit. Magna commodo culpa sit id.\n",
						Gender: "male",
					},

					{
						Id:     4,
						Name:   "Owen Lynn",
						Age:    30,
						About:  "Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n",
						Gender: "male",
					},
				},
				NextPage: false,
			},
			isError: false,
		},
		TestCase{
			Name: "check order age, asc",
			Req: golang_2019_1.SearchRequest{
				Limit:      20,
				Offset:     0,
				Query:      "we",
				OrderField: "Age",
				OrderBy:    golang_2019_1.OrderByAsc,
			},
			Ans: &golang_2019_1.SearchResponse{
				Users: []golang_2019_1.User{
					{
						Id:     20,
						Name:   "Lowery York",
						Age:    27,
						About:  "Dolor enim sit id dolore enim sint nostrud deserunt. Occaecat minim enim veniam proident mollit Lorem irure ex. Adipisicing pariatur adipisicing aliqua amet proident velit. Magna commodo culpa sit id.\n",
						Gender: "male",
					},
					{
						Id:     4,
						Name:   "Owen Lynn",
						Age:    30,
						About:  "Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n",
						Gender: "male",
					},
					{
						Id:     10,
						Name:   "Henderson Maxwell",
						Age:    30,
						About:  "Ex et excepteur anim in eiusmod. Cupidatat sunt aliquip exercitation velit minim aliqua ad ipsum cillum dolor do sit dolore cillum. Exercitation eu in ex qui voluptate fugiat amet.\n",
						Gender: "male",
					},
				},
				NextPage: false,
			},
			isError: false,
		},
		TestCase{
			Name: "check order empty string(Name), asc",
			Req: golang_2019_1.SearchRequest{
				Limit:      20,
				Offset:     0,
				Query:      "we",
				OrderField: "",
				OrderBy:    golang_2019_1.OrderByAsc,
			},
			Ans: &golang_2019_1.SearchResponse{
				Users: []golang_2019_1.User{
					{
						Id:     10,
						Name:   "Henderson Maxwell",
						Age:    30,
						About:  "Ex et excepteur anim in eiusmod. Cupidatat sunt aliquip exercitation velit minim aliqua ad ipsum cillum dolor do sit dolore cillum. Exercitation eu in ex qui voluptate fugiat amet.\n",
						Gender: "male",
					},
					{
						Id:     20,
						Name:   "Lowery York",
						Age:    27,
						About:  "Dolor enim sit id dolore enim sint nostrud deserunt. Occaecat minim enim veniam proident mollit Lorem irure ex. Adipisicing pariatur adipisicing aliqua amet proident velit. Magna commodo culpa sit id.\n",
						Gender: "male",
					},

					{
						Id:     4,
						Name:   "Owen Lynn",
						Age:    30,
						About:  "Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n",
						Gender: "male",
					},
				},
				NextPage: false,
			},
			isError: false,
		},
		TestCase{
			Name: "check bad order field",
			Req: golang_2019_1.SearchRequest{
				Limit:      15,
				Offset:     0,
				Query:      "Boyd",
				OrderField: "Bad",
				OrderBy:    golang_2019_1.OrderByAsc,
			},
			isError: true,
		},
		TestCase{
			Name: "check order name, AsIs",
			Req: golang_2019_1.SearchRequest{
				Limit:      20,
				Offset:     0,
				Query:      "we",
				OrderField: "Name",
				OrderBy:    golang_2019_1.OrderByAsIs,
			},
			Ans: &golang_2019_1.SearchResponse{
				Users: []golang_2019_1.User{
					{
						Id:     4,
						Name:   "Owen Lynn",
						Age:    30,
						About:  "Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n",
						Gender: "male",
					},
					{
						Id:     10,
						Name:   "Henderson Maxwell",
						Age:    30,
						About:  "Ex et excepteur anim in eiusmod. Cupidatat sunt aliquip exercitation velit minim aliqua ad ipsum cillum dolor do sit dolore cillum. Exercitation eu in ex qui voluptate fugiat amet.\n",
						Gender: "male",
					},
					{
						Id:     20,
						Name:   "Lowery York",
						Age:    27,
						About:  "Dolor enim sit id dolore enim sint nostrud deserunt. Occaecat minim enim veniam proident mollit Lorem irure ex. Adipisicing pariatur adipisicing aliqua amet proident velit. Magna commodo culpa sit id.\n",
						Gender: "male",
					},
				},
				NextPage: false,
			},
			isError: false,
		},
		TestCase{
			Name: "check limit < 0",
			Req: golang_2019_1.SearchRequest{
				Limit:      -1,
				Offset:     0,
				Query:      "Boyd",
				OrderField: "Id",
				OrderBy:    golang_2019_1.OrderByDesc,
			},
			isError: true,
		},
		TestCase{
			Name: "check offset < 0",
			Req: golang_2019_1.SearchRequest{
				Limit:      3,
				Offset:     -1,
				Query:      "Boyd",
				OrderField: "Id",
				OrderBy:    golang_2019_1.OrderByDesc,
			},
			isError: true,
		},
		TestCase{
			Name: "check bad orderBy",
			Req: golang_2019_1.SearchRequest{
				Limit:      15,
				Offset:     0,
				Query:      "Boyd",
				OrderField: "Name",
				OrderBy:    23,
			},
			isError: true,
		},
		TestCase{
			Name: "check offset > data",
			Req: golang_2019_1.SearchRequest{
				Limit:      15,
				Offset:     3,
				Query:      "Boyd",
				OrderField: "Name",
				OrderBy:    golang_2019_1.OrderByDesc,
			},
			isError: true,
		},
		TestCase{
			Name:    "empty request",
			isError: true,
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(golang_2019_1.SearchServer))
	defer ts.Close()
	for caseNum, item := range cases {
		client := golang_2019_1.SearchClient{
			AccessToken: golang_2019_1.GoodToken,
			URL:         ts.URL,
		}
		result, err := client.FindUsers(item.Req)
		if err != nil && !item.isError {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && item.isError {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !reflect.DeepEqual(result, item.Ans) {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
			t.Error("expected ", item.Ans, "\nreturned ", result)
		}
	}
}

func TestUninitializedUrl(t *testing.T) {
	cases := []TestCase{
		{
			isError: true,
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(golang_2019_1.SearchServer))
	defer ts.Close()
	for caseNum, item := range cases {
		client := golang_2019_1.SearchClient{
			AccessToken: golang_2019_1.GoodToken,
		}
		result, err := client.FindUsers(item.Req)
		if err != nil && !item.isError {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && item.isError {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !reflect.DeepEqual(result, item.Ans) {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
			t.Error("expected ", item.Ans, "\nreturned ", result)
		}
	}
}

func TestBadToken(t *testing.T) {
	cases := []TestCase{
		{
			isError: true,
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(golang_2019_1.SearchServer))
	defer ts.Close()
	for caseNum, item := range cases {
		client := golang_2019_1.SearchClient{
			AccessToken: "1",
			URL:         ts.URL,
		}
		result, err := client.FindUsers(item.Req)
		if err != nil && !item.isError {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && item.isError {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !reflect.DeepEqual(result, item.Ans) {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
			t.Error("expected ", item.Ans, "\nreturned ", result)
		}
	}
}

func TestTimeout(t *testing.T) {
	cases := []TestCase{
		{
			isError: true,
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(golang_2019_1.SearchServer))
	defer ts.Close()
	golang_2019_1.client.Timeout = time.Nanosecond
	defer func() { golang_2019_1.client.Timeout = time.Second }()
	for caseNum, item := range cases {
		client := golang_2019_1.SearchClient{
			AccessToken: golang_2019_1.GoodToken,
			URL:         ts.URL,
		}
		result, err := client.FindUsers(item.Req)
		if err != nil && !item.isError {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && item.isError {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !reflect.DeepEqual(result, item.Ans) {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
			t.Error("expected ", item.Ans, "\nreturned ", result)
		}
	}
}

func TestInternalError(t *testing.T) {
	cases := []TestCase{
		{
			isError: true,
		},
	}
	golang_2019_1.Filename = "bad"
	defer func() { golang_2019_1.Filename = "dataset.xml" }()
	ts := httptest.NewServer(http.HandlerFunc(golang_2019_1.SearchServer))
	defer ts.Close()
	for caseNum, item := range cases {
		client := golang_2019_1.SearchClient{
			AccessToken: golang_2019_1.GoodToken,
			URL:         ts.URL,
		}
		result, err := client.FindUsers(item.Req)
		if err != nil && !item.isError {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && item.isError {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !reflect.DeepEqual(result, item.Ans) {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
			t.Error("expected ", item.Ans, "\nreturned ", result)
		}
	}
}
