package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Kurztrip/driver/data"
)

func (h *DriverHandler) MiddlewareDriverValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		dvr := &data.Driver{}

		err := dvr.FromJson(r.Body)
		if err != nil {
			h.l.Println("[ERROR] deserializing driver", err)
			http.Error(rw, "Error reading driver", http.StatusBadRequest)
			return
		}

		//validate driver
		err = dvr.Validate()
		if err != nil {
			h.l.Println("[ERROR] validating driver", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating driver: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		//add driver to context
		ctx := context.WithValue(r.Context(), KeyDriver{}, dvr)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}

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
