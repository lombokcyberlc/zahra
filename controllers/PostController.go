package controllers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/labstack/echo"

	uuid "github.com/google/uuid"
)

type ArtikelKursus struct {
	ID           string `form:"id" json:"id"`
	UserID       string `form:"user_id" json:"user_id"`
	NamaUser     string `form:"nama_lengkap" json:"nama_lengkap"`
	KategoriID   string `form:"kategori_id" json:"kategori_id"`
	NamaKategori string `form:"nama_kategori" json:"nama_kategori"`
	Judul        string `form:"judul" json:"judul" validate:"required"`
	Konten       string `form:"konten" json:"konten" validate:"required"`
	Gambar       string `form:"gambar" json:"gambar"`
	Slug         string `form:"slug" json:"slug"`
	KategoriSlug string `form:"kategori_slug" json:"kategori_slug"`
}

// Membuat Konten Kursus
func TambahArtikel(c echo.Context) (err error) {

	artikel := new(ArtikelKursus)
	if err = c.Bind(artikel); err != nil {
		return
	}

	tolower := strings.ToLower(artikel.Judul)
	// Make a Regex to say we only want letters and numbers
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	processedString := reg.ReplaceAllString(tolower, " ")

	minChar := processedString[:len(processedString)-1]
	slug := strings.ReplaceAll(minChar, " ", "-")

	//** Start of File Upload
	//------------
	// Read files
	//------------

	// Multipart form
	// Get avatar
	avatar, err := c.FormFile("gambar")
	if err != nil {
		return err
	}

	// Source
	src, err := avatar.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Generate Random File Name
	file, err := ioutil.TempFile("/var/www/vhosts/zcomeducation.com/httpdocs/upload/blog/", "blog-*.png")
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

	strFile := strings.ReplaceAll(file.Name(), "/var/www/vhosts/zcomeducation.com/httpdocs/upload/blog/", "")

	//** End of File Upload

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	id := uuid.New()

	_, err = db.Exec(`
				INSERT INTO 
					posts (
						id, 
						user_id,
						kategori_id, 
						judul, 
						konten,
						gambar, 
						slug
				) VALUES (?, ?, ?, ?, ?, ?, ?)`, id.String(), &artikel.UserID, &artikel.KategoriID, &artikel.Judul, &artikel.Konten, strFile, slug)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":      "berhasil",
		"dataArtikel": artikel,
	})

}

// Menampilkan data artikel
func GetArtikelKursus(c echo.Context) (err error) {

	// init user
	artikel := new(ArtikelKursus)
	if err = c.Bind(artikel); err != nil {
		return
	}

	// init db
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	// scan mysql_rows db query
	mysql_rows, err := db.Query(`
			SELECT 
				posts.id, 
				posts.user_id,
				users.nama_lengkap,
				posts.kategori_id, 
				post_kategori.nama_kategori,
				posts.judul, 
				posts.konten, 
				posts.gambar, 
				posts.slug 
			FROM 
				posts,
				users,
				post_kategori
			WHERE 
				posts.user_id = users.id 
			AND 
				posts.kategori_id = post_kategori.id`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer mysql_rows.Close()

	// definisi variabel slice result
	var hasil []ArtikelKursus

	for mysql_rows.Next() {

		// init single kursus struct object
		var each = ArtikelKursus{}
		var err = mysql_rows.Scan(&each.ID, &each.UserID, &each.NamaUser, &each.KategoriID, &each.NamaKategori, &each.Judul, &each.Konten, &each.Gambar, &each.Slug)

		if err != nil {
			return err
		}

		// tambah data slice ke variable hasil
		hasil = append(hasil, each)
	}

	if err = mysql_rows.Err(); err != nil {
		fmt.Println(err.Error())
		return
	}

	// kembalikan response json
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":      "berhasil",
		"dataArtikel": hasil,
	})
}

// Get data artikel by slug
func GetArtikelBySlug(c echo.Context) (err error) {

	// definisi variabel kursus dengan struct Kursus
	artikel := new(ArtikelKursus)

	// binding data inputan ke struct Kursus
	if err = c.Bind(artikel); err != nil {
		return err
	}

	// get id by parameter
	slug := c.Param("slug")

	// koneksi database
	db, err := connect()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// tutup koneksi database
	defer db.Close()

	err = db.QueryRow(`
					SELECT 
						posts.id, 
						posts.user_id,
						users.nama_lengkap,
						posts.kategori_id, 
						post_kategori.nama_kategori,
						posts.judul, 
						posts.konten, 
						posts.gambar, 
						posts.slug 
					FROM 
						posts,
						users,
						post_kategori
					WHERE 
						posts.user_id = users.id 
					AND 
						posts.kategori_id = post_kategori.id
					AND
						posts.slug = ?`, slug).Scan(&artikel.ID, &artikel.UserID, &artikel.NamaUser, &artikel.KategoriID, &artikel.NamaKategori, &artikel.Judul, &artikel.Konten, &artikel.Gambar, &artikel.Slug)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":      "berhasil",
		"dataArtikel": artikel,
	})

}

