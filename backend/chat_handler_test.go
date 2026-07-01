package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestChatHandler_StreamsReplyThenFormData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)

		if stream, _ := body["stream"].(bool); stream {
			w.Header().Set("Content-Type", "text/event-stream")
			w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"Sure, \"}}]}\n\n"))
			w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"let's start.\"}}]}\n\n"))
			w.Write([]byte("data: [DONE]\n\n"))
			return
		}

		formJSON := `{"party1":{"name":"","title":"","company":"","address":""},` +
			`"party2":{"name":"","title":"","company":"","address":""},` +
			`"effectiveDate":"","mndaTermType":"expires","mndaTermYears":1,` +
			`"confidentialityTermType":"years","confidentialityTermYears":3,` +
			`"purpose":"","governingLaw":"","jurisdiction":"","modifications":""}`
		resp := map[string]any{
			"choices": []map[string]any{{"message": map[string]any{"content": formJSON}}},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := newOpenRouterClient("test-key")
	client.baseURL = server.URL

	e := echo.New()
	e.POST("/api/chat", newChatHandler(client))

	reqBody := `{"messages":[{"role":"user","content":"hi"}]}`
	req := httptest.NewRequest(http.MethodPost, "/api/chat", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	out := rec.Body.String()
	if !strings.Contains(out, "event: message") {
		t.Errorf("expected message events in output, got: %s", out)
	}
	if !strings.Contains(out, "Sure, ") || !strings.Contains(out, "let's start.") {
		t.Errorf("expected streamed reply chunks in output, got: %s", out)
	}
	if !strings.Contains(out, "event: replyDone") {
		t.Errorf("expected a replyDone event marking the end of the streamed reply, got: %s", out)
	}
	if !strings.Contains(out, "event: formData") {
		t.Errorf("expected a formData event in output, got: %s", out)
	}
	if !strings.Contains(out, `"mndaTermType":"expires"`) {
		t.Errorf("expected extracted form data JSON in output, got: %s", out)
	}

	replyDoneIdx := strings.Index(out, "event: replyDone")
	formDataIdx := strings.Index(out, "event: formData")
	lastMessageIdx := strings.LastIndex(out, "event: message")
	if !(lastMessageIdx < replyDoneIdx && replyDoneIdx < formDataIdx) {
		t.Errorf("expected event order message -> replyDone -> formData, got: %s", out)
	}
}

func TestChatHandler_EncodesMessageChunksContainingNewlines(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)

		if stream, _ := body["stream"].(bool); stream {
			w.Header().Set("Content-Type", "text/event-stream")
			// A delta chunk that is itself a line break, as real models emit
			// between paragraphs.
			w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"line one\\nline two\"}}]}\n\n"))
			w.Write([]byte("data: [DONE]\n\n"))
			return
		}

		formJSON := `{"party1":{"name":"","title":"","company":"","address":""},` +
			`"party2":{"name":"","title":"","company":"","address":""},` +
			`"effectiveDate":"","mndaTermType":"expires","mndaTermYears":1,` +
			`"confidentialityTermType":"years","confidentialityTermYears":3,` +
			`"purpose":"","governingLaw":"","jurisdiction":"","modifications":""}`
		resp := map[string]any{
			"choices": []map[string]any{{"message": map[string]any{"content": formJSON}}},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := newOpenRouterClient("test-key")
	client.baseURL = server.URL

	e := echo.New()
	e.POST("/api/chat", newChatHandler(client))

	req := httptest.NewRequest(http.MethodPost, "/api/chat", strings.NewReader(`{"messages":[{"role":"user","content":"hi"}]}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	// The SSE protocol requires one logical event per "\n\n"-delimited
	// block; a raw, un-encoded newline inside a chunk would split it into
	// two blocks and corrupt the stream. Find the message event's data
	// line and confirm it's valid, single-line JSON that decodes back to
	// the original two-line content.
	out := rec.Body.String()
	blocks := strings.Split(out, "\n\n")
	var found bool
	for _, block := range blocks {
		if !strings.Contains(block, "event: message") {
			continue
		}
		lines := strings.Split(block, "\n")
		if len(lines) != 2 {
			t.Fatalf("expected message event to be exactly 2 lines (event + single data line), got %d: %q", len(lines), block)
		}
		data := strings.TrimPrefix(lines[1], "data: ")
		var decoded string
		if err := json.Unmarshal([]byte(data), &decoded); err != nil {
			t.Fatalf("expected data line to be valid JSON string, got %q: %v", data, err)
		}
		if decoded == "line one\nline two" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected a message event whose JSON-decoded data equals %q, got: %s", "line one\nline two", out)
	}
}

func TestChatHandler_EmitsErrorEventOnUpstreamFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client := newOpenRouterClient("test-key")
	client.baseURL = server.URL

	e := echo.New()
	e.POST("/api/chat", newChatHandler(client))

	req := httptest.NewRequest(http.MethodPost, "/api/chat", strings.NewReader(`{"messages":[{"role":"user","content":"hi"}]}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	out := rec.Body.String()
	if !strings.Contains(out, "event: error") {
		t.Errorf("expected an error event in output, got: %s", out)
	}
}

func TestChatHandler_RejectsInvalidJSON(t *testing.T) {
	client := newOpenRouterClient("test-key")

	e := echo.New()
	e.POST("/api/chat", newChatHandler(client))

	req := httptest.NewRequest(http.MethodPost, "/api/chat", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}
