package handlers

import (
	"mk-pos-billing/internal/service"
	"mk-pos-billing/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CacheTestHandler struct {
	cacheMaster service.CacheMasterService
}

func NewCacheTestHandler(cacheMaster service.CacheMasterService) *CacheTestHandler {
	return &CacheTestHandler{cacheMaster: cacheMaster}
}

func (h *CacheTestHandler) GetTillCache(c *gin.Context) {
	storeID, err := strconv.Atoi(c.Query("store_id"))
	if err != nil {
		response.ValidationError(c, "Validation Error", map[string][]string{"store_id": {"invalid store_id"}})
		return
	}

	data, err := h.cacheMaster.GetTillCache(c.Request.Context(), storeID, false)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", data)
}

func (h *CacheTestHandler) GetUserAuthCache(c *gin.Context) {
	token := c.Query("token")
	
	if token == "" {
		response.ValidationError(c, "Validation Error", map[string][]string{"token": {"invalid token"}})
		return
	}

	data, err := h.cacheMaster.GetUserAuthCache(c.Request.Context(), token)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", data)
}

func (h *CacheTestHandler) GetStoreCache(c *gin.Context) {
	storeID, err := strconv.Atoi(c.Query("store_id"))
	if err != nil {
		response.ValidationError(c, "Validation Error", map[string][]string{"store_id": {"invalid store_id"}})
		return
	}

	data, err := h.cacheMaster.GetStoreCache(c.Request.Context(), storeID, false)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", data)
}

func (h *CacheTestHandler) GetProductCache(c *gin.Context) {
	productID, err := strconv.Atoi(c.Query("product_id"))
	if err != nil {
		response.ValidationError(c, "Validation Error", map[string][]string{"product_id": {"invalid product_id"}})
		return
	}

	data, err := h.cacheMaster.GetProductCache(c.Request.Context(), productID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", data)
}
