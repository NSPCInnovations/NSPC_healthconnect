package handlers

import (
	"HospitalAppointment_booking/config"
	"HospitalAppointment_booking/models"
	"HospitalAppointment_booking/services"
	"encoding/json"
	"fmt"
	"net/http"
)

func PaymentWebhook(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		fmt.Println("Webhook Body Decode Failed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	event, _ := body["event"].(string)
	fmt.Printf("\n--- WEBHOOK RECEIVED: %s ---\n", event)

	var razorpayOrderID string
	payload, _ := body["payload"].(map[string]interface{})

	// --- FUTURE-PROOF EXTRACTION LOGIC (FALLBACKS) ---
	if event == "payment_link.paid" {
		// Path 1: Payload -> Order -> Entity -> Notes (Your current successful path)
		if order, ok := payload["order"].(map[string]interface{}); ok {
			if entity, ok := order["entity"].(map[string]interface{}); ok {
				if notes, ok := entity["notes"].(map[string]interface{}); ok {
					razorpayOrderID = fmt.Sprintf("%v", notes["order_id"])
				}
			}
		}

		// Path 2: Fallback to direct Payment Link Entity (In case Path 1 fails)
		if razorpayOrderID == "" || razorpayOrderID == "<nil>" {
			if pLink, ok := payload["payment_link"].(map[string]interface{}); ok {
				if entity, ok := pLink["entity"].(map[string]interface{}); ok {
					// Direct search in entity's order_id field
					if entity["order_id"] != nil {
						razorpayOrderID = fmt.Sprintf("%v", entity["order_id"])
					} else if notes, ok := entity["notes"].(map[string]interface{}); ok {
						// Search in entity's notes field
						razorpayOrderID = fmt.Sprintf("%v", notes["order_id"])
					}
				}
			}
		}

		// Path 3: Fallback to Payment Entity (Just in case Razorpay changes behavior)
		if razorpayOrderID == "" || razorpayOrderID == "<nil>" {
			if payment, ok := payload["payment"].(map[string]interface{}); ok {
				if entity, ok := payment["entity"].(map[string]interface{}); ok {
					if notes, ok := entity["notes"].(map[string]interface{}); ok {
						razorpayOrderID = fmt.Sprintf("%v", notes["order_id"])
					} else if entity["order_id"] != nil {
						razorpayOrderID = fmt.Sprintf("%v", entity["order_id"])
					}
				}
			}
		}

	}

	// --- LOGGING & VALIDATION ---
	if razorpayOrderID == "" || razorpayOrderID == "<nil>" {
		fmt.Println("CRITICAL ERROR: Could not extract OrderID from any path!")
		// Full payload print chestundi so nuvvu future lo ventane structure chudochu
		fmt.Printf("FAILED PAYLOAD: %+v\n", body)
		w.WriteHeader(http.StatusOK) // Return 200 to stop Razorpay retries
		return
	}

	fmt.Printf("Final Search ID: %s\n", razorpayOrderID)

	// --- DB UPDATE SECTION ---
	var booking models.Booking
	result := config.DB.Where("order_id = ?", razorpayOrderID).First(&booking)

	if result.Error != nil {
		fmt.Printf(" DB Error: No record found for ID %s. Check if DB was updated on creation.\n", razorpayOrderID)
	} else {
		// Prevent double confirmation
		if booking.Status == "confirmed" {
			fmt.Println("ℹBooking already confirmed, skipping...")
			w.WriteHeader(http.StatusOK)
			return
		}

		fmt.Println("Record Found! Attempting to Update...")

		// 1. Amount Conversion (50000 -> 500)
		newAmount := booking.Amount / 100

		// 2. Database Status Update
		err := config.DB.Model(&booking).Updates(map[string]interface{}{
			"status": "confirmed",
			"amount": newAmount,
		}).Error

		if err != nil {
			fmt.Printf(" DB Update Failed: %v\n", err)
		} else {
			fmt.Println(" DB Updated Successfully! Preparing Email...")

			// Update local object for PDF content
			booking.Status = "confirmed"
			booking.Amount = newAmount

			// 3. Email and PDF Generation
			pdfPath, pdfErr := services.GenerateBookingPDF(booking)
			if pdfErr != nil {
				fmt.Printf("PDF Generation Error: %v\n", pdfErr)
			} else {
				go services.SendEmailWithPDF(booking.PatientEmail, pdfPath, booking.PatientName)
				fmt.Printf("Email Dispatched to: %s\n", booking.PatientEmail)
			}

			fmt.Printf("SUCCESS: Booking %s Confirmed!\n", booking.AppointmentID)
		}
	}
	w.WriteHeader(http.StatusOK)
}
