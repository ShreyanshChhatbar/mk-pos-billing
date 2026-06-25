package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"mk-pos-billing/internal/infrastructure/cache"
	"mk-pos-billing/internal/infrastructure/config"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type DuplicateRequestMiddleware struct {
	cache *cache.Service
	cfg   config.SalesInvoiceConfig
}

func NewDuplicateRequestMiddleware(cache *cache.Service, cfg config.SalesInvoiceConfig) *DuplicateRequestMiddleware {
	return &DuplicateRequestMiddleware{cache: cache, cfg: cfg}
}

func (m *DuplicateRequestMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodPost && c.Request.Method != http.MethodPut {
			c.Next()
			return
		}

		raw, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "type": "Bad Request", "message": "Validation Error", "errors": gin.H{"body": []string{"invalid request body"}}})
			c.Abort()
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(raw))

		payload := buildDuplicatePayload(c, raw)
		canonicalPayload, _ := json.Marshal(payload)

		userID, _ := POSUserID(c)
		key := buildDuplicateKey(userID, m.cfg.DuplicateCheckCacheKey, c.Request.URL.Path)
		hash := hashBody(canonicalPayload)

		added, err := m.cache.Add(c.Request.Context(), key, hash, time.Duration(m.cfg.DuplicateRequestExpirySeconds)*time.Second)
		if err == nil && !added {
			cached, getErr := m.cache.Get(c.Request.Context(), key)
			if getErr == nil && cached == hash {
				c.JSON(http.StatusTooManyRequests, gin.H{"code": 429, "type": "Not Acceptable", "message": "Duplicate Request Received"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

func buildDuplicatePayload(c *gin.Context, raw []byte) map[string]interface{} {
	payload := map[string]interface{}{}
	_ = json.Unmarshal(raw, &payload)
	if userID, ok := POSUserID(c); ok {
		payload["user_id"] = userID
	}
	if storeID, ok := POSStoreID(c); ok {
		payload["store_id"] = storeID
	}
	if deviceMasterID, ok := POSDeviceMasterID(c); ok && deviceMasterID != nil {
		payload["device_master_id"] = *deviceMasterID
	}
	if tillID, ok := POSTillID(c); ok && tillID != nil {
		payload["till_id"] = *tillID
	}
	if tillTransactionID, ok := POSTillTransactionID(c); ok && tillTransactionID != nil {
		payload["till_transaction_id"] = *tillTransactionID
	}
	return payload
}

func buildDuplicateKey(userID uint64, prefix, path string) string {
	segments := strings.Split(strings.Trim(path, "/"), "/")
	seg3 := ""
	seg4 := ""
	if len(segments) > 2 {
		seg3 = segments[2]
	}
	if len(segments) > 3 {
		seg4 = segments[3]
	}
	return strconv.FormatUint(userID, 10) + "_" + prefix + "_" + seg3 + "_" + seg4
}

func hashBody(raw []byte) string {
	h := sha256.Sum256(raw)
	return hex.EncodeToString(h[:])
}
