package middlewares

import (
	"net/http"
	"strings"

	"github.com/rs/cors"
)

func NewCorsMiddleware(options ...cors.Options) *cors.Cors {
	if len(options) == 0 {
		return cors.Default()
	}

	userOptions := options[0]

	defaultOptions := cors.Options{
		AllowOriginFunc: func(origin string) bool {
			// Check wildcard patterns
			for _, allowedOrigin := range userOptions.AllowedOrigins {
				if matchOrigin(allowedOrigin, origin) {
					return true
				}
			}
			return false
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders:             []string{"*"},
		AllowCredentials:           true,
		Debug:                      userOptions.Debug,
		AllowOriginVaryRequestFunc: userOptions.AllowOriginVaryRequestFunc,
		ExposedHeaders:             userOptions.ExposedHeaders,
		MaxAge:                     userOptions.MaxAge,
		OptionsPassthrough:         userOptions.OptionsPassthrough,
		AllowPrivateNetwork:        userOptions.AllowPrivateNetwork,
		OptionsSuccessStatus:       userOptions.OptionsSuccessStatus,
		Logger:                     userOptions.Logger,
	}

	if userOptions.AllowOriginFunc != nil {
		defaultOptions.AllowOriginFunc = userOptions.AllowOriginFunc
	}

	if len(userOptions.AllowedMethods) > 0 {
		defaultOptions.AllowedMethods = userOptions.AllowedMethods
	}

	if len(userOptions.AllowedHeaders) > 0 {
		defaultOptions.AllowedHeaders = userOptions.AllowedHeaders
	}

	if userOptions.AllowCredentials != defaultOptions.AllowCredentials {
		defaultOptions.AllowCredentials = userOptions.AllowCredentials
	}

	return cors.New(defaultOptions)
}

func matchOrigin(pattern, origin string) bool {
	if pattern == origin {
		return true
	}
	if strings.Contains(pattern, "*") {
		// Handle wildcard matching
		// Split both pattern and origin by "://" to separate protocol
		patternParts := strings.SplitN(pattern, "://", 2)
		originParts := strings.SplitN(origin, "://", 2)

		if len(patternParts) != 2 || len(originParts) != 2 {
			return false
		}

		// Protocol must match exactly
		if patternParts[0] != originParts[0] {
			return false
		}

		// Check host part (remove port from origin if present)
		patternHost := patternParts[1]
		originHost := originParts[1]

		// Remove port from origin host
		if portIndex := strings.LastIndex(originHost, ":"); portIndex != -1 {
			originHost = originHost[:portIndex]
		}

		// Simple wildcard matching
		if strings.HasPrefix(patternHost, "*.") {
			suffix := patternHost[2:] // Remove "*."
			if strings.HasSuffix(originHost, "."+suffix) || originHost == suffix {
				return true
			}
		}
	}

	return false
}
