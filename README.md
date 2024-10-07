# **Jots: A Simple Twitter Clone**

## **Overview**
Jots is a rudimentary Twitter clone designed to explore basic functionalities of a social media platform. Users can post short messages, follow specific "channels" (akin to topics), and receive real-time notifications about new content in the channels they follow. This project serves as a learning tool for backend development, database interactions, and real-time communication using WebSockets.

---

## **Table of Contents**
- [Project Overview](#overview)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Setup and Installation](#setup-and-installation)
- [How to Run Locally](#how-to-run-locally)
- [Next Steps](#next-steps)

---

## **Features**
- Post "jots" (short messages)
- Follow/unfollow channels
- View jots based on channels
- Real-time notifications for new posts in followed channels
- User authentication (login/signup)

---

## **Tech Stack**
- **Backend:** Go
- **Database:** MySQL
- **Real-time Communication:** WebSockets, Redis
- **Frontend:** Basic HTML/CSS
- **Libraries:** 
  - `github.com/go-redis/redis/v8`
  - `github.com/gorilla/websocket`
  - `github.com/go-sql-driver/mysql`

---

## **Setup and Installation**

### **1. Clone the Repository**
```bash
git clone https://github.com/your-username/jots-app.git
cd jots-app
```
### **2. Install Dependencies**

Ensure Go is installed on your machine. Install the necessary Go modules by running:

```bash
go mod download
```

### **3. Setup MySQL Database**

Create a MySQL database and update the dsn in the db.go file with your database credentials. Example:

```bash
dsn := "your_user:your_password@tcp(127.0.0.1:3306)/jots_db"
```

### **4. Setup Redis**

Ensure Redis is running on your machine or use a remote Redis instance. Update the Redis connection settings in the redis.go file if necessary.

```bash
go mod download
```

### **5. Run the Application**

Once all the dependencies are set up, run the application using the following command:
```bash
go run main.go db.go models.go handlers.go redis.go
```

## **Next Steps**
	•	Enhance Frontend: Add more user-friendly design and UI features.
	•	User Profiles: Implement individual user profile pages.
	•	Direct Messaging: Introduce a direct messaging feature between users.
	•	Likes and Comments: Add functionality for users to like and comment on jots.
	•	Search: Implement a search functionality to find jots or users.
