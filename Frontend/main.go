package main

import (
	"fmt"
	"net/http"
	"html/template"
	"strings"
	"os"
	"log"
	"github.com/joho/godotenv"
	"encoding/json"
)

type TagObject struct {
	Tags []string
}
type User struct {
	Username string
	Password string
	FavoriteTags []TagObject
}
type JWTObject struct {
	Token string
}
type Product struct {
	ID string
	Name string
	Price float64
	Tags []string
	Images []string
	Description string
	Rating float64
	URL string
}
type IndexInput struct {
	Products []Product
	User User
	AllTags []TagObject
}

var USER_SERVICE_URL string
var CATALOG_SERVICE_URL string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	USER_SERVICE_URL = os.Getenv("USER_SERVICE_URL")
	if USER_SERVICE_URL == "" {
		log.Fatalf("USER_SERVICE_URL not set in .env file")
	}
	CATALOG_SERVICE_URL = os.Getenv("CATALOG_SERVICE_URL")
	if CATALOG_SERVICE_URL == "" {
		log.Fatalf("CATALOG_SERVICE_URL not set in .env file")
	}
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

	resp , err := http.Get(USER_SERVICE_URL+"/user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	var products []Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		http.Error(w, "Failed to decode response", http.StatusInternalServerError)
		return
	}

	// TODO: Make The IndexInput object!!!
	
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
	resp, err := http.Post(USER_SERVICE_URL+"/login", "application/json", strings.NewReader(`{"username":"`+username+`","password":"`+password+`"}`))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	var jwtObj JWTObject
	if err := json.NewDecoder(resp.Body).Decode(&jwtObj); err != nil {
		http.Error(w, "Failed to decode response", http.StatusInternalServerError)
		return
	}

	// Set JWT token in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "jwt",
		Value: jwtObj.Token,
		Path:  "/",
	})
	// ----------------------------CATALOG PART--------------------------------
	resp, err = http.Get(CATALOG_SERVICE_URL+"/catalog")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var products []Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		http.Error(w, "Failed to decode response", http.StatusInternalServerError)
		return
	}
	renderTemplate(w, "catalog.html", products)
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	http.Post(USER_SERVICE_URL+"/register", "application/json", strings.NewReader(`{"username":"`+username+`","password":"`+password+`"}`))
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