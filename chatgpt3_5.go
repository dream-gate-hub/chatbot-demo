package main

import (
	"bytes"
	"encoding/json"
	"errors"
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

type ChatGPTMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type MessagesRequest struct {
	Messages []ChatGPTMessage `json:"messages"`
}

func ChatWithGPT3_5(chatBotSystemPrompt string, msgReq MessagesRequest) (string, error) {
	// API URL
	url := "https://api.openai.com/v1/chat/completions"

	// 获取环境变量中的API密钥
	//apiKey := os.Getenv("OPENAI_API_KEY")
	apiKey := "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"

	// 构建请求体
	var chatGPTReqMsg = []map[string]string{
		{
			"role":    "system",
			"content": chatBotSystemPrompt,
		},
	}
	for _, msg := range msgReq.Messages {
		chatGPTReqMsg = append(chatGPTReqMsg, map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		})
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"model":    "gpt-3.5-turbo",
		"messages": chatGPTReqMsg,
	})
	if err != nil {
		return "", err
	}

	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	// 设置请求头部
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 用于存储解析后的数据
	var chatResponse ChatResponse

	// 解析JSON
	err = json.Unmarshal(body, &chatResponse)
	if err != nil {
		log.Fatalf("Error occurred during unmarshalling. Error: %s", err.Error())
		return "", err
	}

	if len(chatResponse.Choices) > 0 {
		return chatResponse.Choices[0].Message.Content, nil
	} else {
		return "", errors.New("no choices available")
	}

}
