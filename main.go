package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// 404 not found
func notFoundHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, "404 Not Found")
}

//article handler
func articleHandler(w http.ResponseWriter, req *http.Request) {

	db, err := sql.Open("mysql", "tim:102030@tcp(localhost:8889)/go")

	rows, err := db.Query("SELECT * FROM users")

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var autoID string
		var account string
		var pwd string
		var email string

		if err := rows.Scan(&autoID, &account, &pwd, &email); err != nil {
			panic(err)
		}

		io.WriteString(w, account)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

	vars := mux.Vars(req)
	//get url path param
	articleName := vars["name"]
	//get url query
	page := req.URL.Query().Get("p")
	//get request header
	io.WriteString(w, "header"+req.Header.Get("User-Agent")+"\n")
	io.WriteString(w, "article:"+articleName+page)

}

func jsonHandler(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	//parse json
	v := struct {
		Type string
		Name string
		Age  int
	}{}

	if err := json.Unmarshal(b, &v); err != nil {
		panic(err)
	}

	//respone
	res := struct {
		Name string
		Age  int
	}{
		v.Name,
		v.Age,
	}

	b2, _ := json.Marshal(res)
	w.Header().Set("CONTENT-TYPE", "application/json; charset=utf-8")
	w.Write(b2)

}

func handler(w http.ResponseWriter, req *http.Request) {
	w.Write(([]byte)("hello world"))
}

func middleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				w.WriteHeader(http.StatusInternalServerError)
				io.WriteString(w, fmt.Sprintf("internal error %v", r))
			}
		}()
		next.ServeHTTP(w, req)

	})
}

func sendToSlack(w http.ResponseWriter, req *http.Request) {

	deviceToken := req.URL.Query().Get("deviceToken")
	url := "https://hooks.slack.com/services/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"
	fmt.Println("URL:>", url)
	io.WriteString(w, deviceToken)
	s := `{"text":"` + deviceToken + `"}`
	var jsonStr = []byte(s)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}
	client.Do(req)

	if err != nil {
		panic(err)
	}

	/*resp, err := client.Do(www)
	    if err != nil {
	        panic(err)
		}
	    defer resp.Body.Close()
	*/
}

func main() {

	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	// //get path param
	r.HandleFunc("/article/{name}", articleHandler)
	r.HandleFunc("/sendToSlack", sendToSlack)
	r.HandleFunc("/json", jsonHandler).Methods(http.MethodPost)

	http.ListenAndServe(":8080", middleWare(r))
}
