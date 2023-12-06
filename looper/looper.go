package looper

import (
	"expvar"
	"fmt"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/thalesfsp/status"
	"github.com/thalesfsp/sypl"
	"github.com/thalesfsp/sypl/level"
)

//////
// Consts, vars, and types.
//////

// Type of the entity in the framework.
const Type = "looper"

const defaultTimeout = 10 * time.Second

// Func is the function which is executed by the looper.
type Func func(IMetrics)

// IMetrics is the interface of the metrics.
type IMetrics interface {
	// GetLogger returns the logger of the looper.
	GetLogger() sypl.ISypl

	// IncrementFailedCounter increments the failed counter.
	IncrementFailedCounter()

	// IncrementSuccessCounter increments the success counter.
	IncrementSuccessCounter()

	// IncrementTotalCounter increments the total counter.
	IncrementTotalCounter()
}

// Metrics contains the metrics of a loop.
type Metrics struct {
	// Failed will be incremented every time a it fails execution.
	FailedCounter *expvar.Int `json:"-" validate:"omitempty,gte=0"`

	// Succeeded will be incremented every time a it succeeds execution.
	SuccessCounter *expvar.Int `json:"-" validate:"omitempty,gte=0"`

	// Total will be incremented every time a it is executed.
	TotalCounter *expvar.Int `json:"-" validate:"omitempty,gte=0"`
}

// ILooper is the interface of the looper.
type ILooper interface {
	// Start starts the loop.
	Start()

	// StartAsync starts the loop asynchronously.
	StartAsync()

	// Stop stops the loop.
	Stop()
}

// Looper is the individual, concurrent/non-blocking. The primary
// goal is to check the status of the API.
type Looper struct {
	// Func is the function to be executed.
	Func Func `json:"-" validate:"required"`

	// Interval is the time between each loop.
	Interval time.Duration `json:"interval"`

	// Logger is the logger.
	Logger sypl.ISypl `json:"-" validate:"required"`

	// Name of the looper.
	Name string `json:"name" validate:"required,lowercase,gte=1"`

	*Metrics

	s *gocron.Scheduler
}

//////
// Methods.
//////

// Start the loop.
func (l *Looper) Start() {
	// Starts the scheduler asynchronously.
	l.s.StartBlocking()
}

// StartAsync the loop.
func (l *Looper) StartAsync() {
	l.s.StartAsync()
}

// Stop the loop.
func (l *Looper) Stop() {
	l.s.Stop()
}

// GetLogger returns the logger of the looper.
func (l *Looper) GetLogger() sypl.ISypl {
	return l.Logger
}

// GetMetrics returns the metrics of the looper.
func (l *Looper) GetMetrics() any {
	return l.Metrics
}

// IncrementFailedCounter increments the failed counter.
func (l *Looper) IncrementFailedCounter() {
	if l.Metrics != nil {
		l.FailedCounter.Add(1)
	}
}

// IncrementSuccessCounter increments the success counter.
func (l *Looper) IncrementSuccessCounter() {
	if l.Metrics != nil {
		l.SuccessCounter.Add(1)
	}
}

// IncrementTotalCounter increments the total counter.
func (l *Looper) IncrementTotalCounter() {
	if l.Metrics != nil {
		l.TotalCounter.Add(1)
	}
}

//////
// Factory.
//////

// New creates a new Looper.
func New(
	name string,
	interval time.Duration,
	enableMetrics bool,
	loggingLevel string,
	metricsServerPort string,
	fn Func,
) (ILooper, error) {
	// Optionally enable logging.
	finalLoggingLevel := level.None

	if loggingLevel != "" {
		lvlInternal, err := level.FromString(loggingLevel)
		if err != nil {
			return nil, err
		}

		finalLoggingLevel = lvlInternal
	}

	l := &Looper{
		Func:     fn,
		Interval: interval,
		Logger:   sypl.NewDefault(name, finalLoggingLevel),
		Name:     name,
	}

	// Optionally enable metrics.
	if enableMetrics {
		l.Metrics = &Metrics{}

		l.FailedCounter = expvar.NewInt(fmt.Sprintf("%s.%s.%s", Type, name, status.Failed.String()))
		l.SuccessCounter = expvar.NewInt(fmt.Sprintf("%s.%s.%s", Type, name, status.Succeeded.String()))
		l.TotalCounter = expvar.NewInt(fmt.Sprintf("%s.%s.%s", Type, name, status.Total.String()))

		// Automatically expose the metrics.
		go func() {
			server := &http.Server{
				Addr:         fmt.Sprintf(":%s", metricsServerPort),
				ReadTimeout:  defaultTimeout,
				WriteTimeout: defaultTimeout,
			}

			if err := server.ListenAndServe(); err != nil {
				l.GetLogger().Fatalln(err)
			}
		}()

		// Set up function to be exposed when the metrics are called.
		expvar.Publish("", expvar.Func(l.GetMetrics))
	}

	// Set up the scheduler.
	s := gocron.NewScheduler(time.UTC)

	if _, err := s.Every(l.Interval).Do(func() {
		l.IncrementTotalCounter()

		l.Func(l)
	}); err != nil {
		return nil, err
	}

	l.s = s

	return l, nil
}
