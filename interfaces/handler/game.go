package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"
	"github.com/tusmasoma/go-tech-dojo/usecase"
)

type GameHandler interface {
	FinishGame(w http.ResponseWriter, r *http.Request)
	DrawGacha(w http.ResponseWriter, r *http.Request)
}

type gameHandler struct {
	guc usecase.GameUseCase
}

func NewGameHandler(guc usecase.GameUseCase) GameHandler {
	return &gameHandler{
		guc: guc,
	}
}

type FinishGameRequest struct {
	Score int `json:"score"`
}

type FinishGameResponse struct {
	Coin int `json:"coin"`
}

func (gh *gameHandler) FinishGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var requestBody FinishGameRequest
	defer r.Body.Close()
	if !gh.isValidFinishGameRequest(r.Body, &requestBody) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	coin, err := gh.guc.FinishGame(ctx, requestBody.Score)
	if err != nil {
		log.Error("Failed to finish game", log.Ferror(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(FinishGameResponse{
		Coin: coin,
	}); err != nil {
		log.Error("Failed to encode response to JSON", log.Ferror(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (gh *gameHandler) isValidFinishGameRequest(body io.ReadCloser, requestBody *FinishGameRequest) bool {
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Error("Failed to decode request body: %v", err)
		return false
	}
	if requestBody.Score < 0 {
		log.Warn("Invalid request body: %v", requestBody)
		return false
	}
	return true
}

type DrawGachaRequest struct {
	Times int `json:"times"`
}

type DrawGachaResponse struct {
	Results []struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Rarity int    `json:"rarity"`
		IsNew  bool   `json:"is_new"`
	} `json:"results"`
}

func (gh *gameHandler) DrawGacha(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var requestBody DrawGachaRequest
	defer r.Body.Close()
	if !gh.isValidDrawGachaRequest(r.Body, &requestBody) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	gachaResults, err := gh.guc.DrawGacha(ctx, requestBody.Times)
	if err != nil {
		log.Error("Failed to draw gacha", log.Ferror(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := gh.convertToDrawGachaResponse(gachaResults)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Failed to encode response to JSON", log.Ferror(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (gh *gameHandler) isValidDrawGachaRequest(body io.ReadCloser, requestBody *DrawGachaRequest) bool {
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Error("Failed to decode request body: %v", err)
		return false
	}
	if requestBody.Times < 0 || requestBody.Times > 10 {
		log.Warn("Invalid request body: %v", requestBody)
		return false
	}
	return true
}

func (gh *gameHandler) convertToDrawGachaResponse(gachaResults []*usecase.GachaResult) DrawGachaResponse {
	var results []struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Rarity int    `json:"rarity"`
		IsNew  bool   `json:"is_new"`
	}
	for _, item := range gachaResults {
		results = append(results, struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Rarity int    `json:"rarity"`
			IsNew  bool   `json:"is_new"`
		}{
			ID:     item.ID,
			Name:   item.Name,
			Rarity: item.Rarity,
			IsNew:  !item.Has,
		})
	}
	return DrawGachaResponse{
		Results: results,
	}
}
