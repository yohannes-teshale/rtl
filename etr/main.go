package etr

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func Main(testDir string, projDir string, outDir string) {

	err := os.MkdirAll("generated-RTL", 0755)
	if err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}

	lerr := godotenv.Load()
	if lerr != nil {
		log.Fatal("Error loading .env file")
	}

	client := openai.NewClient(os.Getenv("OPENAI_KEY"))

	err = filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".js" {
			enzymeCode, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			rtlCode, err := convertToRTL(client, string(enzymeCode))
			if err != nil {
				return err
			}

			outputPath := filepath.Join("generated-RTL", info.Name())
			err = ioutil.WriteFile(outputPath, []byte(rtlCode), 0644)
			if err != nil {
				return err
			}

			fmt.Printf("Converted %s to RTL format and saved to %s\n", path, outputPath)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func convertToRTL(client *openai.Client, enzymeCode string) (string, error) {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Convert the following Enzyme test code to React Testing Library format:\n\n" + enzymeCode + ". you should not generate anyother output except for the output code. No explanation or any form of text are allowed.",
				},
			},
			MaxTokens:   4096,
			Temperature: 1,
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
