package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"database/sql"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// 404 not found
func notFoundHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, "Not Found")
}

//article handler
func articleHandler(w http.ResponseWriter, req *http.Request) {
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

func db() {

	//init db
	db, err := sql.Open("postgres", "user=tim password=a dbname=test sslmode=disable")
	defer db.Close()
	panic(err)

	//insert
	stmt, err := db.Prepare("INSERT INTO my(use,pw) VALUES($1,$2);")
	// panicErr(err)
	stmt.Exec("apple", "iphone")
	// panicErr(err)
	fmt.Println("insert data")

	//query
	rows, err := db.Query("SELECT * FROM my")
	// panicErr(err)
	var use, pw string
	for rows.Next() {
		// panicErr(err)
		fmt.Println("\t", use, pw)
	}

}

func main() {

	// db()

	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	//get path param
	r.HandleFunc("/article/{name}", articleHandler)

	r.HandleFunc("/json", jsonHandler).Methods(http.MethodPost)

	http.ListenAndServe(":8080", middleWare(r))
}
