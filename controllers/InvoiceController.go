package controllers

import (
	"github.com/labstack/echo"
	"fmt"
	"net/http"

    "io/ioutil"
    "log"
    "os"
    "strings"

	"gopkg.in/gomail.v2"
    uuid "github.com/google/uuid"
)

type UserOrder struct  {
    ID              string      `form:"id"                  json:"invoice_id"`
    UserID          string      `form:"user_id"             json:"user_id"`
    NamaLengkap     string      `form:"nama_lengkap"        json:"nama_lengkap"`
    Alamat          string      `form:"alamat"              json:"alamat"`
    NomorHP         string      `form:"hp"                  json:"hp"`
    KursusID        string      `form:"kursus_id"           json:"kursus_id"`
    NamaKursus      string      `form:"nama_kursus"         json:"nama_kursus"`
    Harga           int         `form:"harga"               json:"harga"`
    Diskon          int         `form:"diskon"              json:"diskon"`
    HargaDiskon     int         `form:"harga_diskon"        json:"harga_diskon"`
    Foto            string      `form:"foto"                json:"foto"`
    PaymentID       int         `form:"metode_pembayaran"   json:"metode_pembayaran"`
    NamaPayment     string      `form:"nama_payment"        json:"pembayaran_via"`
    Status          int         `form:"status"              json:"status"`
    CreatedAt       string      `form:"tanggal_pembelian"   json:"tanggal_pembelian"`
}

// Pembelian
func BeliKursus(c echo.Context) (err error) {

    order := new(UserOrder)
    if err = c.Bind(order); err != nil {
        fmt.Println("Gagal binding objek")
        return
    }

    db, err := connect()
    if err != nil {
        fmt.Println("Tidak dapat terkoneksi ke database !")
        return err
    }

    defer db.Close()

    id := uuid.New()

    var (
        harga           int
        hargaDiskon     int
        diskon          int  
    )

    db.QueryRow(`
                SELECT
                   harga,
                   harga_diskon
                FROM 
                    kursus_konten
                WHERE
                    id = ?
                ORDER BY id DESC LIMIT 1
    `,&order.KursusID).Scan(&harga, &diskon)

    //* Generate Diskon
        potongan := harga * diskon / 100
        hargaDiskon = harga - potongan
    //* Generate Diskon

    _, err = db.Exec(`
                INSERT INTO 
                    orders (
                        id,
                        user_id,
                        kursus_id,
                        harga_diskon,
                        payment_id,
                        status
                    ) VALUES (?, ?, ?, ?, ?, 1)
    `, id.String(), &order.UserID, &order.KursusID, &hargaDiskon, &order.PaymentID)
    if err != nil {
        fmt.Println("gagal insert ke database")
        return err
    }
    err = db.QueryRow(`
                SELECT
                    orders.id,
                    orders.user_id,
                    users.nama_lengkap,
                    users.hp,
                    users.alamat,
                    orders.kursus_id,
                    kursus_konten.nama_kursus,
                    kursus_konten.harga,
                    kursus_konten.harga_diskon,
                    orders.harga_diskon,
                    orders.foto,
                    orders.payment_id,
                    payment.nama_payment,
                    orders.status,
                    orders.created_at
                FROM 
                    orders,
                    users,
                    kursus_konten,
                    payment
                WHERE
                    orders.user_id = users.id 
                AND
                    orders.payment_id = payment.id
                AND 
                    orders.kursus_id = kursus_konten.id
                ORDER BY orders.id DESC limit 1
    `).Scan(&order.ID, &order.UserID, &order.NamaLengkap, &order.NomorHP, &order.Alamat, &order.KursusID, &order.NamaKursus, &order.Harga, &order.Diskon, &order.HargaDiskon, &order.Foto, &order.PaymentID, &order.NamaPayment, &order.Status, &order.CreatedAt)

    if err != nil {
        fmt.Println("gagal get data setelah insert")
        return err
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "status": true,
        "pesan": "Kelas berhasil dibeli !",
    })
}

