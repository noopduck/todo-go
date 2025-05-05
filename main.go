package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

type InputData struct {
	Task string `json:"task"`
}

type Db struct {
	DB *sql.DB
}

// Setup db connection
func setupDb() *Db {

	//os.Remove("./todo.db")

	var err error

	db, err := sql.Open("sqlite3", "./todo.db")
	if err != nil {
		log.Fatal(err)
	}

	app := &Db{DB: db}

	return app
	//stmt := `
	//create table todo (id integer not null primary key AUTOINCREMENT, task text);
	//`

	//_, err = db.Exec(stmt)
	//if err != nil {
	//	log.Printf("%q: %s\n", err, stmt)
	//	return
	//}
}

// Adding the task into the global list
func (app *Db) addTask(task string) {

	stmt := `
  insert into todo (task) values(?)
  `
	_, err := app.DB.Exec(stmt, task)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
	}
}

// Takes POST requests from web for adding a task
func (connection *Db) addTaskHandler(writer http.ResponseWriter, request *http.Request) {
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

func (app *Db) listTasks(writer http.ResponseWriter, request *http.Request) {

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

func returnCORS(writer http.ResponseWriter) http.ResponseWriter {

	writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	return writer
}

func main() {

	app := setupDb() // Start the database file connection

	http.HandleFunc("/", app.listTasks)
	http.HandleFunc("/addtask", app.addTaskHandler)
	http.HandleFunc("/listtasks", app.listTasks)

	var port = "8080"
	if http.ListenAndServe(":"+port, nil) != nil {
		log.Println("Prolly already listening on port:", port)
	}
}
