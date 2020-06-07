package controllers

import (
	"fmt"
	"log"
	"net/http"
	// "time"

	"github.com/labstack/echo"
	uuid "github.com/google/uuid"
)

type ProgresBelajar struct {
	ID        string	 `form:"id_progres" json:"id_progres"`
	KelasID   string	 `form:"id_kelas" json:"id_kelas"`
	KursusID  string	 `form:"kursus_id" json:"kursus_id"`
	VideoID   string   	 `form:"video_id" json:"video_id"`
	VideoURL  string	 `form:"video_url" json:"video_url"`
	NamaVideo string	 `form:"nama_video" json:"nama_video"`
	Progres   int   	 `form:"status_tonton" json:"status_tonton"`
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []ProgresBelajar
}

// Menampilkan data soal
func GetProgresByKelas(c echo.Context) (err error) {

	progres := new(ProgresBelajar)
	if err = c.Bind(progres); err != nil {
		return
	}

	id := c.Param("id")

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	rows, err := db.Query(`
					SELECT 
						IFNULL(progres_belajar.id_progres, ""),
						kelas.idkelas, 
						kursus_video.kursus_id, 
						kursus_video.id, 
						kursus_video.video_url, 
						kursus_video.nama_video, 
						ifnull(progres_belajar.status_tonton, 0) 
					FROM 
						kelas 
					INNER JOIN 
						kursus_video 
					ON 
						kelas.kursus_id = kursus_video.kursus_id 
					LEFT JOIN 
						progres_belajar 
					ON 
						kelas.idkelas = progres_belajar.id_kelas 
					AND 
						kursus_video.id = progres_belajar.video_id 
					WHERE 
						kelas.idkelas = ? ORDER BY kursus_video.created_at ASC`, id)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer rows.Close()

	var result []ProgresBelajar
	var response Response

	for rows.Next() {

		var each = ProgresBelajar{}
		var err = rows.Scan(&each.ID, &each.KelasID, &each.KursusID, &each.VideoID, &each.VideoURL, &each.NamaVideo, &each.Progres)

		if err != nil {
			log.Fatal(err.Error())
			return err
		}

		result = append(result, each)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
		return
	}

	response.Status = 1
	response.Message = "berhasil"
	response.Data = result

	return c.JSON(http.StatusOK, response)
}

// Status Tonton
func PostProgres(c echo.Context) (err error) {

	progres := new(ProgresBelajar)
	if err = c.Bind(progres); err != nil {
		return
	}

	progresID := uuid.New()

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	err = db.QueryRow("SELECT id_progres, id_kelas, video_id, status_tonton FROM progres_belajar WHERE id_kelas = ? AND video_id = ? AND status_tonton = 1", &progres.KelasID, &progres.VideoID).Scan(&progres.ID, &progres.KelasID, &progres.VideoID, &progres.Progres)
	if err != nil {
		_, err = db.Exec("INSERT INTO progres_belajar (id_progres, id_kelas, video_id, status_tonton) VALUES (?, ?, ?, 1)", progresID, &progres.KelasID, &progres.VideoID)
		if err != nil {
			return err
		}
	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": false,
			"pesan": "Gagal",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": true,
		"pesan": "Berhasil",
	})
}
