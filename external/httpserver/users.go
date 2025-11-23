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

func (httpServer *HttpServer) setIsActive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Msgf("Couldn't get body from request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var user *models.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Error().Msgf("Couldn't unmarshal body: %v", err)
		errByte, err := json.Marshal(helpers.GetError(constants.BAD_BODY))
		if err != nil {
			log.Error().Msgf("Couldn't marshal body error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.WriteResponse(w, errByte)
		return
	}

	err = httpServer.storage.UpdateUserActivity(user.UserID, user.IsActive)
	if err != nil {
		log.Error().Msgf("Couldn't update user activity: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userNew, err := httpServer.storage.GetUser(user.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errByte, err := json.Marshal(helpers.GetError(constants.NOT_FOUND))
			if err != nil {
				log.Error().Msgf("Couldn't marshal error %v:%v", constants.NOT_FOUND, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNotFound)
			helpers.WriteResponse(w, errByte)
			return
		}
		log.Error().Msgf("Couldn't get updated user %v: %v", user.UserID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userByte, err := json.Marshal(&models.UserResponse{User: userNew})
	if err != nil {
		log.Error().Msgf("Couldn't marshal user json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	helpers.WriteResponse(w, userByte)
}

func (httpServer *HttpServer) getUserReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		log.Error().Msgf("Empty user id")
		errByte, err := json.Marshal(helpers.GetError(constants.NOT_FOUND))
		if err != nil {
			log.Error().Msgf("Couldn't marshal error %v:%v", constants.NOT_FOUND, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.WriteResponse(w, errByte)
		return
	}

	user, err := httpServer.storage.GetUserPR(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errByte, err := json.Marshal(helpers.GetError(constants.NOT_FOUND))
			if err != nil {
				log.Error().Msgf("Couldn't marshal error %v:%v", constants.NOT_FOUND, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNotFound)
			helpers.WriteResponse(w, errByte)
			return
		}
		log.Error().Msgf("Couldn't get user from db: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userByte, err := json.Marshal(user)
	if err != nil {
		log.Error().Msgf("Couldn't marshal user json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	helpers.WriteResponse(w, userByte)
}
