package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	City      string `json:"city"`
}

var users = map[int]User{
	1: {
		ID:        1,
		FirstName: "Даниял",
		LastName:  "Булатов",
		City:      "Москва",
	},
	2: {
		ID:        2,
		FirstName: "Иван",
		LastName:  "Иванов",
		City:      "Казань",
	},
}

// Ниже напишите обработчики для каждого эндпоинта
func getUsers(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user_id, _ := strconv.Atoi(id)
	val, ok := users[user_id]
	if !ok {
		http.Error(w, "Пользователь не найден", http.StatusNoContent)
		return
	}
	resp, err := json.Marshal(val)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func createTasks(w http.ResponseWriter, r *http.Request) {
	var user User
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = json.Unmarshal(buf.Bytes(), &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	users[user.ID] = user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/users", getUsers)
	r.Post("/users", createTasks)
	r.Get("/users/{id}", getTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}

}
