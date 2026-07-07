package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var todos = []Todo{
	{ID: 1, Title: "Learn Go", Done: false},
	{ID: 2, Title: "Build a REST API", Done: false},
	{ID: 3, Title: "Deploy to production", Done: false},
}

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func main() {
	http.HandleFunc("/health", health)
	http.HandleFunc("/todos", todosHandler)
	http.HandleFunc("PUT /todos/{id}", toggleTodo)
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Ok")
}

func todosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var newTodo Todo
		err := json.NewDecoder(r.Body).Decode(&newTodo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newTodo.ID = len(todos) + 1
		todos = append(todos, newTodo)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTodo)

	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todos)
	}
}

func toggleTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	for i, todo := range todos {
		if todo.ID == id {
			todos[i].Done = !todos[i].Done
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(todos[i])
			return
		}
	}
	http.Error(w, "Todo not found", http.StatusNotFound)
}
