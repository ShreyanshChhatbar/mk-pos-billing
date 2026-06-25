package middleware

import (
	"mk-pos-billing/internal/service"
	"mk-pos-billing/pkg/response"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type MapTillMiddleware struct {
	cacheMaster service.CacheMasterService
}

func NewMapTillMiddleware(cacheMaster service.CacheMasterService) *MapTillMiddleware {
	return &MapTillMiddleware{cacheMaster: cacheMaster}
}

func (m *MapTillMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		storeHeader := strings.TrimSpace(c.GetHeader("store"))
		if storeHeader == "" {
			c.Next()
			return
		}

		storeID, err := strconv.ParseUint(storeHeader, 10, 64)

		// print("store_id", storeID);
		if err != nil || storeID == 0 {
			response.Error(c, http.StatusBadRequest, "Invalid Store")
			c.Abort()
			return
		}

		SetPOSStoreID(c, storeID)

		print("-------------------------------------------------------------------------------------------------\n")

		tillData, err := m.cacheMaster.GetTillCache(c.Request.Context(), int(storeID), true)
		if err != nil || len(tillData) == 0 {
			c.Next()
			return
		}

		if tillID, ok := toUint64(tillData["id"]); ok {
			SetPOSTillID(c, tillID)
		}

		c.Next()
	}
}
