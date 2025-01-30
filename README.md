# Todo App
A simple yet feature-rich Todo application built with Go,
utilizing MongoDB for database storage, JWT for authentication, and bcrypt for password hashing.
The project also includes a docker-compose setup to easily deploy the application along with Nginx for frontend support
and Mongo-Express for database management.

It is a complete rewrite for similar project of [todo-cpp](https://github.com/aufam/todo)

## Features
- User authentication with JWT
- Password hashing using bcrypt
- CRUD operations for todos
- Deployment using Docker Compose
- Reverse proxy with Nginx
- Mongo-Express for database management

## Tech Stack
- **Backend:** Go (Fiber)
- **Database:** MongoDB
- **Authentication:** JWT (JSON Web Token)
- **Password Hashing:** bcrypt
- **Containerization:** Docker, Docker Compose
- **Reverse Proxy:** Nginx
- **Database UI:** Mongo-Express

## Getting Started
### Prerequisites
Make sure you have the following installed:
- [Go](https://go.dev/dl/)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Clone the Repository
```bash
git clone https://github.com/aufam/todo-go.git
cd todo-go
```

### Run Locally
1. Copy the example `.env` file and modify as needed:
   ```bash
   cp template.env .env
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Start the application:
   ```bash
   go run main.go
   ```

### Run with Docker Compose
To deploy the entire stack using Docker Compose, run:
```bash
docker-compose up -d
```
This will start the following services:
- **Go API** (Backend)
- **MongoDB** (Database)
- **Mongo-Express** (DB Admin UI)
- **Nginx** (Reverse Proxy)

### API Endpoints
#### Authentication
- **POST** `/api/v1/user/signup` – Register a new user
- **POST** `/api/v1/user/login` – Login and receive a JWT token

#### Todo Management (Requires JWT)
- **GET** `/api/v1/todos` – Get all todos
- **POST** `/api/v1/todo` – Create a new todo
- **PUT** `/api/v1/todo/:id` – Update a todo
- **DELETE** `/api/v1/todo/:id` – Delete a todo

### Accessing the Services
- API: http://localhost:8000/api/v1
- Mongo-Express: http://localhost:8081
- Nginx (Frontend): http://localhost

### Stopping the Application
```bash
docker-compose down
```

## License
This project is open-source and available under the [MIT License](LICENSE).

