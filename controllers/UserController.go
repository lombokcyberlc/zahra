package controllers

import (
	"github.com/labstack/echo"
	"fmt"
	"net/http"
    "strconv"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
	"crypto/tls"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"io/ioutil"
	"strings"

	uuid "github.com/google/uuid"
)

type UserProfile struct {
	ID          	string  `json:"id" form:"id"`
	NamaLengkap 	string  `json:"nama_lengkap" form:"nama_lengkap" query:"nama_lengkap"`
	Email       	string  `json:"email" form:"email" query:"email" validate:"required,email"`
	NomorHP			string	`json:"hp" form:"hp" query:"hp" validate:"required"`
	Alamat			string	`json:"alamat" form:"alamat" query:"alamat"`
	JenisKelamin	string	`json:"jenis_kelamin" form:"jenis_kelamin" query:"jenis_kelamin"`
	StatusPekerjaan	string	`json:"status_pekerjaan" form:"status_pekerjaan" query:"status_pekerjaan"`
	TempatKerja		string	`json:"tempat_kerja" form:"tempat_kerja" query:"tempat_kerja"`
	Foto			string 	`json:"foto" form:"foto"`
}


// Returns an int >= min, < max
func randomInt(min, max int) int {
    return min + rand.Intn(max-min)
}

// Menampilkan data semua user
func GetAllUsers(c echo.Context) (err error) {
	
	u := new(User)
	if err = c.Bind(u); err != nil {
		return
	}
	
	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    defer db.Close()

    rows, err := db.Query("select id, group_id, nama_lengkap, email, token, verifikasi from users")
    
	if err != nil {
        fmt.Println(err.Error())
        return
    }
    
	defer rows.Close()

    var result []User

    for rows.Next() {
        var each = User{}
        var err = rows.Scan(&each.ID, &each.GroupID, &each.NamaLengkap ,&each.Email, &each.Token, &each.Verified)

        if err != nil {
			return err
        }

        result = append(result, each)
    }

    if err = rows.Err(); err != nil {
        fmt.Println(err.Error())
        return
    }
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "berhasil",
		"dataUser": result,
	})
}

// Menampilkan data user berdasarkan ID user
func GetUsersById(c echo.Context) (err error) {
	
	user := new(User)
	if err = c.Bind(user); err != nil {
		return
	}
	
	id := c.Param("id")
	
	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    defer db.Close()

    rows, err := db.Query("select users.id, users.group_id, users.nama_lengkap, users.email, users.token, users.verifikasi from users JOIN user_group ON users.group_id = user_group.id where users.id = ?", id)
    
	if err != nil {
        fmt.Println(err.Error())
        return
    }
    
	defer rows.Close()

    var result User

    for rows.Next() {
        var each = User{}
        var err = rows.Scan(&each.ID, &each.GroupID, &each.NamaLengkap ,&each.Email, &each.Token, &each.Verified)

        if err != nil {
			return err
        }

        result = each
    }

    if err = rows.Err(); err != nil {
        fmt.Println(err.Error())
        return
    }
		
	if result.ID != "" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "gagal",
			"dataUser": result,
		})
	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "berhasil",
			"dataUser": result,
		})
	}

}

