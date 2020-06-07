package main

import (
	"zahra/config"
	"zahra/controllers"

	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/go-playground/validator.v9"

	log "github.com/sirupsen/logrus"
)

// Custom Validator
type CustomValidator struct {
	validator *validator.Validate
}

// custom validator methode
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func makeLogEntry(c echo.Context) *log.Entry {
	if c == nil {
		return log.WithFields(log.Fields{
			"at": time.Now().Format("2006-01-02 15:04:05"),
		})
	}

	return log.WithFields(log.Fields{
		"at":     time.Now().Format("2006-01-02 15:04:05"),
		"method": c.Request().Method,
		"uri":    c.Request().URL.String(),
		"ip":     c.Request().RemoteAddr,
	})
}

func middlewareLogging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		makeLogEntry(c).Info("incoming request")
		return next(c)
	}
}

func errorHandler(err error, c echo.Context) {
	report, ok := err.(*echo.HTTPError)
	if ok {
		report.Message = fmt.Sprintf("http error %d - %v", report.Code, report.Message)
	} else {
		report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	makeLogEntry(c).Error(report.Message)
	c.HTML(report.Code, report.Message.(string))
}

var appConfig = config.Config.App
var jwtConfig = config.Config.JWT

func main() {
	e := echo.New()

	e.Use(middlewareLogging)
	e.HTTPErrorHandler = errorHandler

	e.Validator = &CustomValidator{validator: validator.New()}

	// generate cache untuk https
	// e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	/*=====================================
	| User Registration And Authentication|
	======================================*/
	e.POST("/login", controllers.PostLogin)
	e.POST("/register", controllers.Register)
	e.PUT("/aktivasi/:id", controllers.VerifikasiEmail)

	// User handler
	e.POST("/users", controllers.CreateUser)
	e.GET("/users", controllers.GetAllUsers)
	e.GET("/users/:id", controllers.GetUsersById)
	e.GET("/users/bygroupid/:id", controllers.GetUsersByGroupId)
	e.PUT("/users/:id", controllers.UpdateUser)
	e.DELETE("/users/:id", controllers.DeleteUser)

	e.GET("/profile", controllers.GetUserProfile)
	e.GET("/profile/:id", controllers.GetProfileById)
	e.PUT("/profile/:id", controllers.UpdateProfile)
	// e.DELETE("/profile/:id", controllers.DeleteProfile)

	// User Group handler
	e.GET("/usergroup", controllers.GetUsersGroup)
	e.GET("/usergroup/:id", controllers.GetUserGroupById)
	e.POST("/usergroup", controllers.TambahUserGroup)

	// Kategori Kursus /kategorikursus
	e.GET("/kategorikursus", controllers.GetAllKursusKategori)
	e.GET("/kategorikursus/:slug", controllers.GetKursusKategoriBySlug)
	e.POST("/kategorikursus", controllers.TambahKursusKategori)
	e.PUT("/kategorikursus/:id", controllers.UpdateKursusKategori)
	e.DELETE("/kategorikursus/:id", controllers.DeleteKursusKategori)

	// Kursus Konten /kursus
	e.GET("/kursus", controllers.GetAllKursus)
	e.GET("/kursus/:slug", controllers.GetKursusBySlug)
	e.GET("/kursus/bykategorislug/:slug", controllers.GetKursusByKategoriSlug)
	e.GET("/kursus/byuserid/:id", controllers.GetKursusByUserId)
	e.PUT("/kursus/:slug", controllers.UpdateKursus)
	e.POST("/kursus", controllers.TambahKursus)
	e.DELETE("/kursus/:slug", controllers.DeleteKursus)

	// Kursus Video /video
	e.POST("/video", controllers.TambahVideo)
	e.GET("/video", controllers.GetAllVideo)
	e.GET("/video/:id", controllers.GetVideoByID)
	e.GET("/video/bykursusslug/:slug", controllers.GetVideoByKursusSlug)
	e.GET("/video/bykursusid/:id", controllers.GetVideoByKursusID)
	e.PUT("/video/:id", controllers.UpdateVideo)
	e.DELETE("/video/:id", controllers.DeleteVideo)

	// Kategori Post /kategoripost
	e.GET("/kategoripost", controllers.GetAllPostKategori)
	e.GET("/kategoripost/:slug", controllers.GetPostKategoriBySlug)
	e.POST("/kategoripost", controllers.TambahPostKategori)
	e.PUT("/kategoripost/:slug", controllers.UpdatePostKategori)
	e.DELETE("/kategoripost/:slug", controllers.DeletePostKategori)

	// Post Kursus /posts
	e.GET("/artikel", controllers.GetArtikelKursus)
	e.GET("/artikel/:slug", controllers.GetArtikelBySlug)
	e.GET("/artikel/bykategorislug/:slug", controllers.GetArtikelByKategoriSlug)
	e.POST("/artikel", controllers.TambahArtikel)
	e.PUT("/artikel/:slug", controllers.UpdateArtikel)
	e.DELETE("/artikel/:slug", controllers.DeleteArtikel)

	// Promo kursus /promokursus
	e.POST("/promo", controllers.TambahPromo)
	e.GET("/promo", controllers.GetAllPromo)
	e.GET("/promo/:id", controllers.GetPromoById)
	e.PUT("/promo/:id", controllers.UpdatePromo)
	e.DELETE("/promo/:id", controllers.DeletePromo)

	// Soal Ujian
	e.POST("/soal", controllers.TambahUjianSoal)
	e.GET("/soal", controllers.GetAllUjianSoal)
	e.GET("/soal/:id", controllers.GetSoalUjianByID)
	e.GET("/soal/byidkelas/:id", controllers.GetSoalUjianByKursus)
	e.PUT("/soal/:id", controllers.UpdateUjianSoal)
	e.DELETE("/soal/:id", controllers.DeleteUjianSoal)

	// Post Ujian
	e.POST("/ujian", controllers.PostJawaban)
	// e.POST("/nilai", controllers.HitungNilai)
	// Ujian Kursus
	e.GET("/ujian", controllers.GetAllUjian)
	// e.GET("/ujian/hasil", controllers.GetHasilUjian)

	// Kelas Kursus
	e.GET("/kelas", controllers.GetAllKelas)
	e.GET("/kelas/:id", controllers.GetKelasByUserId)
	e.GET("/kelasbyslug/:slug", controllers.GetKelasBySlug)

	// Pembelian Kursus
	e.POST("/orders", controllers.BeliKursus)
	e.PUT("/orders/user/cancel/:id", controllers.CancelOrder)

	// Konfirmasi Admin for Order
	e.PUT("/orders/konfirmasi/:id", controllers.ConfirmOrder)

	// Pemesanan Kursus
	e.GET("/orders", controllers.GetAllOrder)
	e.GET("/orders/byuserid/:id", controllers.GetOrderByUserID)
	e.PUT("/orders/:id", controllers.UpdateOrder)

	e.POST("/invoice/mail", controllers.InvoiceMail)

	// Quis Kursus
	e.POST("/pertanyaan", controllers.PostQuestions)
	e.GET("/pertanyaan", controllers.GetPertanyaan)
	e.GET("/pertanyaan/:id", controllers.GetPertanyaanById)

	// Progres Belajar
	e.POST("/progres", controllers.PostProgres)
	e.GET("/progres/:id", controllers.GetProgresByKelas)

	// Data KUPON
	e.GET("/kupon/:kursusID", controllers.GetKuponByKursusID)
	e.POST("/kupon", controllers.InputKupon)
	e.POST("/kupon/generate", controllers.GenerateKupon)

	// Team
	e.GET("/team", controllers.GetAllTeam)

	// Keunggulan Lembaga
	e.GET("/keunggulanlembaga", controllers.GetAllKeunggulan)

	// Partner Lembaga
	e.GET("/partnerlembaga", controllers.GetPartnerLembaga)

	// Testimonial
	e.GET("/testimonial", controllers.GetTestimoniLembaga)

	// group url /academy/*
	ac := e.Group("/academy")
	admin := e.Group("/admin")
	user := e.Group("/user")

	// * Start of jwt Secret
	// Router dengan JWT
	ac.Use(middleware.JWT([]byte(jwtConfig.Secret)))

	// List kursus dengan jwt
	ac.POST("/kursus", controllers.TambahKursus)
	ac.GET("/kursus", controllers.GetAllKursus)
	// ac.GET("/kursuskategori/:id", controllers.GetKursusByKategoriId)
	ac.GET("/kursus/:id", controllers.GetKursusById)
	ac.PUT("/kursus/:id", controllers.UpdateKursus)
	ac.DELETE("/kursus/:id", controllers.DeleteKursus)

	// Users /academy/users
	ac.GET("/users", controllers.GetAllUsers)
	ac.GET("/users/:id", controllers.GetUsersById)
	ac.GET("/users/bygroupid/:id", controllers.GetUsersByGroupId)
	// ac.PUT("/users/:id", controllers.UpdateUser)
	// ac.DELETE("/users/:id", controllers.DeleteUser)

	// Kategori Kurus /academy/kursus
	ac.GET("/kategorikursus", controllers.GetAllKursusKategori)
	ac.GET("/kategorikursus/:slug", controllers.GetKursusKategoriBySlug)

	// Articles /academy/articles
	// ac.GET("/articles", controllers.AllArticles)
	// ac.POST("/articles", controllers.CreateArticle)
	// ac.GET("/articles/:id", controllers.ShowArticle)
	// ac.PUT("/articles/:id", controllers.UpdateArticle)
	// ac.DELETE("/articles/:id", controllers.DeleteArticle)

	// Soal Handler
	admin.GET("/soal/:id", controllers.GetSoalUjianByID)

	user.GET("soal/:id", controllers.GetSoalUjianByID)

	// * End of jwt secret

	// start server dengan TLS/HTTPS
	// e.Logger.Fatal(e.StartAutoTLS(":1945"))

	// e.Logger.Fatal(e.Start(":1945"))

	lock := make(chan error)
	go func(lock chan error) { lock <- e.Start(":1945") }(lock)

	time.Sleep(1 * time.Millisecond)
	makeLogEntry(nil).Warning("application started without ssl/tls enabled")

	err := <-lock
	if err != nil {
		makeLogEntry(nil).Panic("failed to start application")
	}
}
