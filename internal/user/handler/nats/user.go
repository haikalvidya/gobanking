package nats

import (
	"context"
	"gobanking/internal/user/dto"
	"gobanking/pkg/serializer"

	"github.com/nats-io/nats.go"
)

type UserHandler handler

func (h *UserHandler) GetAccountByID(msg *nats.Msg) {
	ctx := context.Background()
	req := msg.Data
	var dtoGetUser dto.GetUserByIDRequest
	err := serializer.Unmarshal(req, &dtoGetUser)
	if err != nil {
		h.Logger.Errorf("error unmarshal request: %v", err)
		return
	}

	dataResp, err := h.Usecase.User.GetByID(ctx, dtoGetUser.ID)
	if err != nil {
		h.Logger.Errorf("failed to get user: %v", err)
		return
	}

	resp := dto.GetUserByIDResponse{
		Data: *dataResp,
	}

	respByte, err := serializer.Marshal(resp)
	if err != nil {
		h.Logger.Errorf("failed to marshal response: %v", err)
		return
	}

	msg.Respond(respByte)
}
