package controllers

import (
	"net/http"
	"fmt"
	"github.com/labstack/echo"
)

type TestimoniLembaga struct {
	ID			int	   `form:"id" json:"id"`
	Nama		string `form:"nama" json:"nama"`
	Testimoni	string `form:"testimoni" json:"testimoni"`
	Foto 		string `form:"foto" json:"foto"`
}

func GetTestimoniLembaga(c echo.Context) (err error) {
	
	tl := new(TestimoniLembaga)
	if err = c.Bind(tl); err != nil {
		return
	}

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	rows, err := db.Query("SELECT id, nama, testimoni, foto FROM testimoni_lembaga")
	if err != nil {
		return
	}

	defer rows.Close()

	var result []TestimoniLembaga

	for rows.Next() {
		var each = TestimoniLembaga{}
		var err = rows.Scan(&each.ID, &each.Nama, &each.Testimoni, &each.Foto)
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
		"dataTestimoni": result,
	})
}