package main

import (
	"database/sql"
	_ "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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

func handleFunc() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/create", create)
	http.HandleFunc("/not_filled_data", not_filled_data)
	http.HandleFunc("/save_article", save_article)

	http.ListenAndServe(":8080", nil)
}

func main() {
	handleFunc()
}
