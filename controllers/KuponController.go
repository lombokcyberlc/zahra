package controllers

import (
	"github.com/labstack/echo"
	"fmt"
	"net/http"
	"math/rand"

	uuid "github.com/google/uuid"
)

type KuponKursus struct {
	IDKupon		string  `form:"id_kupon" json:"id"`
	KodeKupon	string  `form:"kupon" json:"kupon"`
	KursusID	string 	`form:"kursus_id" json:"kursus_id"`
	UserID 		string  `form:"user_id" json:"user_id"`
	IDKelas		string 	`form:"idkelas" json:"idkelas"`
	StatusKupon int 	`form:"status_kupon" json:"status_kupon"`
	Jumlah 		int 	`form:"jumlah" json:"-"`
}

// Random String
func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func GetKuponByKursusID(c echo.Context) (err error) {

	// init user
	kupon := new(KuponKursus)
	if err = c.Bind(kupon); err != nil {
		return
	}

	// init db
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

    kursusId := c.Param("kursusID")

	// scan mysql_rows db query
	mysql_rows, err := db.Query("SELECT id_kupon as id, kupon, status_kupon, kursus_id FROM kupon_kursus WHERE kursus_id = ?" , kursusId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer mysql_rows.Close()

	// definisi variabel slice result
	var hasil []KuponKursus

	for mysql_rows.Next() {

		// init single kursus struct object
		var objekKuponKursus = KuponKursus{}
		var err = mysql_rows.Scan(&objekKuponKursus.IDKupon, &objekKuponKursus.KodeKupon,&objekKuponKursus.StatusKupon, &objekKuponKursus.KursusID)

		if err != nil {
			return err
		}

		// tambah data slice ke variable hasil
		hasil = append(hasil, objekKuponKursus)
	}

	if err = mysql_rows.Err(); err != nil {
		fmt.Println(err.Error())
		return
	}

	// kembalikan response json
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":     "berhasil",
		"dataKupon": hasil,
	})
}

// Verifikasi Kupon
func InputKupon(c echo.Context) (err error) {
	
	// Binding Kupon
	kupon := new(KuponKursus)
	if err = c.Bind(kupon); err != nil {
		return
	}
	
	kursusId := c.FormValue("kursus_id")
	kode := c.FormValue("kupon")
	userId := c.FormValue("user_id")
	
	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return err
    }
	defer db.Close()

	// var result = KuponKursus{}

	// // Check kupon dan id kursus
    err = db.QueryRow(`
				SELECT 
					id_kupon, 
					kupon, 
					kursus_id
				FROM 
					kupon_kursus 
				WHERE 
					kupon = ? 
				AND 
					kursus_id = ? 
				AND 
					status_kupon = 1 `, kode, kursusId).
	Scan(&kupon.IDKupon, &kupon.KodeKupon, &kupon.KursusID)
		
	if err != nil {
		fmt.Println("gagal cek kupon")
		fmt.Println(err.Error())
		return err
	}

	// // SELECT Before Insert Kupon
	// err = db.QueryRow(`
	// 			SELECT
	// 				idkelas,
	// 				user_id,
	// 				kursus_id
	// 			FROM 
	// 				kelas
	// 			WHERE
	// 				user_id = ?
	// 			AND
	// 				kursus_id = ?`, userId, kursusId).
	// Scan(&kupon.IDKelas, &kupon.UserID, &kupon.KursusID)
	// if err != nil {
	// 	fmt.Println("Gagal")
	// 	return err
	// } 

	// Define id dengan random string
		id := uuid.New()

		_, err = db.Exec("INSERT INTO kelas (idkelas, user_id, kursus_id) VALUES (?, ?, ?)", id.String(), userId, kursusId)
		
		if err != nil {
			fmt.Println("gagal insert kelas")
			fmt.Println(err.Error())
			return err
		}

		_, err = db.Exec("UPDATE kupon_kursus SET status_kupon = 0 WHERE kupon = ? AND kursus_id = ?", kode, kursusId)

		if err != nil {
			fmt.Println("gagal update kupon")
			fmt.Println(err.Error())
			return err
		}
		
	err = db.QueryRow("select idkelas, user_id, kursus_id FROM kelas order by idkelas desc limit 1").
	Scan(&kupon.IDKelas, &kupon.UserID, &kupon.KursusID)

	if err != nil {
		fmt.Println("gagal get data kelas")
		fmt.Println(err.Error())
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "berhasil",
		"kode" : kupon,
	})
	
}

// Generate Kupon
func GenerateKupon(c echo.Context) (err error) {

	kupon := new(KuponKursus)
	if err = c.Bind(kupon); err != nil {
		return
	}

	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return
    }

	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO kupon_kursus (id_kupon, kupon, kursus_id, status_kupon) VALUES (?, ?, ?, ?)")
	for i := 1; i <= kupon.Jumlah; i++ {
		kode := c.FormValue("kode_kursus")
		id := uuid.New()
		stmt.Exec(id.String(), kode+RandomString(5), &kupon.KursusID, 1)
	}
    
	if err != nil {
        fmt.Println(err.Error())
        return
    }
    
    rows, err := db.Query("SELECT id_kupon, kupon, kursus_id, status_kupon FROM kupon_kursus WHERE kursus_id = ? ", &kupon.KursusID)
    if err != nil {
        fmt.Println(err.Error())
        return
	}
	
	defer rows.Close()

	var hasil []KuponKursus
	for rows.Next() {
		var objek = KuponKursus{}
		var err = rows.Scan(&objek.IDKupon, &objek.KodeKupon, &objek.KursusID, &objek.StatusKupon)

		if err != nil {
			return err
		}

		hasil = append(hasil, objek)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
		return
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "berhasil",
		"dataKupon": hasil,
	})
}

// Update Kupon
// func UpdateKupon(c echo.Context) (err error) {
// 	kupon := new(KuponKursus)
// 	if err = c.Bind(kupon)
// }