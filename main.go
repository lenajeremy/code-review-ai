package main

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
	"strings"
)

func init() {
	if godotenv.Load() != nil {
		log.Fatal("Failed to setup .env")
	}
}

var ctx = context.Background()

func main() {
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))

	defer func(client *genai.Client) {
		err := client.Close()
		if err != nil {
			log.Fatal("Failed to close client")
		}
	}(client)

	if err != nil {
		log.Fatal(err)
	}

	model := client.GenerativeModel("gemini-1.5-flash")

	settings := genai.SafetySetting{Category: genai.HarmCategorySexuallyExplicit, Threshold: genai.HarmBlockNone}
	model.SafetySettings = []*genai.SafetySetting{&settings}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
		return
	})

	r.POST("/ai", func(c *gin.Context) {
		message := c.PostForm("chat-prompt")
		resp, err := model.GenerateContent(ctx, genai.Text(message))

		if err != nil {
			log.Println("Error with model")
			log.Fatal(err)
			return
		}

		if story, err := json.Marshal(resp.Candidates[0].Content.Parts[0]); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			//log.Println(strings.Split(string(story), "\n"))
			c.HTML(http.StatusOK, "ai-response.html", gin.H{
				"story": strings.Join(strings.Split(string(story), "\\n"), "<br />"),
			})
		}

		return
	})

	log.Fatal(r.Run(":3000"))
}
