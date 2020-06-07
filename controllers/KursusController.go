package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	uuid "github.com/google/uuid"

	"log"
	// "io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// menambahkan validator (validasi input) fungsi ada di AuthController
type Kursus struct {
	IdKursus         string `form:"id" json:"id" query:"id"`
	UserId           string `form:"user_id" json:"user_id" validate:"required"`
	NamaInstruktur   string `form:"nama_instruktur" json:"nama_instruktur"`
	KursusKategoriId string `form:"kursus_kategori_id" json:"kursus_kategori_id" validate:"required"`
	NamaKursus       string `form:"nama_kursus" json:"nama_kursus" query:"nama_kursus" validate:"required"`
	Deskripsi        string `form:"deskripsi" json:"deskripsi" query:"deskripsi" validate:"required"`
	Gambar           string `form:"gambar" json:"gambar"`
	Harga 			 int 	`form:"harga" json:"harga"`
	Diskon			 int 	`form:"harga_diskon" json:"harga_diskon"`
	Slug 			 string `form:"slug" json:"slug"`
	SlugKategori	 string `form:"slug_kategori" json:"slug_kategori"`
}

// fungsi hapus file
func hapusFile(lokasi string) error {
	path := lokasi
	err := os.Remove(path)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// Menampilkan data kursus
func GetAllKursus(c echo.Context) (err error) {

	// init user
	kursus := new(Kursus)
	if err = c.Bind(kursus); err != nil {
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
	mysql_rows, err := db.Query(`SELECT 
									kursus_konten.id, 
									kursus_konten.user_id, 
									kursus_konten.kursus_kategori_id, 
									kursus_konten.nama_kursus, 
									kursus_konten.deskripsi, 
									kursus_konten.gambar, 
									kursus_konten.harga, 
									kursus_konten.harga_diskon, 
									kursus_konten.slug, 
									kursus_kategori.slug, 
									users.nama_lengkap 
								FROM 
									kursus_konten, 
									users, 
									kursus_kategori 
								WHERE
									kursus_konten.user_id = users.id 
								AND 
									kursus_konten.kursus_kategori_id = kursus_kategori.id
								ORDER BY kursus_konten.id ASC`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	
	defer mysql_rows.Close()

	// definisi variabel slice result
	var hasil []Kursus

	for mysql_rows.Next() {

		// init single kursus struct object
		var objekKursus = Kursus{}
		var err = mysql_rows.Scan(&objekKursus.IdKursus, &objekKursus.UserId, &objekKursus.KursusKategoriId, &objekKursus.NamaKursus, &objekKursus.Deskripsi, &objekKursus.Gambar, &objekKursus.Harga, &objekKursus.Diskon, &objekKursus.Slug, &objekKursus.SlugKategori, &objekKursus.NamaInstruktur)

		if err != nil {
			return err
		}

		// tambah data slice ke variable hasil
		hasil = append(hasil, objekKursus)
	}

	if err = mysql_rows.Err(); err != nil {
		fmt.Println(err.Error())
		return
	}

	// kembalikan response json
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":     "berhasil",
		"dataKursus": hasil,
	})
}

// Menampilkan data Kursus Konten berdasarkan ID
func GetKursusById(c echo.Context) (err error) {

	// definisi variabel kursus dengan struct Kursus
	kursus := new(Kursus)

	// binding data inputan ke struct Kursus
	if err = c.Bind(kursus); err != nil {
		return err
	}

	// get id by parameter
	idKursus := c.Param("id")

	// koneksi database
	db, err := connect()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// tutup koneksi database
	defer db.Close()

	err = db.QueryRow(`SELECT
					kursus_konten.id, 
					kursus_konten.user_id, 
					kursus_konten.kursus_kategori_id, 
					kursus_konten.nama_kursus, 
					kursus_konten.deskripsi, 
					kursus_konten.gambar, 
					kursus_konten.harga, 
					kursus_konten.harga_diskon, 
					kursus_konten.slug, 
					users.nama_lengkap
		 		 FROM 
					kursus_konten, 
					users 
				 WHERE 
					kursus_konten.user_id = users.id 
				 AND 
					kursus_konten.id = ?`, idKursus).
	Scan(&kursus.IdKursus, &kursus.UserId, &kursus.KursusKategoriId, &kursus.NamaKursus, &kursus.Deskripsi, &kursus.Gambar, &kursus.Harga, &kursus.Diskon, &kursus.Slug, &kursus.NamaInstruktur)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if kursus.IdKursus != "" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":     "gagal",
			"dataKursus": kursus,
		})

	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":     "berhasil",
			"dataKursus": kursus,
		})
	}

}

