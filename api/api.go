package api

import (
	"fmt"
	"net/http"

	"github.com/floetenleague/floetenleague/api/apigen"
	"github.com/floetenleague/floetenleague/config"
	"github.com/floetenleague/floetenleague/conv/convgen"
	"github.com/floetenleague/floetenleague/database"
	"github.com/floetenleague/floetenleague/ui"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

const sessionID = "session_id"

func New(cfg *config.Config, db *database.DB) *echo.Echo {

	auth := &oauth2.Config{
		ClientID:     cfg.POEClientID,
		ClientSecret: cfg.POEClientSecret,
		RedirectURL:  "https://floetenleague.de/oauth2/poe/callback",
		Scopes:       []string{"account:profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://www.pathofexile.com/oauth/authorize",
			TokenURL:  "https://www.pathofexile.com/oauth/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	store := sessions.NewCookieStore([]byte(cfg.SessionKey))
	store.Options.HttpOnly = true

	app := echo.New()
	app.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		message := ""
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = fmt.Sprint(he.Message)
		} else {
			message = err.Error()
		}
		c.JSON(code, &apigen.ApiError{
			Error:       http.StatusText(code),
			Description: message,
		})
	}

	api := &api{
		cfg:   cfg,
		db:    db,
		auth:  auth,
		store: store,
	}
	apigen.RegisterHandlers(app, api)
	app.StaticFS("/", ui.Files)

	return app
}

type api struct {
	cfg   *config.Config
	db    *database.DB
	auth  *oauth2.Config
	store *sessions.CookieStore
	conv  *convgen.ConverterImpl
}
