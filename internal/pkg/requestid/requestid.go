package requestid

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/xiaowuzai/payroll/pkg/uuid"
)

var RequestId = "RequestId"

// 将 requestId 从 gin.Context 放到 context中，没有则新建
func WithRequestId(gc *gin.Context) context.Context {
	requestId, has := gc.Get(RequestId)
	if !has {
		requestId = uuid.CreateUUID()
	}

	gctx := gc.Request.Context()

	return context.WithValue(gctx, RequestId, requestId)
}
