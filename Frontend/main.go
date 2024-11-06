package main

import (
	"fmt"
	"html/template"
	"log"
	"main/helper"
	"main/objects"
	"net/http"
	"os"
	"strconv"
	"strings"
	"main/logic"
	"github.com/joho/godotenv"
)

// ----------------- Object Structs ----------------- //

// ----------------- Varibles ----------------- //
var USER_SERVICE_URL string
var CATALOG_SERVICE_URL string
var JWTTIMEOUT int

const PORT = "8081"

// ----------------- Setup functions ----------------- //
func init() {
	err := godotenv.Load("../.env")
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
	JWTTIMEOUT, err = strconv.Atoi(strings.TrimSuffix(os.Getenv("JWT_TIMEOUT"), "h"))
	if err != nil {
		log.Fatalf("Invalid JWT_TIMEOUT value in .env file")
	}
}

func main() {
	fmt.Println("\n\n\n\nServer is running on port", PORT) // This will print before the server starts
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/login", GetLogin)
	http.HandleFunc("/register", GetRegister)
	http.HandleFunc("/logout", HandleLogout)
	http.HandleFunc("/handle-login", HandleLogin)
	http.HandleFunc("/handle-register", HandleRegister)
	http.HandleFunc("/catalog", GetCatalog)
	http.HandleFunc("/catalog/", GetCatalogStyleSearch)
	http.HandleFunc("/favorite/tag/", HandleFavoriteTag)
	http.HandleFunc("/unfavorite/tag/", HandleUnfavoriteTag)
	http.HandleFunc("/product/", GetProduct)
	http.HandleFunc("/report/tag/", HandleReportTag)

	if err := http.ListenAndServe(fmt.Sprint(":", PORT), nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

// ----------------- Endpoint functions ----------------- //
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	indexInput, err := logic.IndexHandler(w, r)
	if err != nil {
		return
	}
	renderTemplate(w, "index.html", indexInput)
}

func GetLogin(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login.html", nil)
}

func GetRegister(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "register.html", nil)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, POST only", http.StatusMethodNotAllowed)
		return
	}
	var loginObject objects.LoginObject = objects.LoginObject{Email: r.FormValue("email"), Password: r.FormValue("password")}
	resp, err := helper.MakehttpPostRequest(USER_SERVICE_URL+"/login", "", strings.NewReader(`{"email":"`+loginObject.Email+`","password":"`+loginObject.Password+`"}`))
	if err != nil {
		renderTemplate(w, "login.html", objects.LoginPageData{Login: &loginObject, Error: &objects.HtmlError{Message: "Failed to login", StatusCode: http.StatusInternalServerError, Error: true}})
		return
	}
	defer resp.Body.Close()
	var jwtObj objects.JWTObject
	helper.ResponseToObj(resp, &jwtObj)

	resp, err = helper.MakehttpGetRequest(USER_SERVICE_URL+"/me", jwtObj.Token)
	if err != nil || resp.StatusCode == http.StatusUnauthorized {
		renderTemplate(w, "login.html", objects.LoginPageData{Login: &loginObject, Error: &objects.HtmlError{Message: "Incorrect JWT? (This shoulden't happen)", StatusCode: http.StatusInternalServerError, Error: true}})
		return
	}
	var user objects.User
	helper.ResponseToObj(resp, &user)
	// Set users cookie's
	helper.SetUsersCookies(w, user, jwtObj.Token)
	w.Header().Set("HX-Redirect", "/")
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	helper.ClearUsersCookies(w)
	w.Header().Set("HX-Redirect", "/")
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	var page, registerOject, err = logic.HandleRegister(w, r)
	if err != nil {
		return
	}
	renderTemplate(w, page, registerOject)
}

func GetCatalog(w http.ResponseWriter, r *http.Request) {
	var pagedata, err = logic.GetCatalog(w, r)
	if err != nil {
		return
	}
	renderTemplate(w, "catalog.html", pagedata)
}

func GetCatalogStyleSearch(w http.ResponseWriter, r *http.Request) {
	var pagedata, err = logic.GetCatalogStyleSearch(w, r)
	if err != nil {
		return
	}
	renderTemplate(w, "catalog.html", pagedata)
}

func HandleFavoriteTag(w http.ResponseWriter, r *http.Request) {
	var pagedata, err = logic.HandleFavoriteTag(w, r)
	if err != nil {
		return
	}
	renderTemplate(w, "favoriteTag.html", pagedata)
}

func HandleUnfavoriteTag(w http.ResponseWriter, r *http.Request) {
	var pagedata, err = logic.HandleUnfavoriteTag(w, r)
	if err != nil {
		return
	}
	renderTemplate(w, "normalTag.html", pagedata)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	var pagedata, err = logic.GetProduct(w, r)
	if err != nil {
		return
	}
	renderTemplate(w, "product.html", pagedata)
}

func HandleReportTag(w http.ResponseWriter, r *http.Request) {
	var pagedata, err = logic.HandleReportTag(w, r)
	if err != nil {
		return
	}
	renderTemplate(w, "feedback.html", pagedata)
}

// ----------------- Helper functions ----------------- //
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
