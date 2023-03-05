package mid

import (
	"context"
	"net/http"

	"github.com/deliveranceTechSolutions/erp/foundation/web"
	"go.uber.org/zap"
)

// Logger ...
func Logger(log *zap.SugaredLogger) web.Middleware {

	// These two closures are ensuring that we have a proper
	// ctx when making these api calls from a client.
	// This allows us to maintain the idea of testing in the core
	// while wrapping the core with code, funcs, and ultimately middleware
	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// If the context is missing this value, request the service
			// to be shutdown gracefully.
			v, err := web.GetValues(ctx)
			if err != nil {
				return web.NewShutdownError("web value missing from context")
			}

			// If we cannot extract the ctx properly
			// then we want to exit this middleware so
			// we can determine where in the call stack
			// an error occurred without needing a log parser
			log.Infow("request started", "trace_id", v.TraceID, "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr)
	
			err = handler(ctx, w, r)

			// Log regardless of error because
			// we exit directly after the handler call
			// if err != nil
			log.Infow("request completed", "trace_id", v.TraceID, "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr, "statuscode", v.StatusCode, "since", v.Now)

			return err
		}
		
		return h
	}

	return m
}