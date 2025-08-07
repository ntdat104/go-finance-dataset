package response

import (
	"encoding/base64"

	"github.com/gin-gonic/gin"
	"github.com/ntdat104/go-finance-dataset/internal/application/constants"
	"github.com/ntdat104/go-finance-dataset/pkg/config"
	"github.com/ntdat104/go-finance-dataset/pkg/datetime"
	"github.com/ntdat104/go-finance-dataset/pkg/http"
	"github.com/ntdat104/go-finance-dataset/pkg/json"
	"github.com/ntdat104/go-finance-dataset/pkg/logger"
	"github.com/ntdat104/go-finance-dataset/pkg/rsa"
	"github.com/ntdat104/go-finance-dataset/pkg/uuid"
)

type Meta struct {
	MessageID string `json:"message_id"`
	Timestamp int64  `json:"timestamp"`
	Datetime  string `json:"datetime"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Token     string `json:"token,omitempty"`
}

type Response struct {
	Meta Meta `json:"meta"`
	Data any  `json:"data,omitempty"`
}

func getMessageID(ctx *gin.Context) string {
	messageID := ctx.GetHeader(constants.X_Message_ID)
	if messageID == "" {
		messageID = uuid.NewShortUUID()
	}
	return messageID
}

func buildResponse(ctx *gin.Context, code int, obj any) {
	now := datetime.GetCurrentMiliseconds()

	response := Response{
		Meta: Meta{
			MessageID: getMessageID(ctx),
			Timestamp: now,
			Datetime:  datetime.ConvertMillisecondsToString(now, datetime.YYYY_MM_DD_HH_MM_SS),
			Code:      code,
			Message:   http.StatusText(code),
		},
		Data: obj,
	}

	// Set custom response header
	privateKey, err := base64.StdEncoding.DecodeString(config.GetGlobalConfig().App.PrivateKey)
	if err != nil {
		logger.Error("base64.StdEncoding.DecodeString has error: " + err.Error())
	}
	responseStr, err := json.ToJSON(response)
	if err != nil {
		logger.Error("json.ToJSON has error: " + err.Error())
	}
	signature, err := rsa.SignMessage(string(privateKey), responseStr, "SHA256")
	if err != nil {
		logger.Error("rsa.SignMessage has error: " + err.Error())
	}
	logger.Debug("Signature: " + signature)
	logger.Warn("response: " + responseStr)
	ctx.Header(constants.Signature, signature)

	ctx.JSON(code, response)
}

func JSON(ctx *gin.Context, code int, obj any) {
	buildResponse(ctx, code, obj)
}

func Success(ctx *gin.Context, obj any) {
	buildResponse(ctx, http.StatusOK, obj)
}

func Failed(ctx *gin.Context, obj any) {
	buildResponse(ctx, http.StatusBadRequest, obj)
}
