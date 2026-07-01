package main

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

type chatCompletionChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

// streamOpenRouterChunks reads an OpenRouter SSE response body, calling
// onChunk for each non-empty content delta, and returns the full
// accumulated text once the stream ends.
func streamOpenRouterChunks(body io.Reader, onChunk func(chunk string)) (string, error) {
	var full strings.Builder

	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		line := scanner.Text()
		data, ok := strings.CutPrefix(line, "data: ")
		if !ok {
			continue
		}
		if data == "[DONE]" {
			break
		}

		var chunk chatCompletionChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			return "", err
		}
		if len(chunk.Choices) == 0 {
			continue
		}

		content := chunk.Choices[0].Delta.Content
		if content == "" {
			continue
		}
		full.WriteString(content)
		onChunk(content)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return full.String(), nil
}
