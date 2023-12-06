package looper

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/thalesfsp/httpclient"
	"github.com/thalesfsp/sypl"
)

// Name of the looper.
const Name = "requester"

// ErrHTTPVerbNotImplemented is the error returned when the http verb is not
// implemented by the underlying http client.
var ErrHTTPVerbNotImplemented = errors.New("http verb not implemented")

// Requester is the requester looper.
type Requester struct {
	ILooper
}

// Handle requests errors.
func handleError(metrics IMetrics, url string, err error) {
	// Increment failed counter.
	metrics.IncrementFailedCounter()

	// Log.
	metrics.GetLogger().Errorlnf("error while requesting %s: %s", url, err.Error())
}

// Handle requests success.
func handleSuccess(metrics IMetrics, url string, statusCode int) {
	// Increment success counter.
	metrics.IncrementSuccessCounter()

	// Log.
	metrics.GetLogger().Infolnf("request to %s returned %d", url, statusCode)
}

// Runs a GET request.
func getRequester(
	ctx context.Context,
	hc *httpclient.Client,
	url string,
	o ...httpclient.Func,
) (*http.Response, error) {
	return hc.Get(ctx, url, o...)
}

// Runs a POST request.
func postRequester(
	ctx context.Context,
	hc *httpclient.Client,
	url string,
	o ...httpclient.Func,
) (*http.Response, error) {
	return hc.Post(ctx, url, o...)
}

// Runs a PUT request.
func putRequester(
	ctx context.Context,
	hc *httpclient.Client,
	url string,
	o ...httpclient.Func,
) (*http.Response, error) {
	return hc.Put(ctx, url, o...)
}

// Runs a DELETE request.
func deleteRequester(
	ctx context.Context,
	hc *httpclient.Client,
	url string,
	o ...httpclient.Func,
) (*http.Response, error) {
	return hc.Delete(ctx, url, o...)
}

// Runs a PATCH request.
func patchRequester(
	ctx context.Context,
	hc *httpclient.Client,
	url string,
	o ...httpclient.Func,
) (*http.Response, error) {
	return hc.Patch(ctx, url, o...)
}

// Switcher determines which request to run.
func Switcher(
	ctx context.Context,
	logger sypl.ISypl,
	httpVerb string,
	url string,
	hc *httpclient.Client,
	o ...httpclient.Func,
) (*http.Response, error) {
	// switch case.
	switch httpVerb {
	case http.MethodGet:
		return getRequester(ctx, hc, url, o...)
	case http.MethodPost:
		return postRequester(ctx, hc, url, o...)
	case http.MethodPut:
		return putRequester(ctx, hc, url, o...)
	case http.MethodDelete:
		return deleteRequester(ctx, hc, url, o...)
	case http.MethodPatch:
		return patchRequester(ctx, hc, url, o...)
	default:
		logger.Fatalln(ErrHTTPVerbNotImplemented)

		return nil, ErrHTTPVerbNotImplemented
	}
}

// NewRequester creates a new requester.
func NewRequester(
	httpVerb string,
	url string,
	interval time.Duration,
	enableMetrics bool,
	loggingLevel string,
	metricsServerPort string,
	o ...httpclient.Func,
) (ILooper, error) {
	hc, err := httpclient.NewDefault(Name)
	if err != nil {
		return nil, err
	}

	l, err := New(
		Name,
		interval,
		enableMetrics,
		loggingLevel,
		metricsServerPort,
		func(metrics IMetrics) {
			resp, err := Switcher(context.Background(), metrics.GetLogger(), httpVerb, url, hc, o...)
			if err != nil {
				handleError(metrics, url, err)

				return
			}

			defer resp.Body.Close()

			handleSuccess(metrics, url, resp.StatusCode)
		},
	)
	if err != nil {
		return nil, err
	}

	return &Requester{
		ILooper: l,
	}, nil
}
