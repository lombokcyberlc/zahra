package controllers

import (
	"github.com/labstack/echo"
	"fmt"
	"net/http"
	"strconv"
	"os"
	"io"
	"path/filepath"

)

type PromoKursus struct {
	ID			int `form:"id" json:"id"`
	NamaPromo	string `form:"nama_promo json:"nama_promo" validate:"required"`
	Gambar		string `form:"gambar" json:"gambar" validate:"required"`
}

// Menampilkan semua data Promo
func GetAllPromo(c echo.Context) (err error) {
	
	promo := new(PromoKursus)
	if err = c.Bind(promo); err != nil {
		return
	}
	
	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return
	}
	
    defer db.Close()

    rows, err := db.Query("SELECT id, nama_promo, gambar FROM promo")
    
	if err != nil {
        fmt.Println(err.Error())
        return
    }
    
	defer rows.Close()

    var result []PromoKursus

    for rows.Next() {

        var each = PromoKursus{}
        var err = rows.Scan(&each.ID, &each.NamaPromo, &each.Gambar)

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
		"data": result,
	})
}

// Menampilkan data Promo Kursus berdasarkan ID
func GetPromoById(c echo.Context) (err error) {
	
	promo := new(PromoKursus)

	if err = c.Bind(promo); err != nil {
		return
	}
	
	id := c.Param("id")
	
	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    defer db.Close()

    rows, err := db.Query("SELECT id, nama_promo, gambar FROM promo WHERE id = ?", id)
    
	if err != nil {
        fmt.Println(err.Error())
        return
    }
    
	defer rows.Close()

    var result PromoKursus

    for rows.Next() {
        var each = PromoKursus{}
        var err = rows.Scan(&each.ID, &each.NamaPromo, &each.Gambar)

        if err != nil {
			return err
        }

        result = each
    }

    if err = rows.Err(); err != nil {
        fmt.Println(err.Error())
        return
	}

	if strconv.Itoa(result.ID) == "" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "gagal",
			"data": result,
		})
	
	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "berhasil",
			"data": result,
		})
	}

}

// Tambah Promo
func TambahPromo(c echo.Context) (err error) {

	promo := new(PromoKursus)
	if err = c.Bind(promo); err != nil {
		return
	}

	namaPromo := c.FormValue("nama_promo")

	//** Start of File Upload
	//------------
	// Read files
	//------------

	// Multipart form
	// Get avatar
	imagePromo, err := c.FormFile("gambar")
	if err != nil {
		return err
	}

	// Source
	src, err := imagePromo.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	// ambil lokasi folder saat ini
	dir, err := os.Getwd()

	if err != nil {
		return err
	}

	// menentukan letak upload folder
	fileLocation := filepath.Join(dir, "/home/zahra/go/src/zcomeducation.com/zahra/upload/promo", imagePromo.Filename)
	//fileLocation := filepath.Join(dir, "/upload/foto", avatar.Filename)
	dst, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		return err
	}

	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	//** End of File Upload

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	_, err = db.Exec("INSERT INTO promo (nama_promo, gambar) VALUES (?, ?)", namaPromo, imagePromo.Filename)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = db.
		QueryRow("SELECT id, nama_promo, gambar FROM promo ORDER BY id DESC LIMIT 1").
		Scan(&promo.ID, &promo.NamaPromo, &promo.Gambar)

	if err != nil {
		// fmt.Println(err.Error())
		fmt.Println("gagal get data")
		return err
	} else {
		fmt.Println("berhasil get data")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":     "berhasil",
		"dataPromo": promo,
	})

}

