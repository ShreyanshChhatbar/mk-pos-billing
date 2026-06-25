package middleware

import (
	"mk-pos-billing/internal/service"
	"mk-pos-billing/pkg/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CheckTillStatusMiddleware struct {
	cacheMaster service.CacheMasterService
}

func NewCheckTillStatusMiddleware(cacheMaster service.CacheMasterService) *CheckTillStatusMiddleware {
	return &CheckTillStatusMiddleware{cacheMaster: cacheMaster}
}

func (m *CheckTillStatusMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		storeID, ok := POSStoreID(c)
		if !ok {
			response.Error(c, http.StatusBadRequest, "Store Not Found")
			c.Abort()
			return
		}

		tillData, err := m.cacheMaster.GetTillCache(c.Request.Context(), int(storeID), true)
		if err != nil || len(tillData) == 0 {
			response.Error(c, http.StatusBadRequest, "Please Generate Till Number for this store till")
			c.Abort()
			return
		}

		if tillID, ok := toUint64(tillData["id"]); ok {
			SetPOSTillID(c, tillID)
		}

		tillTransaction, _ := tillData["till_transaction"].(map[string]interface{})
		if len(tillTransaction) == 0 {
			response.Error(c, http.StatusTemporaryRedirect, "Please open the till", map[string]interface{}{"data": map[string]interface{}{"is_till_open": false}})
			c.Abort()
			return
		}

		status := toString(tillTransaction["status"])
		if status != "OPEN" {
			response.Error(c, http.StatusTemporaryRedirect, "Please open the till", map[string]interface{}{"data": map[string]interface{}{"is_till_open": false}})
			c.Abort()
			return
		}

		if transactionDate, ok := toTime(tillTransaction["date"]); ok && !sameLocalDay(transactionDate) && !isDraftUpdate(c) {
			response.Error(c, http.StatusTemporaryRedirect, "Please open the till", map[string]interface{}{"data": map[string]interface{}{"is_till_open": false}})
			c.Abort()
			return
		}

		if transactionID, ok := toUint64(tillTransaction["id"]); ok {
			SetPOSTillTransactionID(c, transactionID)
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
