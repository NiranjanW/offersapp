package models

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

var (
	tokenSecret = []byte(os.Getenv("TOKEN_SECRET"))
)

type User struct {
	ID              uuid.UUID `json:"id"`
	CreatedAt       time.Time `json:"_"`
	UpdatedAt       time.Time `json:"_"`
	Email           string    `json:"email"`
	PasswordHash    string    `json:"-"`
	Password        string    `json:"password"`
	PasswordConfirm string    `json:"password_confirm"`
}

func (u *User) Register(conn *pgx.Conn) error {
	if len(u.Password) < 4 || len(u.PasswordConfirm) < 4 {
		return fmt.Errorf("Password must be aleast 4 char long")

	}

	if u.Password != u.PasswordConfirm {
		return fmt.Errorf("Passwords do not match")
	}
	return nil
}

func (u *User) GetAuthToken() (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = u.ID
	claims["expires"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodES256, &claims)
	authtoken, err := token.SignedString(tokenSecret)
	return authtoken, err
}

func (u *User) isAuthenticated(conn *pgx.Conn) error {
	row := conn.QueryRow(context.Background(), "SELECT id, password_hash from user_account WHERE email = $1", u.Email)
	err := row.Scan(&u.ID, &u.PasswordHash)

	if err == pgx.ErrNoRows {
		fmt.Println("User with email not found")
		return fmt.Errorf("Invalid login credentials")
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))
	if err != nil {
		return fmt.Errorf("Invalid login credentials")
	}

	return nil
}

func isTOkenValid(tokenString string) (bool, string) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// fmt.Printf("Parsing: %v \n", token)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok == false {
			return nil, fmt.Errorf("Token signing method is not valid: %v", token.Header["alg"])
		}

		return tokenSecret, nil
	})
	if err != nil {
		fmt.Printf("Err %v \n", err)
		return false, ""
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// fmt.Println(claims)
		userID := claims["user_id"]
		return true, userID.(string)
	} else {
		fmt.Printf("The alg header %v \n", claims["alg"])
		fmt.Println(err)
		return false, "uuid.UUID{}"
	}
}
