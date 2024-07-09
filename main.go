package main

import (
	"fmt"
	"log"
	"main/server"
	"main/task"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	storage, err := task.NewTaskStorage()
	if err != nil {
		log.Fatalf("Failed to create task storage: %v", err)
	}

	webDir := "./web"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(http.Dir(webDir)).ServeHTTP(w, r)
	})
	port := server.GetPort()

	http.HandleFunc("/api/nextdate", server.HandleNextDate)
	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			server.AddTaskHandler(storage)(w, r)
		case http.MethodGet:
			server.GetTaskHandler(storage)(w, r)
		case http.MethodPut:
			server.UpdateTaskHandler(storage)(w, r)
		case http.MethodDelete:
			server.DeleteTaskHandler(storage)(w, r)
		default:
			http.Error(w, "Неверный метод запроса", http.StatusBadRequest)
		}
	})
	http.HandleFunc("/api/tasks", server.GetTasksHandler(storage))
	http.HandleFunc("/api/task/done", server.MarkTasksAsDoneHandler(storage))

	fmt.Printf("Запуск сервера на порту %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
