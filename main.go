package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"todo-go/internal/auth"
)

type InputData struct {
	Task string `json:"task"`
}

type App struct {
	DB   *sql.DB
	Auth *auth.Auth
}

// Setup db connection
func setupDb() *sql.DB {
	//os.Remove("./todo.db")
	var err error

	db, err := sql.Open("sqlite3", "./todo.db")
	if err != nil {
		log.Fatal(err)
	}

	// TODO: this stuff should be tested for create if exists..
	stmt := `
	create table if not exists todo (id integer not null primary key AUTOINCREMENT, task text);
	`

	_, err = db.Exec(stmt)

	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return nil
	}

	return db
}

// Adding the task into the global list
func (app *App) addTask(task string) {

	stmt := `
  insert into todo (task) values(?)
  `
	_, err := app.DB.Exec(stmt, task)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
	}
}

// Takes POST requests from web for adding a task
func (connection *App) addTaskHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Only POST supported here", http.StatusMethodNotAllowed)
		return
	}

	var data InputData
	if err := json.NewDecoder(request.Body).Decode(&data); err != nil {
		http.Error(writer, "Bad JSON form", http.StatusBadRequest)
		return
	}

	fmt.Println(data.Task)
	connection.addTask(data.Task) // Write task into db

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Thank you, come again!"))

	fmt.Println(request.Header)
}

func (app *App) listTasks(writer http.ResponseWriter, request *http.Request) {

	// Handle CORS
	writer = returnCORS(writer)

	fmt.Println("#### Welcome to our Todolist app! ####")

	stmt := `
  select id, task from todo;
  `
	rows, err := app.DB.Query(stmt)

	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
	}

	for rows.Next() {
		var id int
		var task string
		if err := rows.Scan(&id, &task); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		fmt.Println("Task:", id, task)
		fmt.Fprintln(writer, "Task:", id, task)
	}
}

func (app *App) login(writer http.ResponseWriter, request *http.Request) {
	// Generer evt. random state senere
	state := "static-state"
	authURL := app.Auth.Config.AuthCodeURL(state)

	http.Redirect(writer, request, authURL, http.StatusFound) // Generer evt. random state senere
	fmt.Println(request.Body)

}

func (app *App) callback(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	code := request.URL.Query().Get("code")
	if code == "" {
		http.Error(writer, "Missing code", http.StatusBadRequest)
		return
	}

	// Exchange code for tokens
	token, err := app.Auth.Config.Exchange(ctx, code)
	if err != nil {
		http.Error(writer, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract ID token
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(writer, "No id_token field in token", http.StatusInternalServerError)
		return
	}

	// Verify ID token
	idToken, err := app.Auth.IDTokenVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		http.Error(writer, "Invalid ID token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Decode claims
	var claims struct {
		Email string `json:"email"`
	}
	if err := idToken.Claims(&claims); err != nil {
		http.Error(writer, "Failed to parse claims: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Skriv ut til bruker
	fmt.Fprintf(writer, "You are logged in as: %s", claims.Email)
}

func (app *App) logout(writer http.ResponseWriter, request *http.Request) {
	// Clear the session or token information to log out the user
	http.SetCookie(writer, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Delete the cookie
		HttpOnly: true,
	})

	// Redirect the user to the homepage or login page after logging out
	http.Redirect(writer, request, "/login", http.StatusSeeOther)
}

func returnCORS(writer http.ResponseWriter) http.ResponseWriter {

	writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	return writer
}

func main() {
	authenticator := auth.NewAuth()

	db := setupDb() // Start the database file connection
	app := &App{
		Auth: authenticator,
		DB:   db,
	}

	http.HandleFunc("/", app.listTasks)
	http.HandleFunc("/addtask", app.addTaskHandler)
	http.HandleFunc("/listtasks", app.listTasks)
	http.HandleFunc("/login", app.login)
	http.HandleFunc("/callback", app.callback)
	http.HandleFunc("/logout", app.logout)

	var port = "8080"
	if http.ListenAndServe(":"+port, nil) != nil {
		log.Println("Prolly already listening on port:", port)
	}
}
