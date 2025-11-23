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
		w.WriteHeader(http.StatusBadRequest)
		helpers.WriteResponse(w, errByte)
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
			w.WriteHeader(http.StatusNotFound)
			helpers.WriteResponse(w, errByte)
			return
		}
		log.Error().Msgf("Couldn't get pr author %v: %v", pr.AuthorID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pr.AssignedReviewers, err = httpServer.storage.GetTeamReviewers(author.TeamName, author.UserID)
	if err != nil {
		log.Error().Msgf("Couldn't get team users: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tNow := time.Now()
	pr.CreatedAt = &tNow

	err = httpServer.storage.CreatePR(pr)
	if err != nil {
		if helpers.IsAlreadyExists(err) {
			errByte, err := json.Marshal(helpers.GetError(constants.PR_EXISTS))
			if err != nil {
				log.Error().Msgf("Couldn't marshal error %v:%v", constants.PR_EXISTS, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusConflict)
			helpers.WriteResponse(w, errByte)
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

	prByte, err := json.Marshal(&models.PullRequestResponse{PR: prNew})
	if err != nil {
		log.Error().Msgf("Couldn't marshal pr json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	helpers.WriteResponse(w, prByte)
}

func (httpServer *HttpServer) mergePR(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusBadRequest)
		helpers.WriteResponse(w, errByte)
		return
	}

	tNow := time.Now()
	err = httpServer.storage.MergePR(pr.PullRequestID, &tNow)
	if err != nil {
		log.Error().Msgf("Couldn't merge pr %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	prNew, err := httpServer.storage.GetPR(pr.PullRequestID)
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
		log.Error().Msgf("Couldn't get merged pr from db: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	prByte, err := json.Marshal(&models.PullRequestResponse{PR: prNew})
	if err != nil {
		log.Error().Msgf("Couldn't marshal pr json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	helpers.WriteResponse(w, prByte)
}

func (httpServer *HttpServer) reassignPR(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Msgf("Couldn't get body from request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var newRev *models.NewPRReviewer
	err = json.Unmarshal(body, &newRev)
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

	//user, err := httpServer.storage.GetUser(newRev.OldReviewerID)
	//if err != nil {
	//	if errors.Is(err, gorm.ErrRecordNotFound) {
	//		errByte, err := json.Marshal(helpers.GetError(constants.NOT_FOUND))
	//		if err != nil {
	//			log.Error().Msgf("Couldn't marshal error %v:%v", constants.NOT_FOUND, err)
	//			w.WriteHeader(http.StatusInternalServerError)
	//			return
	//		}
	//		ok := helpers.WriteResponse(w, errByte)
	//		if !ok {
	//			return
	//		}
	//		w.WriteHeader(http.StatusNotFound)
	//		return
	//	}
	//	log.Error().Msgf("Couldn't get user from db: %v", err)
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}

	pr, err := httpServer.storage.GetPR(newRev.PullRequestID)
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
		log.Error().Msgf("Couldn't get pr from db: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if pr.Status == constants.MERGED_STATUS {
		log.Error().Msgf("Try changed reviewer in merged pr")
		errByte, err := json.Marshal(helpers.GetError(constants.PR_MERGED))
		if err != nil {
			log.Error().Msgf("Couldn't marshal error %v:%v", constants.PR_MERGED, err)
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
	//var usersList = []string{pr.AuthorID}
	//isReviewer := false
	//var i int
	//for i = 0; i < len(pr.Reviewers); i++ {
	//	usersList = append(usersList, pr.Reviewers[i].UserID)
	//	isReviewer = isReviewer || (pr.Reviewers[i].UserID == newRev.OldReviewerID)
	//}
	//if !isReviewer {
	//	log.Error().Msgf("%v is not reviewer of pr %v", newRev.OldReviewerID, newRev.PullRequestID)
	//	errByte, err := json.Marshal(helpers.GetError(constants.NOT_ASSIGNED))
	//	if err != nil {
	//		log.Error().Msgf("Couldn't marshal error %v:%v", constants.NOT_ASSIGNED, err)
	//		w.WriteHeader(http.StatusInternalServerError)
	//		return
	//	}
	//	ok := helpers.WriteResponse(w, errByte)
	//	if !ok {
	//		return
	//	}
	//	w.WriteHeader(http.StatusConflict)
	//	return
	//}

	//cand, err := httpServer.storage.GetTeamActiveUser(user.TeamName, usersList...)
	//if err != nil {
	//	if errors.Is(err, gorm.ErrRecordNotFound) {
	//		errByte, err := json.Marshal(helpers.GetError(constants.NO_CANDIDATE))
	//		if err != nil {
	//			log.Error().Msgf("Couldn't marshal error %v:%v", constants.NO_CANDIDATE, err)
	//			w.WriteHeader(http.StatusInternalServerError)
	//			return
	//		}
	//		ok := helpers.WriteResponse(w, errByte)
	//		if !ok {
	//			return
	//		}
	//		w.WriteHeader(http.StatusNotFound)
	//		return
	//	}
	//	log.Error().Msgf("Couldn't find new candidate: %v", err)
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}

	//pr.Reviewers[i] = cand
	err = httpServer.storage.UpdatePR(pr)
	if err != nil {
		log.Error().Msgf("Couldn't update pr reviewers: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	prNew, err := httpServer.storage.GetPR(pr.PullRequestID)
	if err != nil {
		log.Error().Msgf("Couldn't get pr with new reviewers from db: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//for i := 0; i < len(prNew.Reviewers); i++ {
	//	prNew.AssignedReviewers[i] = prNew.Reviewers[i].UserID
	//}

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
}
