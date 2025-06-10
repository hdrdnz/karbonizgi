package controllers

import (
	"carbonfootprint/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func PostComp(c *gin.Context) {
	db := model.GetDB()
	categories := []string{
		"industry",
		"energy_transportation",
		"service_trade",
		"public_waste_management",
	}
	for _, cat := range categories {
		quest := &model.QuestionTypes{}
		quest.QuestionKey = cat
		quest.QuestionType = "company"
		if err := db.Save(quest).Error; err != nil {
			fmt.Println("hata:", cat)
		}
	}
}

func PostSubComp(c *gin.Context) {
	db := model.GetDB()
	categories := []string{"waste_recycling_rate", "waste_separation_facility", "recycling_initiatives", "renewable_energy_projects"}
	for _, cat := range categories {
		quest := &model.QuestionSubhead{}
		quest.QuestionKey = cat
		quest.QuestionTypesId = 28
		if err := db.Save(quest).Error; err != nil {
			fmt.Println("hata:", cat)
		}
	}
}

// @Description  Şirket test kısmı
// @Tags         Test
// @Accept       json
// @Produce      json
// @Success      200 {object} QuesInf
// @Router       /company-questions [get]
func GetCompanyQues(c *gin.Context) {
	file, err := os.Open("./data/company2.json")
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
