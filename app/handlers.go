package app

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/RobleDev498/spaces/model"
	"gorm.io/gorm"
)

type Handler struct {
	App            *App
	HandleFunc     func(*Context, http.ResponseWriter, *http.Request)
	RequireSession bool
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &Context{}
	c.App = h.App
	c.path = r.URL.Path
	c.userAgent = r.UserAgent()
	c.ipAddress = r.RemoteAddr

	c.Params = ParamsFromRequest(r)

	if h.RequireSession {
		token, tokenLocation := ParseToken(r)
		if token != "" {
			session, err := c.App.Store.GetSession(token)

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					c.RemoveSessionCookie(w, r)
					c.Err = &model.AppErr{
						Msg:        fmt.Sprintf("Invalid or expired token. token=%s", token),
						AppErrCode: model.ERROR_INVALID_EXPIRED_SESSION,
						StatusCode: http.StatusUnauthorized,
					}
				} else {
					c.Err = &model.AppErr{
						StatusCode: http.StatusInternalServerError,
					}
				}
			} else {
				c.SetSession(session)
			}

			h.checkCSRF(c, r, token, tokenLocation, session)
		} else {
			c.Err = &model.AppErr{
				Msg:        "Token isn't present.",
				AppErrCode: model.ERROR_INVALID_EXPIRED_SESSION,
				StatusCode: http.StatusUnauthorized,
			}
		}
	}

	if c.Err == nil {
		h.HandleFunc(c, w, r)
	}

	if c.Err != nil {
		w.WriteHeader(c.Err.StatusCode)
		w.Write([]byte(c.Err.ToJson()))
	}

}

// Performs a CSRF check on the provided request with the given CSRF token.
// Returns whether or not a CSRF check occurred and whether or not it succeeded.
func (h *Handler) checkCSRF(c *Context, r *http.Request, token string, tokenLocation int, session *model.Session) (checked, passed bool) {
	checked = session != nil && c.Err != nil && tokenLocation == TOKEN_LOCATION_COOKIE && r.Method != "GET"
	passed = false

	if checked {
		csrfHeader := r.Header.Get(model.HEADER_CSRF_TOKEN)

		if csrfHeader == session.GetCSRF() {
			passed = true
		} else if r.Header.Get(model.HEADER_REQUESTED_WITH) == model.HEADER_REQUESTED_WITH_XML {
			csrfErrorMessage := "CSRF Header check failed for request - Please upgrade your web application or custom app to set a CSRF Header"
			log.Println(csrfErrorMessage)
		}

		if !passed {
			c.SetSession(&model.Session{})
			c.Err = &model.AppErr{
				Msg:        "token=" + token + " Appears to be a CSRF attempt",
				AppErrCode: model.ERROR_INVALID_EXPIRED_SESSION,
				StatusCode: http.StatusUnauthorized,
			}
		}
	}

	return
}
