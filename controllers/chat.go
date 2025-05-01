package controllers

import (
	"carbonfootprint/config"
	"carbonfootprint/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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
	Filter  string
	Key     string
}
type Chat struct {
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
// @Param        PersonInf body PersonInf true "Kullanıcının user id bilgisini ve ilk başta score bilgisini girmelisin.Score bilgisi /score endpointinde data kısmında dönüyor.Eğer kullanıcı konuşmayı devam ettirirse message kısmında kullanıcının mesajını gönderebilirsin. Filter kısmında main ya da detail ifadelerini göndermelisin. main için key kısmını girmene gerek yok ama detail kısmında key kısmına alt başlığın key bilgisini, score kısmına da alt başlıkta aldığı skoru göndermelisin."
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
	fmt.Println("userıd:", Claims["userId"])
	fmt.Println("girilken:", inf.UserId)

	if inf.UserId != int(Claims["userId"].(float64)) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "geçersiz userId",
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
	var key string
	if inf.Filter == "detail" {
		key = userId + ":" + inf.Key
	} else if inf.Filter == "main" {
		key = userId + ":messages"
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Geçersiz key",
		})
		return
	}
	filter := inf.Filter

	//önceki mesajların kontrolü
	pastMsgs, err := rdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("len(pastMsgs):", len(pastMsgs), pastMsgs)
	if len(pastMsgs) != 0 {
		var messages []openai.ChatCompletionMessage
		for _, msgJSON := range pastMsgs {
			fmt.Println("----JSON", msgJSON)
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

		//mesajı kaydetme
		rdb.RPush(ctx, key, msgByte)

		response, control, err := chatWithCarbonExpert(client, messages, user, filter)
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
		fmt.Println("ilk kısım")
		if inf.Score == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Skor boş bırakılamaz.",
			})
			return
		}
		var firstMes string
		if filter == "detail" {
			// [DETAY SKOR MODU] burak, key: house_type, skor: 330"
			firstMes = "[DETAY SKOR MODU] " + user.Firstname + ", key:" + inf.Key + ",skor:" + inf.Score
			// firstMes = inf.Key + ":" + inf.Score
		} else {
			firstMes = "Kullanıcı ismi:" + user.Firstname + "Karbon ayak izi değeri:" + inf.Score

		}
		messages := []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem,
				Content: firstMes,
			},
		}
		response, control, err := chatWithCarbonExpert(client, messages, user, filter)
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

func chatWithCarbonExpert(client *openai.Client, userMessages []openai.ChatCompletionMessage, user *model.User, filter string) (string, bool, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	moduleName := config.GetModulName()
	var content string
	if filter == "detail" {
		if user.UserType == "person" {
			content = "Kullanıcı bireysel bir kişi. Kullanıcı belirli bir konuda karbon etkisini değerlendirmeni istiyor."
		} else {
			content = "Kullanıcı bir şirket temsilcisi. Kullanıcı belirli bir konuda karbon etkisini değerlendirmeni istiyor."
		}
	} else {
		if user.UserType == "person" {
			content = "Kullanıcı bireysel bir kişi. Kullanıcının karbon ayak izi skoru ile ilgili genel bir yorum yap."
		} else {
			content = "Kullanıcı bir şirket temsilcisi. Kullanıcının karbon ayak izi skoru ile ilgili genel bir yorum yap."
		}
	}
	prompt := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: os.Getenv("chatControl")},
	}
	messages := append(prompt, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: content,
	})
	// Kullanıcı mesajları (önceki mesajlar dahil)
	messages = append(messages, userMessages...)
	fmt.Println("messages:", messages)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    moduleName,
			Messages: messages,
		},
	)
	if err != nil {
		return "", true, err
	}
	// if strings.Contains(resp.Choices[0].Message.Content, os.Getenv("ControlPrompt")) {
	// 	cont := controlMessage(strconv.FormatInt(int64(user.Id), 10))
	// 	if !cont {
	// 		return os.Getenv("EndPrompt"), false, nil
	// 	}
	// }
	return resp.Choices[0].Message.Content, true, nil
}

// func controlMessage(userId string) bool {
// 	rdb, ctx := config.GetRedis()
// 	resp, _ := rdb.Get(ctx, userId+"-control").Result()
// 	if resp != "" {
// 		if resp == "2" {
// 			return false
// 		} else {
// 			rdb.Set(ctx, userId+"-control", 2, 5*time.Minute)
// 		}

// 	} else {
// 		rdb.Set(ctx, userId+"-control", 1, 5*time.Minute)
// 	}
// 	return true

// }

func GeneralChat(c *gin.Context) {
	fmt.Println("deneme:", os.Getenv("chatControl"))
	db := model.GetDB()
	client := config.GetClient()
	chat := Chat{}
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	if chat.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Geçersiz mesaj.",
		})
		return
	}
	user := model.User{}
	if err := db.Where("id=?", Claims["userId"]).First(&user).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
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

	rdb, ctx := config.GetRedis()
	key := strconv.FormatInt(int64(user.Id), 10) + ":chat"
	pastMsgs, err := rdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		panic(err)
	}

	var content string
	if user.UserType == "company" {
		content = "Kullanıcı bir şirket. Kullanıcı karbon ayak izi ile ilgili sorular soruyor."
	} else {
		content = "Kullanıcı bireysel bir kişi. Kullanıcı karbon ayak izi ile ilgili sorular soruyor."

	}
	// content += os.Getenv("chatControl")
	// chatMessages := []openai.ChatCompletionMessage{
	// 	{Role: openai.ChatMessageRoleSystem, Content: os.Getenv("chatControl")},
	// }
	chatMessages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: content},
	}
	// chatMessages := append(chatMessages, openai.ChatCompletionMessage{
	// 	Role:    openai.ChatMessageRoleSystem,
	// 	Content: content,
	// })
	// chatMessages := []openai.ChatCompletionMessage{
	// 	{Role: openai.ChatMessageRoleSystem, Content: content},
	// }
	moduleName := config.GetModulName()
	userMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: "[GENEL CHAT MODU] " + chat.Message,
		//Content: chat.Message,
	}
	if len(pastMsgs) != 0 {
		var messages []openai.ChatCompletionMessage
		for _, msgJSON := range pastMsgs {
			fmt.Println("----JSON", msgJSON)
			if msgJSON != "" {
				var cm openai.ChatCompletionMessage
				if err := json.Unmarshal([]byte(msgJSON), &cm); err != nil {
					panic(err)
				}
				messages = append(messages, cm)
			}
		}

		// Kullanıcı mesajları (önceki mesajlar dahil)
		chatMessages = append(chatMessages, messages...)
		fmt.Println("messages:", chatMessages)
	}
	chatMessages = append(chatMessages, userMessage)
	fmt.Println("chatMessages:", chatMessages)
	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    moduleName,
			Messages: chatMessages,
		},
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	respChat := response.Choices[0].Message.Content
	allMessages := []string{chat.Message, respChat}
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
		"message": respChat,
	})
}
