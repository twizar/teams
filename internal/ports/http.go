package ports

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/twizar/teams/internal/application/service"
	"github.com/twizar/teams/internal/ports/converter"
)

const (
	minRatingQueryParam = "min-rating"
	leaguesParam        = "leagues"
	orderByQueryParam   = "order-by"
	limitQueryParam     = "limit"
	idsQueryParam       = "ids"

	queryParamsSeparator = ","

	defaultMinRatingQueryParamValue = 5
	defaultOrderByQueryParamValue   = "rating"
	defaultLimitQueryParamValue     = 0
)

type HTTPServer struct {
	serviceTeams *service.Teams
}

func NewHTTPServer(serviceTeams *service.Teams) *HTTPServer {
	return &HTTPServer{serviceTeams: serviceTeams}
}

func (s HTTPServer) Teams(writer http.ResponseWriter, request *http.Request) {
	teams, err := s.serviceTeams.AllTeams(request.Context())
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	teamDTOs := converter.EntitiesToDTOs(teams)

	err = json.NewEncoder(writer).Encode(teamDTOs)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (s HTTPServer) TeamsByID(writer http.ResponseWriter, request *http.Request) {
	ids := make([]string, 0)
	if idsParam := request.URL.Query().Get(idsQueryParam); idsParam != "" {
		ids = strings.Split(idsParam, queryParamsSeparator)
	}

	teams, err := s.serviceTeams.TeamsByID(request.Context(), ids)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	teamDTOs := converter.EntitiesToDTOs(teams)

	err = json.NewEncoder(writer).Encode(teamDTOs)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (s HTTPServer) SearchTeams(writer http.ResponseWriter, request *http.Request) {
	var (
		minRating float64 = defaultMinRatingQueryParamValue
		limit     int64   = defaultLimitQueryParamValue
		orderBy   string
		err       error
	)

	if minRatingParam := request.URL.Query().Get(minRatingQueryParam); minRatingParam != "" {
		bitSize := 64

		if minRating, err = strconv.ParseFloat(minRatingParam, bitSize); err != nil {
			log.Printf("error parsing %s param to float: %s", minRatingParam, err.Error())
			writer.WriteHeader(http.StatusBadRequest)

			_, err = fmt.Fprintf(writer, "bad param %s", minRatingParam)
			if err != nil {
				log.Printf("error writing response %s", err.Error())
			}

			return
		}
	}

	leagues := resolveLeagues(request)

	if orderBy = request.URL.Query().Get(orderByQueryParam); orderBy == "" {
		orderBy = defaultOrderByQueryParamValue
	}

	if limitParam := request.URL.Query().Get(limitQueryParam); limitParam != "" {
		base, bitSize := 10, 64

		if limit, err = strconv.ParseInt(limitParam, base, bitSize); err != nil {
			log.Printf("error parsing %s param to int: %s", limitParam, err.Error())
			writer.WriteHeader(http.StatusBadRequest)

			_, err = fmt.Fprintf(writer, "bad param %s", limitParam)
			if err != nil {
				log.Printf("error writing response %s", err.Error())
			}

			return
		}
	}

	teams, err := s.serviceTeams.FilterTeams(request.Context(), minRating, leagues, orderBy, int(limit))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	teamDTOs := converter.EntitiesToDTOs(teams)

	err = json.NewEncoder(writer).Encode(teamDTOs)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	writer.WriteHeader(http.StatusOK)
}

func resolveLeagues(request *http.Request) (leagues []string) {
	leagues = make([]string, 0)
	if leaguesParamQuery := request.URL.Query().Get(leaguesParam); leaguesParamQuery != "" {
		leagues = strings.Split(leaguesParamQuery, queryParamsSeparator)
	}

	return
}

func ConfigureRouter(server *HTTPServer) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/teams", server.TeamsByID).Queries(idsQueryParam, "").Methods(http.MethodGet)
	router.HandleFunc("/teams", server.Teams).Methods(http.MethodGet)
	router.HandleFunc("/teams/search", server.SearchTeams).Methods(http.MethodGet)

	return router
}
