# Zync Chat

Zync Chat is a real-time chat application that allows users to register, log in, and communicate with each other instantly. Built with Go for the backend and a simple web interface for the frontend, it provides a seamless chat experience.

## Features

- User registration and authentication
- Real-time messaging using WebSocket
- Secure password storage with bcrypt
- JWT for session management

## Technologies Used

- **Backend**: Go, SQLite, WebSocket
- **Frontend**: HTML, CSS, JavaScript
- **Libraries**:
  - [gorilla/websocket](https://github.com/gorilla/websocket)
  - [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)
  - [dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go)
  - [golang.org/x/crypto/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)

## Setup Instructions

To run the project locally, follow these steps:

1. **Clone the repository**:

   ```bash
   git clone https://github.com/MRQ67/zync-chat.git
   cd zync-chat
   ```

2. **Install dependencies**:

   ```bash
   go mod tidy
   ```

3. **Set up the database**:

   - The application uses SQLite, so no additional setup is required. The database file `zyncchat.db` will be created automatically in the project directory.

4. **Run the server**:

   ```bash
   go run main.go db.go
   ```

5. **Access the application**:

   - Open your browser and navigate to `http://localhost:8080`.

## Usage

- **Register**: Create a new account by providing a username and password.
- **Login**: Use your credentials to log in.
- **Chat**: Once logged in, you can send and receive messages in real-time.

## Project Structure

```
zync-chat/
├── main.go          # Handles HTTP requests, WebSocket connections, and authentication
├── db.go            # Manages database operations for user registration and authentication
└── static/          # Contains frontend files
    ├── index.html
    ├── style.css
    └── chat.js
```

## Building the Project

To build an executable, run:

```bash
go build -o zync-chat
```

Then, run the executable:

```bash
./zync-chat
```

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your changes.

## License

This project is licensed under the MIT License.