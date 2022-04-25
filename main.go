package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/claygod/coffer"
)

var tpl = template.Must(template.ParseFiles("index.html"))

const curDir = "save/"

func main() {
	message := "Hello World!"
	message1 := "1"
	message2 := "23"

	arrayGo := [5]int{1, 5, 75, 33, 5}

	fmt.Println(message)
	fmt.Println(message1)
	fmt.Println(message2)
	fmt.Println(arrayGo)
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

// Обработчик главной странице.
func home(w http.ResponseWriter, r *http.Request) {
	// STEP init
	db, err, wrn := coffer.Db(curDir).Create()
	switch {
	case err != nil:
		fmt.Println("Error:", err)
		return
	case wrn != nil:
		fmt.Println("Warning:", err)
		return
	}
	if !db.Start() {
		fmt.Println("Error: not start")
		return
	}
	defer db.Stop()

	keys := db.RecordsList()
	rep := db.ReadList(keys.Data)
	rep.IsCodeError()
	if rep.IsCodeError() {
		fmt.Sprintf("Read error: code `%v` msg `%v`", rep.Code, rep.Error)
		return
	}

	tpl.Execute(w, rep.Data)
}

// Создание задачи
func create(w http.ResponseWriter, r *http.Request) {
	task := r.URL.Query().Get("task")
	// STEP init
	db, err, wrn := coffer.Db(curDir).Create()
	switch {
	case err != nil:
		fmt.Println("Error:", err)
		return
	case wrn != nil:
		fmt.Println("Warning:", err)
		return
	}
	if !db.Start() {
		fmt.Println("Error: not start")
		return
	}
	defer db.Stop()

	// STEP write
	uuid := string(time.Now().Unix())

	if rep := db.Write(uuid, []byte(task)); rep.IsCodeError() {
		fmt.Sprintf("Write error: code `%d` msg `%s`", rep.Code, rep.Error)
		return
	}

	w.Write([]byte(uuid))
}

// Удаление задачи
func delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	w.Write([]byte("Создание задачи"))
	w.Write([]byte("\n"))
	w.Write([]byte(id))
}
