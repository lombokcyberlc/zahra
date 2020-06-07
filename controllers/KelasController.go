package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type KelasKursus struct {
	IdKelas        string `form:"idkelas" json:"idkelas"`
	UserID         string `form:"user_id" json:"user_id"`
	KursusID       string `form:"kursus_id" json:"kursus_id"`
	Status         string `form:"status_selesai" json:"status_selesai"`
	Deskripsi      string `form:"deskripsi" json:"deskripsi"`
	NamaKursus     string `form:"nama_kursus" json:"nama_kursus"`
	Gambar         string `form:"gambar" json:"gambar"`
	NamaInstruktur string `form:"nama_instruktur" json:"nama_instruktur"`
	Slug           string `form:"slug" json:"slug"`
}

// Menampilkan data semua kelas
func GetAllKelas(c echo.Context) (err error) {

	kelas := new(KelasKursus)
	if err = c.Bind(kelas); err != nil {
		return
	}

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	rows, err := db.Query("SELECT kelas.idkelas, kelas.user_id, kelas.kursus_id, kelas.status_selesai, kursus_konten.deskripsi, kursus_konten.nama_kursus, kursus_konten.gambar, users.nama_lengkap, kursus_konten.slug FROM kelas INNER JOIN kursus_konten ON kelas.kursus_id = kursus_konten.id INNER JOIN users ON kursus_konten.user_id = users.id")

	defer rows.Close()

	var result []KelasKursus

	for rows.Next() {
		var each = KelasKursus{}
		var err = rows.Scan(&each.IdKelas, &each.UserID, &each.KursusID, &each.Status, &each.Deskripsi, &each.NamaKursus, &each.Gambar, &each.NamaInstruktur, &each.Slug)

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
		"status":    "berhasil",
		"dataKelas": result,
	})
}

// Menampilkan data kelas by user id
func GetKelasByUserId(c echo.Context) (err error) {

	kelas := new(KelasKursus)
	if err = c.Bind(kelas); err != nil {
		return
	}

	id := c.Param("id")

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	rows, err := db.Query("SELECT kelas.idkelas, kelas.user_id, kelas.kursus_id, kelas.status_selesai, kursus_konten.deskripsi, kursus_konten.nama_kursus, kursus_konten.gambar, users.nama_lengkap, kursus_konten.slug FROM kelas INNER JOIN kursus_konten ON kelas.kursus_id = kursus_konten.id INNER JOIN users ON kursus_konten.user_id = users.id WHERE kelas.user_id = ?", id)

	defer rows.Close()

	var result []KelasKursus

	for rows.Next() {
		var each = KelasKursus{}
		var err = rows.Scan(&each.IdKelas, &each.UserID, &each.KursusID, &each.Status, &each.Deskripsi, &each.NamaKursus, &each.Gambar, &each.NamaInstruktur, &each.Slug)

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
		"status":      "berhasil",
		"dataKelasku": result,
	})
}

// Menampilkan data kelas by slug
func GetKelasBySlug(c echo.Context) (err error) {

	kelas := new(KelasKursus)
	if err = c.Bind(kelas); err != nil {
		return
	}

	slug := c.Param("slug")

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	rows, err := db.Query("SELECT kelas.idkelas, kelas.user_id, kelas.kursus_id, kelas.status_selesai, kursus_konten.deskripsi, kursus_konten.nama_kursus, kursus_konten.gambar, users.nama_lengkap, kursus_konten.slug FROM kelas INNER JOIN kursus_konten ON kelas.kursus_id = kursus_konten.id INNER JOIN users ON kursus_konten.user_id = users.id WHERE kursus_konten.slug = ?", slug)

	defer rows.Close()

	var result KelasKursus

	for rows.Next() {
		var each = KelasKursus{}
		var err = rows.Scan(&each.IdKelas, &each.UserID, &each.KursusID, &each.Status, &each.Deskripsi, &each.NamaKursus, &each.Gambar, &each.NamaInstruktur, &each.Slug)

		if err != nil {
			return err
		}

		result = each
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":      "berhasil",
		"dataKelasku": result,
	})
}
