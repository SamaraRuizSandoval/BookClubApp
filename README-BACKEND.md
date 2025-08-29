# BookClubApp Backend - HTTP Server with Chi

This is the BookClubApp backend, implemented in Go using the Chi framework to create a robust and scalable HTTP server.

## 🚀 Features

- **HTTP Server** running on port 8080
- **RESTful API** with routes for Books, Chapters, Comments and Users
- **Middlewares** for logging, CORS, timeout and error recovery
- **Health Check** for monitoring
- **Organized structure** and modular
- **Complete functionalities** of BookClubApp (books, chapters, comments)

## 📋 Prerequisites

- Go 1.23.3 or higher
- Git

## 🛠️ Installation

1. **Clone the repository** (if you haven't already):
```bash
git clone https://github.com/SamaraRuizSandoval/BookClubApp.git
cd BookClubApp
```

2. **Install dependencies**:
```bash
go mod tidy
```

3. **Run the server**:
```bash
go run main.go
```

## 🌐 API Endpoints

### Root Route
- `GET /` - Server information

### Health Check
- `GET /health` - Server health status

### Books API
- `GET /api/books` - List all books
- `GET /api/books/{id}` - Get book by ID
- `POST /api/books` - Create new book
- `PUT /api/books/{id}` - Update existing book
- `DELETE /api/books/{id}` - Remove book

### Chapters API
- `GET /api/books/{bookId}/chapters` - List chapters of a book
- `GET /api/books/{bookId}/chapters/{chapterId}` - Get chapter by ID
- `POST /api/books/{bookId}/chapters` - Create new chapter
- `PUT /api/books/{bookId}/chapters/{chapterId}` - Update existing chapter
- `DELETE /api/books/{bookId}/chapters/{chapterId}` - Remove chapter

### Comments API
- `GET /api/comments` - List all comments
- `GET /api/comments/{id}` - Get comment by ID
- `POST /api/comments` - Create new comment
- `PUT /api/comments/{id}` - Update existing comment
- `DELETE /api/comments/{id}` - Remove comment
- `GET /api/comments/chapter/{chapterId}` - List comments of a chapter

### Users API
- `GET /api/users` - List all users
- `GET /api/users/{id}` - Get user by ID
- `POST /api/users` - Create new user
- `PUT /api/users/{id}` - Update existing user
- `DELETE /api/users/{id}` - Remove user

## 🏗️ Project Structure

```
BookClubApp/
├── main.go                 # Application entry point
├── go.mod                  # Go dependencies
├── internal/
│   └── app/
│       └── app.go         # Server logic and handlers
├── config/
│   └── config.go          # Configuration system
├── README-BACKEND.md       # Complete documentation
├── Makefile               # Development commands
├── Dockerfile             # Containerization
├── docker-compose.yml     # Orchestration
└── test-api.http          # API tests
```

## 🔧 Configuration

The server uses default configurations that can be customized:

```go
config := &app.Config{
    Port:         ":8080",
    ReadTimeout:  15 * time.Second,
    WriteTimeout: 15 * time.Second,
    IdleTimeout:  60 * time.Second,
}
```

## 📝 Usage Examples

### Testing with curl

```bash
# Health check
curl http://localhost:8080/health

# List books
curl http://localhost:8080/api/books

# Get specific book
curl http://localhost:8080/api/books/123

# List chapters of a book
curl http://localhost:8080/api/books/123/chapters

# List comments
curl http://localhost:8080/api/comments

# List users
curl http://localhost:8080/api/users
```

### Testing in browser

- **Server**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **API Books**: http://localhost:8080/api/books
- **API Chapters**: http://localhost:8080/api/books/123/chapters
- **API Comments**: http://localhost:8080/api/comments
- **API Users**: http://localhost:8080/api/users

## 🚦 Active Middlewares

- **Logger**: Logs all HTTP requests
- **Recoverer**: Recovers from panics and errors
- **RequestID**: Adds unique ID for each request
- **RealIP**: Gets real client IP
- **Timeout**: 60 second timeout for requests
- **CORS**: Configuration for Cross-Origin Resource Sharing

## 🔮 Next Steps

- [ ] Implement database
- [ ] Add JWT authentication
- [ ] Implement data validation
- [ ] Add unit tests
- [ ] Implement rate limiting
- [ ] Add Swagger documentation

## 📚 Technologies Used

- **Go 1.23.3** - Programming language
- **Chi v5** - Minimalist HTTP framework
- **Chi CORS** - CORS middleware
- **Chi Middleware** - Standard middlewares

## 🤝 Contributing

1. Fork the project
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is under MIT license. See the `LICENSE` file for more details.

## 🆘 Support

If you encounter any problems or have questions, open an issue on GitHub.

---

**Developed with ❤️ for BookClubApp**
