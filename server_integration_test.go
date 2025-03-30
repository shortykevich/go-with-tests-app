package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	storage := &InMemStorage{
		Scores:   map[string]int{},
		winCalls: []string{},
	}
	server := &PlayersScoreServer{Storage: storage}
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostRequest(player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreRequest(player))

	assertStatus(t, response.Code, http.StatusOK)
	assertResponseBody(t, response.Body.String(), "3")
}