// Get Kursus by Kategori ID
func GetKursusByKategoriSlug(c echo.Context) (err error) {

	// definisi variabel kursus dengan struct Kursus
	kursus := new(Kursus)

	// binding data inputan ke struct Kursus
	if err = c.Bind(kursus); err != nil {
		return err
	}

	// get id by parameter
	slugKategori := c.Param("slug")

	// koneksi database
	db, err := connect()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// tutup koneksi database
	defer db.Close()

	rows, err := db.
		Query(`SELECT 
					kursus_konten.id, 
					kursus_konten.user_id, 
					kursus_konten.kursus_kategori_id, 
					kursus_konten.nama_kursus, 
					kursus_konten.deskripsi, 
					kursus_konten.gambar, 
					kursus_konten.harga, 
					kursus_konten.harga_diskon, 
					kursus_konten.slug, 
					kursus_kategori.slug, 
					users.nama_lengkap 
				FROM 
					kursus_konten, 
					users, 
					kursus_kategori 
				WHERE 
					kursus_konten.user_id = users.id 
				AND 
					kursus_konten.kursus_kategori_id = kursus_kategori.id 
				AND 
					kursus_kategori.slug = ?`, slugKategori)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		
		defer rows.Close()
	
		var result []Kursus
	
		for rows.Next() {
			var objekKursus = Kursus{}
			var err = rows.Scan(&objekKursus.IdKursus, &objekKursus.UserId, &objekKursus.KursusKategoriId, &objekKursus.NamaKursus, &objekKursus.Deskripsi, &objekKursus.Gambar, &objekKursus.Harga, &objekKursus.Diskon, &objekKursus.Slug, &objekKursus.SlugKategori, &objekKursus.NamaInstruktur)
	
			if err != nil {
				return err
			}
	
			result = append(result, objekKursus)
		}
	

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":     "berhasil",
		"dataKursus": result,
	})

}

// Get Kursus by User ID
func GetKursusByUserId(c echo.Context) (err error) {

	// definisi variabel kursus dengan struct Kursus
	kursus := new(Kursus)

	// binding data inputan ke struct Kursus
	if err = c.Bind(kursus); err != nil {
		return err
	}

	// get id by parameter
	idUser := c.Param("id")

	// koneksi database
	db, err := connect()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// tutup koneksi database
	defer db.Close()

	rows, err := db.
		Query("SELECT kursus_konten.id, kursus_konten.nama_kursus, kursus_konten.deskripsi, kursus_konten.gambar, kursus_konten.harga, kursus_konten.harga_diskon, kursus_konten.slug, users.nama_lengkap FROM kursus_konten, users WHERE kursus_konten.user_id = users.id AND kursus_konten.user_id = ?", idUser)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		
		defer rows.Close()
	
		var result []Kursus
	
		for rows.Next() {
			var objekKursus = Kursus{}
			var err = rows.Scan(&objekKursus.IdKursus, &objekKursus.NamaKursus, &objekKursus.Deskripsi, &objekKursus.Gambar, &objekKursus.Harga, &objekKursus.Diskon, &objekKursus.Slug, &objekKursus.NamaInstruktur)
	
			if err != nil {
				return err
			}
	
			result = append(result, objekKursus)
		}
	

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":     "berhasil",
		"dataKursus": result,
	})

}

// Menampilkan data Kursus Konten berdasarkan ID
func GetKursusBySlug(c echo.Context) (err error) {

	// definisi variabel kursus dengan struct Kursus
	kursus := new(Kursus)

	// binding data inputan ke struct Kursus
	if err = c.Bind(kursus); err != nil {
		return err
	}

	// get nama parameter
	slug := c.Param("slug")


	// koneksi database
	db, err := connect()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// tutup koneksi database
	defer db.Close()

	err = db.
		QueryRow(`SELECT 
						kursus_konten.id, 
						kursus_konten.user_id, 
						kursus_konten.kursus_kategori_id, 
						kursus_konten.nama_kursus, 
						kursus_konten.deskripsi, 
						kursus_konten.gambar, 
						kursus_konten.harga, 
						kursus_konten.harga_diskon, 
						kursus_konten.slug, 
						users.nama_lengkap 
					FROM 
						kursus_konten, 
						users 
					WHERE 
						kursus_konten.user_id = users.id 
					AND 
						kursus_konten.slug = ?`, slug).
		Scan(&kursus.IdKursus, &kursus.UserId, &kursus.KursusKategoriId, &kursus.NamaKursus, &kursus.Deskripsi, &kursus.Gambar, &kursus.Harga, &kursus.Diskon, &kursus.Slug, &kursus.NamaInstruktur)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if kursus.Slug == "" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":     "gagal",
			"dataKursus": kursus,
		})

	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":     "berhasil",
			"dataKursus": kursus,
		})
	}

}

