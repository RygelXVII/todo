package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var tpl = template.Must(template.ParseFiles("index.html"))

func main() {
	fmt.Println("Hello World!")

	startServer()

}

func startServer() {
	//Настройка сервера
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/create", create)
	mux.HandleFunc("/delete", delete)
	fs := http.FileServer(http.Dir("src"))
	mux.Handle("/src/", http.StripPrefix("/src/", fs))

	log.Println("Запуск веб-сервера на http://127.0.0.1:4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

type task struct {
	id      int
	text    string
	checked bool
}

type templateData struct {
	Tasks []task
}

// Обработчик главной странице.
func home(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "task.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("select * from tasks where checked = false")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	tasks := []task{}

	for rows.Next() {
		p := task{}
		err := rows.Scan(&p.id, &p.text, &p.checked)
		if err != nil {
			fmt.Println(err)
			continue
		}
		tasks = append(tasks, p)
	}

	data := &templateData{Tasks: tasks}

	tpl.Execute(w, data)
}

// Создание задачи
func create(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	task_text := r.Form.Get("task")
	if len(task_text) < 1 {
		panic("error")
	}

	db, err := sql.Open("sqlite3", "task.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	result, err := db.Exec("insert into tasks (text, checked) values ( $1, false)", task_text)
	if err != nil {
		panic(err)
	}

	id, err := result.LastInsertId()
	fmt.Println(id)
	w.Write([]byte(strconv.FormatInt(id, 10)))
}

// Удаление задачи
func delete(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.Form.Get("id")

	db, err := sql.Open("sqlite3", "task.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// удаляем строку с id=1
	result, err := db.Exec("delete from tasks where id = $1", id)
	if err != nil {
		panic(err)
	}

	count, err := result.RowsAffected() // количество удаленных строк
	fmt.Println(count)
	w.Write([]byte(strconv.FormatInt(count, 10)))

}
