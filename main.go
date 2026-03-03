package main

import (
	"fmt"
	"log"
	"net/http"

	"NSPC_HEALTHCONNECT/appointment"
	"NSPC_HEALTHCONNECT/appointment/handlers"
	"NSPC_HEALTHCONNECT/database"
)

func main() {
	database.ConnectDB()
	database.DB.AutoMigrate(&appointment.Booking{})

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/booking", handlers.CreateBooking)
	mux.HandleFunc("/api/v1/available-slots", handlers.GetAvailableSlots)
	mux.HandleFunc("/api/v1/booking/history", handlers.GetMyBookings)
	mux.HandleFunc("/api/v1/webhook/payment", handlers.PaymentWebhook)
	// Wrap mux with middlewares from the middleware package
	finalHandler := database.RecoveryMiddleware(
		database.LoggingMiddleware(
			database.JSONContentMiddleware(mux),
		),
	)

	fmt.Println("Server running at :8080...")
	log.Fatal(http.ListenAndServe(":8080", finalHandler))
}
