package controllers

import (
	"github.com/labstack/echo"
	"fmt"
	"net/http"
	// "strings"
	"os"
	"io"
	// "io/ioutil"
	"path/filepath"
	"context"
	"log"
	"net/url"

	"github.com/jwplayer/jwplatform-go"
	uuid "github.com/google/uuid"
)

type KursusVideo struct {
	ID			string `form:"id" json:"id"`
	KursusID 	string `form:"kursus_id" json:"kursus_id"`
	NamaVideo	string `form:"nama_video" json:"nama_video"`
	VideoUrl 	string `form:"video_url" json:"video_url"`
	Slug		string `json:"-"`
}

// Menampilkan data Video
func GetAllVideo(c echo.Context) (err error) {
	
	video := new(KursusVideo)
	if err = c.Bind(video); err != nil {
		return
	}
	
	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return
	}
	
    defer db.Close()

	rows, err := db.Query(`SELECT 
								kursus_video.id, 
								kursus_video.kursus_id, 
								kursus_video.video_url, 
								kursus_video.nama_video, 
								kursus_konten.slug 
							FROM 
								kursus_video 
							JOIN 
								kursus_konten 
							ON 
								kursus_video.kursus_id = kursus_konten.id`)
    
	if err != nil {
        fmt.Println(err.Error())
        return
    }
    
	defer rows.Close()

    var result []KursusVideo

    for rows.Next() {

        var each = KursusVideo{}
        var err = rows.Scan(&each.ID, &each.KursusID, &each.VideoUrl, &each.NamaVideo, &each.Slug)

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
		"dataVideo": result,
	})
}

// Menampilkan data Video berdasarkan ID
func GetVideoByKursusSlug(c echo.Context) (err error) {
	
	video := new(KursusVideo)

	if err = c.Bind(video); err != nil {
		return
	}
	
	db, err := connect()

    if err != nil {
        fmt.Println(err.Error())
        return
	}
	
    defer db.Close()

	slug := c.Param("slug")

	// SELECT video dengan SUB query dengan mengambil id slug dari kursus_konten untuk direlasikan ke video
    rows, err := db.Query(`
					SELECT 
						kursus_video.id, 
						kursus_video.kursus_id, 
						kursus_video.video_url, 
						kursus_video.nama_video, 
						kursus_konten.slug 
					FROM 
						kursus_video 
					JOIN 
						kursus_konten 
					ON 
						kursus_video.kursus_id = kursus_konten.id 
					WHERE 
						kursus_konten.slug = ?`, slug)
    
	if err != nil {
        fmt.Println(err.Error())
        return
    }
    
	defer rows.Close()

    var result []KursusVideo

    for rows.Next() {

        var each = KursusVideo{}
        var err = rows.Scan(&each.ID, &each.KursusID, &each.VideoUrl, &each.NamaVideo, &each.Slug)

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
		"dataVideo": result,
	})

}

// Menampilkan data Video berdasarkan ID
func GetVideoByKursusID(c echo.Context) (err error) {
	
	video := new(KursusVideo)

	if err = c.Bind(video); err != nil {
		return
	}
	
	db, err := connect()

    if err != nil {
        fmt.Println(err.Error())
        return
    }
	
	defer db.Close()

	id := c.Param("id")

	// SELECT video dengan SUB query dengan mengambil id slug dari kursus_konten untuk direlasikan ke video
    rows, err := db.Query(`
					SELECT 
						kursus_video.id, 
						kursus_video.kursus_id, 
						kursus_video.video_url, 
						kursus_video.nama_video, 
						kursus_konten.slug 
					FROM 
						kursus_video 
					JOIN 
						kursus_konten 
					ON 
						kursus_video.kursus_id = kursus_konten.id 
					WHERE 
						kursus_konten.id = ?
					ORDER BY kursus_video.created_at ASC`, id)
    
	if err != nil {
        fmt.Println(err.Error())
        return
    }
    
	defer rows.Close()

    var result []KursusVideo

    for rows.Next() {

        var each = KursusVideo{}
        var err = rows.Scan(&each.ID, &each.KursusID, &each.VideoUrl, &each.NamaVideo, &each.Slug)

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
		"dataVideo": result,
	})

}

