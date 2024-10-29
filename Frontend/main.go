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
	Role 	   string
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
	User    User
	AllTags []Tag
}
type HtmlError struct {
	Message    string
	StatusCode int
	Endpoint   string
}

var USER_SERVICE_URL string
var CATALOG_SERVICE_URL string

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
}

func main() {
	fmt.Println("Server is running on port", PORT) // This will print before the server starts
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/login", GetLogin)
	http.HandleFunc("/register", GetRegister)
	http.HandleFunc("/handle-login", HandleLogin)
	http.HandleFunc("/handle-register", HandleRegister)

	if err := http.ListenAndServe(fmt.Sprint(":", PORT), nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	var user User = GetUserFromCookie(r)
	if user.Email == "" {
		var jwt = r.Header.Get("authorization")
		if jwt != "" {
			resp, err := MakehttpGetRequest(USER_SERVICE_URL+"/me", jwt)
			if err != nil {
				if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
					log.Fatalf("Failed to decode response")
				}
				user.Username = strings.Split(user.Email, "@")[0]
				r.Header.Set("email", user.Email)
				favtagNames := []string{}
				for _, tag := range user.FavoriteTags {
					favtagNames = append(favtagNames, tag.Name)
				}
				r.Header.Set("favorite_tags", strings.Join(favtagNames, ","))
			}
		}
	}
	var alltags []Tag = GetAllTags()
	if user.Email != "" {
		alltags = RemoveFavoriteTagsFromAllTags(alltags, user.FavoriteTags)
	}
	alltags = sortTagsByFavoriteCount(alltags)

	var indexInput IndexInput = IndexInput{user, alltags}

	renderTemplate(w, "index.html", indexInput)
}

func GetLogin(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login.html", LoginObject{"", ""})
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
		renderTemplate(w, "error.html", HtmlError{"Failed to login", http.StatusInternalServerError, "/login"})
		return
	}
	defer resp.Body.Close()
	var jwtObj JWTObject
	ResponseToObj(resp, &jwtObj)

	resp, err = MakehttpGetRequest(USER_SERVICE_URL+"/me", jwtObj.Token)
	if err != nil || resp.StatusCode == http.StatusUnauthorized {
		renderTemplate(w, "error.html", HtmlError{"Failed to get you", http.StatusInternalServerError, "/login"})
		return
	}
	var user User
	ResponseToObj(resp, &user)
	// Set JWT token in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "jwt",
		Value: jwtObj.Token,
		Path:  "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "email",
		Value: user.Email,
		Path:  "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "username",
		Value: strings.Split(user.Email, "@")[0],
		Path:  "/",
	})
	// Set favorite tags in a cookie
	if len(user.FavoriteTags) != 0 {
		favtagNames := []string{}
		for _, tag := range user.FavoriteTags {
			favtagNames = append(favtagNames, tag.Name)
		}
		favtagCounts := []string{}
		for _, tag := range user.FavoriteTags {
			favtagCounts = append(favtagCounts, strconv.Itoa(tag.FavoriteCount))
		}
		http.SetCookie(w, &http.Cookie{
			Name:  "favorite_tags",
			Value: strings.Join(favtagNames, ","),
			Path:  "/",
		})
		http.SetCookie(w, &http.Cookie{
			Name:  "favorite_tag_counts",
			Value: strings.Join(favtagCounts, ","),
			Path:  "/",
		})
	}
	// ----------------------------Tags PART--------------------------------
	var alltags []Tag = GetAllTags()
	if len(user.FavoriteTags) != 0 {
		alltags = RemoveFavoriteTagsFromAllTags(alltags, user.FavoriteTags)
	}
	indexInput := IndexInput{user, alltags}

	renderTemplate(w, "index.html", indexInput)
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, POST only", http.StatusMethodNotAllowed)
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	// _, err := http.Post(USER_SERVICE_URL+"/register", "application/json", strings.NewReader(`{"email":"`+email+`","password":"`+password+`"}`))
	_, err := MakehttpPostRequest(USER_SERVICE_URL+"/register", "", strings.NewReader(`{"email":"`+email+`","password":"`+password+`"}`))
	if err != nil {
		renderTemplate(w, "error.html", HtmlError{"Failed to register", http.StatusInternalServerError, "/register"})
		return
	}
	renderTemplate(w, "login.html", LoginObject{email, password})
}

func GetCatalog() []Product {
	resp, err := http.Get(CATALOG_SERVICE_URL + "/catalog")
	if err != nil {
		log.Fatalf("Error getting catalog from catalog service")
	}
	defer resp.Body.Close()
	var products []Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		log.Fatalf("Failed to decode response")
	}
	return products
}

func GetUserFromCookie(r *http.Request) User {
	var user User
	if (http.Cookie{Name: "email"}.Value != "") {
		user = User{
			Username:     http.Cookie{Name: "username"}.Value,
			Email:        http.Cookie{Name: "email"}.Value,
			FavoriteTags: []Tag{},
		}
		tagnames := strings.Split(http.Cookie{Name: "favorite_tags"}.Value, ",")
		tagfavoritecounts := strings.Split(http.Cookie{Name: "favorite_tag_counts"}.Value, ",")
		for i, tagname := range tagnames {
			tag := Tag{
				Name: tagname,
				FavoriteCount: func() int {
					count, err := strconv.Atoi(tagfavoritecounts[i])
					if err != nil {
						return 0
					}
					return count
				}(),
			}
			user.FavoriteTags = append(user.FavoriteTags, tag)
		}
	} else {
		user = User{}
	}
	return user
}

func GetAllTags() []Tag {
	resp, err := http.Get(CATALOG_SERVICE_URL + "/tags")
	if err != nil {
		log.Fatalf("Error getting tags from catalog service")
	}
	defer resp.Body.Close()
	var tags []Tag
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		log.Fatalf("Failed to decode response")
	}
	return tags
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

func GetCurrentUser(r *http.Request) User {
	jwt := r.Header.Get("Authorization")
	if jwt == "" {
		return User{}
	}
	resp, err := http.Get(USER_SERVICE_URL + "/me")
	if err != nil {
		log.Fatalf("Error getting user from user service")
	}
	defer resp.Body.Close()
	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		log.Fatalf("Failed to decode response")
	}
	return user
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

func ResponseToObj(resp *http.Response, obj interface{}) {
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
