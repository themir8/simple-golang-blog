package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/jmoiron/sqlx"

	"github.com/gorilla/mux"
)

var DB *gorm.DB

type Article struct {
	gorm.Model
	Id       int    `gorm:"AUTO_INCREMENT" form:"id" json:"id"`
	Title    string `gorm:"not null" form:"title" json:"title"`
	Anons    string `gorm:"not null" form:"anons" json:"anons"`
	FullText string `gorm:"not null" form:"full_text" json:"full_text"`
}

// var posts = []Article{}
var showPost = Article{}

func GetAllPosts(article *[]Article) (err error) {
	if err = DB.Find(article).Error; err != nil {
		return err
	}
	return nil
}

func index(w http.ResponseWriter, r *http.Request) {
	// log.Println(r.URL)
	t, err := template.ParseFiles(
		"templates/index.html",
		"templates/header.html",
		"templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, _ := gorm.Open("sqlite3", "./gorm.db")
	defer db.Close()

	var articles = []Article{}
	db.Find(&articles)
	for _, u := range articles {
		fmt.Println(u)
	}

	t.ExecuteTemplate(w, "index", articles)
}

func create(w http.ResponseWriter, r *http.Request) {
	// URI := r.RequestURI
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

		db, _ := gorm.Open("sqlite3", "./gorm.db")
		defer db.Close()

		el := Article{Title: title, Anons: anons, FullText: full_text}

		db.Create(&el)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

func post_detail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// fmt.Println(vars["id"])
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

func post_delete(w http.ResponseWriter, r *http.Request) {
	// var err error
	db, _ := gorm.Open("sqlite3", "./gorm.db")
	defer db.Close()

	vars := mux.Vars(r)

	db.Delete(&Article{}, vars["id"])
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleFunc() {

	r := mux.NewRouter()
	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/create/", create).Methods("GET")
	r.HandleFunc("/save_article/", save_article).Methods("POST")
	r.HandleFunc("/post/{id:[0-9]+}/", post_detail).Methods("GET")
	r.HandleFunc("/delete/{id:[0-9]+}/", post_delete).Methods("GET")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.Handle("/", r)

	log.Println("Server running")
	http.ListenAndServe(":8000", nil)
}

func main() {
	db, _ := gorm.Open("sqlite3", "./gorm.db")
	defer db.Close()

	// db.AutoMigrate(&Article{})
	handleFunc()
}
