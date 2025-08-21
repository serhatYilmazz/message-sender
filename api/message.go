package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/serhatYilmazz/message-sender/internal/message"
	"github.com/serhatYilmazz/message-sender/pkg/model"
	"github.com/sirupsen/logrus"
)

type Message struct {
	MessageService message.Service
	logger         *logrus.Logger
}

func (m Message) NewMessageHandler(messageService *message.Service) {
	app := fiber.New()

	api := app.Group("/api/messages")

	api.Get("", m.FindAllMessages)
	api.Post("")
	api.Post("/process-message-sender", m.ProcessMessageSender)

	err := app.Listen(":8080")
	if err != nil {
		m.logger.WithError(err).Errorf("error while listening port 8080")
		return
	}
}

func (m Message) FindAllMessages(ctx *fiber.Ctx) error {
	messages, err := m.MessageService.FindAllMessages(ctx.Context())
	if err != nil {
		return nil
	}

	return ctx.Status(fiber.StatusOK).JSON(messages)
}

func (m Message) ProcessMessageSender(ctx *fiber.Ctx) error {
	var messageSenderRequest model.MessageSenderRequest
	if err := ctx.BodyParser(&messageSenderRequest); err != nil {
		errorString := err.Error()
		m.logger.Error(errorString)
		errorResponse := &model.Response{
			Code:    fiber.StatusBadRequest,
			Message: errorString,
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	err := m.MessageService.ProcessMessageSender(ctx.Context(), messageSenderRequest)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&model.Response{
			Code:    500,
			Message: "internal server error",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&model.Response{
		Code:    200,
		Message: "message sender is changed as desired.",
	})
}

func (m Message) AddMessage(ctx *fiber.Ctx) error {
	var addMessageRequest model.AddMessageRequest
	if err := ctx.BodyParser(&addMessageRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if err := model.Validator.Struct(addMessageRequest); err != nil {
		errors := make([]string, 0)
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Field()+" failed on the "+err.Tag()+" tag")
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errors,
		})
	}

	return m.MessageService.SaveMessage(ctx.Context(), addMessageRequest)
}
