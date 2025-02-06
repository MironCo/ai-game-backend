package api

import (
	"fmt"

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
	// Log incoming request for debugging
	fmt.Printf("Received SMS webhook: %+v\n", c.Request.PostForm)

	// Return 200 immediately to acknowledge receipt
	c.JSON(200, gin.H{
		"status": "received",
	})

	// Process asynchronously
	go func() {
		from := c.PostForm("From")
		body := c.PostForm("Body")

		err := h.SendSMS(
			"+18885103459",
			from,
			fmt.Sprintf("Got your message: %s", body),
		)

		if err != nil {
			fmt.Printf("Error sending response SMS: %v\n", err)
		}
	}()
}
