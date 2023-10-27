package nats

import (
	"gobanking/internal/user/usecase"
	"gobanking/pkg/logger"
	natsPkg "gobanking/pkg/nats"

	"github.com/nats-io/nats.go"
)

type Handler struct {
	User *UserHandler
}

type handler struct {
	Usecase  *usecase.Usecase
	Logger   logger.Logger
	natsConn *nats.Conn
}

func NewNatsHandler(usecase *usecase.Usecase, logger logger.Logger, natsConn *nats.Conn) *Handler {
	handler := &handler{
		Usecase:  usecase,
		Logger:   logger,
		natsConn: natsConn,
	}

	h := &Handler{
		User: (*UserHandler)(handler),
	}

	// user
	_, err := handler.natsConn.Subscribe(natsPkg.UserGetUserByIdReqRep, h.User.GetAccountByID)
	if err != nil {
		handler.Logger.Errorf("failed to subscribe user-service.get-user-by-id-req-rep: %v", err)
		return nil
	}

	return h
}