// Membuat User baru
func Register(c echo.Context) (err error) {
	
	user := new(User)
	if err = c.Bind(user); err != nil {
		return
	}
	// Define id sebagai random string
	id := uuid.New()
	
	// Generate Token
	rand.Seed(time.Now().UnixNano())
	token := randomInt(1000, 4000)

    hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	
	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return
    }

	defer db.Close()

    err = db.QueryRow("SELECT email FROM users WHERE email = ?", &user.Email).Scan(&user.Email)
    if err != nil {
        _, err = db.Exec("INSERT INTO users (id, group_id, nama_lengkap, email, password, token, verifikasi) VALUES (?, ?, ?, ?, ?, ?, ?)", id.String(), 4, &user.NamaLengkap, &user.Email, hashPassword, token, 0)
    
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		
		err = db.QueryRow("SELECT id, group_id, nama_lengkap, email, password, token, verifikasi FROM users ORDER BY id DESC limit 1").Scan(&user.ID, &user.GroupID, &user.NamaLengkap, &user.Email, &user.Password, &user.Token, &user.Verified)
		if err != nil {
			fmt.Println(err.Error())
			return
		} else {
			from := mail.Address{"", "Zahra Computer Education"}
			to := mail.Address{"", user.Email}
			subj := "Aktivasi Akun"
			body := "=> Klik Link untuk Aktivasi Akun : http://zcomeducation.com/aktivasi/"+user.Email+"/"+strconv.Itoa(user.Token)+"/"+user.ID+" <="
		
			// Setup headers
			headers := make(map[string]string)
			headers["From"] = from.String()
			headers["To"] = to.String()
			headers["Subject"] = subj
		
			// Setup message
			message := ""
			for k, v := range headers {
				message += fmt.Sprintf("%s: %s\r\n", k, v)
			}
			message += "\r\n" + body
		
			// Connect to the SMTP Server
			servername := "smtp.googlemail.com:465"
		
			host, _, _ := net.SplitHostPort(servername)
		
			auth := smtp.PlainAuth("", "your.email.aja@gmail.com", "Makannasi", host)
		
			// TLS config
			tlsconfig := &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         host,
			}
		
			// Here is the key, you need to call tls.Dial instead of smtp.Dial
			// for smtp servers running on 465 that require an ssl connection
			// from the very beginning (no starttls)
			conn, err := tls.Dial("tcp", servername, tlsconfig)
			if err != nil {
				log.Panic(err)
			}
		
			c, err := smtp.NewClient(conn, host)
			if err != nil {
				log.Panic(err)
			}
		
			// Auth
			if err = c.Auth(auth); err != nil {
				log.Panic(err)
			}
		
			// To && From
			if err = c.Mail(from.Address); err != nil {
				log.Panic(err)
			}
		
			if err = c.Rcpt(to.Address); err != nil {
				log.Panic(err)
			}
		
			// Data
			w, err := c.Data()
			if err != nil {
				log.Panic(err)
			}
		
			_, err = w.Write([]byte(message))
			if err != nil {
				log.Panic(err)
			}
		
			err = w.Close()
			if err != nil {
				log.Panic(err)
			}
		
			c.Quit()
		}	
	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":false,
			"pesan":"Email sudah digunakan !",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"pesan":"Akun berhasil dibuat ! Silahkan Aktivasi.",
		"data": user,
	})
}

// Verifikasi Email
func VerifikasiEmail(c echo.Context) (err error) {
	user := new(User)
	if err = c.Bind(user); err != nil {
		return
	}

	id := c.Param("id")

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	_, err = db.Exec("UPDATE users SET verifikasi = 1 WHERE email = ? AND token = ? AND id = ?", &user.Email, &user.Token, id)
    
	if err != nil {
        fmt.Println(err.Error())
        return
	}
	

	rows, err := db.Query("select users.id, users.group_id, users.nama_lengkap, users.email, users.password, users.token, users.verifikasi from users JOIN user_group ON users.group_id = user_group.id where users.id = ?", id)
    
	if err != nil {
        fmt.Println(err.Error())
        return
    }
    
	defer rows.Close()

    var result User

    for rows.Next() {
        var each = User{}
        var err = rows.Scan(&each.ID, &each.GroupID, &each.NamaLengkap ,&each.Email, &each.Password, &each.Token, &each.Verified)

        if err != nil {
			return err
        }

        result = each
    }

    if err = rows.Err(); err != nil {
        fmt.Println(err.Error())
        return
    }
		
	if result.ID != "" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "Verifikasi Gagal !",
			// "dataUser": result,
		})
	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "Verifikasi Berhasil !",
			// "dataUser": result,
		})
	}
}

