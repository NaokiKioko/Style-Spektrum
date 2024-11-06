package helper

import (
	"encoding/json"
	"fmt"
	"log"
	"main/objects"
	"net/http"
	"sort"
	"strings"
	"time"
	"os"
	"strconv"
	"github.com/joho/godotenv"
)

var USER_SERVICE_URL string
var CATALOG_SERVICE_URL string
var JWTTIMEOUT int
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

func SetUsersCookies(w http.ResponseWriter, user objects.User, jwt string) {
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

func GetUserFromCookies(r *http.Request) (objects.User, string, error) {
	var user objects.User
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
		user.FavoriteTags = append(user.FavoriteTags, objects.Tag{Name: name})
	}
	for x, tag := range user.FavoriteTags {
		if tag.Name == "" {
			user.FavoriteTags = append(user.FavoriteTags[:x], user.FavoriteTags[x+1:]...)
		}
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

func RemoveFavoriteTagsFromAllTags(alltags []objects.Tag, favtags []objects.Tag) []objects.Tag {
	for x := 0; x < len(alltags); {
		removed := false
		for _, favtag := range favtags {
			if alltags[x].Name == favtag.Name {
				alltags = append(alltags[:x], alltags[x+1:]...)
				removed = true
				break
			}
		}
		if !removed {
			x++
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
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
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

func SortTagsByFavoriteCount(tags []objects.Tag) []objects.Tag {
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].FavoriteCount > tags[j].FavoriteCount
	})
	return tags
}
