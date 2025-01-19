# Go-Auth-App

Go-Auth-App is a modern authentication and authorization system built using **Golang**, **Gin-Gonic**, **PostgreSQL**, and **JWT**. It provides a robust, secure, and efficient way to manage user authentication and secure API access.

---

## 🚀 Features

- **User Authentication**: Sign up, log in, and log out seamlessly.
- **JWT-Based Authorization**: Secure your API routes with JSON Web Tokens.
- **Password Reset**: Simple flow for resetting forgotten passwords.
- **Session Management**: Handle user sessions with security and efficiency.
- **Scalable Design**: Built for scalability and performance.

---

## 📋 Prerequisites

Before running the project, ensure you have the following installed on your system:

- **Go**: Version 1.15 or higher ([Download Go](https://go.dev/dl/))
- **PostgreSQL**: ([Download PostgreSQL](https://www.postgresql.org/download/))
- **Gin-Gonic**: Go web framework ([Gin Documentation](https://github.com/gin-gonic/gin))

---

## 📦 Installation

Follow these steps to set up the project on your local machine:

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/atulguptag/go-auth-app
   ```

2. **Navigate to the Project Directory**:

   ```bash
   cd go-auth-app
   ```

3. **Download Dependencies**:

   ```bash
   go mod tidy
   ```

4. **Configure PostgreSQL**:

   - Create a PostgreSQL database for the project.
   - Update the database credentials in the `config` file.

5. **Run the Project**:

   ```bash
   go run main.go
   ```

---

## 🔧 Configuration

To configure the project, update the database and JWT settings in the `config` directory. Ensure the following environment variables are set:

- `DB_HOST`: Database host
- `DB_PORT`: Database port
- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `JWT_SECRET`: Secret key for JWT

---

## 🛠 API Endpoints

| HTTP Method | Endpoint          | Description                    |
| ----------- | ----------------- | ------------------------------ |
| `POST`      | `/signup`         | Create a new user account      |
| `POST`      | `/login`          | Log in to an existing account  |
| `GET`       | `/logout`         | Log out of the current session |
| `GET`       | `/home`           | Access the home page           |
| `POST`      | `/reset-password` | Reset your password            |
| `POST`      | `/generate-jokes` | Generate AI-Powered Jokes      |

---

## 🖥️ Running in Development Mode

To run the project in development mode:

1. Export the required environment variables:

   ```bash
   export GIN_MODE=debug
   ```

2. Start the server:
   ```bash
   go run main.go
   ```

Access the application at [http://localhost:8080](http://localhost:8080).

---

## 🤝 Contributing

We welcome contributions to enhance the project! Here’s how you can help:

1. **Report Bugs**: Open an issue to let us know.
2. **Request Features**: Suggest new features via an issue.
3. **Contribute Code**: Submit a pull request with your improvements.

---

## 📜 License

Go-Auth-App is licensed under the [MIT License](LICENSE). You are free to use, modify, and distribute this project with proper attribution.

---

## 💬 Support

If you have questions or need help, feel free to reach out:

- **Email**: [atulguptag111@gmail.com](mailto:atulguptag111@gmail.com)
- **GitHub Issues**: [Create an Issue](https://github.com/atulguptag/go-auth-app/issues)

---

## 🌟 Acknowledgments

Special thanks to:

- The [Gin-Gonic](https://github.com/gin-gonic/gin) team for the excellent web framework.
- The PostgreSQL community for their reliable database solutions.
- The Go community for their robust programming language.
