package controllers

import (
	"github.com/labstack/echo"
	"fmt"
	"net/http"
    "strings"

	uuid "github.com/google/uuid"

)

type KursusKategori struct {
	ID_Kategori 	string `form:"id" json:"id"`
    NamaKategori 	string `form:"nama_kategori_kursus" json:"nama_kategori_kursus"`
    Slug            string `json:"slug"`
}

// Menampilkan data semua user
func GetAllKursusKategori(c echo.Context) (err error) {
	
	kategori := new(KursusKategori)
	if err = c.Bind(kategori); err != nil {
		return
	}
	
	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    defer db.Close()

    rows, err := db.Query("SELECT id, nama_kategori_kursus, slug FROM kursus_kategori")
    
	if err != nil {
        fmt.Println(err.Error())
        return
    }
    
	defer rows.Close()

    var result []KursusKategori

    for rows.Next() {
        var each = KursusKategori{}
        var err = rows.Scan(&each.ID_Kategori, &each.NamaKategori, &each.Slug)

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
		"dataKategori": result,
	})
}

// Menampilkan data Kursus Kategori berdasarkan Id
func GetKursusKategoriBySlug(c echo.Context) (err error) {
	
	kategori := new(KursusKategori)
	if err = c.Bind(kategori); err != nil {
		return
	}
	
	slug := c.Param("slug")
	
	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    
    defer db.Close()

    err = db.QueryRow("SELECT id, nama_kategori_kursus, slug FROM kursus_kategori WHERE slug = ?", slug).Scan(&kategori.ID_Kategori, &kategori.NamaKategori, &kategori.Slug)
    
	if err != nil {
        fmt.Println(err.Error())
        return
    }
		
	if kategori.Slug == "" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "gagal",
			"dataKategori": kategori,
		})
	
	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "berhasil",
			"dataKategori": kategori,
		})
	}
}

// Membuat data Kategori kursus
func TambahKursusKategori(c echo.Context) (err error) {
	
	tkk := new(KursusKategori)
	if err = c.Bind(tkk); err != nil {
		return
    }

	// Define id sebaga uuid
	id := uuid.New()
    
	// Generate Slug
    tolower := strings.ToLower(tkk.NamaKategori)
    slug := strings.ReplaceAll(tolower, " ", "-")
	
	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    defer db.Close()

    _, err = db.Query("INSERT INTO kursus_kategori (id, nama_kategori_kursus, slug) VALUES (?, ?, ?)", id.String(), tkk.NamaKategori, slug)
    
	if err != nil {
        fmt.Println(err.Error())
        return
    }
    
	var result = KursusKategori{}
    err = db.QueryRow("select id, nama_kategori_kursus, slug from kursus_kategori order by id desc limit 1").
	Scan(&result.ID_Kategori, &result.NamaKategori, &result.Slug)
    
	if err != nil {
        fmt.Println(err.Error())
        return
    }
		
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "berhasil",
		"data": result,
	})
}

// Update Kategori Kursus
func UpdateKursusKategori(c echo.Context) (err error) {
	
	ukk := new(KursusKategori)
	if err = c.Bind(ukk); err != nil {
		return
	}
    
	slug := c.Param("slug")
    tolower := strings.ToLower(ukk.NamaKategori)
    formSlug := strings.ReplaceAll(tolower, " ", "-")
	
	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    defer db.Close()

    _, err = db.Exec("UPDATE kursus_kategori SET nama_kategori_kursus = ?, slug = ? WHERE slug = ?", ukk.NamaKategori, formSlug, slug)
    
	if err != nil {
        fmt.Println(err.Error())
        return
    }
    
	var result = KursusKategori{}
    err = db.QueryRow("select id, nama_kategori_kursus, slug from kursus_kategori order by created_at desc limit 1").
	Scan(&result.ID_Kategori, &result.NamaKategori, &result.Slug)
	
    if err != nil {
        fmt.Println(err.Error())
        return
    }
		
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "berhasil",
		"dataKategori": result,
	})
}

// Hapus kategori kursus
func DeleteKursusKategori(c echo.Context) (err error) {

	dk := new(KursusKategori)
	if err = c.Bind(dk); err != nil {
		return
	}

	slug := c.Param("slug")
	
	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return err
	}
	
    defer db.Close()

    _, err = db.Exec("DELETE FROM kursus_kategori WHERE slug = ?", slug)
    
	if err != nil {
        fmt.Println(err.Error())
        return c.JSON(http.StatusOK, map[string]interface{}{
			"status":"gagal",
			"slug": slug,
			"error": err,
		})
	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":"berhasil",
			"slug": slug,
		})
	}
}