package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
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
	FavoriteTags []Tag
}
type LoginObject struct {
	Username string
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
	User     User
	AllTags  []Tag
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
	if err := http.ListenAndServe(fmt.Sprint(":", PORT), nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	var user User = GetUserFromCookie(r)
	if user.Username == "" {
		var jwt = r.Header.Get("authorization")
		if jwt != "" {
			resp := MakehttpGetRequest(USER_SERVICE_URL+"/me", jwt)
			if resp.StatusCode == http.StatusOK {
				if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
					log.Fatalf("Failed to decode response")
				}
				r.Header.Set("username", user.Username)
				favtagNames := []string{}
				for _, tag := range user.FavoriteTags {
					favtagNames = append(favtagNames, tag.Name)
				}
				r.Header.Set("favorite_tags", strings.Join(favtagNames, ","))
			}
		}
	}
	var alltags []Tag = getAllTags()
	if user.Username != "" {
		alltags = RemoveFavoriteTagsFromAllTags(alltags, user.FavoriteTags)
	}

	var indexInput IndexInput = IndexInput{user, alltags}

	renderTemplate(w, "index.html", indexInput)
}

func GetLogin(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login.html", nil)
}

func GetRegister(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "register.html", nil)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var loginObject LoginObject = LoginObject{r.FormValue("username"), r.FormValue("password")}
	resp, err := http.Post(USER_SERVICE_URL+"/login", "application/json", strings.NewReader(`{"username":"`+loginObject.Username+`","password":"`+loginObject.Password+`"}`))
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
	resp, err = http.Get(CATALOG_SERVICE_URL + "/catalog")
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
	if (http.Cookie{Name: "username"}.Value != "") {
		user = User{
			Username:     http.Cookie{Name: "username"}.Value,
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

func getAllTags() []Tag {
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

func MakehttpGetRequest(url string, jwt string) *http.Response {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request")
	}
	req.Header.Set("authorization", jwt)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error making request")
	}
	return resp
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