package order

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/Nahbox/streamed-order-viewer/service/internal/db"
	"github.com/Nahbox/streamed-order-viewer/service/internal/handler"
	"github.com/Nahbox/streamed-order-viewer/service/internal/models"
)

const orderIDKey = "orderID"

func (h *Handler) OrderContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := chi.URLParam(r, "orderId")
		if userId == "" {
			render.Render(w, r, handler.ErrorRenderer(fmt.Errorf("user ID is required")))
			return
		}
		id, err := strconv.Atoi(userId)
		if err != nil {
			render.Render(w, r, handler.ErrorRenderer(fmt.Errorf("invalid user ID")))
		}
		ctx := context.WithValue(r.Context(), orderIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// TODO: убрать структуру Out и создать в models
func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(orderIDKey).(int)
	var err error
	orderDetails, orderIsFound := h.CacheManager.Get(strconv.Itoa(userID))

	if !orderIsFound {
		orderDetails, err = h.DB.GetOrderById(userID)
		if err == nil {
			h.CacheManager.Set(strconv.Itoa(orderDetails.(models.OrderDetails).Order.Id), orderDetails.(models.OrderDetails), -1)
		}
	}

	if err != nil {
		if errors.Is(err, db.ErrNoMatch) {
			render.Render(w, r, handler.ErrNotFound)
		} else {
			render.Render(w, r, handler.ErrorRenderer(err))
		}
		return
	}

	if err != nil {
		render.Render(w, r, handler.ServerErrorRenderer(err))
		return
	}
	render.JSON(w, r, orderDetails)
}
