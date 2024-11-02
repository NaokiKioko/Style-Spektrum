package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Tag struct {
	Name          string
	FavoriteCount int
}
type User struct {
	Username     string
	Email        string
	FavoriteTags []Tag
	Role         string
}
type LoginObject struct {
	Email    string
	Password string
}
type JWTObject struct {
	Token string
}
type Product struct {
	ID          string
	Name        string
	Price       float64
	Tags        []string
	Images      []string
	Description string
	Rating      float64
	URL         string
}
type IndexInput struct {
	User         User
	FavoriteTags []Tag
	AllTags      []Tag
}
type HtmlError struct {
	Message    string
	StatusCode int
	Error      bool
}
type CatalogPageData struct {
	Products []Product
	Error    *HtmlError
}
type LoginPageData struct {
	Login *LoginObject
	Error *HtmlError
}

var USER_SERVICE_URL string
var CATALOG_SERVICE_URL string
var JWTTIMEOUT int

const PORT = "8081"

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
	http.HandleFunc("/favorite/tag/", FavoriteTag)
	http.HandleFunc("/unfavorite/tag/", UnfavoriteTag)

	if err := http.ListenAndServe(fmt.Sprint(":", PORT), nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	var user, jwt, err = GetUserFromCookies(r)
	if err != nil {
		user = User{}
	} else if jwt == "" {
		ClearUsersCookies(w)
		w.Header().Set("HX-Redirect", "/")
		return
	}
	if len(user.FavoriteTags) == 1 && user.FavoriteTags[0].Name == "" {
		user.FavoriteTags = []Tag{}
	}
	resp, err := MakehttpGetRequest(CATALOG_SERVICE_URL+"/tags", "")
	if err != nil {
		log.Fatalf("Error getting tags from catalog service")
	}
	alltags := []Tag{}
	ResponseToObj(resp, &alltags)

	if user.Email != "" {
		alltags = RemoveFavoriteTagsFromAllTags(alltags, user.FavoriteTags)
	}
	alltags = sortTagsByFavoriteCount(alltags)

	var indexInput IndexInput = IndexInput{user, user.FavoriteTags, alltags}

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
	var loginObject LoginObject = LoginObject{r.FormValue("email"), r.FormValue("password")}
	resp, err := MakehttpPostRequest(USER_SERVICE_URL+"/login", "", strings.NewReader(`{"email":"`+loginObject.Email+`","password":"`+loginObject.Password+`"}`))
	if err != nil {
		renderTemplate(w, "login.html", LoginPageData{&loginObject, &HtmlError{"Failed to login", http.StatusInternalServerError, true}})
		return
	}
	defer resp.Body.Close()
	var jwtObj JWTObject
	ResponseToObj(resp, &jwtObj)

	resp, err = MakehttpGetRequest(USER_SERVICE_URL+"/me", jwtObj.Token)
	if err != nil || resp.StatusCode == http.StatusUnauthorized {
		renderTemplate(w, "login.html", LoginPageData{&loginObject, &HtmlError{"Incorrect JWT? (This shoulden't happen)", http.StatusInternalServerError, true}})
		return
	}
	var user User
	ResponseToObj(resp, &user)
	// Set users cookie's
	SetUsersCookies(w, user, jwtObj.Token)
	w.Header().Set("HX-Redirect", "/")
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	ClearUsersCookies(w)
	w.Header().Set("HX-Redirect", "/")
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, POST only", http.StatusMethodNotAllowed)
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmpassword := r.FormValue("confirmpassword")

	if password != confirmpassword {
		renderTemplate(w, "register.html", LoginPageData{&LoginObject{email, ""}, &HtmlError{"Passwords do not match", http.StatusBadRequest, true}})
		return
	}
	// _, err := http.Post(USER_SERVICE_URL+"/register", "application/json", strings.NewReader(`{"email":"`+email+`","password":"`+password+`"}`))
	_, err := MakehttpPostRequest(USER_SERVICE_URL+"/register", "", strings.NewReader(`{"email":"`+email+`","password":"`+password+`"}`))
	if err != nil {
		renderTemplate(w, "error.html", HtmlError{"Failed to register", http.StatusInternalServerError, true})
		return
	}
	renderTemplate(w, "login.html", LoginPageData{&LoginObject{email, password}, nil})
}

func GetCatalog(w http.ResponseWriter, r *http.Request) {
	resp, err := MakehttpGetRequest(CATALOG_SERVICE_URL+"/catalog", "")
	if err != nil {
		log.Fatalf("Error getting catalog from catalog service")
	}
	defer resp.Body.Close()
	var products []Product
	ResponseToObj(resp, &products)
	renderTemplate(w, "catalog.html", CatalogPageData{products, nil})
}

func GetCatalogStyleSearch(w http.ResponseWriter, r *http.Request) {
	resp, err := MakehttpGetRequest(CATALOG_SERVICE_URL+"/catalog/tags/"+r.URL.Path[len("/catalog/"):], "")
	if err != nil {
		log.Fatalf("Error getting catalog from catalog service")
	}
	defer resp.Body.Close()
	var products []Product
	ResponseToObj(resp, &products)
	renderTemplate(w, "catalog.html", CatalogPageData{products, nil})
}

