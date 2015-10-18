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
    //"github.com/gorilla/mux"
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

func mysqlInsert(Url string)(mySqlInsertId string) {
db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/urlConverter?charset=utf8")
checkErr(err)
stmt, err := db.Prepare("INSERT URL_LIST SET LONG_URL=?")
checkErr(err)
res, err := stmt.Exec(Url)
checkErr(err)
 id, err := res.LastInsertId()
checkErr(err)
    fmt.Println(id)

mySqlInsertId = strconv.Itoa(int(id))
return
}

func sendShortUrl(id string)(shortUrl string) {
shortUrl = base64.StdEncoding.EncodeToString([]byte(id))
fmt.Println(shortUrl)

return	
}

func getMysqlID(shortUrl string)(mySqlInsertId int) {
temp := strings.Split(shortUrl, "/");
fmt.Println(temp);
var id []byte
err := base64.StdEncoding.Decode(id, []byte(shortUrl))
return
}

func mysqlGetLongUrl(id int)(longUrl string){
longUrl = "wkdnfjasdhchjd"
return
}

func handleShortRequest(rw http.ResponseWriter, req *http.Request) {
    decoder := json.NewDecoder(req.Body)

	var t longUrlStruct
	err := decoder.Decode(&t)

	checkErr(err)
	fmt.Println(t.Url)
	if len(t.Url)>0{
	mySqlInsertID := mysqlInsert(t.Url)
	shortUrl := sendShortUrl(mySqlInsertID)
	
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
	mySqlInsertID := getMysqlID(t.Short)
	longUrl := mysqlGetLongUrl(mySqlInsertID)
	fmt.Println(longUrl)
	//temp := []string{"http://localhost/", shortUrl}

	//var responseUrl = strings.Join(temp, "")
		//fmt.Println(responseUrl)
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
