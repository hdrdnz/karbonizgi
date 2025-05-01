package controller

import (
	"carbonfootprint/controllers"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func GetQuestions(c *gin.Context) {
	fmt.Println("girdi")
	questionType := c.Query("type")
	if questionType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Lütfen kullanıcı tipi seçiniz",
		})
		return
	}
	fmt.Println("questionType:", questionType)
	var file *os.File
	var err error
	if questionType == "person" {
		file, err = os.Open("./data/person.json")

	} else {
		file, err = os.Open("./data/company2.json")
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Error opening file",
		})
		return
	}
	byteFile, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Error reading file",
		})
		return
	}
	ques := make(map[string]interface{})
	if err := json.Unmarshal(byteFile, &ques); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"eror":    err,
			"message": "Something went worong.",
		})
		return
	}

	fmt.Println("somuçlar gitti")
	c.JSON(http.StatusOK, ques)
}

func GetArticle(c *gin.Context) {
	var infos []controllers.Data
	file, err := os.ReadFile("./data/data.json")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	if err := json.Unmarshal(file, &infos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	rand.Shuffle(len(infos), func(i, j int) {
		infos[i], infos[j] = infos[j], infos[i]
	})
	c.JSON(http.StatusOK, infos)
}

func GetComment(c *gin.Context) {
	data := c.Query("user_type")
	var file *os.File
	var err error
	if data != "" {
		if data == "person" {
			file, err = os.Open("./data/person-question.json")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": "Bir hata oluştu",
				})
				return
			}
		} else if data == "company" {
			file, err = os.Open("./data/company-question.json")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": "Bir hata oluştu",
				})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Geçersiz kullanıcı tipi.",
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Kullanıcı tipi girilmelidir.",
		})
		return
	}

	byteFile, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Error reading file",
		})
		return
	}
	var comments []controllers.Comment
	if err := json.Unmarshal(byteFile, &comments); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"eror":    err,
			"message": "Bir hata oluştu.",
		})
		return
	}

	c.JSON(http.StatusOK, comments)
}

type CommentReq struct {
	Comment  Comment `json:"comment"`
	UserType string  `json:"user_type"`
}

type Comment struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func AddComment(c *gin.Context) {
	fmt.Println("girdi.")
	commentReq := CommentReq{}
	if err := c.ShouldBindJSON(&commentReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  err.Error(),
			"message": "Bir hata oluştu",
		})
		return
	}
	fmt.Println("commentReq:", commentReq)
	if commentReq.Comment.Answer == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Cevap kısmı boş bırakılamaz.",
		})
		return
	}
	if commentReq.Comment.Question == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Soru kısmı boş bırakılamaz.",
		})
		return
	}
	if commentReq.UserType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Kullanıcı tipi boş bırakılamaz.",
		})
		return
	}
	var fileName string
	if commentReq.UserType == "person" {
		fmt.Println("person kısmına girdi.")
		fileName = "./data/person-question.json"
	} else {
		fileName = "./data/company-question.json"
	}
	file, err := os.ReadFile(fileName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Dosya okunamadı",
		})
		return
	}
	var comments []Comment
	if err := json.Unmarshal(file, &comments); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu2",
		})
		return
	}
	fmt.Println("bu ksııam girdis")
	comments = append(comments, commentReq.Comment)
	updatedData, err := json.MarshalIndent(comments, "", "  ")
	if err != nil {
		log.Fatal("JSON dönüştürme hatası:", err)
	}
	if err := os.WriteFile(fileName, updatedData, 0644); err != nil {
		log.Fatal("Yazma hatası:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "başarılı",
	})

}
