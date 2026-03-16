package ai

import (
	"context"
	"encoding/json"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type Client struct {
	client *openai.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		client: openai.NewClient(apiKey),
	}
}

func (c *Client) GenerateWeeklyPlan(ctx context.Context, req WeeklyPlanRequest) (*WeeklyPlan, error) {
	prompt := BuildWeeklyPrompt(req)

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: "gpt-4o-mini",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a helpful meal planning assistant that generates creative, diverse, and practical meal plans. Always respond with valid JSON only, no additional text. CRITICAL: Every shopping_items entry MUST have item_name, quantity (number only), and unit (separate field) - never omit the unit field.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.8,
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	content := resp.Choices[0].Message.Content

	// Debug: log the raw AI response
	fmt.Printf("DEBUG: AI Response (first 500 chars): %s\n", content[:min(500, len(content))])

	var plan WeeklyPlan
	if err := json.Unmarshal([]byte(content), &plan); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w. Response: %s", err, content)
	}

	// Debug: log first shopping item to verify unit is present
	if len(plan.Monday.ShoppingItems) > 0 {
		fmt.Printf("DEBUG: First Monday item: %+v\n", plan.Monday.ShoppingItems[0])
	}

	return &plan, nil
}

func (c *Client) GenerateDayOptions(ctx context.Context, req DayOptionsRequest) (*DayOptions, error) {
	prompt := BuildDayOptionsPrompt(req)

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: "gpt-4o-mini",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a helpful meal planning assistant that generates creative, diverse meal options. Always respond with valid JSON only, no additional text. CRITICAL: Every shopping_items entry MUST have item_name, quantity (number only), and unit (separate field) - never omit the unit field.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.8,
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	content := resp.Choices[0].Message.Content

	var options DayOptions
	if err := json.Unmarshal([]byte(content), &options); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w. Response: %s", err, content)
	}

	if len(options.Options) != 3 {
		return nil, fmt.Errorf("expected 3 options but got %d", len(options.Options))
	}

	return &options, nil
}
