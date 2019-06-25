package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные

func NewDbExplorer(db *sql.DB) (http.Handler, error) {
	type Info struct {
		Field   string
		Type    string
		Null    string
		Key     string
		Default sql.NullString
		Extra   sql.NullString
	}
	var table string
	var tables []string
	var fullInfo = make(map[string][]Info)

	rows, err := db.Query("SHOW tables")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&table)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	rows.Close()

	for _, tab := range tables {
		var tableInfo Info
		var tablesInfo []Info

		rows2, err := db.Query("SHOW COLUMNS FROM " + tab)
		if err != nil {
			return nil, err
		}

		for rows2.Next() {
			err := rows2.Scan(&tableInfo.Field, &tableInfo.Type, &tableInfo.Null, &tableInfo.Key, &tableInfo.Default, &tableInfo.Extra)
			if err != nil {
				return nil, err
			}
			tablesInfo = append(tablesInfo, tableInfo)
		}
		rows2.Close()
		fullInfo[tab] = tablesInfo
	}

	findTable := func(name string) bool {
		_, exits := fullInfo[name]
		return exits
	}

	findPkAndNotNulls := func(Info []Info) (string, map[string]bool) {
		var pk string
		NotNulls := make(map[string]bool)
		for _, v := range Info {
			if v.Key == "PRI" {
				pk = v.Field
			}
			if v.Null == "NO" {
				NotNulls[v.Field] = true
			}
		}
		return pk, NotNulls
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			path := strings.TrimPrefix(r.URL.Path, "/")
			var resp = make(map[string]interface{})
			if path != "" {
				var table string
				var id string
				var args []string
				params := r.URL.Query()

				if strings.Contains(path, "/") {
					args = strings.Split(path, "/")
					table = args[0]
					id = args[1]
					var types []string
					var names []string

					for _, val := range fullInfo[table] {
						types = append(types, val.Type)
						names = append(names, val.Field)
					}

					if !findTable(table) {
						w.WriteHeader(http.StatusNotFound)
						resp["error"] = "unknown table"
						data, _ := json.Marshal(resp)
						w.Write(data)
						return
					}

					primaryKey, _ := findPkAndNotNulls(fullInfo[table])
					getOneRecord(db, table, primaryKey, id, len(fullInfo[table]), types, names, resp)
					jsonData, _ := json.Marshal(resp)
					if _, exist := resp["error"]; exist {
						w.WriteHeader(http.StatusNotFound)
					}
					w.Write(jsonData)
					return
				}

				var limit int
				var offset int
				table = path
				var types []string
				var names []string

				for _, val := range fullInfo[table] {
					types = append(types, val.Type)
					names = append(names, val.Field)
				}

				if len(params) != 0 {
					limit, _ = strconv.Atoi(params.Get("limit"))
					offset, _ = strconv.Atoi(params.Get("offset"))
				}

				if !findTable(table) {
					w.WriteHeader(http.StatusNotFound)
					resp["error"] = "unknown table"
					jsonData, _ := json.Marshal(resp)
					w.Write(jsonData)
					return
				}

				err = getAllRecords(db, table, len(fullInfo[table]), limit, offset, types, names, resp)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
				jsonData, _ := json.Marshal(resp)
				w.Write(jsonData)
				return
			}

			var m = make(map[string]interface{})
			m["tables"] = tables
			resp["response"] = m
			data, _ := json.Marshal(resp)
			w.Write(data)

		case http.MethodPut:
			path := strings.TrimPrefix(r.URL.Path, "/")
			path = strings.TrimSuffix(path, "/")
			respBody, err := ioutil.ReadAll(r.Body)
			var body interface{}
			err = json.Unmarshal(respBody, &body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			mapBody := body.(map[string]interface{})
			var m = make(map[string]interface{})
			var resp = make(map[string]interface{})
			table = path

			if !findTable(table) {
				w.WriteHeader(http.StatusNotFound)
				resp["error"] = "unknown table"
				jsonData, _ := json.Marshal(resp)
				w.Write(jsonData)
				return
			}

			primaryKey, _ := findPkAndNotNulls(fullInfo[table])
			var fields []string
			var values []interface{}
			var questionSign []string
			var myType string

			for _, v := range fullInfo[table] {
				if _, exist := mapBody[v.Field]; !exist && v.Null == "NO" {
					myType = v.Type
					switch myType {
					case "text", "varchar(255)":
						values = append(values, "")
					case "int(11)":
						values = append(values, 0)
					default:
						values = append(values, nil)
					}
					fields = append(fields, v.Field)
					questionSign = append(questionSign, "?")
				}
			}

			for key, val := range mapBody {
				keyExist := false
				for _, v := range fullInfo[table] {
					if key == v.Field {
						myType = v.Type
						keyExist = true
					}
				}

				if !keyExist || key == primaryKey {
					continue
				}

				switch myType {
				case "text", "varchar(255)":
					values = append(values, val.(string))
				case "int(11)":
					values = append(values, strconv.Itoa(int(val.(float64))))
				default:
					values = append(values, val)
				}
				fields = append(fields, key)
				questionSign = append(questionSign, "?")
			}

			insertStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(fields, ","), strings.Join(questionSign, ","))
			result, err := db.Exec(insertStr, values...)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			id, err := result.LastInsertId()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			m[primaryKey] = id
			resp["response"] = m
			jsonData, _ := json.Marshal(resp)
			w.Write(jsonData)

		case http.MethodPost:
			path := strings.TrimPrefix(r.URL.Path, "/")
			respBody, err := ioutil.ReadAll(r.Body)
			var body interface{}
			err = json.Unmarshal(respBody, &body)
			mapBody := body.(map[string]interface{})
			var m = make(map[string]interface{})
			var resp = make(map[string]interface{})
			args := strings.Split(path, "/")
			table := args[0]
			id := args[1]

			if !findTable(table) {
				w.WriteHeader(http.StatusNotFound)
				resp["error"] = "unknown table"
				data, _ := json.Marshal(resp)
				w.Write(data)
				return
			}

			var values []interface{}
			var fields []string
			var myType string
			for key, val := range mapBody {
				primaryKey, notNulls := findPkAndNotNulls(fullInfo[table])

				keyExist := false
				for _, v := range fullInfo[table] {
					if key == v.Field {
						keyExist = true
						myType = v.Type
					}
				}

				if !keyExist {
					continue
				}

				_, exist := notNulls[key]
				if val == nil && exist || key == primaryKey {
					w.WriteHeader(http.StatusBadRequest)
					resp["error"] = "field " + key + " have invalid type"
					data, _ := json.Marshal(resp)
					w.Write(data)
					return
				}

				if val == nil {
					fields = append(fields, key + " = ?")
					values = append(values, nil)
					continue
				}

				switch myType {
				case "text", "varchar(255)":
					if reflect.ValueOf(val).Kind() != reflect.String {
						w.WriteHeader(http.StatusBadRequest)
						resp["error"] = "field " + key + " have invalid type"
						data, _ := json.Marshal(resp)
						w.Write(data)
						return
					}
					fields = append(fields, key + " = ?")
					values = append(values, val.(string))
				case "int(11)":
					if reflect.ValueOf(val).Kind() != reflect.Int {
						w.WriteHeader(http.StatusBadRequest)
						resp["error"] = "field " + key + " have invalid type"
						data, _ := json.Marshal(resp)
						w.Write(data)
						return
					}
					fields = append(fields, key + " = ?")
					values = append(values, strconv.Itoa(int(val.(float64))))
				default:
					fields = append(fields, key + " = ?")
					values = append(values, nil)
				}

			}

			primaryKey, _ := findPkAndNotNulls(fullInfo[table])
			values=append(values, id)
			insertStr := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?", table, strings.Join(fields, ","), primaryKey)
			result, err := db.Exec(insertStr,values...)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			affectedCount, err := result.RowsAffected()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			m["updated"] = affectedCount
			resp["response"] = m
			jsonData, _ := json.Marshal(resp)
			w.Write(jsonData)

		case http.MethodDelete:
			path := strings.TrimPrefix(r.URL.Path, "/")
			var m = make(map[string]interface{})
			var resp = make(map[string]interface{})
			args := strings.Split(path, "/")
			table := args[0]
			id := args[1]

			if !findTable(table) {
				w.WriteHeader(http.StatusNotFound)
				resp["error"] = "unknown table"
				jsonData, _ := json.Marshal(resp)
				w.Write(jsonData)
				return
			}

			primaryKey, _ := findPkAndNotNulls(fullInfo[table])
			insertStr := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", table, primaryKey)
			result, err := db.Exec(insertStr, id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			affectedCount, err := result.RowsAffected()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			m["deleted"] = affectedCount
			resp["response"] = m
			jsonData, _ := json.Marshal(resp)
			w.Write(jsonData)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})
	return mux, nil
}

