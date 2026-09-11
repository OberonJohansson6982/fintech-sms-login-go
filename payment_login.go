package main

import "fmt"

type PaymentEvent struct {
	ID, Phone, Action string
	Approved          bool
	Audit             string
}

func verifyPayment(c *Infrai, phone, code, eventID string) (PaymentEvent, error) {
	_, err := c.Verify(phone, code, "verify-"+eventID)
	e := PaymentEvent{ID: eventID, Phone: phone, Action: "login_payment", Audit: "otp verification"}
	if err != nil {
		return e, fmt.Errorf("verification failed: %w", err)
	}
	e.Approved = true
	return e, nil
}
