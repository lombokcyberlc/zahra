package controllers

import (
	"github.com/labstack/echo"
	"fmt"
	"net/http"
	"strconv"

)

type UserGroup struct {
	ID			int `form:"id" json:"id"`
	NamaGroup	string `form:"nama_group" json:"nama_group"`
}

// Menampilkan data semua user group
func GetUsersGroup(c echo.Context) (err error) {
	
	ug := new(User)
	if err = c.Bind(ug); err != nil {
		return
	}
	
	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    defer db.Close()

    rows, err := db.Query("SELECT id, nama_group FROM user_group")
    
	if err != nil {
        fmt.Println(err.Error())
        return
    }
    
	defer rows.Close()

    var result []UserGroup

    for rows.Next() {
        var each = UserGroup{}
        var err = rows.Scan(&each.ID, &each.NamaGroup)

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
		"data": result,
	})
}

// Menampilkan data user group berdasarkan ID
func GetUserGroupById(c echo.Context) (err error) {
	
	ug := new(UserGroup)
	if err = c.Bind(ug); err != nil {
		return
	}
	
	id := c.Param("id")
	
	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    defer db.Close()

    rows, err := db.Query("SELECT id, nama_group FROM user_group WHERE id = ?", id)
    
	if err != nil {
        fmt.Println(err.Error())
        return
    }
    
	defer rows.Close()

    var result UserGroup

    for rows.Next() {
        var each = UserGroup{}
        var err = rows.Scan(&each.ID, &each.NamaGroup)

        if err != nil {
			return err
        }

        result = each
    }

    if err = rows.Err(); err != nil {
        fmt.Println(err.Error())
        return
    }
		
	if strconv.Itoa(result.ID) == "" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "gagal",
			"data": result,
		})
	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "berhasil",
			"data": result,
		})
	}

}

// Membuat User Group baru
func TambahUserGroup(c echo.Context) (err error) {
	
	ug := new(UserGroup)
	if err = c.Bind(ug); err != nil {
		return
	}
	
	nama := c.FormValue("nama_group")
	
	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    defer db.Close()

    _, err = db.Query("INSERT INTO user_group (nama_group) VALUES (?)", nama)
    
	if err != nil {
        fmt.Println(err.Error())
        return
    }
    
	var result = UserGroup{}
    err = db.QueryRow("select id, nama_group from user_group order by id desc limit 1").Scan(&result.ID, &result.NamaGroup)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
		
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "berhasil",
		"data": result,
	})
}
