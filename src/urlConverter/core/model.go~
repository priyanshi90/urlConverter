package main

import (
	"encoding/json"
	"encoding/base64"
	"net/http"
	"fmt"
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

func mysqlInsert(longUrl string)(mysqlInsertId string) {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/urlConverter?charset=utf8")
	checkErr(err)
	stmt, err := db.Prepare("INSERT URL_LIST SET LONG_URL = ?")
	checkErr(err)
	res, err := stmt.Exec(longUrl)
	checkErr(err)
	id, err := res.LastInsertId()
	checkErr(err)
	fmt.Println(id)
	mysqlInsertId = strconv.Itoa(int(id))
	return
}

func sendShortUrl(id string)(shortUrl string) {
	shortUrl = base64.StdEncoding.EncodeToString([]byte(id))
	fmt.Println(shortUrl)
	return	
}

func getMysqlID(shortUrl string)(mysqlInsertId int) {
	temp := strings.Split(shortUrl, "/")
	lastIndex := len(temp) -1
	fmt.Println(temp[lastIndex]);
	id, err := base64.StdEncoding.DecodeString(temp[lastIndex])
	checkErr(err)
	mysqlInsertId, err = strconv.Atoi(string(id))
	checkErr(err)
	fmt.Println(mysqlInsertId)
	return
}

func mysqlGetLongUrl(id int)(longUrl string){
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/urlConverter?charset=utf8")
	checkErr(err)
	rows, err := db.Query("SELECT LONG_URL FROM URL_LIST WHERE ID = ?", id)
	checkErr(err)
	rows.Next()
	var LONG_URL string
	err = rows.Scan(&LONG_URL)
	fmt.Println(LONG_URL)
	longUrl = LONG_URL
	return
}

func handleShortRequest(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t longUrlStruct
	err := decoder.Decode(&t)
	checkErr(err)
	fmt.Println(t.Url)
	if len(t.Url)>0{
		mysqlInsertId := mysqlInsert(t.Url)
		shortUrl := sendShortUrl(mysqlInsertId)
		fmt.Println(shortUrl)
		temp := []string{"http://localhost/", shortUrl}
		var responseUrl = strings.Join(temp, "")
		fmt.Println(responseUrl)
		response := shortUrlStruct{Short: responseUrl}
		json.NewEncoder(rw).Encode(response)
	} else {
		fmt.Println("Empty")
	}

}

func handleOriginalRequest(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t shortUrlStruct
	err := decoder.Decode(&t)
	checkErr(err)
	fmt.Println(t.Short)
	if len(t.Short)>0{
		mysqlInsertId := getMysqlID(t.Short)
		longUrl := mysqlGetLongUrl(mysqlInsertId)
		fmt.Println(longUrl)
		response := originalUrlStruct{Original: longUrl}
		json.NewEncoder(rw).Encode(response)
	} else {
	 fmt.Println("Empty")
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
