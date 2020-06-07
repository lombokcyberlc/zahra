package controllers

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

type Pertanyaan struct {
	ID		int `form:"id_tanya" json:"id_tanya"`
	UserID	int `form:"id_user" json:"id_user"`
	Tanya 	string `form:"pertanyaan" json:"pertanyaan"`
}

func PostQuestions(c echo.Context) (err error) {

	quis := new(Pertanyaan)
	if err = c.Bind(quis); err != nil {
		return
	}

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	err = db.
	QueryRow("SELECT tanya.id_tanya, tanya.id_user, tanya.pertanyaan FROM tanya, kursus_konten WHERE tanya.id_user = kursus_konten.user_id ORDER BY CURRENT_TIME() DESC LIMIT 1").
	Scan(&quis.ID, &quis.UserID, &quis.Tanya)

	if err != nil {
		// fmt.Println(err.Error())
		fmt.Println("gagal get data")
		return err
	} else {
		fmt.Println("berhasil get data")
	}

	_, err = db.Exec("INSERT INTO tanya (pertanyaan) VALUES (?)", &quis.Tanya)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":     "berhasil",
		"dataQuis": quis,
	})
}	


// Get Pertanyaan
func GetPertanyaan(c echo.Context) (err error) {
	
	quis := new(Pertanyaan)
	if err = c.Bind(quis); err != nil {
		return
	}

	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return
	}
	
    defer db.Close()

	rows, err := db.Query("SELECT tanya.id_tanya, tanya.id_user, tanya.pertanyaan FROM tanya, kursus_konten WHERE tanya.id_user = kursus_konten.user_id")

	defer rows.Close()

	var result []Pertanyaan
	
	for rows.Next() {
		var each = Pertanyaan{}
		var err = rows.Scan(&each.ID, &each.UserID, &each.Tanya)
	
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
		"dataQuis": result,
	})
}

// Get Pertanyaan By ID
func GetPertanyaanById(c echo.Context) (err error) {
	
	quis := new(Pertanyaan)
	if err = c.Bind(quis); err != nil {
		return
	}

	id := c.Param("id")

	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return
	}
	
    defer db.Close()

	rows, err := db.Query("SELECT tanya.id_tanya, tanya.id_user, tanya.pertanyaan FROM tanya, kursus_konten WHERE tanya.id_user = kursus_konten.user_id AND tanya.id_user = ?", id)

	defer rows.Close()

	var result Pertanyaan
	
	for rows.Next() {
		var each = Pertanyaan{}
		var err = rows.Scan(&each.ID, &each.UserID, &each.Tanya)
	
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
		"status":  "berhasil",
		"dataQuis": result,
	})
}