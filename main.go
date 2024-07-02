package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var msgHistory []MSG

func main() {

	err := os.Setenv("apiKey", "YOUR OPENAI KEY HERE")
	if err != nil {
		return
	}

	sysMsg := MSG{
		Role: "system",
		Content: `You are a helpful agent with the ability to access to the internet via the tool provided to you. 
			The tool provided is called searchWeb and you should use it whenever you need information that is real-time
			or current that is not part of your training data. If you do search the web, make sure to always include 
			the source link in your response. The final output should be in formatted markdown.`,
	}

	msgHistory = append(msgHistory, sysMsg)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\nAsk a question: ")
		if scanner.Scan() {
			input := scanner.Text()
			input = strings.TrimSpace(input)
			if input == "end" {
				fmt.Println("Goodbye!")
				os.Exit(0)
			}
			if input != "" {
				newMsg := MSG{
					Role:    "user",
					Content: input,
				}
				msgHistory = append(msgHistory, newMsg)
				askLLMStream("google/gemini-flash-1.5")
			}
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}

	}
}
