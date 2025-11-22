package httpserver

import (
	"avito-trainee/common/constants"
	"avito-trainee/domains/models"
	"avito-trainee/helpers"
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"io"
	"net/http"
)

func (httpServer *HttpServer) addTeam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Msgf("Couldn't get body from request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var team *models.Team
	err = json.Unmarshal(body, &team)
	if err != nil {
		log.Error().Msgf("Couldn't unmarshal body: %v", err)
		errByte, err := json.Marshal(helpers.GetError(constants.BAD_BODY))
		if err != nil {
			log.Error().Msgf("Couldn't marshal body error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ok := helpers.WriteResponse(w, errByte)
		if !ok {
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	operationError, err := httpServer.storage.CreateTeam(team)
	if err != nil {
		log.Error().Msgf("Couldn't create team: %v", err)
		if operationError == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		errorByte, err := json.Marshal(&models.ErrorResponse{Error: operationError})
		if err != nil {
			log.Error().Msgf("Couldn't marshal error json: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ok := helpers.WriteResponse(w, errorByte)
		if !ok {
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	teamDB, err := httpServer.storage.GetTeam(team.TeamName)
	if err != nil {
		log.Error().Msgf("Couldn't get created team from db: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	teamByte, err := json.Marshal(&models.TeamResponse{Team: teamDB})
	if err != nil {
		log.Error().Msgf("Couldn't marshal team json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ok := helpers.WriteResponse(w, teamByte)
	if !ok {
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (httpServer *HttpServer) getTeam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		log.Error().Msgf("Empty team name")
		errByte, err := json.Marshal(helpers.GetError(constants.NOT_FOUND))
		if err != nil {
			log.Error().Msgf("Couldn't marshal error %v:%v", constants.NOT_FOUND, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ok := helpers.WriteResponse(w, errByte)
		if !ok {
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	team, err := httpServer.storage.GetTeam(teamName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errByte, err := json.Marshal(helpers.GetError(constants.NOT_FOUND))
			if err != nil {
				log.Error().Msgf("Couldn't marshal error %v:%v", constants.NOT_FOUND, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			ok := helpers.WriteResponse(w, errByte)
			if !ok {
				return
			}
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Error().Msgf("Couldn't get team from db: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	teamByte, err := json.Marshal(team)
	if err != nil {
		log.Error().Msgf("Couldn't marshal team json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ok := helpers.WriteResponse(w, teamByte)
	if !ok {
		return
	}
}
