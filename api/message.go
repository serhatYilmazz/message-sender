package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/serhatYilmazz/message-sender/internal/cache"
	"github.com/serhatYilmazz/message-sender/internal/message"
	"github.com/serhatYilmazz/message-sender/internal/scheduler"
	"github.com/serhatYilmazz/message-sender/pkg/model"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/fiber-swagger"
)

type MessageHandler struct {
	MessageService          message.Service
	SchedulerControlService scheduler.ControlService
	CacheService            cache.Service
	logger                  *logrus.Logger
}

func NewMessageHandler(messageService message.Service, schedulerControlService scheduler.ControlService, cacheService cache.Service, logger *logrus.Logger) {
	messageHandler := MessageHandler{
		MessageService:          messageService,
		SchedulerControlService: schedulerControlService,
		CacheService:            cacheService,
		logger:                  logger,
	}
	app := fiber.New()

	api := app.Group("/api/messages")

	api.Get("", messageHandler.FindAllMessages)
	api.Post("", messageHandler.AddMessage)
	api.Post("/process-message-sender", messageHandler.ProcessMessageSender)
	api.Get("/scheduler-status", messageHandler.GetSchedulerStatus)

	api.Get("/webhook-delivery/:messageId", messageHandler.GetWebhookDelivery)

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
// @Success 200 {array} model.MessageDto
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

	err := m.SchedulerControlService.ProcessMessageSender(ctx.Context(), messageSenderRequest)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&model.Response{
			Code:    500,
			Message: "internal server error",
		})
	}

	statusMessage := "message sender disabled"
	if messageSenderRequest.IsMessageSenderEnabled {
		statusMessage = "message sender enabled"
	}

	return ctx.Status(fiber.StatusOK).JSON(&model.Response{
		Code:    200,
		Message: statusMessage,
	})
}

// GetSchedulerStatus godoc
// @Summary Get scheduler status
// @Description Get the current status of the message scheduler
// @Tags messages
// @Accept json
// @Produce json
// @Success 200 {object} map[string]bool
// @Failure 500 {object} model.Response
// @Router /api/messages/scheduler-status [get]
func (m MessageHandler) GetSchedulerStatus(ctx *fiber.Ctx) error {
	isRunning := m.SchedulerControlService.GetSchedulerStatus(ctx.Context())

	return ctx.Status(fiber.StatusOK).JSON(map[string]bool{
		"isRunning": isRunning,
	})
}

// AddMessage godoc
// @Summary Add a new message
// @Description Create a new message with content and recipient phone number
// @Tags messages
// @Accept json
// @Produce json
// @Param request body model.AddMessageRequest true "Message data"
// @Success 200 {object} model.MessageDto
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

	savedMessage, err := m.MessageService.SaveMessage(ctx.Context(), addMessageRequest)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&model.Response{
			Code:    500,
			Message: "internal server error",
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(savedMessage)
}

// GetWebhookDelivery godoc
// @Summary Get webhook delivery record
// @Description Retrieve webhook delivery record by message ID from cache
// @Tags webhook
// @Accept json
// @Produce json
// @Param messageId path string true "Message ID"
// @Success 200 {object} cache.WebhookDelivery
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/messages/webhook-delivery/{messageId} [get]
func (m MessageHandler) GetWebhookDelivery(ctx *fiber.Ctx) error {
	messageId := ctx.Params("messageId")
	if messageId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(&model.Response{
			Code:    400,
			Message: "message ID is required",
		})
	}

	delivery, err := m.CacheService.GetDeliveryRecord(ctx.Context(), messageId)
	if err != nil {
		m.logger.WithError(err).Errorf("failed to get webhook delivery for message ID: %s", messageId)
		return ctx.Status(fiber.StatusInternalServerError).JSON(&model.Response{
			Code:    500,
			Message: "failed to retrieve webhook delivery record",
		})
	}

	if delivery == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(&model.Response{
			Code:    404,
			Message: "webhook delivery record not found",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(delivery)
}