func GetAllOrder(c echo.Context) (err error) {
    order := new(UserOrder)
    if err = c.Bind(order); err != nil {
        return
    }

    db, err := connect()
    if err != nil {
        fmt.Println("Tidak dapat terkoneksi ke database !")
        return err
    }

    defer db.Close()

    rows, err := db.Query(`
                    SELECT
                        orders.id,
                        orders.user_id,
                        users.nama_lengkap,
                        users.hp,
                        users.alamat,
                        orders.kursus_id,
                        kursus_konten.nama_kursus,
                        kursus_konten.harga,
                        kursus_konten.harga_diskon,
                        orders.harga_diskon,
                        orders.foto,
                        orders.payment_id,
                        payment.nama_payment,
                        orders.status,
                        orders.created_at
                    FROM 
                        orders,
                        users,
                        kursus_konten,
                        payment
                    WHERE
                        orders.user_id = users.id 
                    AND
                        orders.payment_id = payment.id
                    AND 
                        orders.kursus_id = kursus_konten.id`)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    defer rows.Close()

    var result []UserOrder
    
    for rows.Next() {
        var each = UserOrder{}
        var err = rows.Scan(&each.ID, &each.UserID, &each.NamaLengkap, &each.NomorHP, &each.Alamat, &each.KursusID, &each.NamaKursus, &each.Harga, &each.Diskon, &each.HargaDiskon, &each.Foto, &each.PaymentID, &each.NamaPayment, &each.Status, &each.CreatedAt)

        if err != nil {
            fmt.Println(err.Error())
            return err
        }

        result = append(result, each)
    }

    if err = rows.Err(); err != nil {
        fmt.Println(err.Error())
        return
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "status":"berhasil",
        "dataInvoice": result,
    })
}

func GetOrderByUserID(c echo.Context) (err error) {

    order := new(UserOrder)
    if err = c.Bind(order); err != nil {
        return
    }

    db, err := connect()
    if err != nil {
        fmt.Println("Tidak dapat terkoneksi ke database !")
        return err
    }

    defer db.Close()

    userID := c.Param("id")

    rows, err := db.Query(`
                    SELECT
                        orders.id,
                        orders.user_id,
                        users.nama_lengkap,
                        users.hp,
                        users.alamat,
                        orders.kursus_id,
                        kursus_konten.nama_kursus,
                        kursus_konten.harga,
                        kursus_konten.harga_diskon,
                        orders.harga_diskon,
                        orders.foto,
                        orders.payment_id,
                        payment.nama_payment,
                        orders.status,
                        orders.created_at
                    FROM 
                        orders,
                        users,
                        kursus_konten,
                        payment
                    WHERE
                        orders.user_id = users.id 
                    AND
                        orders.payment_id = payment.id
                    AND 
                        orders.kursus_id = kursus_konten.id
                    AND 
                        orders.user_id = ?`, userID)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    defer rows.Close()

    var result []UserOrder
    
    for rows.Next() {
        var each = UserOrder{}
        var err = rows.Scan(&each.ID, &each.UserID, &each.NamaLengkap, &each.NomorHP, &each.Alamat, &each.KursusID, &each.NamaKursus, &each.Harga, &each.Diskon, &each.HargaDiskon, &each.Foto, &each.PaymentID, &each.NamaPayment, &each.Status, &each.CreatedAt)

        if err != nil {
            fmt.Println(err.Error())
            return err
        }

        result = append(result, each)
    }

    if err = rows.Err(); err != nil {
        fmt.Println(err.Error())
        return
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "status":"berhasil",
        "dataInvoice": result,
    })
}

