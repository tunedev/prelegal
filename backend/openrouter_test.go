package main

import (
	"strings"
	"testing"
)

func TestStreamOpenRouterChunks_CallsOnChunkForEachDelta(t *testing.T) {
	body := strings.NewReader(
		"data: {\"choices\":[{\"delta\":{\"content\":\"Hello\"}}]}\n\n" +
			"data: {\"choices\":[{\"delta\":{\"content\":\" there\"}}]}\n\n" +
			"data: [DONE]\n\n",
	)

	var chunks []string
	full, err := streamOpenRouterChunks(body, func(chunk string) {
		chunks = append(chunks, chunk)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(chunks) != 2 || chunks[0] != "Hello" || chunks[1] != " there" {
		t.Errorf("expected chunks [Hello, there], got %v", chunks)
	}
	if full != "Hello there" {
		t.Errorf("expected accumulated text %q, got %q", "Hello there", full)
	}
}

func TestStreamOpenRouterChunks_IgnoresKeepAliveComments(t *testing.T) {
	body := strings.NewReader(
		": OPENROUTER PROCESSING\n\n" +
			"data: {\"choices\":[{\"delta\":{\"content\":\"Hi\"}}]}\n\n" +
			"data: [DONE]\n\n",
	)

	var chunks []string
	full, err := streamOpenRouterChunks(body, func(chunk string) {
		chunks = append(chunks, chunk)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(chunks) != 1 || chunks[0] != "Hi" {
		t.Errorf("expected chunks [Hi], got %v", chunks)
	}
	if full != "Hi" {
		t.Errorf("expected accumulated text %q, got %q", "Hi", full)
	}
}

func TestStreamOpenRouterChunks_SkipsEmptyDeltas(t *testing.T) {
	body := strings.NewReader(
		"data: {\"choices\":[{\"delta\":{}}]}\n\n" +
			"data: {\"choices\":[{\"delta\":{\"content\":\"x\"}}]}\n\n" +
			"data: [DONE]\n\n",
	)

	var chunks []string
	full, err := streamOpenRouterChunks(body, func(chunk string) {
		chunks = append(chunks, chunk)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(chunks) != 1 || chunks[0] != "x" {
		t.Errorf("expected chunks [x] (empty delta skipped), got %v", chunks)
	}
	if full != "x" {
		t.Errorf("expected accumulated text %q, got %q", "x", full)
	}
}
