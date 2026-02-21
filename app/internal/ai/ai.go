package ai

import (
	"bufio"
	"bytes"
	"commmunity/app/zlog"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

const (
	API_KEY = "43f29e6afa964c199fbde7a6d7b57794.768b8xPe3IMpGUB7"     //viper.GetString("api.key")
	API_URL = "https://open.bigmodel.cn/api/paas/v4/chat/completions" //viper.GetString("api.url")
	MODEL   = "glm-4.7-flash"                                         //viper.GetString("api.model")
)

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream"`
	Temperature float64   `json:"temperature"`
	Thinking    struct {
		Type string `json:"type,omitempty"`
	} `json:"thinking,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func AutoSummary(content string) (string, error) {
	requestData := ChatRequest{
		Model: MODEL,
		Messages: []Message{
			{
				Role:    "system",
				Content: "你是一个文章总结员，你的任务是将文章内容进行总结概括，遇到图片地址选择忽略，回复语言要求俏皮可爱，且带有些许傲娇，但概括内容不能有偏差，可以带颜文字和小表情",
			},
			{
				Role:    "user",
				Content: content,
			},
		},
		Stream:      true,
		Temperature: 1.2,
		Thinking: struct {
			Type string `json:"type,omitempty"`
		}{Type: "enabled"},
	}
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		zlog.Error("JSON 编码失败", zap.Error(err))
		return "", err
	}
	req, err := http.NewRequest("POST", API_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		zlog.Error("创建请求失败", zap.Error(err))
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+API_KEY)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 2 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		zlog.Error("网络请求失败", zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)
	var finalContent string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}
		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}
		if len(chunk.Choices) > 0 {
			finalContent += chunk.Choices[0].Delta.Content
		}
	}
	return finalContent, nil
}
