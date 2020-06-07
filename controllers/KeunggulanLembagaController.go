package controllers

import (
	"github.com/labstack/echo"
	"fmt"
	"net/http"
)


type KeunggulanLembaga struct {
	ID 			int `form:"id_keunggulan" json:"id_keunggulan"`
	Icon 		string `form:"icon" json:"icon"`
	Keunggulan 	string `form:"keunggulan" json:"keunggulan"`
}

// Menampilkan data keunggulan lembaga
func GetAllKeunggulan(c echo.Context) (err error) {
	
	kl := new(KeunggulanLembaga)
	if err = c.Bind(kl); err != nil {
		return
	}

	
	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return
	}
	
    defer db.Close()

	rows, err := db.Query("SELECT id_keunggulan, icon, keunggulan FROM keunggulan_lembaga")
    
	if err != nil {
        fmt.Println(err.Error())
        return
    }
    
	defer rows.Close()

    var result []KeunggulanLembaga

    for rows.Next() {

        var each = KeunggulanLembaga{}
        var err = rows.Scan(&each.ID, &each.Icon, &each.Keunggulan)

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
		"dataKeunggulan": result,
	})
}