// Tambah Video
func TambahVideo(c echo.Context) (err error) {
	
	video := new(KursusVideo)
	if err = c.Bind(video); err != nil {
		return
	}

	//** Start of File Upload

	// Multipart form
	// Get avatar
	avatar, err := c.FormFile("video_url")

	if err != nil {
		fmt.Println("Gagal ambil file dari form")
		fmt.Println(err)
	return err
	}
	
	 // Source
	src, err := avatar.Open()

	if err != nil {
		fmt.Println("Gagal membuka file dari src")
		return err
	}

	defer src.Close()

	// menentukan letak upload folder
	fileLocation := filepath.Join("/var/www/vhosts/zcomeducation.com/httpdocs/upload/video", avatar.Filename)
	dst, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		return err
	}

	fmt.Println(fileLocation)
	fmt.Println(dst)
	
	defer dst.Close()
	
	//  Copy
	_, err = io.Copy(dst, src)

	if err != nil {
		return err
	}

	/*===============================================================================
	| UPLOAD KE JWPLAYER DASHBOARD
	================================================================================*/
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// set URL params
	params := url.Values{}
	params.Set("title", video.NamaVideo)

	apiKey := "vy9IrF5y"
	apiSecret := "9PvR1GzIGk6Q2wWHUJNLdPYY"

	client := jwplatform.NewClient(apiKey, apiSecret)

	// declare an empty interface
	var hasil map[string]interface{}

	// upload video using direct upload method
	err = client.Upload(ctx, fileLocation, params, &hasil)

	if err != nil {
		fmt.Println("Gagal Upload ke JWP")
		log.Fatal(err)
	}

	videoKey := hasil["media"].(map[string]interface{})["key"]

	/*===============================================================================
	| UPLOAD KE JWPLAYER DASHBOARD
	================================================================================*/

	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return
	}
	
    defer db.Close()

	// Define Random ID
	videoID := uuid.New()

	_, err = db.Exec(`INSERT INTO 
								kursus_video (id, kursus_id, video_url, nama_video) 
						VALUES (?, ?, ?, ?)`, videoID.String(), &video.KursusID, videoKey, &video.NamaVideo)
    
	if err != nil {
		fmt.Println("gagal insert video_url ke database !")
        fmt.Println(err.Error())
        return
    }

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "berhasil",
		"location": fileLocation,
		"dataVideo": video,
	})

}

// Hapus video
func DeleteVideo(c echo.Context) (err error) {

	// definisi variabel video dengan struct video
	video := new(KursusVideo)

	// binding data inputan ke struct video
	if err = c.Bind(video); err != nil {
		return err
	}

	// koneksi database
	db, err := connect()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// tutup koneksi database
	defer db.Close()

	// get id by parameter
	videoID := c.Param("id")

	if err != nil {
		return err
	}

	// ambil gambar dari database sebelum dihapus
	err = db.
		QueryRow(`
				SELECT id, 
					kursus_id, 
					video_url, 
					nama_video 
				FROM 
					kursus_video 
				WHERE id = ?`, videoID).
		Scan(&video.ID, &video.KursusID, &video.VideoUrl, &video.NamaVideo)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	/*===============================================================================
	| DELETE DARI KE JWPLAYER DASHBOARD
	================================================================================*/
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	apiKey := "vy9IrF5y"
	apiSecret := "9PvR1GzIGk6Q2wWHUJNLdPYY"

	client := jwplatform.NewClient(apiKey, apiSecret)

	// set URL params
	params := url.Values{}
	params.Set("video_key", video.VideoUrl)

	// declare an empty interface
	var result map[string]interface{}

	err = client.MakeRequest(ctx, http.MethodGet, "/videos/delete/", params, &result)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result["status"])  // ok

	/*===============================================================================
	| DELETE DARI JWPLAYER DASHBOARD
	================================================================================*/

	// Menghapus Video 
	// menentukan letak folder
	fileLocation := filepath.Join("/var/www/vhosts/zcomeducation.com/httpdocs/upload/video/", video.VideoUrl)

	// hapus gambar
	hapusFile(fileLocation)
	// End Of Menghapus Gambar

	// hapus data dari db
	_, err = db.Exec("DELETE FROM kursus_video WHERE id = ?", videoID)
	
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "berhasil",
		"datavideo": video,
	})
}

