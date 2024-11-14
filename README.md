# Simple Blogger API

This is a Golang API for a simple blogging platform where users can manage posts, follow each other, like and unlike posts.

### Features

  * User registration and login with email and password (validation included);
  * JWT authentication for secure access;
  * Users can create, edit, and delete their own posts;
  * Users can follow other users;
  * Users can see posts from users they follow and their own ones;
  * Users can like and unlike posts;

### Dependencies

This project uses the following external libraries:

  * **[jwt-go](https://github.com/dgrijalva/jwt-go):**  for JSON Web Token (JWT) authentication;
  * **[godotenv](https://github.com/joho/godotenv):** for loading environment variables from a `.env` file;
  * **[checkmail](https://github.com/badoux/checkmail):** for email address validation;
  * **[Go Cryptography](https://golang.org/x/crypto):** for cryptographic functions;
  * **[MySQL Driver](https://github.com/go-sql-driver/mysql):**  for connecting and interacting with a MySQL database;

### Installation

1.  **Install Go:** Download and install the latest version of Go from [https://go.dev/dl/](https://go.dev/dl/).
2.  **Clone the Repository:**

```
git clone https://github.com/your-username/simple-blogger-api.git
```

3.  **Install Dependencies:**

```
go mod download
```

4.  **Configure Environment:**

  - Create a file named `.env` in the project root directory.
  - Add the following environment variables to `.env`:

```
DB_USER=your_username
DB_PASS=your_password
DB_NAME=blogger

API_PORT=8000

SECRET_KEY=your_secret_key

```

  - Replace the placeholders with your actual database credentials and a strong secret key for JWT.

> You can execute, before uncomment the `init()` function, the `main.go` file to have a suggested strong secret key.

### Running the API

```
go run main.go
```

This will start the API server on the default port (usually 8080).

### API Documentation

Detailed API documentation with endpoints, request/response formats, and authentication details is on the `blogger.postman_collection.json` file.

### About

This project was implemented in order to learn the Go language and develop APIs with it. 
It was largely based on the project developed in the course 
"[Aprenda Golang do Zero! Desenvolva uma APLICAÇÃO COMPLETA!](https://www.udemy.com/course/aprenda-golang-do-zero-desenvolva-uma-aplicacao-completa/)"
 by [Otávio Augusto Gallego](https://github.com/OtavioGallego) (he teach very well, I recommend him). 
 However, I feel the need to apply some improvements, I also used Go's native mux: [ServeMux](https://pkg.go.dev/net/http#ServeMux).

 The application is pretty simple, you can learn a lot of things though. For instance: authentication and security for it, middlewares, repository pattern,
 generic implementation of responses for the API.

### License

This project is licensed under the MIT License.
