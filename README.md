# 📚 Book Club App

A web application to manage, share, and comment your favorite books. The goal is to create a space where you can interact and express your ideas and though as you go through the chapters of the books you are reading. 

---

## 🚀 Features
- Add and browse books
- Track the books you have read and the ones you want to read
- Comment on each chapter of the books that you are reading

---

## 🛠 Tech Stack
- **Frontend**: React (Create React App)
- **Backend**: Go
- **Database**: PostgreSQL
- **Deployment**: Docker

---

## ⚡ Getting Started

### 1. Frontend setup (React)
``bash
cd book-club-app
npm install
npm start``

### 2. Backend setup (Go)
``bash
cd backend
go mod tidy
go run main.go``

By default, the backend will run on port 5000 but you can use the port flag to change the port. For example:
`go run main.go -port 8081`

You can verify that the backend is running by navigating to:
http://localhost:5000/health


