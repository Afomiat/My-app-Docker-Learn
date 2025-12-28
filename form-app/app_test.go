package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Form App API", func() {
	var router = SetupRouter()

	Describe("GET /health", func() {
		It("should return a 200 status", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/health", nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(ContainSubstring("healthy"))
		})
	})

	Describe("POST /update-profile", func() {
		Context("with valid JSON data", func() {
			It("should upsert the user into MongoDB", func() {
				userData := User{
					Name:      "CI Test User",
					Email:     "ci@test.com",
					Interests: "Automation",
				}
				jsonValue, _ := json.Marshal(userData)
				
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/update-profile", bytes.NewBuffer(jsonValue))
				req.Header.Set("Content-Type", "application/json")
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Body.String()).To(ContainSubstring("CI Test User"))
			})
		})
	})
})