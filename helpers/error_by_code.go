package helpers

import (
	"avito-trainee/domains/models"
	"github.com/rs/zerolog/log"
	"net/http"
)

func GetError(code string, optional ...string) *models.ErrorType {

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
