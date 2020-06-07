package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo"

	uuid "github.com/google/uuid"
)

type PostKategori struct {
	ID           string `form:"id" json:"id"`
	NamaKategori string `form:"nama_kategori" json:"nama_kategori"`
	Slug         string `json:"slug"`
}

// Menampilkan data Post Kategori
func GetAllPostKategori(c echo.Context) (err error) {

	post := new(PostKategori)
	if err = c.Bind(post); err != nil {
		return
	}

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	rows, err := db.Query("SELECT id, nama_kategori, slug FROM post_kategori ORDER BY created_at ASC")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer rows.Close()

	var result []PostKategori

	for rows.Next() {

		var each = PostKategori{}
		var err = rows.Scan(&each.ID, &each.NamaKategori, &each.Slug)

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
		"status": "berhasil",
		"data":   result,
	})
}

// Menampilkan data Post Kategori berdasarkan ID
func GetPostKategoriBySlug(c echo.Context) (err error) {

	post := new(PostKategori)

	if err = c.Bind(post); err != nil {
		return
	}

	slug := c.Param("slug")

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	err = db.QueryRow(`
				SELECT 
					id, 
					nama_kategori, 
					slug 
				FROM 
					post_kategori 
				WHERE slug = ?`, slug).Scan(&post.ID, &post.NamaKategori, &post.Slug)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "berhasil",
		"data":   post,
	})
}

// Membuat data Post Kategori
func TambahPostKategori(c echo.Context) (err error) {

	post := new(PostKategori)
	if err = c.Bind(post); err != nil {
		return
	}

	id := uuid.New()
	// Generate Slug to LowerCase
	tolower := strings.ToLower(post.NamaKategori)
	slug := strings.ReplaceAll(tolower, " ", "-")

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	_, err = db.Query(`
				INSERT INTO 
					post_kategori (
						id, 
						nama_kategori, 
						slug
				) VALUES (?, ?, ?)`, id.String(), &post.NamaKategori, slug)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = db.QueryRow(`
				SELECT 
					id, 
					nama_kategori, 
					slug 
				FROM 
					post_kategori 
				ORDER BY id desc limit 1`).Scan(&post.ID, &post.NamaKategori, &post.Slug)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "berhasil",
		"data":   post,
	})
}

// Membuat data Post Kategori
func UpdatePostKategori(c echo.Context) (err error) {

	post := new(PostKategori)
	if err = c.Bind(post); err != nil {
		return
	}

	slug := c.Param("slug")

	// Generate Slug to LowerCase
	tolower := strings.ToLower(post.NamaKategori)
	newSlug := strings.ReplaceAll(tolower, " ", "-")

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	_, err = db.Exec(`
				UPDATE 
					post_kategori 
				SET 
					nama_kategori = ?,
					slug = ? 
				WHERE 
					slug = ?`, &post.NamaKategori, newSlug, slug)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "berhasil",
	})
}

// Hapus Post kategori
func DeletePostKategori(c echo.Context) (err error) {

	post := new(PostKategori)
	if err = c.Bind(post); err != nil {
		return
	}

	slug := c.Param("slug")

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	_, err = db.Exec("DELETE FROM post_kategori WHERE slug = ?", slug)

	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("gagal hapus post kategori")
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "berhasil",
	})
}
