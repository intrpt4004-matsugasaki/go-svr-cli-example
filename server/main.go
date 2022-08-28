package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Identifier string `gorm:"size:255;unique;not null" binding:"required"`
	Password   string `gorm:"size:255;not null" binding:"required"`
}

type Token struct {
	gorm.Model
	Identifier string `gorm:"size:255;not null" binding:"required"`
	Token      string `gorm:"size:255;unique;not null" binding:"required"`
}

type Post struct {
	gorm.Model
	Identifier string `gorm:"size:255;not null" binding:"required"`
	Post       string `gorm:"size:255;not null" binding:"required"`
}

type PostResponse struct {
	Id   string `json:"id"`
	Post string `json:"post"`
}

func main() {
	/* DB接続 */
	db, err := gorm.Open(mysql.Open("root:password@tcp(localhost:3306)/go_svr_cli_sample_db?parseTime=true"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Token{})
	db.AutoMigrate(&Post{})

	// 3分毎に全トークン無効化
	ticker := time.NewTicker(3 * time.Minute)
	go func() {
		for {
			<-ticker.C

			var tkns []Token
			db.Model(&Token{}).Find(&tkns)
			for _, tkn := range tkns {
				db.Unscoped().Delete(&tkn)
			}
		}
	}()

	/* REST API 構築 */
	engine := gin.Default()

	/** アカウント登録 **/
	engine.POST("/signup", func(c *gin.Context) {
		// Id,Pw長さ検査
		if len(c.PostForm("id")) < 4 || len(c.PostForm("pw")) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "id > 3 and pw > 5.",
			})
			return
		}

		// パスワードをハッシュ化
		pwhash, _ := bcrypt.GenerateFromPassword([]byte(c.PostForm("pw")), bcrypt.DefaultCost)

		// DBに登録
		if result := db.Create(&User{
			Identifier: c.PostForm("id"),
			Password:   string(pwhash),
		}); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "the identifier already exists.",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "signup success",
		})
	})

	/** アカウント認証 **/
	engine.POST("/auth", func(c *gin.Context) {
		// 一致ID検索
		var users []User
		if result := db.Model(&User{}).Where("Identifier = ?", c.PostForm("id")).Find(&users); result.RowsAffected == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "bad identifier",
			})
			return
		}

		// パスワードのハッシュ値を比較
		if err := bcrypt.CompareHashAndPassword(
			[]byte(users[0].Password),
			[]byte(c.PostForm("pw")),
		); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "auth failed.",
			})
			return
		}

		// トークン生成
		var tkn string

		// 乱数生成
		src := make([]byte, 10)
		rand.Read(src)

		const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		for _, digit := range src {
			tkn += string(letters[int(digit)%len(letters)])
		}

		// トークンDB登録
		if result := db.Create(&Token{
			Token:      tkn,
			Identifier: c.PostForm("id"),
		}); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "auth again.",
			})
			return
		}

		// トークン払い出し
		c.JSON(http.StatusOK, gin.H{
			"token": tkn,
		})
	})

	/** 時刻取得 **/
	engine.GET("/time", func(c *gin.Context) {
		// トークン認証
		var tokens []Token
		if result := db.Model(&Token{}).Where("Token = ?", c.Request.Header.Get("token")).Find(&tokens); result.RowsAffected == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "illegal token",
			})
			return
		}

		if c.Query("format") == "unix" {
			c.JSON(http.StatusOK, gin.H{
				"time": strconv.FormatInt(time.Now().Unix(), 10),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"time": time.Now(),
		})
	})
	engine.GET("/time/:format", func(c *gin.Context) {
		// トークン認証
		var tokens []Token
		if result := db.Model(&Token{}).Where("Token = ?", c.Request.Header.Get("token")).Find(&tokens); result.RowsAffected == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "illegal token",
			})
			return
		}

		if c.Param("format") == "unixnano" {
			c.JSON(http.StatusOK, gin.H{
				"time": strconv.FormatInt(time.Now().UnixNano(), 10),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})
	})

	/** UA取得 **/
	ua := ""
	engine.Use(func(c *gin.Context) {
		ua = c.GetHeader("User-Agent")
		c.Next()
	})
	engine.GET("/user-agent", func(c *gin.Context) {
		// トークン認証
		var tokens []Token
		if result := db.Model(&Token{}).Where("Token = ?", c.Request.Header.Get("token")).Find(&tokens); result.RowsAffected == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "illegal token",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user-agent": ua,
		})
	})

	/** 文字列反転 **/
	engine.POST("/reverse", func(c *gin.Context) {
		// トークン認証
		var tokens []Token
		if result := db.Model(&Token{}).Where("Token = ?", c.Request.Header.Get("token")).Find(&tokens); result.RowsAffected == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "illegal token",
			})
			return
		}

		// 文字列反転
		var result []rune
		for i := len(c.PostForm("text")) - 1; i >= 0; i-- {
			result = append(result, []rune(c.PostForm("text"))[i])
		}

		c.JSON(http.StatusOK, gin.H{
			"text": string(result),
		})
	})

	/** 呟き書き込み **/
	engine.POST("/post/create", func(c *gin.Context) {
		// トークン認証
		var tokens []Token
		if result := db.Model(&Token{}).Where("Token = ?", c.Request.Header.Get("token")).Find(&tokens); result.RowsAffected == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "illegal token",
			})
			return
		}

		// post長さ検査
		if len(c.PostForm("post")) < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "post > 0.",
			})
			return
		}

		// tokenからid取得
		id := tokens[0].Identifier
		fmt.Println(id)

		// DBに登録
		if result := db.Create(&Post{
			Identifier: id,
			Post:       c.PostForm("post"),
		}); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "internal error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "post success",
		})
	})

	/** 呟き全取得 **/
	engine.GET("/post/all", func(c *gin.Context) {
		// トークン認証
		var tokens []Token
		if result := db.Model(&Token{}).Where("Token = ?", c.Request.Header.Get("token")).Find(&tokens); result.RowsAffected == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "illegal token",
			})
			return
		}

		// DBから取得
		var posts []Post
		if result := db.Model(&Post{}).Find(&posts); result.Error != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "internal error",
			})
			return
		}

		postsResp := make([]PostResponse, len(posts))
		for i, v := range posts {
			postsResp[i].Id = v.Identifier
			postsResp[i].Post = v.Post
		}

		c.JSON(http.StatusOK, gin.H{
			"posts": postsResp,
		})
	})

	engine.Run(":3000")
}
