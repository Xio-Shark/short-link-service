package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"short-link-service/internal/service"
)

type Handler struct {
	Service *service.ShortLinkService
}

type CreateLinkRequest struct {
	URL      string `json:"url"`
	ExpireAt string `json:"expire_at"`
}

func (h *Handler) Health(c *gin.Context) {
	respondOK(c, gin.H{"status": "ok"})
}

func (h *Handler) CreateLink(c *gin.Context) {
	var req CreateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	code, err := h.Service.Create(c.Request.Context(), req.URL, req.ExpireAt)
	if err != nil {
		if isBadRequest(err) {
			respondError(c, http.StatusBadRequest, 40001, err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, 50001, "create failed")
		return
	}

	respondOK(c, gin.H{"code": code, "short_url": "/" + code})
}

func (h *Handler) Redirect(c *gin.Context) {
	code := c.Param("code")
	access := service.AccessInfo{
		IP:        c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	}
	url, err := h.Service.Resolve(c.Request.Context(), code, access)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			respondError(c, http.StatusNotFound, 40401, "not found")
			return
		}
		respondError(c, http.StatusInternalServerError, 50001, "resolve failed")
		return
	}

	c.Redirect(http.StatusFound, url)
}

func (h *Handler) Stats(c *gin.Context) {
	code := c.Param("code")
	pv, uv, err := h.Service.Stats(c.Request.Context(), code)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			respondError(c, http.StatusNotFound, 40401, "not found")
			return
		}
		respondError(c, http.StatusInternalServerError, 50001, "stats failed")
		return
	}
	respondOK(c, gin.H{"pv": pv, "uv": uv})
}

func isBadRequest(err error) bool {
	if errors.Is(err, service.ErrInvalidRequest) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "invalid") || strings.Contains(msg, "required")
}
