package logic

import (
	"log"
	"main/helper"
	"main/objects"
	"net/http"
	"errors"
	"os"
	"strconv"
	"strings"
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

func IndexHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var user, jwt, err = helper.GetUserFromCookies(r)
	if err != nil {
		user = objects.User{}
	} else if jwt == "" {
		helper.ClearUsersCookies(w)
		w.Header().Set("HX-Redirect", "/")
		return nil, errors.New("user not logged in")
	}
	if len(user.FavoriteTags) <= 0 {
		user.FavoriteTags = []objects.Tag{}
	}
	resp, err := helper.MakehttpGetRequest(CATALOG_SERVICE_URL+"/tags", "")
	if err != nil {
		log.Fatalf("Error getting tags from catalog service")
	}
	alltags := []objects.Tag{}
	helper.ResponseToObj(resp, &alltags)

	if user.Email != "" {
		alltags = helper.RemoveFavoriteTagsFromAllTags(alltags, user.FavoriteTags)
	}
	alltags = helper.SortTagsByFavoriteCount(alltags)

	return objects.IndexInput{User: user, FavoriteTags: user.FavoriteTags, AllTags: alltags}, nil
}

func HandleRegister(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	if r.Method != http.MethodPost {
		
		return "", nil, errors.New("method not allowed")
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmpassword := r.FormValue("confirmpassword")

	if password != confirmpassword {
		return "register.html", objects.LoginPageData{Login: &objects.LoginObject{Email: email, Password: ""}, Error: &objects.HtmlError{Message: "Passwords do not match", StatusCode: http.StatusBadRequest, Error: true}}, nil
	}
	_, err := helper.MakehttpPostRequest(USER_SERVICE_URL+"/register", "", strings.NewReader(`{"email":"`+email+`","password":"`+password+`"}`))
	if err != nil {
		return "error.html", objects.HtmlError{Message: "Failed to register", StatusCode: http.StatusInternalServerError, Error: true}, nil
	}
	return "login.html", objects.LoginPageData{Login: &objects.LoginObject{Email: email, Password: password}, Error: nil}, nil
}

func GetCatalog(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	resp, err := helper.MakehttpGetRequest(CATALOG_SERVICE_URL+"/catalog", "")
	if err != nil {
		return nil, errors.New("error getting catalog from catalog service")
	}
	defer resp.Body.Close()
	var products []objects.Product
	helper.ResponseToObj(resp, &products)
	return objects.CatalogPageData{Products: products, Error: nil}, nil
}

func GetCatalogStyleSearch(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	resp, err := helper.MakehttpGetRequest(CATALOG_SERVICE_URL+"/catalog/tags/"+r.URL.Path[len("/catalog/"):], "")
	if err != nil {
		return nil, errors.New("error getting catalog from catalog service")
	}
	defer resp.Body.Close()
	var products []objects.Product
	helper.ResponseToObj(resp, &products)
	return objects.CatalogPageData{Products: products, Error: nil}, nil
}

func HandleFavoriteTag(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var tag string = r.URL.Path[len("/favorite/tag/"):]
	var user, jwt, err = helper.GetUserFromCookies(r)
	// if err != nil {
	// 	// User is not logged in and cant favorite tags
	// }
	if jwt == "" {
		helper.ClearUsersCookies(w)
		w.Header().Set("HX-Redirect", "/")
		return nil, errors.New("user not logged in")
	}
	_, err = helper.MakehttpPostRequest(USER_SERVICE_URL+"/favorite/tag/"+tag, jwt, strings.NewReader(`{}`))
	if err != nil {
		return nil, errors.New("error favoriting tag")
	}
	user.FavoriteTags = append(user.FavoriteTags, objects.Tag{Name: tag})
	helper.SetUsersCookies(w, user, jwt)
	return objects.Tag{Name: tag, FavoriteCount: 0}, nil
}

func HandleUnfavoriteTag(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var tag string = r.URL.Path[len("/unfavorite/tag/"):]
	var user, jwt, err = helper.GetUserFromCookies(r)
	// if err != nil {
	// 	// User is not logged in and cant favorite tags
	// }
	if jwt == "" {
		helper.ClearUsersCookies(w)
		w.Header().Set("HX-Redirect", "/")
		return nil, errors.New("user not logged in")
	}
	_, err = helper.MakehttpDeleteRequest(USER_SERVICE_URL+"/favorite/tag/"+tag, jwt)
	if err != nil {
		log.Print("Tag not found in favorites")
	}
	for x, favtag := range user.FavoriteTags {
		if favtag.Name == tag {
			user.FavoriteTags = append(user.FavoriteTags[:x], user.FavoriteTags[x+1:]...)
		}
	}
	helper.SetUsersCookies(w, user, jwt)
	return objects.Tag{Name: tag, FavoriteCount: 0}, nil
}

func GetProduct(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var productid string = r.URL.Path[len("/product/"):]
	resp, err := helper.MakehttpGetRequest(CATALOG_SERVICE_URL+"/catalog/"+productid, "")
	if err != nil {
		return nil, errors.New("error getting product from catalog service")
	}
	defer resp.Body.Close()
	var product objects.Product
	helper.ResponseToObj(resp, &product)
	return product, nil
}

// /report/tag/:ID/:tagname
func HandleReportTag(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var Varibles string = r.URL.Path[len("/report/tag/"):]
	var productid, tag string = strings.Split(Varibles, "/")[0], strings.Split(Varibles, "/")[1]
	var user, jwt, err = helper.GetUserFromCookies(r)
	if jwt == "" {
		helper.ClearUsersCookies(w)
		w.Header().Set("HX-Redirect", "/")
		return nil, errors.New("user not logged in")
	}
	_, err = helper.MakehttpPostRequest(CATALOG_SERVICE_URL+"/report/"+productid+"/tag/"+tag, "", strings.NewReader(`{"Email": "`+user.Email+`"}`))
	if err != nil {
		return nil, errors.New("error reporting tag")
	}
	return objects.Feedback{Title: "Report Complete", Message: "You Reported the "+tag+" tag for this product! With enough support this will add or remove this tag from this product!"}, nil
}