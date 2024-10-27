package main

import (
	"fmt"
	"net/http"
	"html/template"
)

type User struct {
	Username string
	Password string
	FaboriteTags []string
}

const PORT = "8081"

func main() {
	http.Handle("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir("dist"))))
    fmt.Println("Server is running on port", PORT) // This will print before the server starts
    http.HandleFunc("/", IndexHandler)
    if err := http.ListenAndServe(fmt.Sprint(":", PORT), nil); err != nil {
        fmt.Println("Error starting server:", err)
    }
}
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	jwt := r.Header.Get("Authorization")

	if jwt == "" {
		renderTemplate(w, "index.html", nil)
	} else {
		renderTemplate(w, "index.html", nil)
		fmt.Println(jwt)
	}
}

func GetLogin(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login.html", nil)
}
func GetRegister(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "register.html", nil)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println(username, password)
	renderTemplate(w, "catalog.html", nil)
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println(username, password)
	renderTemplate(w, "login.html", nil)
}

func renderTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	t, err := template.ParseFiles("templates/" + templateName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}