package services

import (
	"HospitalAppointment_booking/models"
	"fmt"
	"os"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

func GenerateBookingPDF(b models.Booking) (string, error) {
	// Folder lekapothe create chesthundhi
	if _, err := os.Stat("receipts"); os.IsNotExist(err) {
		os.Mkdir("receipts", 0755)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(0, 10, "OFFICIAL APPOINTMENT RECEIPT")
	pdf.Ln(15)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, fmt.Sprintf("Appointment ID: %s", b.AppointmentID))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Patient: %s (%d Years, %s)", b.PatientName, b.PatientAge, b.PatientGender))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Date: %s", b.AppointmentDate))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Time Slot: %s", b.TimeSlot))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Hospital/Doctor: %s / %s", b.HospitalID, b.DoctorID))
	pdf.Ln(12)
	pdf.Cell(0, 10, fmt.Sprintf("Payment Method: %s", b.PaymentMethod))
	pdf.Ln(8)
	actualAmount := float64(b.Amount) // Convert paise to rupees
	if strings.ToLower(b.PaymentMethod) == "online" {
		pdf.SetTextColor(0, 128, 0) // Green color for Paid
		pdf.Cell(0, 10, fmt.Sprintf("Amount Paid Online: Rs. %.2f", actualAmount))
		pdf.Ln(8)
		pdf.Cell(0, 10, "Status: SUCCESS / PAID")
	} else {
		pdf.SetTextColor(255, 0, 0) // Red color for Unpaid
		pdf.Cell(0, 10, fmt.Sprintf("Amount to be Paid: Rs. %.2f", actualAmount))
		pdf.Ln(8)
		pdf.Cell(0, 10, "Status: UNPAID (Please pay at hospital counter)")
	}

	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(20)
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 10, "Note: Please carry this PDF or Mobile SMS at the time of visit.")

	filePath := fmt.Sprintf("receipts/%s.pdf", b.AppointmentID)
	err := pdf.OutputFileAndClose(filePath)
	return filePath, err
}
