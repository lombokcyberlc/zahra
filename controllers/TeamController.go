package controllers

import (
	"github.com/labstack/echo"
	"fmt"
	"net/http"
)


type Team struct {
	ID 		int `form:"id" json:"id"`
	Nama 	string `form:"nama" json:"nama"`
	Jabatan	string `form:"jabatan" json:"jabatan"`
	Foto	string `form:"foto" json:"foto"`
	Motto 	string `form:"motto" json:"motto"`
}

// Menampilkan data team
func GetAllTeam(c echo.Context) (err error) {
	
	t := new(Team)
	if err = c.Bind(t); err != nil {
		return
	}

	
	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return
	}
	
    defer db.Close()

	rows, err := db.Query("SELECT id, nama, jabatan, foto, motto FROM team")
    
	if err != nil {
        fmt.Println(err.Error())
        return
    }
    
	defer rows.Close()

    var result []Team

    for rows.Next() {

        var each = Team{}
        var err = rows.Scan(&each.ID, &each.Nama, &each.Jabatan, &each.Foto, &each.Motto)

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
		"dataTeam": result,
	})
}