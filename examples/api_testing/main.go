package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/b3ndoi/factory-go/factory"
)

// UserRequest represents the API request payload
type UserRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password,omitempty"`
	PasswordConfirm string `json:"password_confirm,omitempty"`
	AcceptTerms     bool   `json:"accept_terms"`
}

// UserResponse represents the API response
type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Message  string `json:"message"`
}

func main() {
	fmt.Println("=== Factory-Go API Testing Example ===\n")

	// Create factory with API-specific fields
	requestFactory := factory.New(func(seq int64) UserRequest {
		return UserRequest{
			Username: fmt.Sprintf("user%d", seq),
			Email:    fmt.Sprintf("user%d@example.com", seq),
		}
	}).WithRawDefaults(func(r *UserRequest) {
		// These fields ONLY appear in Raw/RawJSON, not in Make/Create
		r.Password = "ValidPassword123!"
		r.PasswordConfirm = "ValidPassword123!"
		r.AcceptTerms = true
	})

	// Mock API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req UserRequest
		json.NewDecoder(r.Body).Decode(&req)

		// Validation
		if req.Email == "" || req.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "missing fields"})
			return
		}

		if req.Password != req.PasswordConfirm {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "passwords don't match"})
			return
		}

		if !req.AcceptTerms {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "must accept terms"})
			return
		}

		// Success
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(UserResponse{
			ID:       "user-123",
			Username: req.Username,
			Email:    req.Email,
			Message:  "User created successfully",
		})
	}))
	defer server.Close()

	// 1. Test valid registration
	fmt.Println("1. Valid registration:")
	validPayload := requestFactory.MustRawJSON()
	fmt.Printf("   Payload: %s\n", string(validPayload))

	resp, _ := http.Post(server.URL, "application/json", bytes.NewReader(validPayload))
	fmt.Printf("   Status: %d\n", resp.StatusCode)

	var successResp UserResponse
	json.NewDecoder(resp.Body).Decode(&successResp)
	fmt.Printf("   Response: %s\n", successResp.Message)

	// 2. Test invalid email
	fmt.Println("\n2. Invalid email:")
	invalidPayload := requestFactory.MustRawJSON(func(r *UserRequest) {
		r.Email = "" // Empty email
	})
	resp, _ = http.Post(server.URL, "application/json", bytes.NewReader(invalidPayload))
	fmt.Printf("   Status: %d (Expected 400 ✅)\n", resp.StatusCode)

	// 3. Test password mismatch
	fmt.Println("\n3. Password mismatch:")
	mismatchPayload := requestFactory.MustRawJSON(func(r *UserRequest) {
		r.PasswordConfirm = "DifferentPassword"
	})
	resp, _ = http.Post(server.URL, "application/json", bytes.NewReader(mismatchPayload))
	fmt.Printf("   Status: %d (Expected 400 ✅)\n", resp.StatusCode)

	// 4. Test bulk user creation
	fmt.Println("\n4. Bulk API testing:")
	bulkPayload := requestFactory.Count(5).MustRawJSON()
	fmt.Printf("   Created JSON array with %d users\n", 5)

	var requests []UserRequest
	json.Unmarshal(bulkPayload, &requests)
	for i, req := range requests {
		fmt.Printf("   [%d] %s (%s)\n", i+1, req.Username, req.Email)
	}

	// 5. Demonstrate WithRawDefaults separation
	fmt.Println("\n5. WithRawDefaults demonstration:")

	// Make() does NOT include password (for domain models)
	domainUser := requestFactory.Make()
	fmt.Printf("   Make() result: Password='%s' (empty as expected)\n", domainUser.Password)

	// Raw() DOES include password (for API testing)
	apiUser := requestFactory.Raw()
	fmt.Printf("   Raw() result: Password='%s' (has password!)\n", apiUser.Password)

	fmt.Println("\n✅ API testing example complete!")
	fmt.Println("\nKey Takeaway:")
	fmt.Println("  - WithRawDefaults() adds fields ONLY for Raw/RawJSON")
	fmt.Println("  - Perfect for passwords, tokens, or API-specific fields")
	fmt.Println("  - Keeps your domain models clean!")
}
