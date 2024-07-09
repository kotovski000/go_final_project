package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"main/task"
	"net/http"
	"time"
)

func HandleNextDate(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	now, err := time.Parse(task.DateFormat, nowStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка: Неверный формат даты 'now': %s", err), http.StatusBadRequest)
		return
	}

	nextDate, err := task.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка: %s", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", nextDate)
}

func AddTaskHandler(s task.TaskStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t task.Task
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			writeError(w, "Ошибка десериализации JSON: ", http.StatusBadRequest)
			return
		}

		if t.Title == "" {
			writeError(w, "Не указан заголовок задачи", http.StatusBadRequest)
			return
		}

		if t.Date == "" {
			t.Date = time.Now().Format(task.DateFormat)
		} else {
			_, err = time.Parse(task.DateFormat, t.Date)
			if err != nil {
				writeError(w, "Неверный формат даты, ожидается формат "+task.DateFormat, http.StatusBadRequest)
				return
			}
		}

		if t.Repeat == "" {
			if t.Date < time.Now().Format(task.DateFormat) {
				t.Date = time.Now().Format(task.DateFormat)
			}
		} else {
			nextDate, err := task.NextDate(time.Now(), t.Date, t.Repeat)
			if err != nil {
				writeError(w, "Неверный формат правила повторения: "+err.Error(), http.StatusBadRequest)
				return
			}
			if t.Date < time.Now().Format(task.DateFormat) {
				t.Date = nextDate
			}
		}

		id, err := task.AddTask(s, t)
		if err != nil {
			writeError(w, "Ошибка добавления задачи: "+err.Error(), http.StatusInternalServerError)
			return
		}
		writeResponse(w, map[string]int64{"id": id}, http.StatusCreated)
	}
}

func GetTaskHandler(s task.TaskStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			writeError(w, "Не указан ID задачи", http.StatusBadRequest)
			return
		}
		task, err := task.GetTask(s, id)
		if err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeResponse(w, task, http.StatusOK)
	}
}
func UpdateTaskHandler(s task.TaskStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var updatedTask task.Task
		err := json.NewDecoder(r.Body).Decode(&updatedTask)
		if err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}

		updatedTask, err = updatedTask.Check()
		if err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = task.PutTask(s, updatedTask)
		if err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		writeResponse(w, map[string]string{}, http.StatusCreated)
	}
}
func GetTasksHandler(s task.TaskStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("search")
		tasksList, err := task.GetTasks(s, search)
		if err != nil {
			writeError(w, "Ошибка поиска задач: ", http.StatusInternalServerError)
			return
		}

		writeResponse(w, map[string][]task.Task{"tasks": tasksList}, http.StatusOK)
	}
}
func DeleteTaskHandler(s task.TaskStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			writeError(w, "Не указан ID задачи", http.StatusBadRequest)
			return
		}
		err := task.DeleteTask(s, id)
		if err != nil {
			writeError(w, "Ошибка удаления задачи", http.StatusInternalServerError)
			return
		}

		writeResponse(w, struct{}{}, http.StatusOK)
	}
}

func MarkTasksAsDoneHandler(s task.TaskStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			writeError(w, "Не указан ID задачи", http.StatusBadRequest)
			return
		}

		t, err := task.GetTask(s, id)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, "Задача не найдена", http.StatusNotFound)
			} else {
				writeError(w, "Ошибка получения задачи", http.StatusInternalServerError)
			}
			return
		}

		if t.Repeat != "" {
			nextDate, err := task.NextDate(time.Now(), t.Date, t.Repeat)
			if err != nil {
				writeError(w, "Ошибка расчета следующей даты", http.StatusInternalServerError)
				return
			}
			t.Date = nextDate
			err = task.PutTask(s, t)
			if err != nil {
				writeError(w, "Ошибка обновления задачи", http.StatusInternalServerError)
				return
			}
		} else {
			err = task.DeleteTask(s, id)
			if err != nil {
				writeError(w, "Ошибка удаления задачи", http.StatusInternalServerError)
				return
			}
		}

		writeResponse(w, struct{}{}, http.StatusOK)
	}
}
func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func writeResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
