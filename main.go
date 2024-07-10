package main

import (
	"fmt"
	"log"
	"net/http"

	"main/server"
	"main/task"

	"github.com/go-chi/chi/v5"
)

func main() {
	storage, err := task.NewTaskStorage()
	if err != nil {
		log.Fatalf("Failed to create task storage: %v", err)
	}

	r := chi.NewRouter()

	webDir := "./web"
	fs := http.FileServer(http.Dir(webDir))
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})

	r.Route("/api", func(r chi.Router) {
		r.Get("/nextdate", server.HandleNextDate)
		r.Route("/task", func(r chi.Router) {
			r.Post("/", server.AddTaskHandler(storage))
			r.Get("/", server.GetTaskHandler(storage))
			r.Put("/", server.UpdateTaskHandler(storage))
			r.Delete("/", server.DeleteTaskHandler(storage))
		})
		r.Get("/tasks", server.GetTasksHandler(storage))
		r.Post("/task/done", server.MarkTasksAsDoneHandler(storage))
	})

	port := server.GetPort()
	fmt.Printf("Запуск сервера на порту %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