// Get data artikel by kategori slug
func GetArtikelByKategoriSlug(c echo.Context) (err error) {

	// init user
	artikel := new(ArtikelKursus)
	if err = c.Bind(artikel); err != nil {
		return
	}

	// init db
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	slug := c.Param("slug")

	// scan mysql_rows db query
	mysql_rows, err := db.Query(`
			SELECT 
				posts.id, 
				posts.user_id,
				users.nama_lengkap,
				posts.kategori_id, 
				post_kategori.nama_kategori,
				posts.judul, 
				posts.konten, 
				posts.gambar, 
				posts.slug,
				post_kategori.slug
			FROM 
				posts,
				users,
				post_kategori
			WHERE 
				posts.user_id = users.id 
			AND 
				posts.kategori_id = post_kategori.id
			AND 
				post_kategori.slug = ?`, slug)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer mysql_rows.Close()

	// definisi variabel slice result
	var hasil []ArtikelKursus

	for mysql_rows.Next() {

		// init single kursus struct object
		var each = ArtikelKursus{}
		var err = mysql_rows.Scan(&each.ID, &each.UserID, &each.NamaUser, &each.KategoriID, &each.NamaKategori, &each.Judul, &each.Konten, &each.Gambar, &each.Slug, &each.KategoriSlug)

		if err != nil {
			return err
		}

		// tambah data slice ke variable hasil
		hasil = append(hasil, each)
	}

	if err = mysql_rows.Err(); err != nil {
		fmt.Println(err.Error())
		return
	}

	// kembalikan response json
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":      "berhasil",
		"dataArtikel": hasil,
	})
}

// perbaharui konten kursus
func UpdateArtikel(c echo.Context) (err error) {

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer db.Close()

	artikel := new(ArtikelKursus)

	if err = c.Bind(artikel); err != nil {
		return err
	}

	// validasi inputan
	if err := c.Validate(artikel); err != nil {
		return err
	}

	tolower := strings.ToLower(artikel.Judul)
	newSlug := strings.ReplaceAll(tolower, " ", "-")

	// Menggunakan slug sebagai parameter
	slug := c.Param("slug")
	fmt.Println("slug param ", slug)

	//** Start of File Upload
	//------------
	// Read files
	//------------

	// Multipart form
	// Get avatar
	avatar, err := c.FormFile("gambar")

	if err != nil {
		fmt.Println("tidak ada gambar")

		_, err = db.Exec(`
					UPDATE 
						posts 
					SET 
						user_id = ?,
						kategori_id = ?, 
						judul = ?, 
						konten = ?, 
						slug = ? 
					WHERE 
						slug = ?`, &artikel.UserID, &artikel.KategoriID, &artikel.Judul, &artikel.Konten, newSlug, slug)

		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("gagal update")
			return err
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "berhasil",
		})

	} else {
		/* ============================================================
		| Select gambar dari database berdasarkan slug sebelum diupdate |
		   ===========================================================*/
		db.
			QueryRow(`
					SELECT 
						gambar, 
					FROM 
						posts,
					WHERE
						slug = ?`, slug).
			Scan(&artikel.Gambar)

		/* ============================================================
		|  Menentukan letak file upload   						      |
		   ===========================================================*/
		fileLocation := filepath.Join("/var/www/vhosts/zcomeducation.com/httpdocs/upload/blog/", artikel.Gambar)

		/* ============================================================
		|  Hapus gambar lama 						  				  |
		   ===========================================================*/

		hapusFile(fileLocation)

		/* ============================================================
		|  Tentukan sumber upload dari folder temporary 			  |
		   ===========================================================*/

		src, err := avatar.Open()
		if err != nil {
			fmt.Println("gagal buka gambar dari src")
			return err
		}
		defer src.Close()

		/* ============================================================
		|  Update database 						  					   |
		   ===========================================================*/

		// Generate Random File name
		file, err := ioutil.TempFile("/var/www/vhosts/zcomeducation.com/httpdocs/upload/blog/", "blog-*.png")
		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		fileBytes, err := ioutil.ReadAll(src)
		if err != nil {
			fmt.Println("gagal membaca file dari src")
			fmt.Println(err)
		}

		file.Write(fileBytes)

		os.Chown(file.Name(), 10000, 1004)
		os.Chmod(file.Name(), 0644)

		strFile := strings.ReplaceAll(file.Name(), "/var/www/vhosts/zcomeducation.com/httpdocs/upload/blog/", "")
		//** End of File Upload

		// perbaharui database
		_, err = db.Exec(`
						UPDATE 
							posts 
						SET 
							user_id = ?,
							kategori_id = ?, 
							judul = ?, 
							konten = ?, 
							gambar = ?, 
							slug = ? 
						WHERE 
							slug = ?`, &artikel.UserID, &artikel.KategoriID, &artikel.Judul, &artikel.Konten, strFile, newSlug, slug)

		if err != nil {
			fmt.Println("gagal update database dengan gambar")
			fmt.Println(err.Error())
			return err
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "berhasil",
		})
	}
}

// Hapus Kursus
func DeleteArtikel(c echo.Context) (err error) {

	// definisi variabel kursus dengan struct Kursus
	artikel := new(ArtikelKursus)

	// binding data inputan ke struct Kursus
	if err = c.Bind(artikel); err != nil {
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
	slug := c.Param("slug")

	// ambil gambar dari database sebelum dihapus
	err = db.
		QueryRow(`SELECT 
					id, 
					kategori_id, 
					judul, 
					konten, 
					gambar, 
					slug 
				FROM 
					posts 
				WHERE 
					slug = ?`, slug).
		Scan(&artikel.ID, &artikel.KategoriID, &artikel.Judul, &artikel.Konten, &artikel.Gambar, &artikel.Slug)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// menentukan letak upload folder
	fileLocation := filepath.Join("/var/www/vhosts/zcomeducation.com/httpdocs/upload/blog/", artikel.Gambar)

	// hapus gambar
	hapusFile(fileLocation)
	//* End Of Menghapus Gambar *//

	// hapus data dari db
	_, err = db.Exec("DELETE FROM posts WHERE slug = ?", slug)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":      "berhasil",
		"dataArtikel": artikel,
	})
}
