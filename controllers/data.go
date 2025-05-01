package controllers

import (
	"encoding/json"
	"io"
	"math/rand/v2"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type DataResp struct {
	Data []Data `json:"data"`
}

type Data struct {
	Title    string    `json:"title"`
	Image    string    `json:"image"`
	Sections []Section `json:"sections"`
}

type Section struct {
	Subtitle string `json:"subtitle"`
	Type     string `json:"type"` //-paragraph ve list olarak ayrılır. list ise item kısmı paragraph ise content kısmı doldurulur.
	Content  string `json:"content"`
	Items    []Item `json:"items"`
}
type Item struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// @Description  Genel bilgilendirme kısmı
// @Tags         Data
// @Accept       json
// @Produce      json
// @Success      200 {object} DataResp "Anasayfa için güncel bilgiler yer alır. Rastgele üç veri bulunur.Sections kısmında yazılar yer alır. type 'paragraph' ise content içerisinde yazı bulunur. Eğer type 'list' ise items içerisinde title ve content olarak yazılar yer alır."
// @Router       /data [get]
func GetInfo(c *gin.Context) {
	var infos []Data
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

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   infos[:3],
	})

}

type Comment struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// @Description Kullanıcı soru-cevap kısmı
// @Tags         Data
// @Accept       json
// @Param         user_type query string true "person ya da company ifadesi girilir."
// @Produce      json
// @Success      200 {array} Comment  "Kullanıcı eleştiri soru cevap kısımları bulunur."
// @Router       /comments [get]
func GetComments(c *gin.Context) {
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
	var comments []Comment
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

type SuggestResp struct {
	Status string   `json:"status`
	Data   []string `json:"data"`
}

// @Description  Örnek aksiyonlar
// @Tags         Data
// @Accept       json
// @Param         user_type query string true "person ya da company ifadesi girilir."
// @Produce      json
// @Success      200 {object} SuggestResp  "Kullanıcıya aksiyon önerileri için kullanılır.s"
// @Router       /suggested [get]
func GetSuggested(c *gin.Context) {
	data := c.Query("user_type")
	var requested []string
	if data != "" {
		if data == "person" {
			values := []string{
				"Haftada 1 gün toplu taşıma kullan",
				"Et tüketimini haftada 1 gün azalt",
				"Kullanmadığın prizleri fişten çek",
				"Tasarruflu ampul kullan",
				"İkinci el ürün tercih et",
				"Kıyafet alışverişini azalt",
				"Kısa mesafelerde yürümeyi tercih et",
				"Gereksiz e-posta aboneliklerinden çık",
				"E-fatura ve online bankacılık kullan",
				"Gıda israfını önlemek için alışveriş listesi yap",
			}
			requested = append(requested, values...)
		} else if data == "company" {
			values := []string{
				"LED aydınlatmaya geçiş yap",
				"Kâğıtsız ofis uygulamasına geç",
				"Servis araçlarının güzergâhlarını optimize et",
				"Güneş paneli yatırımı planla",
				"Çalışanlara sürdürülebilirlik eğitimi ver",
				"Ofiste atık ayrıştırma sistemleri kur",
				"Enerji tüketimini aylık olarak takip et",
				"Dijital arşiv sistemine geç",
				"Personel ulaşımında elektrikli araçları teşvik et",
				"Ofis içi sıcaklık ve iklimlendirme sistemlerini optimize et",
			}
			requested = append(requested, values...)

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

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   requested,
	})
}

type CalData struct {
	Title    string       `json:"title"`
	Sections []CalSection `json:"sections"`
}

type CalSection struct {
	Subtitle string    `json:"subtitle"`
	Type     string    `json:"type"`
	Content  string    `json:"content"`
	Items    []CalItem `json:"items"`
}
type CalItem struct {
	Content string `json:"content"`
}
type CallDataResp struct {
	Data CalData `json:"data"`
}

// @Description  Hesaplama yöntemi bilgilendirme kısmı
// @Tags         Data
// @Accept       json
// @Produce      json
// @Success      200 {object} CallDataResp "Hesaplama için bilgiler yer alır.Sections kısmında yazı yer alır. type 'paragraph' ise content içerisinde yazı bulunur. Eğer type 'list' ise items içerisinde content olarak yazılar yer alır."
// @Router       /cal-info [get]
func GetCalInfo(c *gin.Context) {
	var infos CalData
	file, err := os.ReadFile("./data/generalinfo.json")
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

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   infos,
	})

}
