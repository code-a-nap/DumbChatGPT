package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	Message struct {
		Name    string `json:"name" validate:"required"`
		Content string `json:"content" validate:"required"`
	}

	CustomValidator struct {
		validator *validator.Validate
	}
)

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Uuid  string `json:"uuid"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

var t = &Template{
	templates: template.Must(template.ParseGlob("public/views/*.html")),
}

// Echo instance
var e = echo.New()

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// Handler for chat
func chat(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	uuid := claims.Uuid
	if doesFileExist("public/views/" + uuid + ".html") {
		return c.Render(http.StatusOK, uuid, c)
	} else {
		return c.Render(http.StatusOK, "chat", claims)
	}
}

// Handler for login
func login(c echo.Context) error {
	if c.Request().Method == "GET" {
		return c.Render(http.StatusOK, "login", "")
	} else if c.Request().Method == "POST" {
		email := c.FormValue("email")
		password := c.FormValue("password")
		if email != "beta-user@test.com" || password != "secret" {
			return c.JSON(http.StatusForbidden, echo.Map{"error": "Wrong username and/or password!"})
		} else {
			claims := &jwtCustomClaims{
				"Beta User",
				uuid.New().String(),
				false,
				jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
				},
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			t, err := token.SignedString(secretKey)
			if err != nil {
				return err
			}
			cookie := http.Cookie{
				Name:     "token",
				Value:    t,
				Path:     "/",
				MaxAge:   3600,
				HttpOnly: false,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
			}
			c.SetCookie(&cookie)
			return c.JSON(http.StatusOK, echo.Map{"error": ""})
		}
	}
	return c.JSON(http.StatusNotImplemented, echo.Map{"error": "Not Implemented"})
}

// Handler for messages
func handleMessage(c echo.Context) error {
	m := new(Message)
	if err := c.Bind(m); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(m); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Create custom template
	if strings.HasPrefix(m.Content, "/customTemplate ") {
		tmplText := m.Content[strings.Index(m.Content, " ")+1:]
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*jwtCustomClaims)
		// Only admin is allowed!
		if !claims.Admin {
			return c.JSON(http.StatusOK, echo.Map{"name": botName, "content": "Sorry, only admin is allowed.."})
		}
		uuid := claims.Uuid
		ok, err := createTemplate(uuid, tmplText)
		if ok {
			return c.JSON(http.StatusOK, echo.Map{"name": botName, "content": "Template successfully created! Reload the page to use it."})
		} else {
			return c.JSON(http.StatusOK, echo.Map{"name": botName, "content": err.Error()})
		}
	}

	// Send random message
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(answers))
	res := new(Message)
	res.Content = answers[randomIndex]
	res.Name = botName
	return c.JSON(http.StatusOK, res)
}

func createTemplate(filename string, content string) (bool, error) {
	if !IsValidUUID(filename) {
		return false, errors.New("invalid UUID")
	}
	tmpl := template.New(filename)
	_, err := tmpl.Parse(content)
	if err != nil {
		return false, err
	}
	err = os.WriteFile("public/views/"+filename+".html", []byte(content), 0644)
	if err != nil {
		return false, err
	}

	// Update templates
	t = &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = t

	return true, nil
}

func main() {

	value, exists := os.LookupEnv("JWTKEY")
	if exists {
		secretKey = []byte(value)
	}

	e.Renderer = t

	f, err := os.OpenFile("logs/log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}
	defer f.Close()

	// Middlewares
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Output: f}))
	e.Use(middleware.Recover())

	// JWT middleware configuration. JWT token is extracted from the cookie "token".
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey:  secretKey,
		TokenLookup: "cookie:token",
	}

	// Static files
	e.Static("/static", "static")

	e.GET("/", login)
	e.POST("/", login)

	e.Validator = &CustomValidator{validator: validator.New()}

	// Protected routes
	r := e.Group("/chat")
	r.Use(echojwt.WithConfig(config))
	r.GET("/", chat)
	r.POST("/message", handleMessage)

	// Start server
	e.Logger.Fatal(e.Start(":8000"))
}
