package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Post struct {
	Id   int
	Name string
	Text string
}

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("mysql", "root:root@tcp(localhost:8889)/posts?charset=utf8")

	if err != nil {
		panic(err)
	}

	err = Db.Ping()

	if err != nil {
		fmt.Println("データベース接続失敗")
		return
	} else {
		fmt.Println("データベース接続成功")
	}
}

func HandlerIndex(w http.ResponseWriter, r *http.Request) {
	template := template.Must(template.ParseFiles("html/index.html"))
	selected, err := Db.Query("SELECT * FROM posts_table")

	if err != nil {
		panic(err.Error())
	}

	data := []Post{}

	for selected.Next() {
		article := Post{}
		err = selected.Scan(&article.Id, &article.Name, &article.Text)

		if err != nil {
			panic(err.Error())
		}

		data = append(data, article)
	}

	selected.Close()

	if err := template.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Fatal(err)
	}
}

func HandlerShow(w http.ResponseWriter, r *http.Request) {
	template := template.Must(template.ParseFiles("html/show.html"))
	id := r.URL.Query().Get("id")
	selected, err := Db.Query("SELECT * FROM posts_table WHERE id=?", id)

	if err != nil {
		panic(err.Error())
	}

	article := Post{}

	for selected.Next() {
		err = selected.Scan(&article.Id, &article.Name, &article.Text)
		if err != nil {
			panic(err.Error())
		}
	}

	selected.Close()

	if err := template.ExecuteTemplate(w, "show.html", article); err != nil {
		log.Fatal(err)
	}
}

func HandlerCreate(w http.ResponseWriter, r *http.Request) {
	template := template.Must(template.ParseFiles("html/create.html"))

	if r.Method == "GET" {
		template.ExecuteTemplate(w, "create.html", nil)
	} else if r.Method == "POST" {
		name := r.FormValue("name")
		text := r.FormValue("text")
		insert, err := Db.Prepare("INSERT INTO posts_table(name, text) VALUES(?,?)")

		if err != nil {
			panic(err.Error())
		}

		insert.Exec(name, text)
		http.Redirect(w, r, "/", 301)
	}
}

func HandlerEdit(w http.ResponseWriter, r *http.Request) {
	template := template.Must(template.ParseFiles("html/edit.html"))

	if r.Method == "GET" {
		id := r.URL.Query().Get("id")
		selected, err := Db.Query("SELECT * FROM posts_table WHERE id=?", id)

		if err != nil {
			panic(err.Error())
		}

		article := Post{}

		for selected.Next() {
			err = selected.Scan(&article.Id, &article.Name, &article.Text)

			if err != nil {
				panic(err.Error())
			}
		}

		selected.Close()

		if err := template.ExecuteTemplate(w, "edit.html", article); err != nil {
			log.Fatal(err)
		}

	} else if r.Method == "POST" {
		name := r.FormValue("name")
		text := r.FormValue("text")
		id := r.FormValue("id")
		insert, err := Db.Prepare("UPDATE posts_table SET name=?, text=? WHERE id=?")

		if err != nil {
			panic(err.Error())
		}

		insert.Exec(name, text, id)
		http.Redirect(w, r, "/", 301)
	}
}

func HandlerDelete(w http.ResponseWriter, r *http.Request) {
	template := template.Must(template.ParseFiles("html/delete.html"))

	if r.Method == "GET" {
		id := r.URL.Query().Get("id")
		selected, err := Db.Query("SELECT * FROM posts_table WHERE id=?", id)

		if err != nil {
			panic(err.Error())
		}

		article := Post{}

		for selected.Next() {
			err = selected.Scan(&article.Id, &article.Name, &article.Text)
			if err != nil {
				panic(err.Error())
			}
		}

		selected.Close()
		if err := template.ExecuteTemplate(w, "delete.html", article); err != nil {
			log.Fatal(err)
		}
	} else if r.Method == "POST" {
		id := r.FormValue("id")
		insert, err := Db.Prepare("DELETE FROM posts_table WHERE id=?")

		if err != nil {
			panic(err.Error())
		}

		insert.Exec(id)
		http.Redirect(w, r, "/", 301)
	}
}

func main() {

	http.HandleFunc("/", HandlerIndex)

	http.HandleFunc("/show", HandlerShow)

	http.HandleFunc("/create", HandlerCreate)

	http.HandleFunc("/edit", HandlerEdit)

	http.HandleFunc("/delete", HandlerDelete)

	http.ListenAndServe(":8080", nil)
}
