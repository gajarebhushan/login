package main

import (
    "context"
    "html/template"
    "log"
    "net/http"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// User represents a user in the database
type User struct {
    Username string
    Password string
}

// Connect to MongoDB
func connectDB() (*mongo.Client, error) {
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        return nil, err
    }
    return client, nil
}

// Handler for login page
func loginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        // Serve the login page
        tmpl, err := template.ParseFiles("login.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        tmpl.Execute(w, nil)
    } else if r.Method == "POST" {
        // Handle login form submission
        username := r.FormValue("username")
        password := r.FormValue("password")

        // Validate credentials from MongoDB
        client, err := connectDB()
        if err != nil {
            http.Error(w, "Error connecting to database", http.StatusInternalServerError)
            return
        }
        defer client.Disconnect(context.Background())

        collection := client.Database("mydatabase").Collection("users")
        var user User
        err = collection.FindOne(context.Background(), User{Username: username, Password: password}).Decode(&user)
        if err != nil {
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
            return
        }

        // Redirect to dashboard upon successful login
        http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
    }
}

// Handler for dashboard page
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
    // Serve the dashboard page
    tmpl, err := template.ParseFiles("dashboard.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, nil)
}

func main() {
    http.HandleFunc("/", loginHandler)
    http.HandleFunc("/dashboard", dashboardHandler)

    log.Println("Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
