package handler

import (
	"database/sql"
	"net/http"

	"github.com/elanq/tinyurl-go/model"
	"github.com/elanq/tinyurl-go/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type URL interface {
	Create(e echo.Context) error
	GetByShortURL(e echo.Context) error
}

type url struct {
	s service.URL
}

// Create implements URL.
func (u *url) Create(e echo.Context) error {
	var urlModel model.URL
	err := e.Bind(&urlModel)
	if err != nil {
		return err
	}
	u.s.Create(e.Request().Context(), urlModel)
	return e.String(http.StatusOK, "ok")
}

// GetByShortURL implements URL.
func (u *url) GetByShortURL(e echo.Context) error {
	shortUrl := e.Param("url")
	model, err := u.s.GetByShortURL(e.Request().Context(), shortUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.String(http.StatusNotFound, "not found")
		}
		log.Error(err)
		return err
	}
	return e.Redirect(http.StatusFound, model.LongURL)
}

func NewURL(s service.URL) URL {
	return &url{
		s: s,
	}
}