// Menampilkan User berdasarkan group id
func GetUsersByGroupId(c echo.Context) (err error) {
	
	user := new(User)
	if err = c.Bind(user); err != nil {
		return
	}
	
	id := c.Param("id")
	
	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    defer db.Close()

    rows, err := db.Query("select users.id, users.group_id, users.nama_lengkap, users.email, users.token, users.verifikasi from users JOIN user_group ON users.group_id = user_group.id where users.group_id = ?", id)
    
	if err != nil {
        fmt.Println(err.Error())
        return
    }
    
	defer rows.Close()

    var result []User

    for rows.Next() {
        var each = User{}
        var err = rows.Scan(&each.ID, &each.GroupID, &each.NamaLengkap ,&each.Email ,&each.Token, &each.Verified)

        if err != nil {
			return err
        }

        result = append(result, each)
    }

    if err = rows.Err(); err != nil {
        fmt.Println(err.Error())
        return
    }
		
	if strconv.Itoa(user.GroupID) == "" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "gagal",
			"dataUser": result,
		})
	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "berhasil",
			"dataUser": result,
		})
	}

}

// Hapus User
func DeleteUser(c echo.Context) (err error) {

	// definisi variabel User dengan struct User
	users := new(User)

	// binding data inputan ke struct User
	if err = c.Bind(users); err != nil {
		return err
	}

	// koneksi database
	db, err := connect()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// tutup koneksi database
	defer db.Close()

	// get id by parameter
	idUser := c.Param("id")

	err = db.
		QueryRow("SELECT users.id, users.group_id, users.nama_lengkap, users.email, users.password, users.token, users.verifikasi FROM users WHERE users.id = ?", idUser).
		Scan(&users.ID, &users.GroupID, &users.NamaLengkap, &users.Email, &users.Password, &users.Token, &users.Verified)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// hapus data dari db
	_, err = db.Exec("DELETE FROM users WHERE id = ?", idUser)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "berhasil",
		"dataUser": users,
    })
}

// Update Soal
func UpdateUser(c echo.Context) (err error) {

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer db.Close()

	user := new(User)
	if err = c.Bind(user); err != nil {
		return
	}

	// validasi inputan
	if err := c.Validate(user); err != nil {
		return err
	}

	id := c.Param("id")

	rand.Seed(time.Now().UnixNano())
	token := randomInt(1000, 4000)

    hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	// Update Query	
	_, err = db.Exec("UPDATE users SET group_id = ?, nama_lengkap = ?, email = ?, password = ?, token = ?, verifikasi = ? WHERE id = ?", &user.GroupID, &user.NamaLengkap, &user.Email, hashPassword, token, 0, id)
	
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("errornya di sini")
			return err
		}
	
	// Select Query setelah di update
	err = db.QueryRow("SELECT id, group_id, nama_lengkap, email, password, token, verifikasi FROM users WHERE id = ?", id).Scan(&user.ID, &user.GroupID, &user.NamaLengkap, &user.Email, &user.Password, &user.Token, &user.Verified)

	if err != nil {
		// fmt.Println(err.Error())
		fmt.Println("gagal get data")
		return err
	} else {
		fmt.Println("berhasil get data")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "berhasil",
		"dataSoal": user,
	})
}

// Membuat User baru
func CreateUser(c echo.Context) (err error) {
	
	user := new(User)
	if err = c.Bind(user); err != nil {
		return
    }
	// Define id sebagai random string
	id := uuid.New()

	// Enkripsi Password
    hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	
	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    defer db.Close()

    _, err = db.Exec("INSERT INTO users (id, group_id, nama_lengkap, email, password) VALUES (?, ?, ?, ?, ?)", id.String(), &user.GroupID, &user.NamaLengkap, &user.Email, hashPassword)
    
	if err != nil {
        fmt.Println(err.Error())
        return
    }
    
    err = db.QueryRow("SELECT id, group_id, nama_lengkap, email, password FROM users ORDER BY id DESC limit 1").Scan(&user.ID, &user.GroupID, &user.NamaLengkap, &user.Email, &user.Password)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
		
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "berhasil",
		"data": user,
	})
}

