package queue

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/app/service/queue"
	"github.com/axlle-com/blog/pkg/user/queue/model"
	"github.com/axlle-com/blog/pkg/user/service"
)

func NewUserQueueHandler(
	userService *service.UserService,
) contract.QueueHandler {
	return &queueHandler{
		userService: userService,
	}
}

type queueHandler struct {
	userService *service.UserService
}

func (qh *queueHandler) Run(data []byte) {
	action, payload, err := models.NewEnvelopeQueue().Convert(data)
	if err != nil {
		logger.Errorf("[user][queue] error: %v", err)
		return
	}

	switch action {
	case queue.Create:
		if err := qh.create(payload); err != nil {
			logger.Errorf("[user][queue][create] error: %v", err)
		}
	default:
		logger.Errorf("[user][queue] unknown action: %s", action)
	}
}

func (qh *queueHandler) create(payload []byte) error {
	var obj model.User

	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.DisallowUnknownFields()

	if err := dec.Decode(&obj); err != nil {
		return fmt.Errorf("incorrect data format: %v", err)
	}

	_, err := qh.userService.CreateFromInterface(obj.Model())
	return err
}
