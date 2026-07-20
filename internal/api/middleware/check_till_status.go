package middleware

import (
	domaincache "mk-pos-billing/internal/domain/cache"
	"mk-pos-billing/pkg/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CheckTillStatusMiddleware struct {
	tillCache domaincache.TillCache
}

func NewCheckTillStatusMiddleware(tillCache domaincache.TillCache) *CheckTillStatusMiddleware {
	return &CheckTillStatusMiddleware{tillCache: tillCache}
}

func (m *CheckTillStatusMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		storeID, ok := POSStoreID(c)
		if !ok {
			response.Error(c, http.StatusBadRequest, "Store Not Found")
			c.Abort()
			return
		}

		tillData, err := m.tillCache.Get(c.Request.Context(), int(storeID))
		if err != nil || tillData == nil {
			response.Error(c, http.StatusBadRequest, "Please Generate Till Number for this store till")
			c.Abort()
			return
		}

		if tillData.ID > 0 {
			SetPOSTillID(c, uint64(tillData.ID))
		}

		if tillData.TillTransaction == nil {
			response.Error(c, http.StatusTemporaryRedirect, "Please open the till", map[string]interface{}{"data": map[string]interface{}{"is_till_open": false}})
			c.Abort()
			return
		}

		if tillData.TillTransaction.Status != "OPEN" {
			response.Error(c, http.StatusTemporaryRedirect, "Please open the till", map[string]interface{}{"data": map[string]interface{}{"is_till_open": false}})
			c.Abort()
			return
		}

		if transactionDate, err := time.Parse("2006-01-02", tillData.TillTransaction.Date); err == nil && !sameLocalDay(transactionDate) && !isDraftUpdate(c) {
			response.Error(c, http.StatusTemporaryRedirect, "Please open the till", map[string]interface{}{"data": map[string]interface{}{"is_till_open": false}})
			c.Abort()
			return
		}

		if tillData.TillTransaction.ID > 0 {
			SetPOSTillTransactionID(c, uint64(tillData.TillTransaction.ID))
		}

		c.Next()
	}
}

func isDraftUpdate(c *gin.Context) bool {
	return c.Request.Method == http.MethodPut && c.FullPath() == "/api/v1/sales/sales-invoice/draft/:id"
}

func sameLocalDay(tillDate time.Time) bool {
	local := tillDate.In(time.Local)
	now := time.Now()
	return local.Year() == now.Year() && local.YearDay() == now.YearDay()
}
