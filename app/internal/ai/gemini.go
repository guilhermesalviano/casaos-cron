package ai

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func AnalyzeWithGemini(prompt string) (string, error) {
	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY not set in environment")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", fmt.Errorf("failed to create Gemini client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-3.1-flash-lite-preview")
	
	model.SetTemperature(0.7)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %v", err)
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates returned from Gemini")
	}

	var result string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				result += fmt.Sprintf("%v", part)
			}
		}
	}

	return result, nil
}
