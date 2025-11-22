package helpers

import (
	"avito-trainee/common/constants"
	"avito-trainee/domains/models"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

func GetError(code string, optional ...string) *models.ErrorType {
	codeError := &models.ErrorType{Code: code}
	switch code {
	case constants.BAD_BODY:
		codeError.Message = constants.BAD_BODY_TEXT
	case constants.USER_EXISTS:
		codeError.Message = constants.USER_EXISTS_TEXT
	case constants.TEAM_EXISTS:
		if len(optional) != 1 {
			codeError.Message = fmt.Sprintf(constants.TEAM_EXISTS_TEXT, "team")
		} else {
			codeError.Message = fmt.Sprintf(constants.TEAM_EXISTS_TEXT, optional[0])
		}
	case constants.NOT_FOUND:
		codeError.Message = constants.NOT_FOUND_TEXT
	case constants.PR_EXISTS:
		codeError.Message = constants.PR_EXISTS_TEXT
	}
	return codeError
}

func WriteResponse(w http.ResponseWriter, response []byte) bool {
	_, err := w.Write(response)
	if err != nil {
		log.Error().Msgf("Couldn't write response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	return true
}
