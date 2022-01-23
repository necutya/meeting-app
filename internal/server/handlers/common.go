package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"server/internal/models"

	log "github.com/sirupsen/logrus"
)
const (
	contentTypeJSON = "application/json"
)

func SendEmptyResponse(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(statusCode)
}

// SendResponse - common method for encoding and writing any json response.
func SendResponse(w http.ResponseWriter, statusCode int, respBody interface{}) {
	w.Header().Set("Content-Type", contentTypeJSON)

	binRespBody, err := json.Marshal(respBody)
	if err != nil {
		statusCode = http.StatusInternalServerError

		log.Error(err)
	}

	w.WriteHeader(statusCode)
	// nolint
	_, err = w.Write(binRespBody)
	if err != nil {
		log.Error(err)
	}
}

func SendHTTPError(w http.ResponseWriter, err error) {
	var (
		statusCode int
		errCode    string
		errMessage string
		errs       []models.FieldError
	)

	switch err {
	case models.ErrUnauthorized:
		statusCode = http.StatusUnauthorized
		errCode = "unauthorized"
	case models.ErrNotFound:
		statusCode = http.StatusNotFound
		errCode = "not_found"
	case models.ErrForbidden:
		statusCode = http.StatusForbidden
		errCode = "forbidden"
	case models.ErrAlreadyExist,
		models.ErrConflict:
		statusCode = http.StatusConflict
		errCode = "conflict"
	default:
		switch v := err.(type) {
		case models.BadRequest:
			statusCode = http.StatusBadRequest
			errs = v.Errors
			errCode = "bad_request"
			errMessage = v.Msg
		case models.InternalError:
			statusCode = http.StatusInternalServerError
			errCode = "internal"
		default:
			statusCode = http.StatusServiceUnavailable
			errCode = "service_unavailable"
		}
	}

	log.WithError(err).Error("handler error")

	SendResponse(w, statusCode, struct {
		Code    string              `json:"code,omitempty"`
		Message string              `json:"message,omitempty"`
		Errors  []models.FieldError `json:"validation_errors,omitempty"`
	}{
		Code:    errCode,
		Message: errMessage,
		Errors:  errs,
	})
}

func UnmarshalRequest(r *http.Request, body interface{}) error {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return models.InternalError(err.Error())
	}
	defer r.Body.Close()

	if err := json.Unmarshal(reqBody, body); err != nil {
		if e, ok := err.(*json.UnmarshalTypeError); ok {
			return models.BadRequest{
				Msg: err.Error(),
				Errors: []models.FieldError{{
					Field: strings.ToLower(e.Field),
					Code:  "validation_is_" + e.Type.String(),
				}},
			}
		}

		return models.BadRequest{Msg: err.Error()}
	}

	return nil
}