// Membuat Konten Kursus
func TambahKursus(c echo.Context) (err error) {

	kursus := new(Kursus)
	if err = c.Bind(kursus); err != nil {
		return
	}

	// Define id sebagai random string
	id := uuid.New()

	// Generate Slug to LowerCase
	tolower := strings.ToLower(kursus.NamaKursus)
	slug := strings.ReplaceAll(tolower, " ", "-")

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

	//** End of File Upload

	// Generate Random File Name
	file, err := ioutil.TempFile("/var/www/vhosts/zcomeducation.com/httpdocs/upload/kursus", "kursus-*.png")
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

	strFile := strings.ReplaceAll(file.Name(), "/var/www/vhosts/zcomeducation.com/httpdocs/upload/kursus/", "")

	// Koneksi Database
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	_, err = db.Exec("INSERT INTO kursus_konten (id, user_id, kursus_kategori_id, nama_kursus, deskripsi, gambar, harga, harga_diskon, slug) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", id.String(), &kursus.UserId, &kursus.KursusKategoriId, &kursus.NamaKursus, &kursus.Deskripsi, strFile, &kursus.Harga, &kursus.Diskon, slug)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = db.
		QueryRow(`SELECT 
						kursus_konten.id, 
						kursus_konten.user_id, 
						kursus_konten.kursus_kategori_id, 
						kursus_konten.nama_kursus, 
						kursus_konten.deskripsi, 
						kursus_konten.gambar, 
						kursus_konten.harga, 
						kursus_konten.harga_diskon,
						kursus_konten.slug, 
						users.nama_lengkap 
				 FROM 
				 		kursus_konten, 
						users 
				 WHERE 
				 		kursus_konten.user_id = users.id 
				 AND 
						kursus_konten.user_id = ? ORDER BY kursus_konten.id DESC LIMIT 1`, &kursus.UserId).
		Scan(&kursus.IdKursus, &kursus.UserId, &kursus.KursusKategoriId, &kursus.NamaKursus, &kursus.Deskripsi, &kursus.Gambar, &kursus.Harga, &kursus.Diskon, &kursus.Slug, &kursus.NamaInstruktur)

	if err != nil {
		fmt.Println("gagal get data")
		return err
	} else {
		fmt.Println("berhasil get data")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":     "berhasil",
		"dataKursus": kursus,
	})

}

