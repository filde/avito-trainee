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
	"time"
)

func (httpServer *HttpServer) createPR(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Msgf("Couldn't get body from request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var pr *models.PullRequest
	err = json.Unmarshal(body, &pr)
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

	pr.Status = constants.OPEN_STATUS
	author, err := httpServer.storage.GetUser(pr.AuthorID)
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
		log.Error().Msgf("Couldn't get pr author %v: %v", pr.AuthorID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = httpServer.storage.TeamExists(author.TeamName)
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

	users, err := httpServer.storage.GetTeamUsers(author.TeamName, author.UserID)
	if err != nil {
		log.Error().Msgf("Couldn't get team users: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pr.Reviewers = users
	tNow := time.Now()
	pr.CreatedAt = &tNow

	err = httpServer.storage.CreatePR(pr)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			errByte, err := json.Marshal(helpers.GetError(constants.PR_EXISTS))
			if err != nil {
				log.Error().Msgf("Couldn't marshal error %v:%v", constants.PR_EXISTS, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			ok := helpers.WriteResponse(w, errByte)
			if !ok {
				return
			}
			w.WriteHeader(http.StatusConflict)
			return
		}
		log.Error().Msgf("Couldn't get team from db: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	prNew, err := httpServer.storage.GetPR(pr.PullRequestID)
	if err != nil {
		log.Error().Msgf("Couldn't get created pr from db: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for i := 0; i < len(prNew.Reviewers); i++ {
		prNew.AssignedReviewers[i] = prNew.Reviewers[i].UserID
	}

	prByte, err := json.Marshal(&models.PullRequestResponse{PR: prNew})
	if err != nil {
		log.Error().Msgf("Couldn't marshal pr json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ok := helpers.WriteResponse(w, prByte)
	if !ok {
		return
	}
	w.WriteHeader(http.StatusCreated)
}
