package handler

import (
	"em-test/cmd/internal/repository"
	"net/http"
)

type SubscriptionHandler struct {
	Repo repository.SubscriptionRepo
}

func (SH *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request)     {}
func (SH *SubscriptionHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {}
func (SH *SubscriptionHandler) GetByID(w http.ResponseWriter, r *http.Request)    {}
func (SH *SubscriptionHandler) GetList(w http.ResponseWriter, r *http.Request)    {}
func (SH *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request)     {}
func (SH *SubscriptionHandler) Search(w http.ResponseWriter, r *http.Request)     {}
