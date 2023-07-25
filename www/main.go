package main

import (
	"database/sql"
	_ "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" //go get -u github.com/go-sql-driver/mysql
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux" //go get -u github.com/gorilla/mux
	"html/template"
	"net/http"
	"path/filepath"
)

var save_from_data []string

type Article struct {
	Id                     uint16
	Title, Anons, FullText string
}

var posts = []Article{}
var showPosts = Article{}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		filepath.Join("templates", "index.html"),
		filepath.Join("templates", "header.html"),
		filepath.Join("templates", "footer.html"))

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	//Подключение к БД
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//Выборка данных из БД
	res, err := db.Query("SELECT * FROM articles")
	if err != nil {
		panic(err)
	}

	posts = []Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)

	}

	t.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		filepath.Join("templates", "create.html"),
		filepath.Join("templates", "header.html"),
		filepath.Join("templates", "footer.html"))

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request) {
	var title_var string
	var anons_var string
	var full_text_var string

	//fmt.Println(save_from_data)

	if len(save_from_data) == 0 {
		title_var = r.FormValue("title")
		anons_var = r.FormValue("anons")
		full_text_var = r.FormValue("full_text")
	} else {
		title_var = save_from_data[0]
		anons_var = save_from_data[1]
		full_text_var = save_from_data[2]
	}
	if title_var == "" || anons_var == "" || full_text_var == "" {
		save_from_data = nil
		//fmt.Println(title_var)
		save_from_data = append(save_from_data, title_var)
		save_from_data = append(save_from_data, anons_var)
		save_from_data = append(save_from_data, full_text_var)
		http.Redirect(w, r, "/not_filled_data", http.StatusSeeOther)

	} else {
		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")
		if err != nil {
			panic(err)
		}
		defer db.Close()
		insert, err := db.Query(fmt.Sprintf("INSERT INTO articles (title,anons,full_text) VALUES ('%s','%s','%s')", title_var, anons_var, full_text_var))
		if err != nil {
			panic(err)
		}
		defer insert.Close()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func not_filled_data(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		filepath.Join("templates", "not_filled_data.html"),
		filepath.Join("templates", "header.html"),
		filepath.Join("templates", "footer.html"))

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "not_filled_data", nil)
}

func show_post(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		filepath.Join("templates", "show.html"),
		filepath.Join("templates", "header.html"),
		filepath.Join("templates", "footer.html"))

	vars := mux.Vars(r)

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query(fmt.Sprintf("SELECT * FROM articles WHERE id = '%s'", vars["post_id"]))
	if err != nil {
		panic(err)
	}

	showPosts = Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}

		showPosts = post
	}

	t.ExecuteTemplate(w, "show_block", showPosts)

}

func handleFunc() {
	//Enabling address processing via gorilla mux
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET", "POST")
	rtr.HandleFunc("/not_filled_data", not_filled_data).Methods("GET")
	rtr.HandleFunc("/save_article", save_article).Methods("POST")

	//the ending can only accept numbers from 0-9 and there can be several of them
	rtr.HandleFunc("/post/{post_id:[0-9]+}", show_post).Methods("GET", "POST")

	http.Handle("/", rtr)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":8080", nil)
}

func main() {
	handleFunc()
}
