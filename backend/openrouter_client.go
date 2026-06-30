package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

var chatModels = []string{
	"openai/gpt-oss-120b:free",
	"qwen/qwen3-next-80b-a3b-instruct:free",
	"openrouter/free",
}

type OpenRouterClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func newOpenRouterClient(apiKey string) *OpenRouterClient {
	return &OpenRouterClient{
		baseURL:    "https://openrouter.ai/api/v1",
		apiKey:     apiKey,
		httpClient: http.DefaultClient,
	}
}

// StreamChatReply streams an unconstrained conversational reply, calling
// onChunk for each text delta, and returns the full accumulated reply.
func (c *OpenRouterClient) StreamChatReply(ctx context.Context, messages []ChatMessage, onChunk func(string)) (string, error) {
	reqBody := map[string]any{
		"model":    chatModels[0],
		"models":   chatModels,
		"messages": messages,
		"stream":   true,
	}

	resp, err := c.post(ctx, reqBody)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return streamOpenRouterChunks(resp.Body, onChunk)
}

// ExtractFormData makes a single structured-output call over the given
// conversation, returning the model's current best-guess NDA field values.
func (c *OpenRouterClient) ExtractFormData(ctx context.Context, messages []ChatMessage) (FormData, error) {
	reqBody := map[string]any{
		"model":       chatModels[0],
		"models":      chatModels,
		"messages":    messages,
		"stream":      false,
		"temperature": 0,
		// Some free models default to a small max output length, which can
		// silently truncate the JSON schema response. The schema is small
		// (a handful of short NDA fields), so this is a generous ceiling.
		"max_tokens": 1000,
		"response_format": map[string]any{
			"type":        "json_schema",
			"json_schema": formDataJSONSchema(),
		},
	}

	resp, err := c.post(ctx, reqBody)
	if err != nil {
		return FormData{}, err
	}
	defer resp.Body.Close()

	var completion struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&completion); err != nil {
		return FormData{}, fmt.Errorf("decoding completion: %w", err)
	}
	if len(completion.Choices) == 0 {
		return FormData{}, fmt.Errorf("no choices in completion response")
	}

	var data FormData
	content := stripMarkdownCodeFence(completion.Choices[0].Message.Content)
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		return FormData{}, fmt.Errorf("parsing extracted form data: %w", err)
	}
	return data.trimmed(), nil
}

// trimmed returns a copy of the form data with leading/trailing whitespace
// removed from every text field. Some models return a single space rather
// than a true empty string for fields the user hasn't mentioned yet, which
// would otherwise defeat the frontend's `value || 'placeholder'` checks.
func (f FormData) trimmed() FormData {
	f.Party1 = f.Party1.trimmed()
	f.Party2 = f.Party2.trimmed()
	f.EffectiveDate = strings.TrimSpace(f.EffectiveDate)
	f.Purpose = strings.TrimSpace(f.Purpose)
	f.GoverningLaw = strings.TrimSpace(f.GoverningLaw)
	f.Jurisdiction = strings.TrimSpace(f.Jurisdiction)
	f.Modifications = strings.TrimSpace(f.Modifications)
	return f
}

func (p Party) trimmed() Party {
	p.Name = strings.TrimSpace(p.Name)
	p.Title = strings.TrimSpace(p.Title)
	p.Company = strings.TrimSpace(p.Company)
	p.Address = strings.TrimSpace(p.Address)
	return p
}

// stripMarkdownCodeFence removes a leading/trailing ``` or ```json fence if
// present. Despite requesting strict json_schema output, some models still
// wrap their JSON response in a markdown code block.
func stripMarkdownCodeFence(s string) string {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "```") {
		return s
	}
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}

func (c *OpenRouterClient) post(ctx context.Context, body map[string]any) (*http.Response, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("encoding request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling OpenRouter: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, fmt.Errorf("OpenRouter returned status %d", resp.StatusCode)
	}
	return resp, nil
}

const unmentionedFieldNote = " Leave as empty string if the user has not explicitly stated this value yet — never guess or infer it."

func formDataJSONSchema() map[string]any {
	textField := func(label string) map[string]any {
		return map[string]any{"type": "string", "description": label + "." + unmentionedFieldNote}
	}

	party := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name":    textField("The party's full name"),
			"title":   textField("The party's job title"),
			"company": textField("The party's company name"),
			"address": textField("The party's mailing/notice address"),
		},
		"required":             []string{"name", "title", "company", "address"},
		"additionalProperties": false,
	}

	return map[string]any{
		"name":   "nda_form_data",
		"strict": true,
		"schema": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"party1":        party,
				"party2":        party,
				"effectiveDate": textField("ISO date YYYY-MM-DD that the user explicitly stated as the effective date"),
				"mndaTermType": map[string]any{
					"type": "string", "enum": []string{"expires", "continues"},
					"description": "\"expires\" if the user wants the MNDA to expire after N years, \"continues\" if they said it continues until terminated. Default to \"expires\" if not yet discussed.",
				},
				"mndaTermYears": map[string]any{
					"type":        "integer",
					"description": "Number of years until the MNDA expires, only if the user stated one. Default to 1 if not yet discussed.",
				},
				"confidentialityTermType": map[string]any{
					"type": "string", "enum": []string{"years", "perpetuity"},
					"description": "\"years\" if confidentiality lasts a number of years, \"perpetuity\" if the user said it lasts forever. Default to \"years\" if not yet discussed.",
				},
				"confidentialityTermYears": map[string]any{
					"type":        "integer",
					"description": "Number of years confidentiality lasts, only if the user stated one. Default to 3 if not yet discussed.",
				},
				"purpose":       textField("Why confidential information is being shared, in the user's own words"),
				"governingLaw":  textField("The state whose law governs the agreement, only if the user stated one"),
				"jurisdiction":  textField("The city/county and state for dispute jurisdiction, only if the user stated one"),
				"modifications": textField("Any modifications to the standard NDA terms the user explicitly asked for"),
			},
			"required": []string{
				"party1", "party2", "effectiveDate", "mndaTermType", "mndaTermYears",
				"confidentialityTermType", "confidentialityTermYears", "purpose",
				"governingLaw", "jurisdiction", "modifications",
			},
			"additionalProperties": false,
		},
	}
}
