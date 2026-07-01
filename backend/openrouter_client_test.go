package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestStreamChatReply_IncludesResponseBodyOnFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":{"message":"No auth credentials found"}}`))
	}))
	defer server.Close()

	client := newOpenRouterClient("")
	client.baseURL = server.URL

	_, err := client.StreamChatReply(context.Background(), []ChatMessage{{Role: "user", Content: "hi"}}, func(string) {})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "No auth credentials found") {
		t.Errorf("expected error to include the response body so operators can diagnose failures, got: %v", err)
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("expected error to include the status code, got: %v", err)
	}
}

func TestStreamChatReply_StreamsChunksAndReturnsFullText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding request body: %v", err)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("expected Authorization header, got %q", r.Header.Get("Authorization"))
		}
		if stream, _ := body["stream"].(bool); !stream {
			t.Errorf("expected stream:true in request body, got %v", body["stream"])
		}
		if _, hasFormat := body["response_format"]; hasFormat {
			t.Errorf("expected no response_format for the chat reply call, got %v", body["response_format"])
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"Hi\"}}]}\n\n"))
		w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\" there\"}}]}\n\n"))
		w.Write([]byte("data: [DONE]\n\n"))
	}))
	defer server.Close()

	client := newOpenRouterClient("test-key")
	client.baseURL = server.URL

	var chunks []string
	full, err := client.StreamChatReply(context.Background(), []ChatMessage{{Role: "user", Content: "hello"}}, func(c string) {
		chunks = append(chunks, c)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if full != "Hi there" {
		t.Errorf("expected full text %q, got %q", "Hi there", full)
	}
	if len(chunks) != 2 {
		t.Errorf("expected 2 chunks, got %v", chunks)
	}
}

func TestExtractFormData_StripsMarkdownCodeFenceIfPresent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		formJSON := `{"party1":{"name":"Alice","title":"","company":"","address":""},` +
			`"party2":{"name":"","title":"","company":"","address":""},` +
			`"effectiveDate":"","mndaTermType":"expires","mndaTermYears":1,` +
			`"confidentialityTermType":"years","confidentialityTermYears":3,` +
			`"purpose":"","governingLaw":"","jurisdiction":"","modifications":""}`

		// Some models wrap structured JSON output in a markdown code fence
		// despite strict json_schema mode being requested.
		fenced := "```json\n" + formJSON + "\n```"

		resp := map[string]any{
			"choices": []map[string]any{{"message": map[string]any{"content": fenced}}},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := newOpenRouterClient("test-key")
	client.baseURL = server.URL

	data, err := client.ExtractFormData(context.Background(), []ChatMessage{{Role: "user", Content: "I'm Alice"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.Party1.Name != "Alice" {
		t.Errorf("expected party1 name Alice, got %+v", data.Party1)
	}
}

func TestExtractFormData_TrimsWhitespaceOnlyFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Some models return a single space instead of a true empty
		// string for fields the user hasn't mentioned yet, which would
		// defeat the frontend's `value || 'placeholder'` fallback logic.
		formJSON := `{"party1":{"name":"Alice","title":" ","company":" ","address":" "},` +
			`"party2":{"name":" ","title":" ","company":" ","address":" "},` +
			`"effectiveDate":" ","mndaTermType":"expires","mndaTermYears":1,` +
			`"confidentialityTermType":"years","confidentialityTermYears":3,` +
			`"purpose":" ","governingLaw":" ","jurisdiction":" ","modifications":" "}`

		resp := map[string]any{
			"choices": []map[string]any{{"message": map[string]any{"content": formJSON}}},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := newOpenRouterClient("test-key")
	client.baseURL = server.URL

	data, err := client.ExtractFormData(context.Background(), []ChatMessage{{Role: "user", Content: "I'm Alice"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.Party1.Name != "Alice" {
		t.Errorf("expected party1 name Alice, got %q", data.Party1.Name)
	}
	if data.Party1.Company != "" || data.Party2.Name != "" || data.Purpose != "" || data.GoverningLaw != "" {
		t.Errorf("expected whitespace-only fields to be trimmed to empty string, got %+v", data)
	}
}

func TestExtractFormData_ExcludesWildcardFreeRouterFromModelList(t *testing.T) {
	// openrouter/free is a wildcard that can route to any free model,
	// including very small/unreliable ones observed in production to
	// truncate structured JSON output (finish_reason "length" after only
	// a couple of fields). The extraction call needs correctness more
	// than the conversational reply does, so it uses a curated list.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)

		models, _ := body["models"].([]any)
		for _, m := range models {
			if m == "openrouter/free" {
				t.Errorf("expected extraction model list to exclude the openrouter/free wildcard router, got %v", models)
			}
		}
		if len(models) == 0 {
			t.Errorf("expected a non-empty models list, got %v", body["models"])
		}

		formJSON := `{"party1":{"name":"","title":"","company":"","address":""},` +
			`"party2":{"name":"","title":"","company":"","address":""},` +
			`"effectiveDate":"","mndaTermType":"expires","mndaTermYears":1,` +
			`"confidentialityTermType":"years","confidentialityTermYears":3,` +
			`"purpose":"","governingLaw":"","jurisdiction":"","modifications":""}`
		resp := map[string]any{"choices": []map[string]any{{"message": map[string]any{"content": formJSON}}}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := newOpenRouterClient("test-key")
	client.baseURL = server.URL

	if _, err := client.ExtractFormData(context.Background(), []ChatMessage{{Role: "user", Content: "hi"}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExtractFormData_RetriesOnceOnMalformedResponse(t *testing.T) {
	extractionRetryDelay = 0
	defer func() { extractionRetryDelay = 2 * time.Second }()

	attempt := 0
	var modelsPerAttempt [][]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt++
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		models, _ := body["models"].([]any)
		modelsPerAttempt = append(modelsPerAttempt, models)
		w.Header().Set("Content-Type", "application/json")

		if attempt == 1 {
			// Simulates a truncated/malformed structured-output response,
			// as seen in production from a small free model.
			resp := map[string]any{
				"choices": []map[string]any{{"message": map[string]any{"content": `{"confidentialityTermType": "years"`}}},
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		formJSON := `{"party1":{"name":"Alice","title":"","company":"","address":""},` +
			`"party2":{"name":"","title":"","company":"","address":""},` +
			`"effectiveDate":"","mndaTermType":"expires","mndaTermYears":1,` +
			`"confidentialityTermType":"years","confidentialityTermYears":3,` +
			`"purpose":"","governingLaw":"","jurisdiction":"","modifications":""}`
		resp := map[string]any{"choices": []map[string]any{{"message": map[string]any{"content": formJSON}}}}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := newOpenRouterClient("test-key")
	client.baseURL = server.URL

	data, err := client.ExtractFormData(context.Background(), []ChatMessage{{Role: "user", Content: "I'm Alice"}})
	if err != nil {
		t.Fatalf("expected the retry to succeed, got error: %v", err)
	}
	if data.Party1.Name != "Alice" {
		t.Errorf("expected party1 name Alice from the retried response, got %+v", data.Party1)
	}
	if attempt != 2 {
		t.Fatalf("expected exactly 2 attempts (1 failure + 1 retry), got %d", attempt)
	}
	for _, m := range modelsPerAttempt[0] {
		if m == "openrouter/free" {
			t.Errorf("expected the first attempt to use the curated model list, got %v", modelsPerAttempt[0])
		}
	}
	found := false
	for _, m := range modelsPerAttempt[1] {
		if m == "openrouter/free" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected the retry to widen the model list to include openrouter/free as a last resort, got %v", modelsPerAttempt[1])
	}
}

func TestExtractFormData_ParsesStructuredResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding request body: %v", err)
		}
		if stream, _ := body["stream"].(bool); stream {
			t.Errorf("expected stream:false for the extraction call, got %v", body["stream"])
		}
		format, ok := body["response_format"].(map[string]any)
		if !ok || format["type"] != "json_schema" {
			t.Errorf("expected response_format.type json_schema, got %v", body["response_format"])
		}
		if temp, ok := body["temperature"].(float64); !ok || temp != 0 {
			t.Errorf("expected temperature 0 for the extraction call to minimize hallucination, got %v", body["temperature"])
		}
		if maxTokens, ok := body["max_tokens"].(float64); !ok || maxTokens < 500 {
			t.Errorf("expected a generous explicit max_tokens so the JSON schema response can't be silently truncated, got %v", body["max_tokens"])
		}

		formJSON := `{"party1":{"name":"Alice","title":"CEO","company":"Acme","address":"1 Main St"},` +
			`"party2":{"name":"Bob","title":"CTO","company":"Beta","address":"2 Side St"},` +
			`"effectiveDate":"2026-01-01","mndaTermType":"expires","mndaTermYears":1,` +
			`"confidentialityTermType":"years","confidentialityTermYears":3,` +
			`"purpose":"Evaluate partnership","governingLaw":"California","jurisdiction":"San Francisco, California",` +
			`"modifications":""}`

		resp := map[string]any{
			"choices": []map[string]any{
				{"message": map[string]any{"content": formJSON}},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := newOpenRouterClient("test-key")
	client.baseURL = server.URL

	data, err := client.ExtractFormData(context.Background(), []ChatMessage{{Role: "user", Content: "I'm Alice, CEO of Acme"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.Party1.Name != "Alice" || data.Party1.Company != "Acme" {
		t.Errorf("expected party1 Alice/Acme, got %+v", data.Party1)
	}
	if data.Party2.Name != "Bob" {
		t.Errorf("expected party2 Bob, got %+v", data.Party2)
	}
	if data.MndaTermYears != 1 || data.ConfidentialityTermYears != 3 {
		t.Errorf("expected term years 1/3, got %d/%d", data.MndaTermYears, data.ConfidentialityTermYears)
	}
	if data.GoverningLaw != "California" {
		t.Errorf("expected governingLaw California, got %q", data.GoverningLaw)
	}
}