func UpdateOrder(c echo.Context) (err error) {

    order := new(UserOrder)
    if err = c.Bind(order); err != nil {
        return
    }

    db, err := connect()
    if err != nil {
        fmt.Println("Tidak dapat terkoneksi ke database")
        return
    }

     //** Start of File Upload
	//------------
	// Read files
	//------------

	// Multipart form
	// Get foto
	foto, err := c.FormFile("foto")
	if err != nil {
		return err
	}

	// Source
	src, err := foto.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	//** End of File Upload

	// Generate Random File Name
	file, err := ioutil.TempFile("/var/www/vhosts/zcomeducation.com/httpdocs/upload/bukti_pembayaran", "bp-*.png")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fileBytes, err := ioutil.ReadAll(src)
    if err != nil {
        fmt.Println(err)
    }
    
    file.Write(fileBytes)

	os.Chown(file.Name(), 10000, 1004)
	os.Chmod(file.Name(), 0644)

	strFile := strings.ReplaceAll(file.Name(), "/var/www/vhosts/zcomeducation.com/httpdocs/upload/bukti_pembayaran/", "")

    idInvoice := c.Param("id")

    // Update foto di database
    _, err = db.Exec(`
            UPDATE orders 
                SET 
                    foto = ?,
                    status = ?
                WHERE 
                    id = ?
    `, strFile, 2, idInvoice)
    if err != nil {
        fmt.Println("gagal upload bukti pembayaran")
        return
    }

    // Select data order berdasarkan data yang telah diupdate
    err = db.QueryRow(`
                SELECT
                    orders.id,
                    orders.user_id,
                    users.nama_lengkap,
                    users.hp,
                    users.alamat,
                    orders.kursus_id,
                    kursus_konten.nama_kursus,
                    kursus_konten.harga,
                    kursus_konten.harga_diskon,
                    orders.harga_diskon,
                    orders.foto,
                    orders.payment_id,
                    payment.nama_payment,
                    orders.status,
                    orders.created_at
                FROM 
                    orders,
                    users,
                    kursus_konten,
                    payment
                WHERE
                    orders.user_id = users.id 
                AND
                    orders.payment_id = payment.id
                AND 
                    orders.kursus_id = kursus_konten.id
                AND 
                    orders.id = ?
    `, idInvoice).Scan(&order.ID, &order.UserID, &order.NamaLengkap, &order.NomorHP, &order.Alamat, &order.KursusID, &order.NamaKursus, &order.Harga, &order.Diskon, &order.HargaDiskon, &order.Foto, &order.PaymentID, &order.NamaPayment, &order.Status, &order.CreatedAt)

    if err != nil {
        fmt.Println("gagal get data setelah insert")
        return err
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "status": true,
        "pesan" : "Bukti pembayaran berhasil diupload !",
        "data" : order,
        "foto": strFile,
    })

}

func CancelOrder(c echo.Context) (err error) {

    order := new(UserOrder)
    if err = c.Bind(order); err != nil {
        return
    }

    db, err := connect()
    if err != nil {
        fmt.Println("Tidak dapat terkoneksi ke database")
        return
    }

    /*
        @Param("id")
    */

    idInvoice := c.Param("id")
    
    // Update pembatalan
    _, err = db.Exec(`
            UPDATE orders 
                SET 
                    status = ?
                WHERE 
                    id = ?
    `, 4, idInvoice)
    if err != nil {
        fmt.Println("tidak bisa dibatalkan")
        return
    }

    // Select data order berdasarkan data yang telah diupdate
    err = db.QueryRow(`
                SELECT
                    orders.id,
                    orders.user_id,
                    users.nama_lengkap,
                    users.hp,
                    users.alamat,
                    orders.kursus_id,
                    kursus_konten.nama_kursus,
                    kursus_konten.harga,
                    kursus_konten.harga_diskon,
                    orders.harga_diskon,
                    orders.foto,
                    orders.payment_id,
                    payment.nama_payment,
                    orders.status,
                    orders.created_at
                FROM 
                    orders,
                    users,
                    kursus_konten,
                    payment
                WHERE
                    orders.user_id = users.id 
                AND
                    orders.payment_id = payment.id
                AND 
                    orders.kursus_id = kursus_konten.id
                AND 
                    orders.id = ?
    `, idInvoice).Scan(&order.ID, &order.UserID, &order.NamaLengkap, &order.NomorHP, &order.Alamat, &order.KursusID, &order.NamaKursus, &order.Harga, &order.Diskon, &order.HargaDiskon, &order.Foto, &order.PaymentID, &order.NamaPayment, &order.Status, &order.CreatedAt)

    if err != nil {
        fmt.Println("gagal get data setelah update")
        return err
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "status": true,
        "pesan" : "Pembelian berhasil dibatalkan !",
        "data" : order,
    })

}

