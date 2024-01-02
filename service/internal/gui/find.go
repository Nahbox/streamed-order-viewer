package gui

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/Nahbox/streamed-order-viewer/service/internal/handler/order"
	"github.com/Nahbox/streamed-order-viewer/service/internal/models"
)

type FindPageData struct {
	OrderID     string
	ShowMessage bool
	Message     string
}

// GET
func FindOrderByIDPage(message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.find.FindOrderByIDPage"

		orderID := r.URL.Query().Get("orderID")
		var showMessage bool
		if message != "" {
			showMessage = true
		}
		data := FindPageData{OrderID: orderID, ShowMessage: showMessage, Message: message}

		lp := filepath.Join("internal", "html", "find.html")
		tmpl, err := template.ParseFiles(lp)
		if err != nil {
			log.Printf("%s: %s\n", op, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			log.Printf("%s: %s\n", op, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		log.Println("Template find.html executed successful!")
	}
}

// POST
func FindOrderByID(orderHandler *order.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.find.FindOrderByID"

		orderID := r.URL.Query().Get("id")
		if orderID == "" {
			FindOrderByIDPage("Field must be filled!")(w, r) // Повторно отображаем страницу с предупреждением
			return
		}

		var parsed models.OrderDetails

		order, ok := orderHandler.CacheManager.Get(orderID)
		if ok {
			if parsed, ok = order.(models.OrderDetails); ok {
				log.Println("Get order from cache successful!")
			} else {
				log.Println("Failed to convert []byte")
				FindOrderByIDPage("Error getting data from cache!")(w, r)
			}
		} else {
			var err error
			intOrderId, err := strconv.Atoi(orderID)
			if err != nil {
				return
			}

			orderDetails, err := orderHandler.DB.GetOrderById(intOrderId)
			if err != nil {
				log.Printf("%s: %s\n", op, err)
				FindOrderByIDPage("Error getting data from database!")(w, r)
				return
			}

			parsed = *orderDetails
		}

		log.Println("Find order by ID page successful!")

		OrderDetailsPage(parsed)(w, r)
	}
}
