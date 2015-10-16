package main

import (
    "encoding/json"
    "encoding/base64"
    "net/http"
    "fmt"
    "strconv"
    "database/sql"
    _"github.com/go-sql-driver/mysql"
)

type long_url_struct struct {
	Url string
}

func mysqlConnect(Url string) {
db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/urlConverter?charset=utf8")
checkErr(err)
stmt, err := db.Prepare("INSERT URL_LIST SET LONG_URL=?")
checkErr(err)
res, err := stmt.Exec(Url)
checkErr(err)
 id, err := res.LastInsertId()
checkErr(err)
    fmt.Println(id)

sendShortUrl(strconv.Itoa(int(id)));
}

func sendShortUrl(id string) {
shortUrl := base64.StdEncoding.EncodeToString([]byte(id))
fmt.Println(shortUrl)

data, err := base64.StdEncoding.DecodeString(shortUrl)
checkErr(err)
fmt.Printf("%q\n", data)
//fmt.Fprintln(w, "short:", data)	
}



func handlePostRequest(rw http.ResponseWriter, req *http.Request) {
    decoder := json.NewDecoder(req.Body)

	var t long_url_struct
	err := decoder.Decode(&t)

	checkErr(err)
	fmt.Println(t.Url)
	if len(t.Url)>0{
	mysqlConnect(t.Url)
	} else {
	 fmt.Println("Empty")
	}


    //LOG: that
}

func handleGetRequest(rw http.ResponseWriter, req *http.Request) {
    decoder := json.NewDecoder(req.Body)

	var t long_url_struct
	err := decoder.Decode(&t)

	checkErr(err)
	fmt.Println(t.Url)
	if len(t.Url)>0{
	mysqlConnect(t.Url)
	} else {
	 fmt.Println("Empty")
	}


    //LOG: that
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}
