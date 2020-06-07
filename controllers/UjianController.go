package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type UjianKursus struct {
	ID           int    `form:"id" 				json:"id"`
	KelasID      string `form:"kelas_id" 		json:"kelas_id"`
	SoalID       int    `form:"soal_id" 		json:"soal_id"`
	KursusID     string `form:"kursus_id"		json:"-"`
	Pertanyaan   string `form:"pertanyaan" 		json:"pertanyaan"`
	Jawaban      string `form:"jawaban" 		json:"jawaban"`
	KunciJawaban string `form:"kunci_jawaban" 	json:"-"`
	Hasil        int    `json:"hasil"`
}

type (
	Ujian struct {
		ID    string `json:"id_kelas"`
		Soals []Soal `json:"soal"`
	}

	Soal struct {
		ID      string `json:"id_soal"`
		Jawaban string `json:"jawaban"`
	}
)

// Menampilkan data soal
func GetAllUjian(c echo.Context) (err error) {
	

	/*
		============================================================
		| Binding data struct ke variabel baru 					   |
		============================================================
	*/

	soal := new(UjianKursus)

	if err = c.Bind(soal); err != nil {
		return
	}

	/*
		============================================================
		| Menghubungkan function sehingga dapat mengakses database |
		============================================================
	*/

	db, err := connect()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	/*
		===============================================================
		| Mengambil banyak data ujian dengan untuk ditampung ke slice |
		===============================================================
	*/

	rows, err := db.Query(`
					SELECT
						ujian.id,
						ujian.kelas_id,
						ujian.soal_id,
						soal.pertanyaan,
						ujian.jawaban,
						soal.kunci_jawaban
					FROM 
						ujian,
						soal,
						kelas
					WHERE
						ujian.kelas_id = kelas.idkelas
					AND 
						ujian.soal_id = soal.id`)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Tutup Koneksi row
	defer rows.Close()

	/*
		=======================================================================
		| Menampung semua hasil ke dalam slice dengan membentuk variabel baru |
		| (@each) sebagai variabel objek 									  |
		=======================================================================
	*/

	// Definisi result menampung Struct dari UjianKursus
	var result []UjianKursus

	// Looping untuk menyisipkan Objek
	for rows.Next() {

		var each = UjianKursus{}
		var err = rows.Scan(&each.ID, &each.KelasID, &each.SoalID, &each.Pertanyaan, &each.Jawaban, &each.KunciJawaban)

		if err != nil {
			return err
		}

		result = append(result, each)
	}

	// Pengecekan error di rows objek
	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
		return
	}

	/*
		===================================================================
		| Mengembalikan data ujian dalam bentuk JSON dari variabel result |
		===================================================================
	*/

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":    "berhasil",
		"dataUjian": result,
	})
}

/*
* ===================================================================
* Input Jawaban 													|
* ===================================================================
 */

func PostJawaban(c echo.Context) (err error) {

	/*
		=========================================
		| Binding data struct ke variabel ujian |
		=========================================
	*/

	ujian := new(Ujian)

	if err = c.Bind(ujian); err != nil {
		return
	}

	/*
		============================================================
		| Menghubungkan function sehingga dapat mengakses database |
		============================================================
	*/

	db, err := connect()

	if err != nil {
		fmt.Println("Tidak dapat terhubung ke database !")
		return
	}

	defer db.Close()

	/*
		============================================================
		| Mendefinisikan variabel sebagai form data 			   |
		============================================================
	*/

	var (
		kursusID string
		Benar    int
	)

	Benar = 0

	db.QueryRow("SELECT kursus_id FROM `kelas` WHERE idkelas = ?", ujian.ID).Scan(&kursusID)

	fmt.Println("==============================================================")
	fmt.Print("| ID\t| Jawaban\t  | Kunci\t\t| Hasil\t| \n")
	fmt.Println("==============================================================")

	for _, value := range ujian.Soals {

		var (
			idSoal       string
			kunciJawaban string
		)

		db.QueryRow("SELECT id, kunci_jawaban FROM soal WHERE id = ?", value.ID).Scan(&idSoal, &kunciJawaban)

		fmt.Print("| ", value.ID)

		if value.Jawaban == "" {
			fmt.Print("\t| --")
		} else {
			fmt.Print("\t| ", value.Jawaban)
		}

		fmt.Print("\t\t| ", kunciJawaban)

		if value.Jawaban == kunciJawaban {
			Benar++
			fmt.Println("\t\t| Benar\t|")
		} else {
			fmt.Println("\t\t| Salah\t|")
		}
	}

	fmt.Println("==============================================================")

	/*
		============================================================
		| Mengembalikan hasil tampungan struct ujian 			   |
		============================================================
	*/
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":    "berhasil",
		"dataUjian": ujian,
		"benar":     Benar,
	})
}
