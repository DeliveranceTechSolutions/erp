package mid

import (
	"context"
	"net/http"

	"github.com/deliveranceTechSolutions/erp/business/sys/validate"
	"github.com/deliveranceTechSolutions/erp/foundation/web"
	"go.uber.org/zap"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Errors(log *zap.SugaredLogger) web.Middleware {

	// This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {

		// Create the handler that will be attacehd in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// If the context is missing this value, request the service
			// to be shutdown gracefully.
			v, err := web.GetValues(ctx)
			if err != nil {
				return web.NewShutdownError("web value missing from context")
			}


			// Run the next handler and catch any propagated error.
			if err := handler(ctx, w, r); err != nil {

				// Log the error.
				log.Errorw("ERROR", "traceid", v.TraceID, "ERROR", err)

				// switch statement for business errors is more efficient
				// than using As and Is, because we want the error to dictate
				// behavior not recognize the error and assign behavior
				// event based programming vs. classic imperative programming
				// Build out the error response.
				var er validate.ErrorResponse
				var status int
				switch act := validate.Cause(err).(type) {
				// FieldError attached to model/struct validation
				// Value semantics because it's a slice (built-in pointer semantics)
				case validate.FieldErrors:
					er = validate.ErrorResponse{
						Error:  "data validation error",
						Fields: act.Error(),
					}
					status = http.StatusBadRequest
				// RequestError attached to net/http validation
				// Pointer semantics because it's a string
				case *validate.RequestError:
					er = validate.ErrorResponse{
						Error: act.Error(),
					}
					status = act.Status
				// ErrorResponse attached to untrusted errors 
				// basically a 500 internal server error check
				default:
					er = validate.ErrorResponse{
						Error: http.StatusText(http.StatusInternalServerError),
					}
					status = http.StatusInternalServerError
				}

				// Respond with the error back to the client.
				if err := web.Respond(ctx, w, er, status); err != nil {
					return err
				}

				// If we receive the shutdown err we need to retun it
				// back to teh base handler to shutdown the service.
				if ok := web.IsShutdown(err); ok {
					return err
				}
			}

			// The error has been handled so this stops propogating.
			return nil
		}
		return h
	}
	return m
}