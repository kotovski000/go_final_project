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
		r.Post("/signin", server.SigninHandler)

		r.Route("/task", func(r chi.Router) {
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				server.Auth(server.AddTaskHandler(storage)).ServeHTTP(w, r)
			})
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				server.Auth(server.GetTaskHandler(storage)).ServeHTTP(w, r)
			})
			r.Put("/", func(w http.ResponseWriter, r *http.Request) {
				server.Auth(server.UpdateTaskHandler(storage)).ServeHTTP(w, r)
			})
			r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
				server.Auth(server.DeleteTaskHandler(storage)).ServeHTTP(w, r)
			})
		})

		r.Get("/tasks", func(w http.ResponseWriter, r *http.Request) {
			server.Auth(http.HandlerFunc(server.GetTasksHandler(storage))).ServeHTTP(w, r)
		})

		r.Post("/task/done", func(w http.ResponseWriter, r *http.Request) {
			server.Auth(http.HandlerFunc(server.MarkTasksAsDoneHandler(storage))).ServeHTTP(w, r)
		})
	})

	port := server.GetPort()
	fmt.Printf("Запуск сервера на порту %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
