package main

import (
	"fmt"
	"log"
	"net/http"

	"HospitalAppointment_booking/config"
	"HospitalAppointment_booking/handlers"
	"HospitalAppointment_booking/middleware"
	"HospitalAppointment_booking/models"
)

func main() {
	config.ConnectDB()
	config.DB.AutoMigrate(&models.Booking{})

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/booking", handlers.CreateBooking)
	mux.HandleFunc("/api/v1/available-slots", handlers.GetAvailableSlots)
	mux.HandleFunc("/api/v1/booking/history", handlers.GetMyBookings)
	mux.HandleFunc("/api/v1/webhook/payment", handlers.PaymentWebhook)
	// Wrap mux with middlewares from the middleware package
	finalHandler := middleware.RecoveryMiddleware(
		middleware.LoggingMiddleware(
			middleware.JSONContentMiddleware(mux),
		),
	)

	fmt.Println("Server running at :8080...")
	log.Fatal(http.ListenAndServe(":8080", finalHandler))
}