func GetUserProfile(c echo.Context) (err error) {
	profile := new(UserProfile)
	if err = c.Bind(profile); err != nil {
		return
	}

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	rows, err := db.Query(`
				SELECT 
					id,
					nama_lengkap,
					email,
					hp,
					alamat,
					jenis_kelamin,
					status_pekerjaan,
					tempat_kerja,
					foto
				FROM users ORDER BY created_at ASC	
	`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer rows.Close()

	var userprof []UserProfile

	for rows.Next() {
		var user = UserProfile{}
		var err = rows.Scan(&user.ID, &user.NamaLengkap, &user.Email, &user.NomorHP, &user.Alamat, &user.JenisKelamin, &user.StatusPekerjaan, &user.TempatKerja, &user.Foto)
		if err != nil {
			return err
		}

		userprof = append(userprof, user)
	}

	if err = rows.Err(); err != nil {
        fmt.Println(err.Error())
        return
    }

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":"berhasil",
		"dataProfile": userprof,
	})
}

func GetProfileById(c echo.Context) (err error) {
	
	// Define new profile
	profile := new(UserProfile)
	if err = c.Bind(profile); err != nil {
		return err
	}

	// Konek database
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	// menggunakan parameter
	id := c.Param("id")

	err = db.QueryRow(`
				SELECT
					id,
					nama_lengkap,
					email,
					hp,
					alamat,
					jenis_kelamin,
					status_pekerjaan,
					tempat_kerja,
					foto
				FROM 
					users
				WHERE id = ?
	`, id).
	Scan(&profile.ID, &profile.NamaLengkap, &profile.Email, &profile.NomorHP, &profile.Alamat, &profile.JenisKelamin, &profile.StatusPekerjaan, &profile.TempatKerja, &profile.Foto)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":"berhasil",
		"dataProfile": profile,
	})
}

func UpdateProfile(c echo.Context) (err error) {
	
	// Define new profile
	profile := new(UserProfile)
	if err = c.Bind(profile); err != nil {
		return err
	}

	// Konek database
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	id := c.Param("id")

	// INSERT PROCESS
	/*============================================================
	  |  Tentukan sumber upload dari folder temporary 			 |
	  ============================================================*/
	
	avatar, err := c.FormFile("foto")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
		
	src, err := avatar.Open()
	if err != nil {
		fmt.Println("gagal buka gambar src")
		return err
	}
	defer src.Close()

	// Generate Random File Name
	file, err := ioutil.TempFile("/var/www/vhosts/zcomeducation.com/httpdocs/upload/profile", "profile-*.png")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fileBytes, err := ioutil.ReadAll(src)
	if err != nil {
		fmt.Println(err)
	}
		
	file.Write(fileBytes)

	os.Chown(file.Name(), 10000, 1004)
	os.Chmod(file.Name(), 0644)

	strFile := strings.ReplaceAll(file.Name(), "/var/www/vhosts/zcomeducation.com/httpdocs/upload/profile/", "")

	/* INSERT DATA PROFILE */
	_, err = db.Exec(`
				UPDATE users SET 
					hp = ?,
					alamat = ?,
					jenis_kelamin = ?,
					status_pekerjaan = ?,
					tempat_kerja = ?,
					foto = ?
				WHERE id = ?
	`, &profile.NomorHP, &profile.Alamat, &profile.JenisKelamin, &profile.StatusPekerjaan, &profile.TempatKerja, strFile, id)
	if err != nil{
		fmt.Println(err.Error())
		return
	}

	//** SELECT DATA AFTER INSERT
	err = db.QueryRow(`
				SELECT
					id,
					nama_lengkap,
					email,
					hp,
					alamat,
					jenis_kelamin,
					status_pekerjaan,
					tempat_kerja,
					foto
				FROM 
					users
				WHERE id = ?
	`, id).
	Scan(&profile.ID, &profile.NamaLengkap, &profile.Email, &profile.NomorHP, &profile.Alamat, &profile.JenisKelamin, &profile.StatusPekerjaan, &profile.TempatKerja, &profile.Foto)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//** END OF SELECTED DATA

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":"berhasil",
		"dataProfile": profile,
	})
}