// Menampilkan data Video berdasarkan ID
func GetVideoByID(c echo.Context) (err error) {
	
	video := new(KursusVideo)
	id := c.Param("id")

	if err = c.Bind(video); err != nil {
		return
	}
	
	db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    defer db.Close()

	err = db.QueryRow(`SELECT 
								id, 
								kursus_id, 
								video_url, 
								nama_video 
							FROM 
								kursus_video 
							WHERE 
								id = ?`, id).
								Scan(&video.ID, &video.KursusID, &video.VideoUrl, &video.NamaVideo)
    
	if err != nil {
        fmt.Println(err.Error())
        return
    }

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "berhasil",
	})

}

// Update video
func UpdateVideo(c echo.Context) (err error) {

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer db.Close()

	id := c.Param("id")
	video := new(KursusVideo)
	
	if err = c.Bind(video); err != nil {
		return err
	}

	// validasi inputan
	if err := c.Validate(video); err != nil {
		return err
	}


	fmt.Println(id)
	//** Start of File Upload

	// Multipart form
	// Get avatar
	avatar, err := c.FormFile("video_url")

	if err != nil {
		
		fmt.Println("tidak ada video")

		_, err = db.Exec(`UPDATE 
								kursus_video 
							SET 
								kursus_id = ?, 
								nama_video = ? 
							WHERE id = ?`, &video.KursusID, &video.NamaVideo, id)
		
		if err != nil {
			fmt.Println("gagal update video tanpa video upload" + err.Error())
			return err
		}

		err = db.
			QueryRow(`SELECT 
							id, 
							kursus_id, 
							nama_video, 
							video_url 
						FROM 
							kursus_video 
						WHERE id = ?`, id).
			Scan(&video.ID, &video.KursusID, &video.NamaVideo, &video.VideoUrl)

		if err != nil {
			fmt.Println("gagal ambil video setelah update ")
			fmt.Println(err)
			return err
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":     "berhasil",
			"dataKursus": video,
		})

	}

	fmt.Println(avatar.Filename)
	fmt.Println(id)
	// hapus video lama
	// select by id
	err = db.
		QueryRow("SELECT id, kursus_id, nama_video, video_url WHERE id = ?", id).
		Scan(&video.ID, &video.KursusID, &video.NamaVideo, &video.VideoUrl)

	if err != nil {
		fmt.Println("gagal get data dengan video upload")
		fmt.Println(err.Error)
		return err
	}

	// menentukan letak upload folder
	fileLocation := filepath.Join("/var/www/vhosts/zcomeducation.com/httpdocs/upload/video", video.VideoUrl)

	// hapus gambar
	hapusFile(fileLocation)

	// kemudian perbaharui gambar
	// sumber upload
	src, err := avatar.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// menentukan letak upload folder
	fileLocation = filepath.Join("/var/www/vhosts/zcomeducation.com/httpdocs/upload/video", avatar.Filename)
	//fileLocation = filepath.Join(dir, "/upload/foto", avatar.Filename)
	dst, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	/*===============================================================================
	| UPLOAD KE JWPLAYER DASHBOARD
	================================================================================*/
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// set URL params
	params := url.Values{}
	params.Set("title", video.NamaVideo)

	apiKey := "vy9IrF5y"
	apiSecret := "9PvR1GzIGk6Q2wWHUJNLdPYY"

	client := jwplatform.NewClient(apiKey, apiSecret)

	// declare an empty interface
	var hasil map[string]interface{}

	// upload video using direct upload method
	err = client.Upload(ctx, fileLocation, params, &hasil)

	if err != nil {
		fmt.Println("Gagal Upload ke JWP")
		log.Fatal(err)
	}

	videoKey := hasil["media"].(map[string]interface{})["key"]

	/*===============================================================================
	| UPLOAD KE JWPLAYER DASHBOARD
	================================================================================*/
	
	//** End of File Upload

	// perbaharui data
	_, err = db.Exec("UPDATE kursus_video SET kursus_id = ?, nama_video = ?, video_url = ? WHERE id = ?", &video.KursusID, &video.NamaVideo, videoKey, id)
	
	if err != nil {
		fmt.Println("gagal update video + upload" + err.Error())
		return err
	}

	// select by id
	err = db.
		QueryRow("SELECT id, kursus_id, nama_video, video_url WHERE id = ?", id).
		Scan(&video.ID, &video.KursusID, &video.NamaVideo, &video.VideoUrl)

	if err != nil {
		fmt.Println("gagal get data" + id)
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":     "berhasil",
		"dataKursus": video,
	})
	
}

