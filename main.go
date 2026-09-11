package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func main() {
	c, err := NewInfrai()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/login/otp", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", 405)
			return
		}
		var in struct {
			Phone string `json:"phone"`
		}
		if json.NewDecoder(r.Body).Decode(&in) != nil || in.Phone == "" {
			http.Error(w, "phone required", 400)
			return
		}
		_, err := c.OTP(in.Phone, "otp-"+strings.NewReplacer("+", "", " ", "").Replace(in.Phone))
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	})
	http.HandleFunc("/login/verify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", 405)
			return
		}
		var in struct {
			Phone   string `json:"phone"`
			Code    string `json:"code"`
			EventID string `json:"eventID"`
		}
		_ = json.NewDecoder(r.Body).Decode(&in)
		e, err := verifyPayment(c, in.Phone, in.Code, in.EventID)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			w.WriteHeader(400)
		}
		_ = json.NewEncoder(w).Encode(e)
	})
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
