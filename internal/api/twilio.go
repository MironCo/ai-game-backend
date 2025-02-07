package api

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TextingHandler struct {
	twilioClient *twilio.RestClient
}

func NewTextingHandler() *TextingHandler {
	return &TextingHandler{
		twilioClient: twilio.NewRestClient(),
	}
}

func (h *TextingHandler) SendSMS(from, to, message string) error {
	params := &openapi.CreateMessageParams{
		To:   &to,
		From: &from,
		Body: &message,
	}

	_, err := h.twilioClient.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("error sending SMS: %v", err)
	}

	return nil
}

func (h *TextingHandler) SendSMSBasic() {
	err := h.SendSMS("+18885103459", "+16074221508", "Hello from your AI friend!")
	if err != nil {
		fmt.Printf("Failed to send SMS: %v\n", err)
		return
	}
	fmt.Println("Message sent successfully!")
}

func (h *TextingHandler) ReceiveSMS(c *gin.Context) {
	// Set content type to XML
	c.Header("Content-Type", "text/xml")

	// Get message details
	from := c.PostForm("From")
	body := c.PostForm("Body")

	// Process the message (replace this with your actual processing logic)
	processedResponse := processMessage(body)

	// Log what happened
	fmt.Printf("Processed message from %s: %s -> %s\n", from, body, processedResponse)

	// Return TwiML response with the processed result
	c.String(200, fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
		<Response>
			<Message>%s</Message>
		</Response>`, processedResponse))
}

// Add your processing logic here
func processMessage(message string) string {
	// Example processing - replace with your actual logic
	switch strings.ToLower(message) {
	case "hi", "hello":
		return "Hello! How can I help you today?"
	case "help":
		return "Available commands: hi, help, status"
	case "status":
		return "All systems operational!"
	default:
		return fmt.Sprintf("You said: %s. What would you like to know?", message)
	}
}
