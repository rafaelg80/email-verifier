package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strings"

    "github.com/AfterShip/email-verifier"
)

// Predefined API key for simplicity, in practice store securely
var apiKey = "ENan5m9w22afdXH85u"

func main() {
    verifier := emailverifier.NewVerifier().
        EnableSMTPCheck().
        FromEmail("barnaby@prodigylead.com")

    http.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
        // Check for API key in request
        key := r.URL.Query().Get("apikey")
        if key == "" || key != apiKey {
            http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
            return
        }

        email := r.URL.Query().Get("email")
        if email == "" {
            http.Error(w, `{"error": "Email parameter is missing"}`, http.StatusBadRequest)
            return
        }

        result, err := verifier.Verify(email)
        if err != nil {
            log.Printf("Error verifying email: %v", err)
            handleVerificationError(w, err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        jsonResult, err := json.Marshal(result)
        if err != nil {
            log.Printf("Error marshaling result: %v", err)
            http.Error(w, fmt.Sprintf(`{"error": "Error marshaling result: %v"}`, err), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Write(jsonResult)
    })

    fmt.Println("Server is running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleVerificationError(w http.ResponseWriter, err error) {
    errMsg := err.Error()
    var response string

    switch {
    case strings.Contains(errMsg, "no such host"):
        response = `{"error": "Reachable: No Host"}`
    case strings.Contains(errMsg, "Try again later"):
        response = `{"error": "Reachable: Temporary Issue"}`
    case strings.Contains(errMsg, "too many errors"):
        response = `{"error": "Reachable: Temporary Issue, Too Many Errors"}`
    default:
        response = fmt.Sprintf(`{"error": "Error verifying email: %v"}`, err)
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(response))
}
