package controllers

import (
	"net/http"
	"fmt"
	"github.com/labstack/echo"
)

type PartnerLembaga struct {
	ID 		int `form:"id_partner" json:"id_partner"`
	Logo	string `form:"logo" json:"logo"`
}

// func partner lembaga
func GetPartnerLembaga(c echo.Context) (err error) {
	
	pl := new(PartnerLembaga)
	if err = c.Bind(pl); err != nil {
		return
	}

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	rows, err := db.Query("SELECT id_partner, logo FROM partner_lembaga")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer rows.Close()

	var result []PartnerLembaga

	for rows.Next() {
		var each = PartnerLembaga{}
		var err = rows.Scan(&each.ID, &each.Logo)

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
		"dataParter" : result,
	})

}