package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"io"
	"log"
	"os"
	"strings"
)

func askLLMStream(model string) {
	endpoint := "https://openrouter.ai/api/v1/chat/completions"

	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Chrome_124),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		log.Println(err)
	}

	tools := []Tool{
		{
			Type: "function",
			Function: llmFunction{
				Name:        "searchWeb",
				Description: "Retrieve realtime and current information e.g. news and weather from a web page based on a query. The returned result is the content of a webpage and the source link",
				Parameters: llmFunctionParam{
					Type: "object",
					Properties: llmFunctionProperties{
						Query: Query{
							Type:        "string",
							Description: "The search query used for getting a link to a webpage.",
						},
					},
					Required: []string{"query"},
				},
			},
		},
	}

	stream := true

	data := T{
		Model:    model,
		Messages: msgHistory,
		Tools:    &tools,
		Stream:   &stream,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(string(jsonData)))
	if err != nil {
		log.Println(err)
	}

	req.Header = http.Header{
		"accept":        {"application/json"},
		"user-agent":    {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"},
		"Authorization": {`Bearer ` + os.Getenv("apiKey")},
		http.HeaderOrderKey: {
			"accept",
			"user-agent",
			"Authorization",
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		return
	}

	reader := bufio.NewReader(resp.Body)
	defer resp.Body.Close()

	var fullResponse strings.Builder
	var chunk Chunk
	isToolCall := false

	for {

		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("Error reading stream:", err)
			return
		}
		// Remove "data: " prefix if present
		line = bytes.TrimPrefix(line, []byte("data: "))

		// Skip empty lines
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}

		// Check for end of stream
		if string(bytes.TrimSpace(line)) == "[DONE]" {
			break
		}

		// Unmarshal the JSON
		err = json.Unmarshal(line, &chunk)
		if err != nil {
			//log.Println("Error unmarshalling JSON:", err)
			continue
		}

		// Process the chunk
		if len(chunk.Choices) > 0 {
			delta := chunk.Choices[0].Delta

			// Check for tool calls
			if delta.ToolCalls != nil && len(*delta.ToolCalls) > 0 {
				isToolCall = true
				break
			}
			// Stream the content
			if delta.Content != nil {
				content, ok := delta.Content.(string)
				if ok && content != "" {
					fmt.Print(content)
					fullResponse.WriteString(content)
				}
			}
		}
	}

	// Update message history
	if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
		newMsg := MSG{
			Role:    chunk.Choices[0].Delta.Role,
			Content: fullResponse.String(),
		}
		delta := chunk.Choices[0].Delta
		if delta.ToolCalls != nil && len(*delta.ToolCalls) > 0 {
			newMsg.ToolCalls = delta.ToolCalls
		}
		msgHistory = append(msgHistory, newMsg)
	}

	if isToolCall {
		// Handle tool call
		handleToolCall(&chunk, model)
	}
}

func handleToolCall(chunk *Chunk, model string) {
	if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.ToolCalls != nil {
		toolCalls := *chunk.Choices[0].Delta.ToolCalls
		if len(toolCalls) > 0 {
			for _, toolCall := range toolCalls {
				fmt.Printf("\nTool Call: %s\n", toolCall.Function.Name)
				fmt.Printf("Arguments: %s\n\n", toolCall.Function.Arguments)

				// Add your tool call handling logic here
				funcName := toolCall.Function.Name
				var funcArg struct {
					Query string `json:"query"`
				}
				err := json.Unmarshal([]byte(toolCall.Function.Arguments), &funcArg)
				if err != nil {
					log.Println("Error unmarshalling JSON:", err)
					continue
				}

				searchRes, err := searchWeb(funcArg.Query)
				if err != nil {
					fmt.Println(err)
				}

				searchWebResult := struct {
					Text   string `json:"text"`
					Source string `json:"source"`
				}{
					Text:   searchRes[0],
					Source: searchRes[1],
				}

				toolResult, err := json.Marshal(searchWebResult)
				if err != nil {
					fmt.Println("Error marshalling JSON:", err)
					return
				}

				toolResultMsg := MSG{
					Role:    "tool",
					Name:    &funcName,
					ToolID:  &chunk.ID,
					Content: string(toolResult),
				}

				msgHistory = append(msgHistory, toolResultMsg)
				askLLMStream(model)
			}
		} else {
			fmt.Println("Tool calls array empty")
		}
	} else {
		fmt.Println("No tool calls found")
	}
}
