package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type QuesInf struct {
	Key      string
	Question string
	Options  Options
}

type Options struct {
	Text     string
	Emission string
}

// @Summary      Test
// @Description  Bireysel test kısmı
// @Tags         Test
// @Accept       json
// @Produce      json
// @Success      200 {object} QuesInf "User test successfully"
// @Failure      400 {object} controllers.Response "Invalid request"
// @Router       /person_questions [get]
func GetPersonQues(c *gin.Context) {

	file, err := os.Open("./data/person.json")
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

	//var questions QuesInf
	ques := make(map[string]interface{})
	fmt.Println("deneme:", string(byteFile))
	if err := json.Unmarshal(byteFile, &ques); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"eror":    err,
			"message": "Something went worong.",
		})
		return
	}

	c.JSON(http.StatusOK, ques)

}

func Test(c *gin.Context) {
	c.JSON(http.StatusOK, "başarılı")
}
func Test2(c *gin.Context) {
	c.JSON(http.StatusOK, "düzgün çalışıyor.")
}
