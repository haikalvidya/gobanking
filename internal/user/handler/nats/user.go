package nats

import (
	"context"
	"gobanking/internal/user/dto"
	"gobanking/pkg/middlewares"
	redisPkg "gobanking/pkg/redis"
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
		msg.Respond([]byte("error unmarshal request"))
		return
	}

	token := dtoGetUser.Token
	// check if token is valid
	if token == "" {
		h.Logger.Error("token is empty")
		msg.Respond([]byte("token is empty"))
		return
	}
	// to struct
	var user *middlewares.UserModel
	key := redisPkg.GetKeyOfTokenUserFromRedis(token, redisPkg.Token)

	user, err = redisPkg.GetDataFromRedis(ctx, h.redis, key, user)
	if err != nil {
		h.Logger.Errorf("failed to get redis user otp: %v", err)
		msg.Respond([]byte("failed to get redis user otp"))
		return
	}

	dataResp, err := h.Usecase.User.GetByID(ctx, user.ID.String())
	if err != nil {
		h.Logger.Errorf("failed to get user: %v", err)
		msg.Respond([]byte("failed to get user"))
		return
	}

	resp := dto.GetUserByIDResponse{
		Data: *dataResp,
	}

	respByte, err := serializer.Marshal(resp)
	if err != nil {
		h.Logger.Errorf("failed to marshal response: %v", err)
		msg.Respond([]byte("failed to marshal response"))
		return
	}

	msg.Respond(respByte)
}