// perbaharui konten kursus
func UpdateKursus(c echo.Context) (err error) {

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	defer db.Close()

	kursus := new(Kursus)
	if err = c.Bind(kursus); err != nil {
		return err
	}

	// validasi inputan
	if err := c.Validate(kursus); err != nil {
		return err
	}

	// Generate Slug To LowerCase
	tolower := strings.ToLower(kursus.NamaKursus)
	slug := strings.ReplaceAll(tolower, " ", "-")

	// Select slug sebagai parameter
	slugKursus := c.Param("slug")
	fmt.Println(slugKursus)

	//** Start of File Upload
	//------------
	// Read files
	//------------

	// Multipart form
	// Get avatar
	avatar, err := c.FormFile("gambar")
	if err != nil {
		fmt.Println("tidak ada gambar")
		_, err = db.Exec(`UPDATE kursus_konten SET 
								user_id = ?, 
								kursus_kategori_id = ?, 
								nama_kursus = ?, 
								deskripsi = ?, 
								harga = ?, 
								harga_diskon = ?, 
								slug = ? 
							WHERE 
								slug = ?`, &kursus.UserId, &kursus.KursusKategoriId, &kursus.NamaKursus, &kursus.Deskripsi, &kursus.Harga, &kursus.Diskon, slug, slugKursus)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("gagal update")
			return err
		}

		/* ==============================================================
		|  Select data terakhir, order by created_at DESC(yang terakhir) |
		   =============================================================*/
		
		err = db.
			QueryRow(`SELECT 
							kursus_konten.id, 
							kursus_konten.user_id, 
							kursus_konten.kursus_kategori_id, 
							kursus_konten.nama_kursus, 
							kursus_konten.deskripsi, 
							kursus_konten.gambar, 
							kursus_konten.harga,
							kursus_konten.harga_diskon, 
							kursus_konten.slug, 
							users.nama_lengkap 
						FROM 
							kursus_konten, 
							users 
						WHERE 
							kursus_konten.user_id = users.id ORDER BY kursus_konten.created_at DESC LIMIT 1`).
			Scan(&kursus.IdKursus, &kursus.UserId, &kursus.KursusKategoriId, &kursus.NamaKursus, &kursus.Deskripsi, &kursus.Gambar, &kursus.Harga, &kursus.Diskon, &kursus.Slug, &kursus.NamaInstruktur)

		if err != nil {
			fmt.Println("gagal get data setelah update")
			return err
		} else {
			fmt.Println("berhasil get data setelah update")
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":     "berhasil",
			"dataKursus": kursus,
		})

	} else {
		
		/* ============================================================
		| Select data dari database berdasarkan slug sebelum diupdate |
		   ===========================================================*/
		
		err = db.
			QueryRow(`SELECT 
							kursus_konten.id, 
							kursus_konten.user_id, 
							kursus_konten.kursus_kategori_id, 
							kursus_konten.nama_kursus, 
							kursus_konten.deskripsi, 
							kursus_konten.gambar, 
							kursus_konten.harga, 
							kursus_konten.harga_diskon,
							kursus_konten.slug, 
							users.nama_lengkap 
						FROM 
							kursus_konten, 
							users 
						WHERE 
							kursus_konten.user_id = users.id 
						AND 
							kursus_konten.slug = ?`, slugKursus).
			Scan(&kursus.IdKursus, &kursus.UserId, &kursus.KursusKategoriId, &kursus.NamaKursus, &kursus.Deskripsi, &kursus.Gambar, &kursus.Harga, &kursus.Diskon, &kursus.Slug, &kursus.NamaInstruktur)

		/* ============================================================
		|  tampilkan error di terminal server 						  |
		   ===========================================================*/
		
		if err != nil {
			fmt.Println("gagal get data sebelum update dengan gambar")
			return err
		} else {
			fmt.Println("berhasil get data sebelum update dengan gambar")
		}

		/* ============================================================
		|  Ambil Lokasi Folder dengan os.Getwd (get working directory) |
		   ===========================================================*/
		
		// dir, err := os.Getwd()
		// if err != nil {
		// 	fmt.Println("gagal ambil folder tempat bekerja")
		// 	return err
		// }

		/* ============================================================
		|  Menentukan letak file upload   						      |
		   ===========================================================*/
		
		fileLocation := filepath.Join("/var/www/vhosts/zcomeducation.com/httpdocs/upload/kursus", kursus.Gambar)
		//fileLocation := filepath.Join(dir, "/upload/foto", kursus.Gambar)

		/* ============================================================
		|  Hapus gambar lama 						  				  |
		   ===========================================================*/
		
		hapusFile(fileLocation)

		/* ============================================================
		|  Tentukan sumber upload dari folder temporary 			  |
		   ===========================================================*/
		
		src, err := avatar.Open()
		if err != nil {
			fmt.Println("gagal buka gambar src")
			return err
		}
		defer src.Close()

		/* ======================================================================
		|  Tentukan folder tujuan upload dengan os.Getwd (get working directory |
		   =====================================================================*/
		
		// dir, err = os.Getwd()
		// if err != nil {
		// 	fmt.Println("gagal ambil folder saat ini")
		// 	return err
		// }

		/* ============================================================
		|  Tentukan lokasi file upload	   						      |
		   ===========================================================*/
		
		// fileLocation = filepath.Join("/var/www/vhosts/zcomeducation.com/httpdocs/upload/kursus", avatar.Filename)
		// //fileLocation = filepath.Join(dir, "/upload/foto", avatar.Filename)
		// dst, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
		// if err != nil {
		// 	fmt.Println("gagal ambil folder dst")
		// 	return err
		// }
		// defer dst.Close()

		/* ============================================================
		|  Copy / Move file 						  				  |
		   ===========================================================*/
		
		// if _, err = io.Copy(dst, src); err != nil {
		// 	fmt.Println("gagal copy file")
		// 	return err
		// }

		
		/* ============================================================
		|  Update database 						  					   |
		   ===========================================================*/

		// Generate Random File Name
		file, err := ioutil.TempFile("/var/www/vhosts/zcomeducation.com/httpdocs/upload/kursus", "kursus-*.png")
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

		strFile := strings.ReplaceAll(file.Name(), "/var/www/vhosts/zcomeducation.com/httpdocs/upload/kursus/", "")

		
		_, err = db.Exec(`UPDATE kursus_konten SET 
							user_id = ?, 
							kursus_kategori_id = ?, 
							nama_kursus = ?, 
							deskripsi = ?, 
							gambar = ?, 
							harga = ?, 
							harga_diskon = ?, 
							slug = ? 
						WHERE 
							slug = ?`, &kursus.UserId, &kursus.KursusKategoriId, &kursus.NamaKursus, &kursus.Deskripsi, strFile, &kursus.Harga, &kursus.Diskon, slug, slugKursus)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("gagal update database")
			return err
		}

		// select by id
		err = db.
			QueryRow(`SELECT 
							kursus_konten.id, 
							kursus_konten.user_id, 
							kursus_konten.kursus_kategori_id, 
							kursus_konten.nama_kursus, 
							kursus_konten.deskripsi, 
							kursus_konten.gambar, 
							kursus_konten.harga, 
							kursus_konten.harga_diskon,
							kursus_konten.slug, 
							users.nama_lengkap 
						FROM 
							kursus_konten, 
							users 
						WHERE 
							kursus_konten.user_id = users.id 
						AND 
							kursus_konten.slug = ?`, slugKursus).
			Scan(&kursus.IdKursus, &kursus.UserId, &kursus.KursusKategoriId, &kursus.NamaKursus, &kursus.Deskripsi, &kursus.Gambar, &kursus.Harga, &kursus.Diskon, &kursus.Slug, &kursus.NamaInstruktur)

		if err != nil {
			// fmt.Println(err.Error())
			fmt.Println("gagal get data return")
			return err
		} else {
			fmt.Println("berhasil get data return")
		}
	}

	// if kursus.IdKursus != "" {
	// 	return c.JSON(http.StatusOK, map[string]interface{}{
	// 		"status":     "gagal",
	// 		"dataKursus": kursus,
	// 	})

	// } else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":     "berhasil",
			"dataKursus": kursus,
		})
	// }

}

