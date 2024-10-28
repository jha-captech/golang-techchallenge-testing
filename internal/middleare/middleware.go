package middleare

import (
	"net/http"
)

// Middleware is a function that wraps an http.Handler.
type Middleware func(next http.Handler) http.Handler
