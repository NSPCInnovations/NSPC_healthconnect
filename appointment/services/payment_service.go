package services

import (
	"github.com/razorpay/razorpay-go"
)

const (
	RazorpayKeyID     = "rzp_test_SBLCJs3P9k6l1b"
	RazorpayKeySecret = "x65V2Uj16mkthxr8tJNZXajo"
)

func CreateRazorpayOrder(amountInPaise int64, receiptID string) (string, error) {
	client := razorpay.NewClient(RazorpayKeyID, RazorpayKeySecret)

	data := map[string]interface{}{
		"amount":   amountInPaise,
		"currency": "INR",
		"receipt":  receiptID,
		"notes": map[string]string{
			"appointment_id": receiptID, // this is useful for webhook identification
		},
	}

	body, err := client.Order.Create(data, nil)
	if err != nil {
		return "", err
	}

	return body["id"].(string), nil
}
func CreatePaymentLink(orderID string, amount float64, patientName, mobile string) (string, error) {
	// Ikkada nee existing keys vaadu (already paina define chesi unte avi vaadu)
	client := razorpay.NewClient(RazorpayKeyID, RazorpayKeySecret)

	data := map[string]interface{}{
		"amount":         int(amount * 100),
		"currency":       "INR",
		"accept_partial": false,
		"description":    "Hospital Appointment - " + patientName,
		"customer": map[string]interface{}{
			"name":    patientName,
			"contact": mobile,
			"email":   "ameermahammad40@gmail.com",
		},
		"notify": map[string]interface{}{
			"sms":   true,
			"email": true,
		},
		"reminder_enable": true,
		"notes": map[string]interface{}{
			"order_id": orderID, // webhook is using this to identify the booking, so use appointment_id or order_id as per your webhook logic
		},
	}

	// Razorpay SDK call
	linkBody, err := client.PaymentLink.Create(data, nil)
	if err != nil {
		return "", err
	}

	return linkBody["short_url"].(string), nil
}
