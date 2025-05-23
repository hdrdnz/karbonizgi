package controllers

import (
	"carbonfootprint/config"
	"carbonfootprint/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
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
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Bir hata oluştu"})
		return
	}

	if inf.UserId == 0 || inf.UserId != int(Claims["userId"].(float64)) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Geçersiz veya boş userId"})
		return
	}

	user := &model.User{}
	if err := db.Where("id=?", inf.UserId).First(&user).Error; err != nil || user.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Kullanıcı bulunamadı."})
		return
	}

	var count int64
	if err := db.Table("user_score").Where("user_id=?", user.Id).Count(&count).Error; err != nil || count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Kullanıcıya ait skor bulunamadı."})
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
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Geçersiz key"})
		return
	}

	pastMsgs, err := rdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		panic(err)
	}

	if len(pastMsgs) != 0 {
		var messages []openai.ChatCompletionMessage
		for _, msgJSON := range pastMsgs {
			var cm openai.ChatCompletionMessage
			if err := json.Unmarshal([]byte(msgJSON), &cm); err == nil {
				messages = append(messages, cm)
			}
		}
		userMessage := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: inf.Message,
		}
		msgByte, _ := json.Marshal(userMessage)
		messages = append(messages, userMessage)
		rdb.RPush(ctx, key, msgByte)

		response, _, err := chatWithCarbonExpert(client, messages, user, inf.Filter)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Bir hata oluştu."})
			return
		}
		rdb.RPush(ctx, key, jsonMust(openai.ChatCompletionMessage{Role: openai.ChatMessageRoleAssistant, Content: response}))
		c.JSON(http.StatusOK, gin.H{"message": response})
		return
	}

	if inf.Score == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Skor boş bırakılamaz."})
		return
	}

	var firstMes string
	if inf.Filter == "detail" {
		firstMes = "[DETAY SKOR MODU] kullanıcı ismi: " + user.Firstname + ", key:" + inf.Key + ", skor:" + inf.Score
	} else {
		firstMes = "Kullanıcı ismi: " + user.Firstname + ", Karbon ayak izi:" + inf.Score + " ton CO₂/yıl"
	}
	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleUser, Content: firstMes},
	}

	response, _, err := chatWithCarbonExpert(client, messages, user, inf.Filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Bir hata oluştu."})
		return
	}
	rdb.RPush(ctx, key, jsonMust(messages[0]))
	rdb.RPush(ctx, key, jsonMust(openai.ChatCompletionMessage{Role: openai.ChatMessageRoleAssistant, Content: response}))
	rdb.Expire(ctx, key, 5*time.Minute)
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": response})
}

func chatWithCarbonExpert(client *openai.Client, userMessages []openai.ChatCompletionMessage, user *model.User, filter string) (string, bool, error) {
	var systemMsg string
	if user.UserType == "person" {
		systemMsg = "Kullanıcı türü: personal"
	} else {
		systemMsg = "Kullanıcı türü: company"

	}
	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: systemMsg},
	}
	messages = append(messages, userMessages...)
	fmt.Println("messages:", messages)

	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:    config.GetModulName(),
		Messages: messages,
	})
	if err != nil {
		return "", true, err
	}
	return resp.Choices[0].Message.Content, true, nil
}

type GeneralReq struct {
	Message string `json:"message"`
}

// @Summary      User genel chat Kısmı
// @Description  Chat Kısmı
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        GeneralReq body GeneralReq true "kullanıcı mesajını girmelisin."
// @Success      200 {object} Response "message içerisinde yanıt döner"
// @Failure      400 {object} Response "Invalid request"
// @Router       /general-chat [post]
func GeneralChat(c *gin.Context) {
	db := model.GetDB()
	client := config.GetClient()
	chat := Chat{}
	if err := c.ShouldBindJSON(&chat); err != nil || chat.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Geçersiz mesaj."})
		return
	}

	user := model.User{}
	if err := db.Where("id=?", Claims["userId"]).First(&user).Error; err != nil || user.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Kullanıcı bulunamadı."})
		return
	}

	rdb, ctx := config.GetRedis()
	key := strconv.FormatInt(int64(user.Id), 10) + ":chat"
	pastMsgs, err := rdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		panic(err)
	}

	systemPrompt := os.Getenv("chat")
	chatMessages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
	}

	if len(pastMsgs) > 0 {
		for _, msgJSON := range pastMsgs {
			if msgJSON != "" {
				var cm openai.ChatCompletionMessage
				if err := json.Unmarshal([]byte(msgJSON), &cm); err == nil {
					chatMessages = append(chatMessages, cm)
				}
			}
		}
	}

	// Yeni gelen kullanıcı mesajı eklenir
	userMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: chat.Message,
	}
	chatMessages = append(chatMessages, userMessage)

	fmt.Println("genel mesaj:", chatMessages)
	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    "ft:gpt-3.5-turbo-0125:personal::BY8WUAsP",
			Messages: chatMessages,
		},
	)
	if err != nil {
		fmt.Println("err:", err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Bir hata oluştu"})
		return
	}

	respChat := response.Choices[0].Message.Content

	rdb.RPush(ctx, key, jsonMust(userMessage))
	rdb.RPush(ctx, key, jsonMust(openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: respChat,
	}))
	rdb.Expire(ctx, key, 5*time.Minute)

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": respChat})
}

func jsonMust(msg openai.ChatCompletionMessage) []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}
