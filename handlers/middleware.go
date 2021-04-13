package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Kurztrip/driver/data"
)

func (h *LocationHandler) MiddlewareLocationValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		loc := &data.Location{}

		err := loc.FromJson(r.Body)
		if err != nil {
			h.l.Println("[ERROR] deserializing location", err)
			http.Error(rw, "Error reading location", http.StatusBadRequest)
			return
		}

		//validate location
		err = loc.Validate()
		if err != nil {
			h.l.Println("[ERROR] validating location", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating location: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		//add location to context
		ctx := context.WithValue(r.Context(), KeyLocation{}, loc)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}
