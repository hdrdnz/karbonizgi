package controllers

import (
	"carbonfootprint/config"
	"carbonfootprint/model"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

type PersonInf struct {
	UserId  int    `json:"userId"`
	Score   string `json:"score"`
	Message string `json:"message"`
}

type ChatResp struct {
	Message string `json:"message"`
}

// @Summary      User chat Kısmı
// @Description  Chat Kısmı
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        PersonInf body PersonInf true "Kullanıcının user id bilgisini ve ilk başta score bilgisini girmelisin.Score bilgisi /score endpointinde data kısmında dönüyor.Eğer kullanıcı konuşmayı devam ettirirse message kısmında kullanıcının mesajını gönderebilirsin."
// @Success      200 {object} Response "Message kısmını kullanıcıya gösterebilirsin.Eğer kullanıcı karbonayak izi dışında sorular sorarsa chat cevap vermez karbon ayak izine yönlendirir.2 kere gereksiz soru sorulursa data kısmında 'false' döner ve burada chat konuşmasını bitir sonra message kısmını kullanıcıya gösterebilirsin."
// @Failure      400 {object} Response "Invalid request"
// @Router       /chat [post]
func PostChat(c *gin.Context) {
	db := model.GetDB()
	client := config.GetClient()
	var inf PersonInf
	if err := c.ShouldBindJSON(&inf); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	if inf.UserId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "userId kısmı boş olamaz.",
		})
		return
	}
	user := &model.User{}
	if err := db.Where("id=?", inf.UserId).First(&user).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	if user.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Kullanıcı bulunamadı.",
		})
		return
	}
	var count int64
	if err := db.Table("user_score").Where("user_id=?", user.Id).Count(&count).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Kullanıcıya ait skor bulunamadı.",
		})
		return
	}
	rdb, ctx := config.GetRedis()
	userId := strconv.FormatInt(int64(user.Id), 10)
	key := userId + ":messages"
	pastMsgs, err := rdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		panic(err)
	}
	if len(pastMsgs) != 0 {
		var messages []openai.ChatCompletionMessage
		for _, msgJSON := range pastMsgs {
			if msgJSON != "" {
				var cm openai.ChatCompletionMessage
				if err := json.Unmarshal([]byte(msgJSON), &cm); err != nil {
					panic(err)
				}
				messages = append(messages, cm)
			}
		}
		message := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: inf.Message,
		}
		msgByte, err := json.Marshal(message)
		if err != nil {
			panic(err)
		}
		messages = append(messages, message)
		rdb.RPush(ctx, key, msgByte)

		response, control, err := chatWithCarbonExpert(client, messages, userId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Bir hata oluştu.",
			})
			return
		}
		if !control {
			c.JSON(http.StatusOK, gin.H{
				"message": response,
				"data":    false,
			})
			return
		}
		resp := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: response,
		}
		respByte, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}
		rdb.RPush(ctx, key, respByte)
		c.JSON(http.StatusOK, gin.H{
			"message": response,
		})
		return
	} else {
		if inf.Score == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Skor boş bırakılamaz.",
			})
			return
		}
		firstMes := "Kullanıcı ismi:" + user.Firstname + "Karbon ayak izi değeri:" + inf.Score
		messages := []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem,
				Content: firstMes,
			},
		}
		response, control, err := chatWithCarbonExpert(client, messages, userId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Bir hata oluştu.",
			})
			return
		}
		if !control {
			c.JSON(http.StatusOK, gin.H{
				"status":  "success",
				"message": response,
				"data":    "false",
			})
			return
		}
		allMessages := []string{firstMes, response}
		for _, msg := range allMessages {
			message := openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: msg,
			}
			msgByte, err := json.Marshal(message)
			if err != nil {
				panic(err)
			}
			rdb.RPush(ctx, key, msgByte)
		}
		rdb.Expire(ctx, key, 5*time.Minute)
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": response,
		})
	}

}

func chatWithCarbonExpert(client *openai.Client, userMessages []openai.ChatCompletionMessage, userId string) (string, bool, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: os.Getenv("SystemPrompt")},
	}

	// Kullanıcı mesajları (önceki mesajlar dahil)
	messages = append(messages, userMessages...)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)
	if err != nil {
		return "", true, err
	}
	if strings.Contains(resp.Choices[0].Message.Content, os.Getenv("ControlPrompt")) {
		cont := controlMessage(userId)
		if !cont {
			return os.Getenv("EndPrompt"), false, nil
		}
	}
	return resp.Choices[0].Message.Content, true, nil
}

func controlMessage(userId string) bool {
	rdb, ctx := config.GetRedis()
	resp, _ := rdb.Get(ctx, userId+"-control").Result()
	if resp != "" {
		if resp == "2" {
			return false
		} else {
			rdb.Set(ctx, userId+"-control", 2, 5*time.Minute)
		}

	} else {
		rdb.Set(ctx, userId+"-control", 1, 5*time.Minute)
	}
	return true

}
