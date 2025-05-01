package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type QuesInf struct {
	Key      string  `json:"key"`
	Question string  `json:"question"`
	Options  Options `json:"options"`
}

type Options struct {
	Text     string `json:"text"`
	Emission string `json:"emission"`
}

// @Description  Bireysel test kısmı
// @Tags         Test
// @Accept       json
// @Produce      json
// @Success      200 {object} QuesInf
// @Router       /person-questions [get]
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
	ques := make(map[string]interface{})
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
