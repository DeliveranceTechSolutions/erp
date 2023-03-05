// Package testgrp holds all of the test handlers.
package testgrp

import (
	"context"
	// "errors"
	// "math/rand"
	"net/http"

	// "github.com/deliveranceTechSolutions/erp/business/sys/validate"
	"github.com/deliveranceTechSolutions/erp/foundation/web"
	"go.uber.org/zap"
)

// Handlers manages the set of check endpoints
type Handlers struct {
	Build string
	Log *zap.SugaredLogger
}

// Test handler is for development
func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// testing error functionality
	// if n := rand.Intn(100); n%2 == 0 {
	// 	return validate.NewRequestError(errors.New("trusted error"), http.StatusBadRequest)
	// 	// panic("testing panic")
	// }
	
	status := struct {
		Status string
	}{
		Status: "OK",
	}

	statusCode := http.StatusOK
	h.Log.Infow("readiness", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)

	return web.Respond(ctx, w, status, http.StatusOK)
}