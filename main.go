package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

type Article struct {
	Id                     uint16
	Title, Anons, FullText string
}

var posts = []Article{}
var showPost = Article{}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"templates/index.html",
		"templates/header.html",
		"templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	posts = []Article{}

	// DATABASE_URL := "postgres://mirzohidov:coder@localhost:6432/mir_db?sslmode=disable"

	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err, " Bla Bla1")
	}
	defer db.Close()

	rows, _ := db.Queryx("SELECT * FROM simpleGolangBlog")
	for rows.Next() {
		var post Article
		err := rows.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			log.Println(err, " Bla Bla2")
		}
		posts = append(posts, post)
	}

	t.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"templates/create.html",
		"templates/header.html",
		"templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request) {
	// var err error
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	if title == "" || anons == "" || full_text == "" {
		fmt.Fprintf(w, "Hammasini toldir dal**yob")
	} else {
		db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			log.Fatalln(err, " Bla Bla1")
		}
		defer db.Close()

		tx := db.MustBegin()
		tx.MustExec("INSERT INTO simpleGolangBlog(title, anons, full_text) values($1, $2, $3)", title, anons, full_text)
		tx.Commit()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

func post_detail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t, err := template.ParseFiles(
		"templates/post.html",
		"templates/header.html",
		"templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err, " Bla Bla1")
	}
	defer db.Close()

	rows, _ := db.Queryx(fmt.Sprintf("SELECT * FROM simpleGolangBlog WHERE id=%s", vars["id"]))

	showPost = Article{}
	for rows.Next() {
		var post Article
		err := rows.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			log.Println(err, " Bla Bla2")
		}
		showPost = post
	}

	t.ExecuteTemplate(w, "post_detail", showPost)
}

func handleFunc() {

	r := mux.NewRouter()
	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/create/", create).Methods("GET")
	r.HandleFunc("/save_article/", save_article).Methods("POST")
	r.HandleFunc("/post/{id:[0-9]+}/", post_detail).Methods("GET")

	http.Handle("/", r)

	log.Println("Server running")
	http.ListenAndServe(":8000", nil)
}

func main() {
	handleFunc()
}