func ConfirmOrder(c echo.Context) (err error) {

    order := new(UserOrder)
    if err = c.Bind(order); err != nil {
        return
    }

    db, err := connect()
    if err != nil {
        fmt.Println("Tidak dapat terkoneksi ke database")
        return
    }

    idInvoice := c.Param("id")

    // Update foto di database
    _, err = db.Exec(`
            UPDATE orders 
                SET 
                    status = ?
                WHERE 
                    id = ?
    `, 3, idInvoice)
    if err != nil {
        fmt.Println("konfirmasi gagal")
        return
    }

    // Select data order berdasarkan data yang telah diupdate
    err = db.QueryRow(`
                SELECT
                    orders.id,
                    orders.user_id,
                    users.nama_lengkap,
                    users.hp,
                    users.alamat,
                    orders.kursus_id,
                    kursus_konten.nama_kursus,
                    kursus_konten.harga,
                    kursus_konten.harga_diskon,
                    orders.harga_diskon,
                    orders.foto,
                    orders.payment_id,
                    payment.nama_payment,
                    orders.status,
                    orders.created_at
                FROM 
                    orders,
                    users,
                    kursus_konten,
                    payment
                WHERE
                    orders.user_id = users.id 
                AND
                    orders.payment_id = payment.id
                AND 
                    orders.kursus_id = kursus_konten.id
                AND 
                    orders.id = ?
    `, idInvoice).Scan(&order.ID, &order.UserID, &order.NamaLengkap, &order.NomorHP, &order.Alamat, &order.KursusID, &order.NamaKursus, &order.Harga, &order.Diskon, &order.HargaDiskon, &order.Foto, &order.PaymentID, &order.NamaPayment, &order.Status, &order.CreatedAt)

    if err != nil {
        fmt.Println("gagal get data setelah update")
        return err
    }

    // Define id dengan random string
	id := uuid.New()

	_, err = db.Exec(`
                INSERT INTO kelas (
                    idkelas, 
                    user_id, 
                    kursus_id
                ) VALUES (?, ?, ?)`, id.String(), &order.UserID, &order.KursusID)
		
	if err != nil {
		fmt.Println("gagal insert kelas")
		fmt.Println(err.Error())
		return err
	}

    /*
    |   SQL Query for SELECT Excecution from Database
                                                    */
    var (
        InvoiceID   string
        IDUser      string
        NamaUser    string
        Email       string 
        Tanggal     string
        IDPay       string
        StatusPay   string
        NamaPay     string
        IDKursus    string
        NamaKurs    string
        Harga       int

    )
    userId := c.Param("id")
    err = db.QueryRow(`
            SELECT 
                orders.id,
                orders.user_id,
                users.nama_lengkap,
                users.email,
                orders.created_at,
                orders.status,
                orders.payment_id,
                payment.nama_payment,
                orders.kursus_id,
                kursus_konten.nama_kursus,
                kursus_konten.harga
            FROM 
                orders,
                kursus_konten,
                payment
            WHERE 
                orders.payment_id = payment.id
            AND 
                kursus_id = kursus_konten.id
            AND
                orders.user_id = ?
    `, &userId).
    Scan(&InvoiceID, &IDUser, &NamaUser, &Email, &Tanggal, &StatusPay, &IDPay, &NamaPay, &IDKursus, &NamaKurs, &Harga)
    if err != nil {
        fmt.Println("gagal get data order")
        return
    }
    /* ================================ #
    |         EMAIL SENDER CONFIG       |
    #  ================================ */
    const CONFIG_SMTP_HOST = "smtp.gmail.com"
    const CONFIG_SMTP_PORT = 587
    const CONFIG_EMAIL = "your.email.aja@gmail.com"
    const CONFIG_PASSWORD = "Makannasi"

    /* =============================== # 
    |   SEND INVOICE PAYMENT TO USER   |
    #  =============================== */
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_EMAIL)
	mailer.SetHeader("To", "{{.Email}}")
	mailer.SetAddressHeader("Cc", "rocker.hunt@gmail.com", "Tra Lala La")
	mailer.SetHeader("Subject", "Test mail")
	mailer.SetBody("text/html", `<div style="width: 80%;margin: auto;">
      <div style="width: 100%;
      height: 200px;
      background-color: #004d40; margin: auto;">
        <img
          src="http://zcomeducation.com/static/media/logo.23e672b2.png"
          alt="logo"
          height="150px"
          style="margin-left:30%;margin-right: 30%;"
        />
      </div>
      <div class="body" style="color: gray;">
        <h1 align="center">Pembayaran Berhasil</h1>
        <h2>Halo {{.NamaUser}},</h2>
        <p>Terima kasih sudah melakukan pembelian kelas di <a href="http://zcomeducation.com">zcomeducation</a> Pembayaran kamu telah <b>sukses</b></p>
        <h2 align="center">Insentif Pemerintah</h2>
        <p>Segera selesaikan kelas kamu untuk mendapatkan insentif berikut
            <ol>
                <li>Insentif pelatihan sebesar Rp. 600.000,-/ bulan ( selama 4 bulan )</li>
                <li>Insentif survey Kebekerjaan sebesar Rp. 50.000/survey ( akan ada 3 survey )</li>
            </ol>
            Insentif akan diberikan setelah menyelesaikan kelas pertama dan hanya dapat di klaim 1x untuk setiap pengguna Kartu Prakerja
        </p>
        <br>
        <h2 align="center">Invoice #INV-"{{.InvoiceID}}"</h2>
        <p align="center">Tanggal Pembelian : {{.Tanggal}}</p>
        <h2>Status Pembayaran</h2>
        <h3>{{.StatusPay}}</h3>
        <h2>Detail Pembelian</h2>
        <table>
            <tr>
                <td align="left">Total yang telah dibayar</td>
                <td>: Rp. {{.Harga}}</td>
            </tr>
            <tr>
                <td align="left">Metode Pembayaran</td>
                <td>: {{.NamaPay}}</td>
            </tr>
            <tr>
                <td align="left">Kelas</td>
                <td>: {{.NamaKurs}}</td>
            </tr>
        </table>
        <div style="width: 100%; position: relative;">
          <div style="width: 45%; float: left;
            padding: 10px;
              border: 1px solid gray;">
              <p>Untuk masalah dan pertanyaan terkait insentif dari pemerintah dan Kartu Prakerja, hubungi <br>
               <br>
                021-25541246 <br>
              <a href="https://info@prakerja.go.id">info@prakerja.go.id</a>
              </p>
          </div>
          <div style="width: 45%;float:right;
                padding: 10px;
                border: 1px solid gray;">
            <p>Untuk masalah dan pertanyaan terkait akses akun, materi, dan video belajar zcomeducation.com, hubungi <br>
                <br>
                021-25541246 <br>
                <a href="https://zcomeducation.com/info">info@zcomeducation.com</a>
            </p>
          </div>
        </div>
        <div style="position: relative;">
          <br>
          <br>
          <br>
          <br>
          <br>
          <br>
          <br>
          <br>
          <br>
          <p>Selamat belajar, {{.NamaUser}}</p>
          <br><br><br>
          <p>Salam hangat,<br>
            Team zcomeducation.com</p>
            <br><br>
            <p align="center">&copy;2020 zcomeducation.com</p>
        </div>
          </div>
    </div>`)
	mailer.Attach("./invoice.html")

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_EMAIL,
		CONFIG_PASSWORD,
	)

	err = dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatal(err.Error())
	}

    log.Println("Mail sent!")
    return c.JSON(http.StatusOK, map[string]interface{}{
        "status": "terkirim",
    })

    return c.JSON(http.StatusOK, map[string]interface{}{
        "status": true,
        "pesan" : "Konfirmasi Berhasil !",
    })

}

