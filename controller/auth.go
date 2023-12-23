package controller

import (
	// "encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"

	// "github.com/rs/xid"
	"github.com/shubhamxg/go-hunger/models"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/message"
)

type AuthHandler struct {
	db *sqlx.DB
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		db: models.Start(),
	}
}

func (handler *AuthHandler) SignUpHandler(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, recipe_response(backend_error))
		return
	}
	username := strings.ToLower(user.Username)
	hashed_bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Errorf("Failed to create user: %w", err)
	}

	password_hash := string(hashed_bytes)
	// user := models.User{
	// 	Username: username,
	// 	Password: password_hash,
	// }

	tx := handler.db.MustBegin()
	executed, err := handler.db.Exec(
		`INSERT INTO users (email, password) VALUES ($1, $2)`,
		username,
		password_hash,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, recipe_response(backend_error))
		return
	}

	if count, _ := executed.RowsAffected(); count == 0 {
		c.JSON(http.StatusNotFound, recipe_response(not_found))
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message": "User registered Successfully",
	})
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type JWTOutput struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

var JWT_SECRET = models.Env(models.JWT_SECRET)

func (handler *AuthHandler) SignInHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var creds []models.User
	if err := handler.db.Select(&creds,
		`SELECT id, password FROM users WHERE email=$1`,
		user.Username,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"db_error": err.Error(),
		})
		return
	}
	// foo, _ := json.Marshal(creds)
	// var bar models.User
	// json.Unmarshal(foo, &bar)

	if err := bcrypt.CompareHashAndPassword([]byte(creds[0].Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, recipe_response(invalid_creds))
		return
	}

	session_token := xid.New().String()
	session := sessions.Default(c)
	session.Set("username", user.Username)
	session.Set("token", session_token)
	session.Save()

	// if user.Username != "foo" || user.Password != "bar" {
	// 	c.JSON(http.StatusUnauthorized, recipe_response(invalid_creds))
	// 	return
	// }

	expiration_time := time.Now().Add(10 * time.Minute)
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration_time.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token_string, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		c.JSON(http.StatusInternalServerError, recipe_response(backend_error))
		return
	}

	jwt_output := JWTOutput{
		Token:   token_string,
		Expires: expiration_time,
	}
	c.JSON(http.StatusOK, jwt_output)
}

func (handler *AuthHandler) SignOutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "Signed out...",
	})
}

func (handler *AuthHandler) AuthMiddlerware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// token_value := c.GetHeader("Authorization")
		session := sessions.Default(c)
		session_token := session.Get("token")
		if session_token == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Not logged in",
			})
			c.Abort()
		}
		// c.Ne
		// claims := &Claims{}
		//
		// token, err := jwt.ParseWithClaims(
		// 	token_value,
		// 	claims,
		// 	func(t *jwt.Token) (interface{}, error) {
		// 		return []byte(JWT_SECRET), nil
		// 	},
		// )
		// if err != nil {
		// 	c.AbortWithStatus(http.StatusUnauthorized)
		// }
		//
		// if token == nil || !token.Valid {
		// 	c.AbortWithStatus(http.StatusUnauthorized)
		// }
		c.Next()
	}
}

func (handler *AuthHandler) RefreshHandler(c *gin.Context) {
	token_value := c.GetHeader("Authorization")
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(token_value, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	if token == nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, recipe_response(invalid_token))
		return
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		c.JSON(http.StatusBadRequest, recipe_response(not_expired))
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token_string, err := token.SignedString(JWT_SECRET)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	jwt_output := JWTOutput{
		Token:   token_string,
		Expires: expirationTime,
	}
	c.JSON(http.StatusOK, jwt_output)
}
