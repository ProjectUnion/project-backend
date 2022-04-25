package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/ProjectUnion/project-backend.git/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Register(c *gin.Context) {
	var inp domain.UserAuth
	if err := c.BindJSON(&inp); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid input body")
		return
	}

	res, err := h.services.Authorization.Register(c, inp)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refreshToken",
		Value:    res.RefreshToken,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
	})

	h.logger.Infof("Register user %s", res.UserID)

	c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken": res.AccessToken,
		//"refreshToken": res.RefreshToken,
		"userID": res.UserID,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var inp domain.UserAuth
	if err := c.BindJSON(&inp); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid input body")
		return
	}

	res, err := h.services.Authorization.Login(c, inp.Email, inp.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refreshToken",
		Value:    res.RefreshToken,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
	})

	h.logger.Infof("Login admin %s", res.UserID)

	c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken": res.AccessToken,
		//"refreshToken": res.RefreshToken,
		"userID": res.UserID,
	})
}

func (h *Handler) Refresh(c *gin.Context) {
	cookie, err := c.Request.Cookie("refreshToken")
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	refreshToken := cookie.Value

	res, err := h.services.Authorization.Refresh(c, refreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refreshToken",
		Value:    res.RefreshToken,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
	})

	h.logger.Infof("Refresh admin %s", res.UserID)

	c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken": res.AccessToken,
		"position":    res.Position,
		//"refreshToken": res.RefreshToken,
		"userID": res.UserID,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	cookie, err := c.Request.Cookie("refreshToken")
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	refreshToken := cookie.Value

	if err := h.services.Authorization.Logout(c, refreshToken); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:   "refreshToken",
		Value:  "",
		MaxAge: -1,
	})

	h.logger.Infof("Logout admin %s", refreshToken)

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}
