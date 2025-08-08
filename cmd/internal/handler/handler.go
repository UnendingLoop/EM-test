package handler

import (
	"em-test/cmd/internal/model"
	"em-test/cmd/internal/repository"
	"em-test/cmd/internal/service"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// SubscriptionHandler provides process to HTTP-requests
type SubscriptionHandler struct {
	Service service.SubscriptionService
}

// Create - хендлер для создания новой подписки в базе
// @Summary      Хендлер для создания новой подписки
// @Description  Создаёт новую подписку из данных в теле запроса
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        subscription  body      model.Subscription  true  "Subscription info"
// @Success      201   {object}  model.Subscription
// @Failure      400   {string}  string  "Incomplete/incorrect data input"
// @Failure      500   {string}  string  "Internal server error"
// @Router       /subscription/create [post]
func (SH *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var newSub model.RawSubscription

	if err := json.NewDecoder(r.Body).Decode(&newSub); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := SH.Service.Create(r.Context(), &newSub); err != nil {
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
	if err := json.NewEncoder(w).Encode(newSub); err != nil {
		http.Error(w, "Failed do encode user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

}

// UpdateBySID - хендлер для обновления данных подписки
// @Summary      Обновление подписки по ее SID
// @Description  Обновляет подписку по ее SID из URL, новые данные берутся из тела запроса
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        sid   path      uint  true  "SID подписки"
// @Success      200  {object}  model.Subscription	"Subscription updated successfully"
// @Failure      404  {string}  string  "Subscription not found"
// @Failure      400  {string}  string  "Bad request"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /subscription/update/{sid}	[put]
func (SH *SubscriptionHandler) UpdateBySID(w http.ResponseWriter, r *http.Request) {
	var newSub model.RawSubscription
	sidStr := chi.URLParam(r, "sid")

	if err := json.NewDecoder(r.Body).Decode(&newSub); err != nil {
		http.Error(w, "Failed to decode subscription from json", http.StatusBadRequest)
		return
	}

	if err := SH.Service.UpdateBySID(r.Context(), &newSub, sidStr); err != nil {
		switch {
		case errors.Is(err, repository.ErrEmptyFields):
			http.Error(w, "At least one field must not be empty", http.StatusBadRequest)
		case errors.Is(err, repository.ErrSubNotFound):
			http.Error(w, "Subscription not found", http.StatusNotFound)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newSub); err != nil {
		http.Error(w, "Failed to encode subscription", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetBySID - хендлер для получения подписки по ее SID
// @Summary      Получение подписки по SID
// @Description  Возвращает подписку в формате JSON по ее SID из URL
// @Tags         subscriptions
// @Produce      json
// @Param        sid   path      uint  true  "SID подписки"
// @Success      200  {object}  model.Subscription
// @Failure      404  {string}  string  "Subscription not found"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /subscription/{sid} [get]
func (SH *SubscriptionHandler) GetBySID(w http.ResponseWriter, r *http.Request) {
	sidStr := chi.URLParam(r, "sid")
	sid, err := strconv.ParseInt(sidStr, 10, 64)
	if err != nil {
		http.Error(w, "Failed to parse subscription SID", http.StatusInternalServerError)
		return
	}

	subscription, err := SH.Service.GetBySID(r.Context(), uint(sid))
	if err != nil {
		http.Error(w, "Failed to find subscription SID", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(subscription); err != nil {
		http.Error(w, "Failed to encode subscription", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetList - хендлер для получения списка всех подписок из базы
// @Summary      Хендлер для получения списка всех подписок из базы
// @Description  Отдает массив из всех пользователей базы
// @Tags         subscriptions
// @Produce      json
// @Success      200   {array}  model.Subscription
// @Failure      500   {string}  string  "Internal server error"
// @Router       /subscription/list [get]
func (SH *SubscriptionHandler) GetList(w http.ResponseWriter, r *http.Request) {
	subscriptions, err := SH.Service.GetList(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch subscriptions", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(subscriptions); err != nil {
		http.Error(w, "Failed to encode users", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

// Delete - хендлер для удаления подписки по SID
// @Summary      Удаление подписки по SID
// @Description  Удаляет подписки по ее SID из URL
// @Tags         subscriptions
// @Param        sid   path      uint  true  "SID подписки"
// @Success      204  {string}  string  "No Content"
// @Failure      404  {string}  string  "Subscription not found"
// @Failure      400  {string}  string  "Bad request"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /subscription/{sid}	[delete]
func (SH *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	sidStr := chi.URLParam(r, "sid")
	sid, err := strconv.ParseInt(sidStr, 10, 64)
	if err != nil {
		http.Error(w, "Failed to parse subscription SID", http.StatusInternalServerError)
		return
	}
	err = SH.Service.Delete(r.Context(), uint(sid))
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
// @Param        period     path      string  true  "Период(конкретный месяц в формате "07-2024") для поиска подписок в активном статусе"
// @Param        uid        path      string  false "UID пользователя"
// @Param        provider   path      string  false "Имя провайдера услуги"
// @Success      200  {string}  string  "Status OK"
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
	}

	var err error
	if result.Total, err = SH.Service.Report(r.Context(), &filter); err != nil {
		switch {
		case errors.Is(err, service.ErrConvertToNorm):
			http.Error(w, "Incorrect input data", http.StatusBadRequest)
		default:
			http.Error(w, fmt.Sprintf("Failed to compose report: %v", err), http.StatusInternalServerError)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode result", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