// Hapus Kursus
func DeleteKursus(c echo.Context) (err error) {

	// definisi variabel kursus dengan struct Kursus
	kursus := new(Kursus)

	// binding data inputan ke struct Kursus
	if err = c.Bind(kursus); err != nil {
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
	err = db.QueryRow(`
		SELECT 
			kursus_konten.id, 
			kursus_konten.user_id, 
			kursus_konten.kursus_kategori_id, 
			kursus_konten.nama_kursus, 
			kursus_konten.deskripsi, 
			kursus_konten.gambar, 
			kursus_konten.harga, 
			kursus_konten.harga_diskon, 
			kursus_konten.slug, 
			users.nama_lengkap 
		FROM 
			kursus_konten, 
			users 
		WHERE 
			kursus_konten.user_id = users.id 
		AND 
			kursus_konten.slug = ?
			`, slug).
		Scan(&kursus.IdKursus, &kursus.UserId, &kursus.KursusKategoriId, &kursus.NamaKursus, &kursus.Deskripsi, &kursus.Gambar,&kursus.Harga, &kursus.Diskon, &kursus.Slug, &kursus.NamaInstruktur)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	//* Menghapus Gambar *//
	// ambil lokasi folder saat ini
	// dir, err := os.Getwd()
	// if err != nil {
	// 	return err
	// }

	// menentukan letak upload folder
	fileLocation := filepath.Join("/var/www/vhosts/zcomeducation.com/httpdocs/upload/kursus", kursus.Gambar)
	//fileLocation := filepath.Join(dir, "/upload/foto", kursus.Gambar)

	// hapus gambar
	hapusFile(fileLocation)
	//* End Of Menghapus Gambar *//

	// hapus data dari db
	_, err = db.Exec(`DELETE FROM kursus_konten WHERE slug = ?`, slug)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":     "berhasil",
		"dataKursus": kursus,
	})
}