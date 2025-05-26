# Todo List Application

This is a simple Todo List application built with Go and SQLite. It allows users to add tasks and list all tasks stored in a SQLite database.

## Features

- Add a new task
- List all tasks

## Setup

1. **Install Go**: Make sure you have Go installed on your machine. You can download it from [golang.org](https://golang.org/dl/).

2. **Clone the repository**:

   ```bash
   git clone <repository-url>
   cd todo
   ```

3. **Install dependencies**:

   ```bash
   go get -u github.com/mattn/go-sqlite3
   ```

4. **Run the application**:

   ```bash
   go run main.go
   ```

5. **Access the application**:
   - Open your browser and go to `http://localhost:8080` to list tasks.
   - Use a tool like Postman to send a POST request to `http://localhost:8080/addtask` with a JSON body to add a task.

## Environment Variables

Create a `.env` file in the root of your project with the following keys. Do not include the actual secrets, just use example values:

```
HOST=example_host
CLIENT_ID=example_client_id
CLIENT_SECRET=example_client_secret
CALLBACK=example_callback_url
```

## API Endpoints

- `GET /`: List all tasks.
- `POST /addtask`: Add a new task. Requires a JSON body with a `task` field.

## Database

The application uses SQLite to store tasks. The database file is named `todo.db` and is created in the project directory.

## License

This project is licensed under the MIT License.
