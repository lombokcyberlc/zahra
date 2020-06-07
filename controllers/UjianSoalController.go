package controllers

import (
	"github.com/labstack/echo"
	"fmt"
	"net/http"
	// "strconv"
)


type UjianSoal struct {
	IdSoal			int 	`form:"id" 				json:"id"`
	KelasID 		string	`form:"kelas_id"		json:"kelas_id"`
	KursusID 		string 	`form:"kursus_id" 		json:"kursus_id"`
	Pertanyaan		string 	`form:"pertanyaan"		json:"pertanyaan"`
	JawabanA	 	string 	`form:"jawaban_a" 		json:"jawaban_a"`
	JawabanB 		string 	`form:"jawaban_b" 		json:"jawaban_b"`
	JawabanC 		string 	`form:"jawaban_c" 		json:"jawaban_c"`
	JawabanD 		string 	`form:"jawaban_d" 		json:"jawaban_d"`
	KunciJawaban	string 	`form:"kunci_jawaban" 	json:"kunci_jawaban"`
}

// Menampilkan data soal
func GetAllUjianSoal(c echo.Context) (err error) {
	
	soal := new(UjianSoal)
	if err = c.Bind(soal); err != nil {
		return
	}
	
	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return
	}
	
    defer db.Close()

	rows, err := db.Query(`
					SELECT 
						soal.id, 
						soal.kursus_id, 
						soal.pertanyaan, 
						soal.jawaban_a, 
						soal.jawaban_b, 
						soal.jawaban_c, 
						soal.jawaban_d, 
						soal.kunci_jawaban 
					FROM 
						kelas 
					INNER JOIN 
						soal 
					ON 
						kelas.kursus_id = soal.kursus_id 
					GROUP BY soal.id`)
    
	if err != nil {
        fmt.Println(err.Error())
        return
    }
    
	defer rows.Close()

    var result []UjianSoal

    for rows.Next() {

        var each = UjianSoal{}
        var err = rows.Scan(&each.IdSoal, &each.KursusID, &each.Pertanyaan, &each.JawabanA, &each.JawabanB, &each.JawabanC, &each.JawabanD, &each.KunciJawaban)

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
		"dataSoal": result,
	})
}

