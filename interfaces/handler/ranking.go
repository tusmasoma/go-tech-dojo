package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tusmasoma/go-tech-dojo/domain/model"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
	"github.com/tusmasoma/go-tech-dojo/usecase"
)

type RankingHandler interface {
	ListRankings(w http.ResponseWriter, r *http.Request)
}

type rankingHandler struct {
	ruc usecase.RankingUseCase
}

func NewRankingHandler(ruc usecase.RankingUseCase) RankingHandler {
	return &rankingHandler{
		ruc: ruc,
	}
}

type ListRankingsResponse struct {
	Rankings []struct {
		Name  string `json:"name"`
		Score int    `json:"score"`
		Rank  int    `json:"rank"`
	} `json:"rankings"`
}

func (rh *rankingHandler) ListRankings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	start, ok := rh.isValidListRankingsRequest(r)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rankings, err := rh.ruc.ListRankings(ctx, start)
	if err != nil {
		log.Error("Failed to list rankings", log.Ferror(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := rh.convertToResponseRankings(rankings)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Failed to encode rankings to JSON", log.Ferror(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (rh *rankingHandler) isValidListRankingsRequest(r *http.Request) (int, bool) {
	startStr := r.URL.Query().Get("start")
	if startStr == "" {
		log.Warn("Missing 'start' parameter")
		return 0, false
	}

	start, err := strconv.Atoi(startStr)
	if err != nil {
		log.Warn("Invalid 'start' parameter", log.Ferror(err))
		return 0, false
	}

	if start < 1 {
		log.Warn("Invalid 'start' parameter", log.Fint("start", start))
		return 0, false
	}
	return start, true
}

func (rh *rankingHandler) convertToResponseRankings(rankings []*model.Ranking) ListRankingsResponse {
	response := ListRankingsResponse{
		Rankings: make([]struct {
			Name  string `json:"name"`
			Score int    `json:"score"`
			Rank  int    `json:"rank"`
		}, 0, len(rankings)),
	}

	for _, r := range rankings {
		response.Rankings = append(response.Rankings, struct {
			Name  string `json:"name"`
			Score int    `json:"score"`
			Rank  int    `json:"rank"`
		}{
			Name:  r.UserName,
			Score: r.Score,
			Rank:  r.Rank,
		})
	}
	return response
}
