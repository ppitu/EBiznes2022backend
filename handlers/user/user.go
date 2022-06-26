package user

import (
	"backend/database/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
)

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file")
	}
	fmt.Println(os.Getenv(key))
	return os.Getenv(key)
}

func Get(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusPreconditionFailed)
	}

	db, _ := c.Get("db").(*gorm.DB)

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, user)
}

func GetAll(c echo.Context) error {
	var users []models.User

	db, _ := c.Get("db").(*gorm.DB)

	if err := db.Model(&models.User{}).Find(&users).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, users)
}

func Create(c echo.Context) error {
	type RequestBody struct {
		Email string `json:"email"`
	}

	var body RequestBody

	if err := c.Bind(&body); err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	db, _ := c.Get("db").(*gorm.DB)

	user := models.User{
		Email: body.Email,
	}

	db.Create(&user)

	return c.JSON(http.StatusOK, user)
}

func Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusPreconditionFailed)
	}

	type RequestBody struct {
		Name string `json:"name"`
	}

	var body RequestBody

	if err := c.Bind(&body); err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	db, _ := c.Get("db").(*gorm.DB)

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	user.Email = body.Name

	db.Save(&user)

	return c.JSON(http.StatusOK, user)
}

func Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusPreconditionFailed)
	}

	db, _ := c.Get("db").(*gorm.DB)

	var user models.User

	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	db.Delete(&user)

	return c.NoContent(http.StatusOK)
}

func GenerateToken(email string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)
	if err != nil {
	}
	fmt.Println("Hash to store:", string(hash))

	return string(hash)

}

// Scopes: OAuth 2.0 scopes provide a way to limit the amount of access that is granted to an access token.
var googleOauthConfig *oauth2.Config

var githubOauthConfig *oauth2.Config

var oauthState = "Random-string"

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func OauthGoogleLogin(c echo.Context) error {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:1323/users/auth/google/callback",
		ClientID:     goDotEnvVariable("GOOGLECLIENTID"),
		ClientSecret: goDotEnvVariable("GOOGLECLIENTSECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	u := googleOauthConfig.AuthCodeURL(oauthState)
	return c.Redirect(http.StatusTemporaryRedirect, u)
}

func OAuthGithubLogin(c echo.Context) error {
	githubOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:1323/users/auth/github/callback",
		ClientID:     goDotEnvVariable("GITHUBCLIENTID"),
		ClientSecret: goDotEnvVariable("GITHUBCLIENTSECRET"),
		Scopes:       []string{"https://github.com/login/oauth/access_token"},
		Endpoint:     github.Endpoint,
	}

	fmt.Println("Test:")
	fmt.Println(githubOauthConfig.ClientID)
	u := githubOauthConfig.AuthCodeURL(oauthState)
	return c.Redirect(http.StatusTemporaryRedirect, u)
}

func OauthGoogleCallback(c echo.Context) error {
	data, err := getUserInfoGoogle(c.FormValue("state"), c.FormValue("code"))

	if err != nil {
		fmt.Println(err.Error())
		return c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/")
	}

	userinfo := new(models.User)
	json.Unmarshal([]byte(string(data)), &userinfo)

	userinfo.Token = GenerateToken(userinfo.Email)

	fmt.Println(userinfo.Email)

	cookieUser := new(http.Cookie)
	cookieUser.Path = "/success"
	cookieUser.Name = "user"
	cookieUser.Value = userinfo.Email
	cookieUser.Expires = time.Now().Add(30 * time.Second)
	cookieUser.HttpOnly = false
	cookieUser.Secure = true
	c.SetCookie(cookieUser)

	return c.Redirect(http.StatusSeeOther, "http://localhost:3000/success")
}

func OauthGithubCallback(c echo.Context) error {
	fmt.Println(c.Request().Body)
	data, err := getUserInfoGithub(c.FormValue("state"), c.FormValue("code"))

	if err != nil {
		fmt.Println(err.Error())
		return c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/")
	}

	userinfo := new(models.User)
	json.Unmarshal([]byte(string(data)), &userinfo)

	userinfo.Token = GenerateToken(userinfo.Email)

	fmt.Println(userinfo.Email)

	cookieUser := new(http.Cookie)
	cookieUser.Path = "/success"
	cookieUser.Name = "user"
	cookieUser.Value = userinfo.Email
	cookieUser.Expires = time.Now().Add(30 * time.Second)
	cookieUser.HttpOnly = false
	cookieUser.Secure = true
	c.SetCookie(cookieUser)

	return c.Redirect(http.StatusSeeOther, "http://localhost:3000/success")
}

func getUserInfoGoogle(state string, code string) ([]byte, error) {
	if state != oauthState {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	return contents, nil
}

func getUserInfoGithub(state string, code string) ([]byte, error) {
	fmt.Println(state)
	fmt.Println(oauthState)
	if state != oauthState {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := githubOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	return contents, nil
}

type item struct {
	id string
}

func calculateOrderAmount(items []item) int64 {
	// Replace this constant with a calculation of the order's amount
	// Calculate the order total on the server to prevent
	// people from directly manipulating the amount on the client
	return 1400
}

func HandleCreatePaymentIntent(c echo.Context) error {
	stripe.Key = "sk_test_51LD9HFEAp8DjvQZDaoDyh5FaY6WJ4n4mUWIB05eSK3JiYhZosDMgu0YQByYoBGAsoeI4uu5x7i16UC1L9KFqYcQ600WTABIIip"

	if c.Request().Method != "POST" {
		//http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return c.NoContent(http.StatusNotFound)
	}

	var req struct {
		Items []item `json:"items"`
	}

	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("json.NewDecoder.Decode: %v", err)
		return c.NoContent(http.StatusNotFound)
	}

	// Create a PaymentIntent with amount and currency
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(calculateOrderAmount(req.Items)),
		Currency: stripe.String(string(stripe.CurrencyEUR)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
	log.Printf("pi.New: %v", pi.ClientSecret)

	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("pi.New: %v", err)
		return c.NoContent(http.StatusNotFound)
	}

	writeJSON(c, struct {
		ClientSecret string `json:"clientSecret"`
	}{
		ClientSecret: pi.ClientSecret,
	})

	return c.NoContent(http.StatusOK)
}

func writeJSON(c echo.Context, v interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("json.NewEncoder.Encode: %v", err)
		return
	}
	c.Response().Header().Set("Content-Type", "application/json")
	if _, err := io.Copy(c.Response(), &buf); err != nil {
		log.Printf("io.Copy: %v", err)
		return
	}
}
