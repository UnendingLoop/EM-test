package handler

import (
	"em-test/cmd/internal/model"
	"em-test/cmd/internal/repository"
	"em-test/cmd/internal/service"
	"em-test/cmd/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// SubscriptionHandler provides process to HTTP-requests
type SubscriptionHandler struct {
	Service *service.SubscriptionService
}

func CreateHandler(db *gorm.DB) *SubscriptionHandler {
	return &SubscriptionHandler{Service: service.CreateService(db)}
}

// Create - хендлер для создания новой подписки в базе
// @Summary      Cоздание новой подписки
// @Description  Создаёт новую подписку из данных в теле запроса
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        subscription  body      model.RawSubscription  true  "Subscription info" example(`{"service_name": "Yandex Plus","price": 400,"user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba","start_date": "07-2025"}`)
// @Success      201   {object}  model.RawSubscription "Subscription successfully created"
// @Failure      400   {string}  string  "Incomplete/incorrect data input"
// @Failure      409   {string}  string  "Subscription already exists"
// @Failure      500   {string}  string  "Internal server error"
// @Router       /subscriptions [post]
func (SH *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var newSub model.RawSubscription

	if err := json.NewDecoder(r.Body).Decode(&newSub); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	if err := SH.Service.CreateSubscription(r.Context(), &newSub); err != nil {
		switch {
		case errors.Is(err, repository.ErrEmptySomeFields):
			http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
			return
		case errors.Is(err, repository.ErrSubExists):
			http.Error(w, fmt.Sprintf("Conflict: %v", err), http.StatusConflict)
			return
		default:
			http.Error(w, fmt.Sprintf("Internal error: %v", err), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newSub); err != nil {
		http.Error(w, "Failed do encode user", http.StatusInternalServerError)
		return
	}

}

// UpdateBySID - Обновление данных существующей подписки
// @Summary      Обновление подписки по ее SID
// @Description  Обновляет подписку по ее SID из URL, новые данные берутся из тела запроса
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        sid   path      int  true  "SID подписки" example(20)
// @Param        subscription  body      model.RawSubscription  true  "Subscription info" example(`{"price": 400,"user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba","end_date": "07-2025"}`)
// @Success      200  {object}  model.RawSubscription	"Subscription updated successfully"
// @Failure      404  {string}  string  "Subscription not found"
// @Failure      400  {string}  string  "Bad request"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /subscriptions/{sid}	[put]
func (SH *SubscriptionHandler) UpdateBySID(w http.ResponseWriter, r *http.Request) {
	var newSub model.RawSubscription
	sidStr := chi.URLParam(r, "sid")

	if err := json.NewDecoder(r.Body).Decode(&newSub); err != nil {
		http.Error(w, "Failed to decode subscription from json", http.StatusBadRequest)
		return
	}

	if err := SH.Service.UpdateBySID(r.Context(), &newSub, sidStr); err != nil {
		switch {
		case errors.Is(err, repository.ErrEmptyAllFields):
			http.Error(w, "At least one field must not be empty", http.StatusBadRequest)
		case errors.Is(err, repository.ErrSubNotFound):
			http.Error(w, "Subscription not found", http.StatusNotFound)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(newSub); err != nil {
		http.Error(w, "Failed to encode subscription", http.StatusInternalServerError)
		return
	}

}

// GetBySID - хендлер для получения подписки по ее SID
// @Summary      Получение подписки по SID
// @Description  Возвращает подписку в формате JSON по ее SID из URL
// @Tags         subscriptions
// @Produce      json
// @Param        sid   path      int  true  "SID подписки" example(20)
// @Success      200  {object}  model.RawSubscription
// @Failure      404  {string}  string  "Subscription not found"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /subscriptions/{sid} [get]
func (SH *SubscriptionHandler) GetBySID(w http.ResponseWriter, r *http.Request) {
	sidStr := chi.URLParam(r, "sid")
	sid, err := strconv.ParseInt(sidStr, 10, 64)
	if err != nil {
		http.Error(w, "Failed to parse subscription SID", http.StatusInternalServerError)
		return
	}

	subscription, err := SH.Service.GetBySID(r.Context(), uint64(sid))
	if err != nil {
		http.Error(w, "Failed to find subscription SID", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(subscription); err != nil {
		http.Error(w, "Failed to encode subscription", http.StatusInternalServerError)
		return
	}

}

// GetList - хендлер для получения списка всех подписок из базы
// @Summary      Получение списка всех подписок из базы
// @Description  Отдает массив из всех подписок в базе; пустой json если подписок нет
// @Tags         subscriptions
// @Produce      json
// @Success      200   {array}  model.RawSubscription
// @Failure      500   {string}  string  "Internal server error"
// @Router       /subscriptions [get]
func (SH *SubscriptionHandler) GetList(w http.ResponseWriter, r *http.Request) {
	subscriptions, err := SH.Service.GetList(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch subscriptions", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(subscriptions); err != nil {
		http.Error(w, "Failed to encode users", http.StatusInternalServerError)
	}

}

// Delete - хендлер для удаления подписки по SID
// @Summary      Удаление подписки по SID
// @Description  Удаляет подписки по ее SID из URL
// @Tags         subscriptions
// @Param        sid   path      int  true  "SID подписки" example(20)
// @Success      204  {string}  string  "No Content"
// @Failure      404  {string}  string  "Subscription not found"
// @Failure      400  {string}  string  "Bad request"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /subscriptions/{sid}	[delete]
func (SH *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	sidStr := chi.URLParam(r, "sid")
	sid, err := strconv.ParseInt(sidStr, 10, 64)
	if err != nil {
		http.Error(w, "Failed to parse subscription SID", http.StatusInternalServerError)
		return
	}
	err = SH.Service.DeleteSubscription(r.Context(), uint(sid))
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrSubNotFound):
			http.Error(w, "Subscription not found", http.StatusNotFound)
		default:
			http.Error(w, fmt.Sprintf("Failed to delete subscription: %v", err), http.StatusBadRequest)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Report - хендлер для формирования отчета по подпискам; результат - сумма стоимости подписок за период(конкретный месяц в формате "07-2024") с фильтрацией
// @Summary      Подсчет суммы подписок удовлетворяющим условиям
// @Description  Выдает сумму стоимости подписок по указанному периоду(конкретному месяцу в формате "07-2024"), пользователю и провайдеру; пользователь и провайдер не являются обязательными полями.
// @Tags         subscriptions
// @Param        period     query      string  true  "Период(конкретный месяц в формате 07-2024) для поиска подписок в активном статусе" example(07-2025)
// @Param        uid        query      string  false "UID пользователя" example(adjhdjfnv-njdfv889)
// @Param        provider   query      string  false "Имя провайдера услуги" example(Yandex)
// @Success      200  {object}  model.Report  "Status OK"
// @Failure      400  {string}  string  "Bad request"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /subscriptions/report	[get]
func (SH *SubscriptionHandler) Report(w http.ResponseWriter, r *http.Request) {
	var result model.Report
	var filter model.RawReportFilter
	filter.Period = r.URL.Query().Get("period")
	filter.UID = r.URL.Query().Get("uid")
	filter.Provider = r.URL.Query().Get("provider")

	if filter.Period == "" {
		http.Error(w, "Empty mandatory period field", http.StatusBadRequest)
		return
	}

	var err error
	if result.Total, err = SH.Service.Report(r.Context(), &filter); err != nil {
		switch {
		case errors.Is(err, utils.ErrConvertToNorm):
			http.Error(w, "Incorrect input data", http.StatusBadRequest)
		default:
			http.Error(w, fmt.Sprintf("Failed to compose report: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode result", http.StatusInternalServerError)
		return
	}
}
