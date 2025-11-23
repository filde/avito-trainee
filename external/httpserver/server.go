package httpserver

import (
	"avito-trainee/domains/models"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

type StorageItf interface {
	CreateTeam(team *models.Team) (*models.ErrorType, error)
	GetTeam(name string) (*models.Team, error)

	UpdateUserActivity(userID string, isActive bool) error
	GetUser(userID string) (*models.UserFull, error)
	GetUserPR(userID string) (*models.UsersPR, error)
	GetTeamReviewers(name string, author string) ([]string, error)
	GetTeamActiveUser(team string, notAllowed ...string) (string, error)

	CreatePR(pr *models.PullRequest) error
	GetPR(id string) (*models.PullRequest, error)
	MergePR(id string, mergeTime *time.Time) error
	ChangeReviewer(oldReviewer *models.NewPRReviewer, newReviewer string) error
}

type HttpServer struct {
	storage     StorageItf
	siteHandler http.Handler
}

func InitAndStart(storage StorageItf) *HttpServer {
	httpServer := Init(storage)
	httpServer.Start()
	return httpServer
}

func Init(storage StorageItf) *HttpServer {
	httpServer := &HttpServer{
		storage: storage,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /team/add", httpServer.addTeam)
	mux.HandleFunc("GET /team/get", httpServer.getTeam)

	mux.HandleFunc("POST /users/setIsActive", httpServer.setIsActive)
	mux.HandleFunc("GET /users/getReview", httpServer.getUserReview)

	mux.HandleFunc("POST /pullRequest/create", httpServer.createPR)
	mux.HandleFunc("POST /pullRequest/merge", httpServer.mergePR)
	mux.HandleFunc("POST /pullRequest/reassign", httpServer.reassignPR)
	mux.HandleFunc("/health", httpServer.health)

	// Middlewares
	httpServer.siteHandler = httpServer.metricsMiddleware(mux)
	httpServer.siteHandler = httpServer.accessControlMiddleware(httpServer.siteHandler)
	httpServer.siteHandler = httpServer.panicMiddleware(httpServer.siteHandler)

	return httpServer
}

func (httpServer *HttpServer) Start() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "80"
	}
	log.Info().Msgf("Start HTTP Server at :%v", port)
	log.Panic().Msgf("HTTP Server ListenAndServe fatal error: %v", http.ListenAndServe(":"+port, httpServer.siteHandler))
}

func (httpServer *HttpServer) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
