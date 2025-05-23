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

type Question struct {
	Key      string `json:"key"`
	Question string `json:"question"`
	Options  Option `json:"options"`
}
type Option struct {
	A OptionKeys `json:"A"`
	B OptionKeys `json:"B"`
	C OptionKeys `json:"C"`
	D OptionKeys `json:"D"`
}
type OptionKeys struct {
	Text     string  `json:"text"`
	Emission float64 `json:"emission"`
}

func GetQuestions(c *gin.Context) {
	fmt.Println("girdi")
	questionType := c.Query("type")
	var fileName string
	if questionType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Lütfen kullanıcı tipi seçiniz",
		})
		return
	}
	fmt.Println("questionType:", questionType)
	if questionType == "person" {
		fileName = "./data/person.json"

	} else {
		fileName = "./data/company2.json"

	}
	file, err := os.Open(fileName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	byteFile, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	fmt.Println("bu kısıma girdi.")
	ques := make(map[string]interface{})
	if err := json.Unmarshal(byteFile, &ques); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	c.JSON(http.StatusOK, ques)
}
func DeleteQuestion(c *gin.Context) {
	questionType := c.PostForm("question_type")
	questionKey := c.PostForm("question_key")
	category := c.PostForm("category")
	var fileName string
	if questionType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Lütfen kullanıcı tipi seçiniz",
		})
		return
	}
	if questionType == "person" {
		fileName = "./data/person.json"

	} else {
		fileName = "./data/company2.json"

	}
	file, err := os.Open(fileName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	byteFile, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	ques := make(map[string]interface{})
	if err := json.Unmarshal(byteFile, &ques); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	var subques []Question
	raw, ok := ques[category]
	if !ok {
		log.Println("diet alanı bulunamadı")
		return
	}
	jsonData, err := json.Marshal(raw)
	if err != nil {
		log.Println("JSON'a çevirirken hata:", err)
		return
	}
	err = json.Unmarshal(jsonData, &subques)
	if err != nil {
		log.Println("JSON çözümlenirken hata:", err)
		return
	}
	found := false
	for subIndex, info := range subques {
		if info.Key == questionKey {
			subques = append(subques[:subIndex], subques[subIndex+1:]...)
			found = true
			break
		}
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Makale bulunamadı",
		})
		return
	}
	ques[category] = subques

	// Güncellenmiş veriyi tekrar JSON'a çevir ve dosyaya yaz
	newData, err := json.MarshalIndent(ques, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "JSON oluşturulamadı",
		})
		return
	}

	if err := os.WriteFile(fileName, newData, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Dosya yazılamadı",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Makale başarıyla silindi",
	})
}
func AddQuestion(c *gin.Context) {
	questionType := c.Param("type")
	key := c.Param("category")
	fmt.Println("key:", key)

	var datas Question
	if err := c.ShouldBindJSON(&datas); err != nil {
		fmt.Println("err:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	var fileName string
	if questionType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Lütfen kullanıcı tipi seçiniz",
		})
		return
	}
	if questionType == "person" {
		fileName = "./data/person.json"
	} else {
		fileName = "./data/company2.json"

	}
	file, err := os.Open(fileName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	byteFile, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	ques := make(map[string]interface{})
	var subques []Question
	if err := json.Unmarshal(byteFile, &ques); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	raw, ok := ques[key]
	if !ok {
		log.Println("diet alanı bulunamadı")
		return
	}
	jsonData, err := json.Marshal(raw)
	if err != nil {
		log.Println("JSON'a çevirirken hata:", err)
		return
	}
	err = json.Unmarshal(jsonData, &subques)
	if err != nil {
		log.Println("JSON çözümlenirken hata:", err)
		return
	}
	subques = append(subques, datas)
	ques[key] = subques
	newData, err := json.MarshalIndent(ques, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "JSON oluşturulamadı",
		})
		return
	}

	if err := os.WriteFile(fileName, newData, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Dosya yazılamadı",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Makale başarıyla silindi",
	})
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

func DeleteArticle(c *gin.Context) {
	var infos []controllers.Data
	title := c.PostForm("title")
	fmt.Println("title:", title)
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
	found := false
	for i, info := range infos {
		if info.Title == title {
			infos = append(infos[:i], infos[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Makale bulunamadı",
		})
		return
	}

	// Güncellenmiş veriyi tekrar JSON'a çevir ve dosyaya yaz
	newData, err := json.MarshalIndent(infos, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "JSON oluşturulamadı",
		})
		return
	}

	if err := os.WriteFile("./data/data.json", newData, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Dosya yazılamadı",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Makale başarıyla silindi",
	})
}

func AddArticle(c *gin.Context) {
	var datas controllers.Data
	if err := c.ShouldBindJSON(&datas); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	file, err := os.ReadFile("./data/data.json")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Dosya okunamadı",
		})
		return
	}
	var comments []controllers.Data
	if err := json.Unmarshal(file, &comments); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu2",
		})
		return
	}
	comments = append(comments, datas)
	updatedData, err := json.MarshalIndent(comments, "", "  ")
	if err != nil {
		log.Fatal("JSON dönüştürme hatası:", err)
	}
	if err := os.WriteFile("./data/data.json", updatedData, 0644); err != nil {
		log.Fatal("Yazma hatası:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "başarılı",
	})
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
	fmt.Println("commentReq.UserType:", commentReq.UserType)
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
func DeleteComment(c *gin.Context) {
	data := c.PostForm("question_type")
	question := c.PostForm("question")
	fmt.Println("data:", data)
	fmt.Println("question:", question)
	var file *os.File
	var fileName string
	var err error
	if data != "" {
		if data == "person" {
			fileName = "./data/person-question.json"
		} else if data == "company" {
			fileName = "./data/company-question.json"
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
	file, err = os.Open(fileName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
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
	found := false
	for i, info := range comments {
		if info.Question == question {
			comments = append(comments[:i], comments[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Soru bulunamadı",
		})
		return
	}

	// Güncellenmiş veriyi tekrar JSON'a çevir ve dosyaya yaz
	newData, err := json.MarshalIndent(comments, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "JSON oluşturulamadı",
		})
		return
	}

	if err := os.WriteFile(fileName, newData, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Dosya yazılamadı",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Soru başarıyla silindi",
	})
}
