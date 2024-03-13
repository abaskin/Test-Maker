package aiclient

import (
	"context"
	"fmt"
	"strings"

	"github.com/daichi-m/go18ds/lists/arraylist"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const (
	reWritePrompt = `
		Rewrite the following, keeping the original format.
		%s
	`
)

// ---- Gemini ----

type AiModel interface {
	Rewrite(string, chan<- *AiResponse)
	Close()
}

type GeminiAi struct {
	ctx    context.Context
	client *genai.Client
	model  *genai.GenerativeModel
	name   *string
}

type AiResponse struct {
	Err       error
	ModelName *string `json:"-"`
	Result    string
}

const (
	aiQueryError = "unable to query AI: %v"
)

func NewGeminiAi(name, apiKey string) (AiModel, error) {
	ai := &GeminiAi{
		ctx:  context.Background(),
		name: &name,
	}

	var err error
	ai.client, err = genai.NewClient(ai.ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to AI: %v", err)
	}

	ai.model = ai.client.GenerativeModel("gemini-pro")
	ai.model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockNone,
		},
	}

	return ai, nil
}

func (g *GeminiAi) Rewrite(original string, responseChan chan<- *AiResponse) {
	resp, err := g.model.GenerateContent(g.ctx, genai.Text(
		fmt.Sprintf(reWritePrompt, string(original))))
	if err != nil {
		responseChan <- &AiResponse{
			Err:    fmt.Errorf(aiQueryError, err),
			Result: "",
		}
		return
	}

	responseChan <- &AiResponse{
		Err:    nil,
		Result: geminiResponseString(resp),
	}
}

func geminiResponseString(resp *genai.GenerateContentResponse) string {
	respParts := arraylist.New[string]()
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				respParts.Add(string(part.(genai.Text)))
			}
		}
	}

	return strings.NewReplacer("```json", "", "```JSON", "", "```", "").
		Replace(strings.Join(respParts.Values(), "\n"))
}

func (g *GeminiAi) Close() {
	g.client.Close()
}
