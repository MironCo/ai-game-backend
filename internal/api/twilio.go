package api

import (
	"fmt"
	"rd-backend/internal/ai"
	"rd-backend/internal/db"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TextingHandler struct {
	twilioClient *twilio.RestClient
	dbHandler    *db.DBHandler
	aiHandler    *ai.AIHandler
}

func NewTextingHandler(dbHandler *db.DBHandler, aiHandler *ai.AIHandler) *TextingHandler {
	return &TextingHandler{
		twilioClient: twilio.NewRestClient(),
		dbHandler:    dbHandler,
		aiHandler:    aiHandler,
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

// get message
// retrieve number
// sql lookup in texts database, based on number received and current number
// create response
// send text message back from that number, based on previous text messages
// therefore : find where rec# and send# are player/ai or ai/player.
// how to get the number?
// ai has it in json, each AI will have a different number

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
	to := c.PostForm("To")
	from := c.PostForm("From")
	body := c.PostForm("Body")

	// Process the message (replace this with your actual processing logic)
	processedResponse := h.processMessage(from, to, body)

	// Log what happened
	fmt.Printf("Processed message from %s: %s -> %s\n", from, body, processedResponse)

	// Return TwiML response with the processed result
	c.String(200, fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
		<Response>
			<Message>%s</Message>
		</Response>`, processedResponse))
}

// Add your processing logic here
func (h *TextingHandler) processMessage(from string, to string, message string) string {
	// FROM is player phone number, TO is AI phone number

	// this is some high grade fuckin annoyance
	if !strings.HasPrefix(from, "+") {
		from = "+" + strings.TrimSpace(strings.TrimPrefix(from, "+"))
	}
	if !strings.HasPrefix(to, "+") {
		to = "+" + strings.TrimSpace(strings.TrimPrefix(to, "+"))
	}

	player, err := h.dbHandler.GetPlayerByPhoneNumber(from)
	if err != nil {
		fmt.Println("Could not find player by phone number: " + err.Error())
		return "Sorry, your number isn't registered in our system."
	}
	if err := h.dbHandler.AddTextToDatabase(player.UnityID, message, from, to, from); err != nil {
		fmt.Println("Could not add text to database.")
		return "Could not add text to database."
	}
	textMessage, err := h.dbHandler.GetLastTextsFromDB(player.UnityID, to, 4)
	if err != nil {
		fmt.Println("Could not get last texts from DB")
		return "Could not get last texts from DB."
	}
	completion, err := h.aiHandler.GetTextCompletion(message, textMessage, to, from)
	if err != nil {
		fmt.Println("Could not get text completion")
		return "Couldn't process completion."
	}
	if err := h.dbHandler.AddTextToDatabase(player.UnityID, *completion, to, from, from); err != nil {
		fmt.Println("Could not add text from AI to player to database.")
	}

	return *completion
}
