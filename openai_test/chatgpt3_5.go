package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func main() {
	// API URL
	url := "https://api.openai.com/v1/chat/completions"

	// 获取环境变量中的API密钥
	//apiKey := os.Getenv("OPENAI_API_KEY")
	apiKey := "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"
	if apiKey == "" {
		fmt.Println("API key is not set. Please set the OPENAI_API_KEY environment variable.")
		return
	}

	chatBotSystemPrompt := `
		[Your Name]
		Miguel O'Hara

		[Your Personality]
		Miguel is a male short tempered, sarcastic asshole that thinks he is always right. He is straight forward and direct. He likes playing around but would never admit it. He knows he is good looking and his voice is a turn on. Features(Fangs+Talons) Powers(can stick to the wall+spider web) Personalty(sarcastic+short tempered) Since he build the spider Society he sits in his office all day. He knows he can't handle a another canon disturbing. He needs to fix the portals and watches.
	
		[Scenario]
		Miguel was looking at the numerous screens in his "office", if you could call it that. He was immersed in his work knowing He needs to fix the portals and watches until he smelled a new scent. Something he had never felt before. Someone new in the headquarters? Without his knowledge? He turned around and, obviously there they were. He unsheathed his talons immediately, visibly tense. "Who are you? What are you doing here?
	
	`

	// 构建请求体
	requestBody, err := json.Marshal(map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": chatBotSystemPrompt,
			},
			{
				"role":    "user",
				"content": "i am your dad",
			},
		},
	})
	if err != nil {
		panic(err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err)
	}

	// 设置请求头部
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 用于存储解析后的数据
	var chatResponse ChatResponse

	// 解析JSON
	err = json.Unmarshal(body, &chatResponse)
	if err != nil {
		log.Fatalf("Error occurred during unmarshalling. Error: %s", err.Error())
	}

	if len(chatResponse.Choices) > 0 {
		fmt.Println(chatResponse.Choices[0].Message.Content)
	} else {
		fmt.Println("No choices available.")
	}

	//fmt.Println(string(body))
}
