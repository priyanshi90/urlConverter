package main

import (
	"encoding/json"
	"encoding/base64"
	"net/http"
	"strconv"
	"database/sql"
	"strings"
	_"github.com/go-sql-driver/mysql"
)

type longUrlStruct struct {
	Url string `json:"url"`
}

type shortUrlStruct struct {
	Short string `json:"short"`
}

type originalUrlStruct struct {
	Original string `json:"original"`
}

type errorStruct struct {
	Exception string `json:"exception"`
}

const(
	sqlConnectError string = "Error while Connecting to Database"
	dbError string = "Error while processing Database"
	shortUrlNotExistsError string = "Short URL does not exist"
	jsonError string = "Wrong JSON Format"
	emptyURLError string = "No URL provided"	
)

func mysqlInsert(longUrl string)(mysqlInsertId string, exception string) {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/urlConverter?charset=utf8")
	if err != nil {
		exception = sqlConnectError
	} else {
		stmt, err := db.Prepare("INSERT URL_LIST SET LONG_URL = ?")
		if err != nil {
			exception = dbError
		}
		res, err := stmt.Exec(longUrl)
		if err != nil {
			exception = dbError
		}
		id, err := res.LastInsertId()
		if err != nil {
			exception = dbError
		}
		mysqlInsertId = strconv.Itoa(int(id))
	}
	return
}

func sendShortUrl(id string)(shortUrl string) {
	shortUrl = base64.StdEncoding.EncodeToString([]byte(id))
	return	
}

func getMysqlID(shortUrl string)(mysqlInsertId int, exception string) {
	temp := strings.Split(shortUrl, "/")
	lastIndex := len(temp) -1
	id, err := base64.StdEncoding.DecodeString(temp[lastIndex])
	if err != nil {
		exception = shortUrlNotExistsError
	}
	mysqlInsertId, err = strconv.Atoi(string(id))
	if err != nil {
		exception = shortUrlNotExistsError
	}
	return
}

func mysqlGetLongUrl(id int)(longUrl string, exception string){
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/urlConverter?charset=utf8")
	if err != nil {
		exception = sqlConnectError
	} else {
		rows, err := db.Query("SELECT LONG_URL FROM URL_LIST WHERE ID = ?", id)
		if err != nil {
			exception = shortUrlNotExistsError
		}
		rows.Next()
		var LONG_URL string
		err = rows.Scan(&LONG_URL)
		if err != nil {
			exception = shortUrlNotExistsError
		}
		longUrl = LONG_URL
	}
	return
}

func handleShortRequest(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t longUrlStruct
	var exception string = ""
	err := decoder.Decode(&t)
	if err != nil {
		exception = jsonError
		response := errorStruct{Exception: exception}
		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode(response)
		return
	}
	if len(t.Url)>0{
		mysqlInsertId, exception := mysqlInsert(t.Url)
		if(len(exception) > 0) {
			response := errorStruct{Exception: exception}
			rw.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(rw).Encode(response)
			return
		}
		shortUrl := sendShortUrl(mysqlInsertId)
		temp := []string{"http://localhost/", shortUrl}
		var responseUrl = strings.Join(temp, "")
		response := shortUrlStruct{Short: responseUrl}
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(response)
		return
	} else {
		exception = emptyURLError
		response := errorStruct{Exception: exception}
		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode(response)
		return
	}

}

func handleOriginalRequest(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t shortUrlStruct
	var exception string = ""
	err := decoder.Decode(&t)
	if err != nil {
		exception = jsonError
		response := errorStruct{Exception: exception}
		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode(response)
		return
	}
	if len(t.Short)>0{
		mysqlInsertId, exception := getMysqlID(t.Short)
		if(len(exception) > 0) {
			response := errorStruct{Exception: exception}
			rw.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(rw).Encode(response)
			return
		}
		longUrl, exception := mysqlGetLongUrl(mysqlInsertId)
		if(len(exception) > 0) {
			response := errorStruct{Exception: exception}
			rw.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(rw).Encode(response)
			return
		}
		response := originalUrlStruct{Original: longUrl}
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(response)
		return
	} else {
		exception = emptyURLError
		response := errorStruct{Exception: exception}
		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode(response)
		return
	}
}

