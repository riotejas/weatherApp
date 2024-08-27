package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

type ContextKey string

const (
	LatKey ContextKey = "latitude"
	LngKey ContextKey = "longitude"
)

// QueryParamsToContext adds all URL params from the request to the context
func QueryParamsToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		query := r.URL.Query()

		ctx = context.WithValue(ctx, LatKey, query.Get("latitude"))
		ctx = context.WithValue(ctx, LngKey, query.Get("longitude"))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type QueryParamRule struct {
	Required  bool
	Validator func(string) bool
}

type QueryParamRules map[ContextKey]QueryParamRule

// ValidateQueryParams rules to validate user input.  Requires parameter to exist and validates value.
func ValidateQueryParams(rules QueryParamRules) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query()
			errors := make(map[ContextKey]string)

			for param, rule := range rules {
				value := query.Get(string(param))

				if rule.Required && value == "" {
					errors[param] = "This parameter is required"
				} else if value != "" && rule.Validator != nil && !rule.Validator(value) {
					errors[param] = "Invalid value for this parameter"
				}
			}

			if len(errors) > 0 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"errors": errors,
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

var ForecastSearchRules = QueryParamRules{
	LatKey: {
		Required: true,
		Validator: func(v string) bool {
			_, err := strconv.ParseFloat(v, 64)
			return err == nil
		},
	},
	LngKey: {
		Required: true,
		Validator: func(v string) bool {
			_, err := strconv.ParseFloat(v, 64)
			return err == nil
		},
	},
}