func InvoiceMail(c echo.Context) (err error) {

    /*                             
    |   Data Binding to New Struct  |
                                    */
    inv := new(UserOrder)
    if err = c.Bind(inv); err != nil {
        return
    }

    /*
    |   Database Binding from connect() function
                                                */
    db, err := connect()
    if err != nil {
        fmt.Println("Tidak dapat terhubung ke database !")
        return
    }

    defer db.Close()

    /*
    |   SQL Query for SELECT Excecution from Database
                                                    */
    var (
        InvoiceID   string
        IDUser      string
        NamaUser    string
        Email       string 
        Tanggal     string
        IDPay       string
        StatusPay   string
        NamaPay     string
        IDKursus    string
        NamaKurs    string
        Harga       int

    )
    userId := c.Param("id")
    err = db.QueryRow(`
            SELECT 
                orders.id,
                orders.user_id,
                users.nama_lengkap,
                users.email,
                orders.created_at,
                orders.status,
                orders.payment_id,
                payment.nama_payment,
                orders.kursus_id,
                kursus_konten.nama_kursus,
                kursus_konten.harga
            FROM 
                orders,
                kursus_konten,
                payment
            WHERE 
                orders.payment_id = payment.id
            AND 
                kursus_id = kursus_konten.id
            AND
                orders.user_id = ?
    `, &userId).
    Scan(&InvoiceID, &IDUser, &NamaUser, &Email, &Tanggal, &StatusPay, &IDPay, &NamaPay, &IDKursus, &NamaKurs, &Harga)
    if err != nil {
        fmt.Println("gagal get data order")
        return
    }


    /* ================================ #
    |         EMAIL SENDER CONFIG       |
    #  ================================ */
    const CONFIG_SMTP_HOST = "smtp.gmail.com"
    const CONFIG_SMTP_PORT = 587
    const CONFIG_EMAIL = "your.email.aja@gmail.com"
    const CONFIG_PASSWORD = "Makannasi"

    /* =============================== # 
    |   SEND INVOICE PAYMENT TO USER   |
    #  =============================== */
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_EMAIL)
	mailer.SetHeader("To", "{{.Email}}")
	mailer.SetAddressHeader("Cc", "rocker.hunt@gmail.com", "Tra Lala La")
	mailer.SetHeader("Subject", "Test mail")
	mailer.SetBody("text/html", `<div style="width: 80%;margin: auto;">
      <div style="width: 100%;
      height: 200px;
      background-color: #004d40; margin: auto;">
        <img
          src="http://zcomeducation.com/static/media/logo.23e672b2.png"
          alt="logo"
          height="150px"
          style="margin-left:30%;margin-right: 30%;"
        />
      </div>
      <div class="body" style="color: gray;">
        <h1 align="center">Pembayaran Berhasil</h1>
        <h2>Halo {{.NamaUser}},</h2>
        <p>Terima kasih sudah melakukan pembelian kelas di <a href="http://zcomeducation.com">zcomeducation</a> Pembayaran kamu telah <b>sukses</b></p>
        <h2 align="center">Insentif Pemerintah</h2>
        <p>Segera selesaikan kelas kamu untuk mendapatkan insentif berikut
            <ol>
                <li>Insentif pelatihan sebesar Rp. 600.000,-/ bulan ( selama 4 bulan )</li>
                <li>Insentif survey Kebekerjaan sebesar Rp. 50.000/survey ( akan ada 3 survey )</li>
            </ol>
            Insentif akan diberikan setelah menyelesaikan kelas pertama dan hanya dapat di klaim 1x untuk setiap pengguna Kartu Prakerja
        </p>
        <br>
        <h2 align="center">Invoice #INV-"{{.InvoiceID}}"</h2>
        <p align="center">Tanggal Pembelian : {{.Tanggal}}</p>
        <h2>Status Pembayaran</h2>
        <h3>{{.StatusPay}}</h3>
        <h2>Detail Pembelian</h2>
        <table>
            <tr>
                <td align="left">Total yang telah dibayar</td>
                <td>: Rp. {{.Harga}}</td>
            </tr>
            <tr>
                <td align="left">Metode Pembayaran</td>
                <td>: {{.NamaPay}}</td>
            </tr>
            <tr>
                <td align="left">Kelas</td>
                <td>: {{.NamaKurs}}</td>
            </tr>
        </table>
        <div style="width: 100%; position: relative;">
          <div style="width: 45%; float: left;
            padding: 10px;
              border: 1px solid gray;">
              <p>Untuk masalah dan pertanyaan terkait insentif dari pemerintah dan Kartu Prakerja, hubungi <br>
               <br>
                021-25541246 <br>
              <a href="https://info@prakerja.go.id">info@prakerja.go.id</a>
              </p>
          </div>
          <div style="width: 45%;float:right;
                padding: 10px;
                border: 1px solid gray;">
            <p>Untuk masalah dan pertanyaan terkait akses akun, materi, dan video belajar zcomeducation.com, hubungi <br>
                <br>
                021-25541246 <br>
                <a href="https://zcomeducation.com/info">info@zcomeducation.com</a>
            </p>
          </div>
        </div>
        <div style="position: relative;">
          <br>
          <br>
          <br>
          <br>
          <br>
          <br>
          <br>
          <br>
          <br>
          <p>Selamat belajar, {{.NamaUser}}</p>
          <br><br><br>
          <p>Salam hangat,<br>
            Team zcomeducation.com</p>
            <br><br>
            <p align="center">&copy;2020 zcomeducation.com</p>
        </div>
          </div>
    </div>`)
	mailer.Attach("./invoice.html")

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_EMAIL,
		CONFIG_PASSWORD,
	)

	err = dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatal(err.Error())
	}

    log.Println("Mail sent!")
    return c.JSON(http.StatusOK, map[string]interface{}{
        "status": "terkirim",
    })
}
