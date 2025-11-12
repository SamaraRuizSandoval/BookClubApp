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
```bash
cd book-club-app
npm install
npm start

### 2. Backend setup (Go)
#### 🐳 Run the backend

From the project root:

```bash
docker compose up --build


You can verify that the backend is running by navigating to:
http://localhost:5000/health

### 3. Building for Web
To create a production-ready web build of the frontend:

``bash
cd book-club-app
npm run build
``

### Running as a Mobile App (Capacitor)
The frontend can be packaged as an iOS/Android app using Capacitor.

``bash
# from project root
npm install @capacitor/core @capacitor/cli
npx cap init "BookClub" com.app.bookclub

# Add platforms
npx cap add ios
npx cap add android
``

Make sure capacitor.config.json has "webDir": "build".

Building & syncing
Whenever you change the frontend:

``bash
cd book-club-app
npm run build      # rebuild frontend
npx cap copy       # copy build into native projects
``

Opening native projects
``bash
npx cap open ios     # opens in Xcode (Mac only)
npx cap open android # opens in Android Studio
``

From there you can run on simulators or devices.