func getOneRecord(db *sql.DB, table string, primaryKey string, id string, colCount int, types []string, names []string, resp map[string]interface{}) {
	var m = make(map[string]interface{})
	insertStr := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", table, primaryKey)
	row := db.QueryRow(insertStr, id)
	m["record"] = make(map[string]interface{})
	var forScan = make([]interface{}, colCount)
	initForScan(forScan, types)

	err := row.Scan(forScan...)
	if err != nil {
		resp["error"] = "record not found"
		return
	}

	var data = make(map[string]interface{})
	extractForScan(forScan, types, names, data)
	m["record"] = data
	resp["response"] = m
	return
}

func getAllRecords(db *sql.DB, table string, colCount int, limit int, offset int, types []string, names []string, resp map[string]interface{}) error {
	var m = make(map[string]interface{})

	if limit == 0 {
		limit = -1
	}

	insertStr := fmt.Sprintf("SELECT * FROM %s LIMIT ?, ?", table)
	rows, err := db.Query(insertStr, offset, limit)
	if err != nil {
		return err
	}

	m["records"] = make([]interface{}, 0)
	var items = make([]interface{}, 0)
	var forScan = make([]interface{}, colCount)
	initForScan(forScan, types)

	for rows.Next() {
		var data = make(map[string]interface{})
		err := rows.Scan(forScan...)
		if err != nil {
			return err
		}
		extractForScan(forScan, types, names, data)
		items = append(items, data)
	}
	rows.Close()

	m["records"] = items
	resp["response"] = m
	return nil
}

func initForScan(forScan []interface{}, types []string) {
	for i := range forScan {
		switch types[i] {
		case "text", "varchar(255)":
			forScan[i] = new(sql.NullString)
		case "int(11)":
			forScan[i] = new(sql.NullInt64)
		default:
			forScan[i] = new(interface{})
		}
	}
}

func extractForScan(forScan []interface{}, types []string, names []string, data map[string]interface{}) {
	for i := range forScan {
		fieldType := types[i]
		name := names[i]

		switch fieldType {
		case "text", "varchar(255)":
			s := string((*forScan[i].(*sql.NullString)).String)
			flag := (*forScan[i].(*sql.NullString)).Valid
			if flag {
				data[name] = s
			} else {
				data[name] = nil
			}
		case "int(11)":
			i := (*forScan[i].(*sql.NullInt64)).Int64
			data[name] = i
		default:
			something := forScan[i]
			data[name] = something
		}
	}
}