func FavoriteTag(w http.ResponseWriter, r *http.Request) {
	var tag string = r.URL.Path[len("/favorite/tag/"):]
	var user, jwt, err = GetUserFromCookies(r)
	if err != nil {
		// User is not logged in and cant favorite tags
	}
	if jwt == "" {
		ClearUsersCookies(w)
		w.Header().Set("HX-Redirect", "/")
		return
	}
	_, err = MakehttpPostRequest(USER_SERVICE_URL+"/favorite/"+tag, jwt, nil)
	if err != nil {
		log.Fatalf("Error favoriting tag")
	}
	user.FavoriteTags = append(user.FavoriteTags, Tag{Name: tag})
	SetUsersCookies(w, user, jwt)
}

func UnfavoriteTag(w http.ResponseWriter, r *http.Request) {
	var tag string = r.URL.Path[len("/unfavorite/tag/"):]
	var user, jwt, err = GetUserFromCookies(r)
	if err != nil {
		// User is not logged in and cant favorite tags
	}
	if jwt == "" {
		ClearUsersCookies(w)
		w.Header().Set("HX-Redirect", "/")
		return
	}
	_, err = MakehttpDeleteRequest(USER_SERVICE_URL+"/favorite/"+tag, jwt)
	if err != nil {
		log.Fatalf("Error unfavoriting tag")
	}
	for x, favtag := range user.FavoriteTags {
		if favtag.Name == tag {
			user.FavoriteTags = append(user.FavoriteTags[:x], user.FavoriteTags[x+1:]...)
		}
	}
	SetUsersCookies(w, user, jwt)
}
// ----------------- Helper functions -----------------
func SetUsersCookies(w http.ResponseWriter, user User, jwt string) {
	// Gather the cookies in a slice
	cookies := []http.Cookie{
		{
			Name:    "JWT",
			Path:    "/",
			Value:   jwt,
			Expires: time.Now().Add(time.Duration(JWTTIMEOUT) * time.Hour),
		},
		{
			Name:  "Email",
			Path:  "/",
			Value: user.Email,
		},
		{
			Name:  "Username",
			Path:  "/",
			Value: strings.Split(user.Email, "@")[0],
		},
	}

	// Join favorite tag names and add as a cookie
	favtagNames := make([]string, len(user.FavoriteTags))
	for i, tag := range user.FavoriteTags {
		favtagNames[i] = tag.Name
	}
	cookies = append(cookies, http.Cookie{
		Name:  "FavoriteTags",
		Path:  "/",
		Value: strings.Join(favtagNames, ","),
	})

	// Set all cookies in a loop
	for _, cookie := range cookies {
		http.SetCookie(w, &cookie)
	}
}

func GetUserFromCookies(r *http.Request) (User, string, error) {
	var user User
	var jwt string

	// Get JWT cookie
	jwtCookie, err := r.Cookie("JWT")
	if err != nil {
		return user, jwt, err
	}
	jwt = jwtCookie.Value

	// Get Email cookie
	emailCookie, err := r.Cookie("Email")
	if err != nil {
		return user, jwt, err
	}
	user.Email = emailCookie.Value

	// Get Username cookie
	usernameCookie, err := r.Cookie("Username")
	if err != nil {
		return user, jwt, err
	}
	user.Username = usernameCookie.Value

	// Get FavoriteTags cookie
	favTagsCookie, err := r.Cookie("FavoriteTags")
	if err != nil {
		return user, jwt, err
	}
	favTagNames := strings.Split(favTagsCookie.Value, ",")
	for _, name := range favTagNames {
		user.FavoriteTags = append(user.FavoriteTags, Tag{Name: name})
	}
	return user, jwt, nil
}

func ClearUsersCookies(w http.ResponseWriter) {
	// Gather the cookies in a slice
	cookies := []http.Cookie{
		{
			Name:   "JWT",
			Path:   "/",
			MaxAge: -1,
		},
		{
			Name:   "Email",
			Path:   "/",
			MaxAge: -1,
		},
		{
			Name:   "Username",
			Path:   "/",
			MaxAge: -1,
		},
		{
			Name:   "FavoriteTags",
			Path:   "/",
			MaxAge: -1,
		},
	}
	for _, cookie := range cookies {
		http.SetCookie(w, &cookie)
	}
}

func RemoveFavoriteTagsFromAllTags(alltags []Tag, favtags []Tag) []Tag {
	for x, tag := range alltags {
		for _, favtag := range favtags {
			if tag.Name == favtag.Name {
				alltags = append(alltags[:x], alltags[x+1:]...)
			}
		}
	}
	return alltags
}

func MakehttpGetRequest(url string, jwt string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request")
	}
	if jwt != "" {
		req.Header.Add("authorization", "Bearer "+jwt)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error making request")
	}
	return resp, err
}

func MakehttpPostRequest(url string, jwt string, body *strings.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	// Add Authorization header if jwt is provided
	if jwt != "" {
		req.Header.Add("authorization", "Bearer "+jwt)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// Optionally, check if the response was successful
	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return resp, nil
}

func MakehttpDeleteRequest(url string, jwt string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	// Add Authorization header if jwt is provided
	if jwt != "" {
		req.Header.Add("authorization", "Bearer "+jwt)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// Optionally, check if the response was successful
	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return resp, nil
}

func ResponseToObj(resp *http.Response, obj interface{}) {
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to get response from:" + resp.Request.URL.String())
	}
	if err := json.NewDecoder(resp.Body).Decode(obj); err != nil {
		log.Fatalf("Failed to decode response")
	}
}

func sortTagsByFavoriteCount(tags []Tag) []Tag {
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].FavoriteCount > tags[j].FavoriteCount
	})
	return tags
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
