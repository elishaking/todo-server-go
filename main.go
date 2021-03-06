package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Todo Struct (Model)
type Todo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	User        *User  `json:"user"`
}

// User Struct (Model)
type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DateCreated string `json:"dateCreated"`
}

var todos []Todo

// Get all Todos
func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// Create a Todo
func createTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todo Todo
	_ = json.NewDecoder(r.Body).Decode(&todo)
	todo.ID = strconv.Itoa(len(todos) + 1)
	todos = append(todos, todo)
	json.NewEncoder(w).Encode(todo)
}

// Get a Todo by id
func getTodoByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, todo := range todos {
		if todo.ID == params["id"] {
			json.NewEncoder(w).Encode(todo)
			return
		}
	}

	fmt.Fprintln(w, "No todos match the id:", params["id"])
}

// Update a Todo by id
func updateTodoByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for idx, todo := range todos {
		if todo.ID == params["id"] {
			json.NewDecoder(r.Body).Decode(&todos[idx])
			todos[idx].ID = params["id"]
			fmt.Fprintln(w, "Updated todo with id:", params["id"])
			return
		}
	}

	fmt.Fprintln(w, "No todos match the id:", params["id"])
}

// Delete a Todo by id
func deleteTodoByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for idx, todo := range todos {
		if todo.ID == params["id"] {
			todos = append(todos[:idx], todos[idx+1:]...)
			fmt.Fprintln(w, "Deleted todo with id:", params["id"])
			return
		}
	}

	fmt.Fprintln(w, "No todos match the id:", params["id"])
}

func main() {
	fmt.Println("Todo Server")

	// Init Router
	router := mux.NewRouter()

	// Mock Data
	todos = append(todos, Todo{ID: "1", Title: "New Todo", Description: "This is a Todo", Date: "now", User: &User{ID: "1", Name: "King", DateCreated: "now"}})
	todos = append(todos, Todo{ID: "2", Title: "New Todo 2", Description: "This is a Todo 2", Date: "now", User: &User{ID: "2", Name: "King", DateCreated: "now"}})

	// Route Handlers / Endpoints
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Golang/Todo API")
	})

	router.HandleFunc("/todos", getTodos).Methods("GET")
	router.HandleFunc("/todos", createTodo).Methods("POST")

	router.HandleFunc("/todos/{id}", getTodoByID).Methods("GET")
	router.HandleFunc("/todos/{id}", updateTodoByID).Methods("PUT")
	router.HandleFunc("/todos/{id}", deleteTodoByID).Methods("DELETE")

	if err := http.ListenAndServe(":7000", router); err != nil {
		log.Fatal(err)
	}
}
