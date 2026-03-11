# 📚 Book Club App

A web application to manage, share, and comment your favorite books. The goal is to create a space where you can interact and express your ideas and though as you go through the chapters of the books you are reading. 

Our link to the backend API documentation
https://bookclub-backend.redwater-26f8bbd2.centralus.azurecontainerapps.io/swagger/index.html

---

## 🚀 Features
- Add and browse books
- Track the books you have read and the ones you want to read
- Comment on each chapter of the books that you are reading

---

## 🏗 Architecture

React Frontend
      │
      ▼
Go REST API
      │
      ▼
PostgreSQL Database

---

## 🛠 Tech Stack
**Frontend**
- React (Create React App)
- Capacitor (mobile builds)

**Backend**
- Go (REST API)
- Swagger API documentation

**Database**
- PostgreSQL

**Infrastructure**
- Docker
- Docker Compose
- Azure deployment

## 📝 Documentation
[Database Schema](https://dbdiagram.io/d/BookClub-68bc7d4961a46d388ec627de)

---

## 🔐 Environment Variables

The backend uses environment variables for configuration such as database credentials and API keys.

For local development, create a .env file in the root of the project.

```DB_HOST=
DB_USER=
DB_PASSWORD=
DB_NAME=
DB_PORT=5432
DB_SSLMODE=require
ALLOWED_ORIGINS=
GOOGLE_BOOKS_API_KEY=
PORT=5000
```

## ⚡ Getting Started

### 1. Frontend setup (React)
```bash
cd book-club-app
npm install
npm start
```

### 2. Backend setup (Go)
#### 🐳 Run the backend

Use the following command to install the dependencies and run the backend:

```bash
docker compose up --build
```


You can verify that the backend is running by navigating to:
http://localhost:5000/health

The internal API documentation can be found using the following link:
http://localhost:5000/swagger/index.html


### 3. Building for Web
To create a production-ready web build of the frontend:

```bash
cd book-club-app
npm run build
```

### Running as a Mobile App (Capacitor)
The frontend can be packaged as an iOS/Android app using Capacitor.

Run from root:

```bash
npm install @capacitor/core @capacitor/cli
npx cap init "BookClub" com.app.bookclub
```

Add platforms:
```
npx cap add ios
npx cap add android
```

Make sure capacitor.config.json has "webDir": "build".

Building & syncing
Whenever you change the frontend:

```bash
cd book-club-app
npm run build      # rebuild frontend
npx cap copy       # copy build into native projects
```

Opening native projects
```bash
npx cap open ios     # opens in Xcode (Mac only)
npx cap open android # opens in Android Studio
```

From there you can run on simulators or devices.