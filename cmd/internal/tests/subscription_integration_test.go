package tests_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"em-test/cmd/internal/handler"
	"em-test/cmd/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/go-chi/chi/v5"
)

// SetupTestDB создает в памяти SQLite базу и мигрирует модель
func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test DB: %v", err)
	}
	if err := db.AutoMigrate(&model.Subscription{}); err != nil {
		t.Fatalf("failed to migrate test DB: %v", err)
	}
	return db
}

func TestSubscriptionAPI(t *testing.T) {
	db := SetupTestDB(t)
	handler := handler.CreateHandler(db)

	var price uint
	// 1. Создать подписку
	price = 400

	newSub := model.RawSubscription{
		Provider: "Yandex Plus",
		Price:    &price,
		UID:      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
		Start:    "07-2025",
	}
	bodyBytes, _ := json.Marshal(newSub)
	createReq := httptest.NewRequest(http.MethodPost, "/subscriptions", bytes.NewReader(bodyBytes))
	createRec := httptest.NewRecorder()
	handler.Create(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("Create subscription: expected status 201, got %d", createRec.Code)
	}

	var createdSub model.RawSubscription
	if err := json.Unmarshal(createRec.Body.Bytes(), &createdSub); err != nil {
		t.Fatalf("Create subscription: failed to parse response: %v", err)
	}

	// 2. Получить по SID
	getReq := httptest.NewRequest(http.MethodGet, "/subscriptions/"+strconv.FormatUint(*createdSub.SID, 10), nil)
	getRec := httptest.NewRecorder()
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("sid", strconv.FormatUint(*createdSub.SID, 10))
	getReq = getReq.WithContext(context.WithValue(getReq.Context(), chi.RouteCtxKey, chiCtx))

	handler.GetBySID(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("GetBySID: expected status 200, got %d", getRec.Code)
	}

	// 3. Получить список всех подписок
	listReq := httptest.NewRequest(http.MethodGet, "/subscriptions", nil)
	listRec := httptest.NewRecorder()
	handler.GetList(listRec, listReq)

	if listRec.Code != http.StatusOK {
		t.Fatalf("GetList: expected status 200, got %d", listRec.Code)
	}

	// 5. Удалить подписку по SID
	deleteReq := httptest.NewRequest(http.MethodDelete, "/subscriptions/"+strconv.FormatUint(*createdSub.SID, 10), nil)
	deleteRec := httptest.NewRecorder()
	chiCtxDelete := chi.NewRouteContext()
	chiCtxDelete.URLParams.Add("sid", strconv.FormatUint(*createdSub.SID, 10))
	deleteReq = deleteReq.WithContext(context.WithValue(deleteReq.Context(), chi.RouteCtxKey, chiCtxDelete))

	handler.Delete(deleteRec, deleteReq)

	if deleteRec.Code != http.StatusNoContent {
		t.Fatalf("Delete: expected status 204, got %d", deleteRec.Code)
	}

	// 6. Тест отчета: создадим пару подписок и запросим сумму за период
	// Создадим первую подписку (цена 300)
	db.Create(&model.Subscription{
		Provider: "Netflix",
		Price:    300,
		UID:      "user1",
		Start:    *mustParseDate("06-2025"),
		End:      nil,
	})
	// Вторая подписка (цена 200), с датой окончания в июле 2025
	End := mustParseDate("07-2025")
	db.Create(&model.Subscription{
		Provider: "Spotify",
		Price:    200,
		UID:      "user1",
		Start:    *mustParseDate("05-2025"),
		End:      End,
	})

	reportReq := httptest.NewRequest(http.MethodGet, "/subscriptions/report?period=07-2025&uid=user1&provider=Netflix", nil)
	reportRec := httptest.NewRecorder()
	handler.Report(reportRec, reportReq)

	if reportRec.Code != http.StatusOK {
		t.Fatalf("Report: expected status 200, got %d", reportRec.Code)
	}

	var report model.Report
	if err := json.Unmarshal(reportRec.Body.Bytes(), &report); err != nil {
		t.Fatalf("Report: failed to parse response: %v", err)
	}

	if report.Total != 300 {
		t.Errorf("Report: expected total 300, got %d", report.Total)
	}
}

// Вспомогательная функция парсинга даты в нужном формате
func mustParseDate(s string) *time.Time {
	tm, err := time.Parse("01-2006", s)
	if err != nil {
		panic(err)
	}
	return &tm
}