// Menampilkan data Soal Ujian untuk admin
func GetSoalUjianByKursus(c echo.Context) (err error) {
	
	soal := new(UjianSoal)

	if err = c.Bind(soal); err != nil {
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
					soal.id, 
					kelas.idkelas,
					soal.kursus_id,
					soal.pertanyaan, 
					soal.jawaban_a, 
					soal.jawaban_b, 
					soal.jawaban_c, 
					soal.jawaban_d, 
					soal.kunci_jawaban
				FROM
					kelas
				INNER JOIN
					soal ON kelas.kursus_id = soal.kursus_id
				LEFT JOIN
					ujian ON kelas.idkelas = ujian.kelas_id
				WHERE kelas.idkelas = ?
				GROUP BY soal.id`, id)
    
	if err != nil {
        fmt.Println(err.Error())
        return
    }
    
	defer rows.Close()

    var result []UjianSoal

    for rows.Next() {
        var each = UjianSoal{}
        var err = rows.Scan(&each.IdSoal, &each.KelasID, &each.KursusID, &each.Pertanyaan, &each.JawabanA, &each.JawabanB, &each.JawabanC, &each.JawabanD, &each.KunciJawaban)

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
		"status": true,
		"pesan":"Data ditemukan !",
		"dataSoal": result,
	})

}

// Menampilkan data Soal Ujian
func GetSoalUjianByID(c echo.Context) (err error) {
	
	soal := new(UjianSoal)

	if err = c.Bind(soal); err != nil {
		return
	}
	
	id := c.Param("id")
	
	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    defer db.Close()

    err = db.QueryRow(`
				SELECT 
					soal.id, 
					soal.kursus_id,
					soal.pertanyaan, 
					soal.jawaban_a, 
					soal.jawaban_b, 
					soal.jawaban_c, 
					soal.jawaban_d, 
					soal.kunci_jawaban 
				FROM 
					kelas 
				INNER JOIN 
					ujian_soal 
				ON 
					kelas.kursus_id = soal.kursus_id 
				WHERE 
					soal.id = ? 
				GROUP BY soal.id`, id).Scan(&soal.IdSoal, &soal.KursusID, &soal.Pertanyaan, &soal.JawabanA, &soal.JawabanB, &soal.JawabanC, &soal.JawabanD, &soal.KunciJawaban)
    
	if err != nil {
        fmt.Println(err.Error())
        return
    }

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "berhasil",
		"dataSoal": soal,
	})

}

// Menambah Soal Ujian
func TambahUjianSoal(c echo.Context) (err error) {
	
	soal := new(UjianSoal)
	if err = c.Bind(soal); err != nil {
		return
	}

	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return
	}
	
    defer db.Close()

    _, err = db.Exec(`
				INSERT INTO soal (
					kursus_id, 
					pertanyaan, 
					jawaban_a, 
					jawaban_b, 
					jawaban_c, 
					jawaban_d, 
					kunci_jawaban
				) VALUES (?, ?, ?, ?, ?, ?, ?)`, &soal.KursusID, &soal.Pertanyaan, &soal.JawabanA, &soal.JawabanB, &soal.JawabanC, &soal.JawabanD, &soal.KunciJawaban)
    
	if err != nil {
		fmt.Println("gagal insert data soal")
		return
	} else {
		fmt.Println("berhasil insert ke database")
	}

	err = db.QueryRow(`
				SELECT 
					id, 
					kursus_id,
					pertanyaan, 
					jawaban_a, 
					jawaban_b, 
					jawaban_c, 
					jawaban_d, 
					kunci_jawaban 
				FROM 
					soal 
				ORDER BY id DESC LIMIT 1`).Scan(&soal.IdSoal, &soal.KursusID, &soal.Pertanyaan, &soal.JawabanA, &soal.JawabanB, &soal.JawabanC, &soal.JawabanD, &soal.KunciJawaban)

	if err != nil {
		fmt.Println("gagal get data")
		return err
	} else {
		fmt.Println("berhasil get data")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "berhasil",
		"dataSoal": soal,
	})

}

// Update Soal
func UpdateUjianSoal(c echo.Context) (err error) {

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer db.Close()

	soal := new(UjianSoal)
	if err = c.Bind(soal); err != nil {
		return
	}

	// validasi inputan
	if err := c.Validate(soal); err != nil {
		return err
	}

	id := c.Param("id")

	// Update Query	
	_, err = db.Exec(`
				UPDATE ujian_soal SET 
					pertanyaan = ?,
					jawaban_a = ?, 
					jawaban_b = ?, 
					jawaban_c = ?, 
					jawaban_d = ?, 
					kunci_jawaban = ? 
				WHERE id = ?`, &soal.Pertanyaan, &soal.JawabanA, &soal.JawabanB, &soal.JawabanC, &soal.JawabanD, &soal.KunciJawaban, id)
	
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("errornya di sini")
			return err
		}
	
	// Select Query setelah di update
	err = db.QueryRow(`
				SELECT 
					id, 
					kursus_id,
					pertanyaan, 
					jawaban_a, 
					jawaban_b, 
					jawaban_c, 
					jawaban_d, 
					kunci_jawaban 
				FROM 
					ujian_soal 
				WHERE id = ?`, id).
		Scan(&soal.IdSoal, &soal.KursusID, &soal.Pertanyaan, &soal.JawabanA, &soal.JawabanB, &soal.JawabanC, &soal.JawabanD, &soal.KunciJawaban)

		if err != nil {
			// fmt.Println(err.Error())
			fmt.Println("gagal get data")
			return err
		} else {
			fmt.Println("berhasil get data")
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "berhasil",
			"dataSoal": soal,
		})
}

// Hapus Soal
func DeleteUjianSoal(c echo.Context) (err error) {

	soal := new(UjianSoal)
	if err = c.Bind(soal); err != nil {
		return
	}

	id := c.Param("id")
	
	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return
	}
	
	defer db.Close()

    _, err = db.Exec("DELETE FROM ujian_soal WHERE id = ?", id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":"berhasil",
		"dataSoal": soal,
	})
}