// perbaharui Konten Promo
func UpdatePromo(c echo.Context) (err error) {

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer db.Close()

	promo := new(PromoKursus)
	if err = c.Bind(promo); err != nil {
		return err
	}

	// validasi inputan
	if err := c.Validate(promo); err != nil {
		return err
	}

	id := c.Param("id")
	namaPromo := c.FormValue("nama_promo")

	//** Start of File Upload
	//------------
	// Read files
	//------------

	// Multipart form
	// Get avatar
	imagePromo, err := c.FormFile("gambar")
	if err != nil {
		fmt.Println("tidak ada gambar")
		_, err = db.Exec("UPDATE promo SET nama_promo = ?, WHERE id = ?", namaPromo, id)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("errornya di sini")
			return err
		}

		err = db.
			QueryRow("SELECT id, nama_promo, gambar FROM promo WHERE id = ?", id).
			Scan(&promo.ID, &promo.NamaPromo, &promo.Gambar)

		if err != nil {
			// fmt.Println(err.Error())
			fmt.Println("gagal get data")
			return err
		} else {
			fmt.Println("berhasil get data")
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":     "berhasil",
			"dataPromo": promo,
		})

	} else {
		fmt.Println(imagePromo.Filename)
		fmt.Println(id)
		// hapus gambar lama
		// select by id
		err = db.
			QueryRow("SELECT id, nama_promo, gambar", id).
			Scan(&promo.ID, &promo.NamaPromo, &promo.Gambar)

		if err != nil {
			// fmt.Println(err.Error())
			fmt.Println("gagal get data")
			return err
		} else {
			fmt.Println("berhasil get data")
		}

		// ambil lokasi folder saat ini
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		// menentukan letak upload folder
		fileLocation := filepath.Join(dir, "/home/zahra/go/src/zcomeducation.com/zahra/upload/promo", promo.Gambar)
		//fileLocation := filepath.Join(dir, "/upload/promo", kursus.Gambar)

		// hapus gambar
		hapusFile(fileLocation)

		// kemudian perbaharui gambar
		// sumber upload
		src, err := imagePromo.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// tujuan upload
		// ambil lokasi folder saat ini
		dir, err = os.Getwd()
		if err != nil {
			return err
		}

		// menentukan letak upload folder
		fileLocation = filepath.Join(dir, "/home/zahra/go/src/zcomeducation.com/zahra/upload/promo", imagePromo.Filename)
		//fileLocation = filepath.Join(dir, "/upload/promo", imagePromo.Filename)
		dst, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		//** End of File Upload

		// perbaharui data
		_, err = db.Exec("UPDATE promo SET nama_promo = ?, gambar = ? WHERE id = ?", namaPromo, imagePromo.Filename, id)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		// select by id
		err = db.
			QueryRow("SELECT id, nama_promo, gambar FROM promo WHERE id = ?", id).
			Scan(&promo.ID, &promo.NamaPromo, &promo.Gambar)

		if err != nil {
			// fmt.Println(err.Error())
			fmt.Println("gagal get data")
			return err
		} else {
			fmt.Println("berhasil get data")
		}
	}

	if promo.ID == 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":     "gagal",
			"dataPromo": promo,
		})

	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":     "berhasil",
			"dataPromo": promo,
		})
	}

}

// Hapus Promo
func DeletePromo(c echo.Context) (err error) {

	// definisi variabel promo dengan struct PromoKursus
	promo := new(PromoKursus)

	// binding data inputan ke struct Promo
	if err = c.Bind(promo); err != nil {
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
	id := c.Param("id")

	// ambil gambar dari database sebelum dihapus
	err = db.
		QueryRow("SELECT id, nama_promo, gambar FROM promo WHERE id = ?", id).
		Scan(&promo.ID, &promo.NamaPromo, &promo.Gambar)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	//* Menghapus Gambar *//
	// ambil lokasi folder saat ini
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	// menentukan letak upload folder
	fileLocation := filepath.Join(dir, "/home/zahra/go/src/zcomeducation.com/zahra/upload/promo", promo.Gambar)
	//fileLocation := filepath.Join(dir, "/upload/promo", promo.Gambar)

	// hapus gambar
	hapusFile(fileLocation)
	//* End Of Menghapus Gambar *//

	// hapus data dari db
	_, err = db.Exec("DELETE FROM promo WHERE id = ?", id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":     "berhasil",
		"dataPromo": promo,
	})
}