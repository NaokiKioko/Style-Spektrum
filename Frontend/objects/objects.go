package objects

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
	ID          string `json:"_id"`
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

