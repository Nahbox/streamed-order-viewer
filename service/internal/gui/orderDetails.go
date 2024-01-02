package gui

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Nahbox/streamed-order-viewer/service/internal/models"
)

func OrderDetailsPage(data models.OrderDetails) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.order.OrderDetailsPage"

		lp := filepath.Join("internal", "html", "orderDetails.html")
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

		log.Println("Template order.html executed successful!")
	}
}
