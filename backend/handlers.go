package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

type chatRequest struct {
	Messages []ChatMessage `json:"messages"`
}

// newChatHandler returns an Echo handler that streams a conversational
// reply as SSE "message" events, then a "replyDone" marker once the reply
// is complete, then one final "formData" event with the current
// best-guess NDA field values extracted from the conversation.
func newChatHandler(client *OpenRouterClient) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req chatRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		}

		res := c.Response()
		res.Header().Set(echo.HeaderContentType, "text/event-stream")
		res.Header().Set("Cache-Control", "no-cache")
		res.Header().Set("Connection", "keep-alive")
		res.WriteHeader(http.StatusOK)

		replyMessages := append([]ChatMessage{{Role: "system", Content: chatSystemPrompt}}, req.Messages...)
		reply, err := client.StreamChatReply(c.Request().Context(), replyMessages, func(chunk string) {
			writeSSEString(res, "message", chunk)
		})
		if err != nil {
			log.Printf("chat: streaming reply failed: %v", err)
			writeSSEString(res, "error", err.Error())
			return nil
		}
		writeSSE(res, "replyDone", "{}")

		transcript := append(req.Messages, ChatMessage{Role: "assistant", Content: reply})
		extractMessages := append([]ChatMessage{{Role: "system", Content: extractSystemPrompt}}, transcript...)
		formData, err := client.ExtractFormData(c.Request().Context(), extractMessages)
		if err != nil {
			log.Printf("chat: extracting form data failed: %v", err)
			writeSSEString(res, "error", err.Error())
			return nil
		}

		formJSON, err := json.Marshal(formData)
		if err != nil {
			return nil
		}
		writeSSE(res, "formData", string(formJSON))

		return nil
	}
}

// writeSSEString JSON-encodes value (so embedded newlines and other
// characters can't break the single-line SSE data field or the
// "\n\n" event-boundary protocol) and writes it as an SSE event.
func writeSSEString(res *echo.Response, event, value string) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return
	}
	writeSSE(res, event, string(encoded))
}

// writeSSE writes data, which must already be safe to place on a single
// line (e.g. JSON), as an SSE event and flushes it to the client.
func writeSSE(res *echo.Response, event, data string) {
	fmt.Fprintf(res, "event: %s\ndata: %s\n\n", event, data)
	res.Flush()
}
