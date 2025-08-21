package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/serhatYilmazz/message-sender/internal/message"
	"github.com/serhatYilmazz/message-sender/pkg/model"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/fiber-swagger"
)

type MessageHandler struct {
	MessageService message.Service
	logger         *logrus.Logger
}

func NewMessageHandler(messageService message.Service, logger *logrus.Logger) {
	messageHandler := MessageHandler{
		MessageService: messageService,
		logger:         logger,
	}
	app := fiber.New()

	api := app.Group("/api/messages")

	api.Get("", messageHandler.FindAllMessages)
	api.Post("", messageHandler.AddMessage)
	api.Post("/process-message-sender", messageHandler.ProcessMessageSender)
	app.Get("/*", fiberSwagger.WrapHandler)

	err := app.Listen(":8080")
	if err != nil {
		messageHandler.logger.WithError(err).Errorf("error while listening port 8080")
		return
	}
}

// FindAllMessages godoc
// @Summary Get all messages
// @Description Retrieve all messages from the database
// @Tags messages
// @Accept json
// @Produce json
// @Success 200 {array} message.Message
// @Failure 500 {object} model.Response
// @Router /api/messages [get]
func (m MessageHandler) FindAllMessages(ctx *fiber.Ctx) error {
	messages, err := m.MessageService.FindAllMessages(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&model.Response{
			Code:    500,
			Message: "Failed to retrieve messages",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(messages)
}

// ProcessMessageSender godoc
// @Summary Process message sender settings
// @Description Enable or disable the message sender functionality
// @Tags messages
// @Accept json
// @Produce json
// @Param request body model.MessageSenderRequest true "Message sender settings"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/messages/process-message-sender [post]
func (m MessageHandler) ProcessMessageSender(ctx *fiber.Ctx) error {
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

// AddMessage godoc
// @Summary Add a new message
// @Description Create a new message with content and recipient phone number
// @Tags messages
// @Accept json
// @Produce json
// @Param request body model.AddMessageRequest true "Message data"
// @Success 200 {object} model.AddMessageResponse
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/messages [post]
func (m MessageHandler) AddMessage(ctx *fiber.Ctx) error {
	var addMessageRequest model.AddMessageRequest
	if err := ctx.BodyParser(&addMessageRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&model.Response{
			Code:    400,
			Message: "invalid request body",
		})
	}

	if err := model.Validator.Struct(addMessageRequest); err != nil {
		errors := make([]string, 0)
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Field()+" failed on the "+err.Tag()+" tag")
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(&model.Response{
			Code:    400,
			Message: "validation failed: " + errors[0],
		})
	}

	return m.MessageService.SaveMessage(ctx.Context(), addMessageRequest)
}
