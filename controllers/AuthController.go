package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"zahra/config"

	_ "github.com/go-sql-driver/mysql"
)

// User struct model
type User struct {
	ID          	string  `json:"id" 	form:"id"`
	GroupID     	int     `json:"group_id" form:"group_id"`
	NamaLengkap 	string  `json:"nama_lengkap" form:"nama_lengkap" query:"nama_lengkap"`
	Email       	string  `json:"email" form:"email" query:"email" validate:"required,email"`
	Password    	string  `json:"-" form:"password" query:"password"`
	Token       	int     `json:"token" form:"token" query:"token"`
	Verified    	int     `json:"verifikasi" form:"verifikasi" query:"verifikasi"`
}

// CustomValidator
type CustomValidator struct {
	validator *validator.Validate
}

var jwtConfig = config.Config.JWT

// validate methode
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "admin_zahra_user:zahra_password1945@tcp(zcomeducation.com:3306)/zahra_web")
	if err != nil {
		return nil, err
	} 
	
	return db, nil
}

// post login
func PostLogin(c echo.Context) (err error) {

	user := new(User)
	if err = c.Bind(user); err != nil {
		return
	}

	// validasi inputan
	if err := c.Validate(user); err != nil {
		return err
	}

	email := fmt.Sprintf("%s", user.Email)
	password := fmt.Sprintf("%s", user.Password)

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	rows, err := db.Query("select id, group_id, nama_lengkap, email, password, token, verifikasi from users where email = ? AND verifikasi = 1", email)

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"status": "gagal",
			"pesan":  "Akun tidak ditemukan...!",
		})
	}

	defer rows.Close()

	var result User

	for rows.Next() {
		var each = User{}
		var err = rows.Scan(&each.ID, &each.GroupID, &each.NamaLengkap, &each.Email, &each.Password, &user.Token, &user.Verified)

		if err != nil {
			return err
		}

		result = each
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
		return
	}

	if user.Verified == 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "gagal",
			"pesan":  "Akun belum diaktivasi !",
		})
	}

	// cocokkan password dengan hash bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password))

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"status": "gagal",
			"pesan":  "Password salah...!",
		})
	}

	// generate token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set data claims ke token berupa data id, email dan waktu kadaluarsa
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = result.ID
	claims["email"] = result.Email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	generatedToken, err := token.SignedString([]byte(jwtConfig.Secret))

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"status": "gagal",
			"pesan":  err,
		})
	}

	err = db.QueryRow("SELECT id, group_id, nama_lengkap, email, password, token, verifikasi FROM users WHERE email = ? AND verifikasi = 1", email).Scan(&user.ID, &user.GroupID, &user.NamaLengkap, &user.Email, &user.Password, &user.Token, &user.Verified)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "berhasil",
		"pesan":  "Selamat datang...!",
		"data":   user,
		"token":  generatedToken,
	})
}
