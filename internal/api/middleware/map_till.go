package middleware

import (
	domaincache "mk-pos-billing/internal/domain/cache"
	"mk-pos-billing/pkg/response"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type MapTillMiddleware struct {
	tillCache domaincache.TillCache
}

func NewMapTillMiddleware(tillCache domaincache.TillCache) *MapTillMiddleware {
	return &MapTillMiddleware{tillCache: tillCache}
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

		tillData, err := m.tillCache.Get(c.Request.Context(), int(storeID))
		if err != nil || tillData == nil {
			c.Next()
			return
		}

		if tillData.ID > 0 {
			SetPOSTillID(c, uint64(tillData.ID))
		}

		c.Next()
	}
}
