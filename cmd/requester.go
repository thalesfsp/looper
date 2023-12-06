package cmd

import (
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/thalesfsp/httpclient"
	"github.com/thalesfsp/looper/looper"
	"github.com/thalesfsp/sypl/level"
)

var (
	enableMetrics     bool
	httpVerb          string
	interval          time.Duration
	logLevel          string
	metricsServerPort string
	reqBodyFilePath   string
	url               string
)

// requesterCmd represents the requester command.
var requesterCmd = &cobra.Command{
	Use:   "requester",
	Short: "Runs a request looper",
	Run: func(cmd *cobra.Command, args []string) {
		var reqBody []byte

		// If req-body-file is set, read the file and set the body.
		if reqBodyFilePath != "" {
			jsonFile, err := os.Open(reqBodyFilePath)
			if err != nil {
				cliLogger.Fatalln(err)
			}

			// defer the closing of our jsonFile so that we can parse it later on
			defer jsonFile.Close()

			b, err := io.ReadAll(jsonFile)
			if err != nil {
				cliLogger.Fatalln(err)
			}

			reqBody = b
		}

		opts := []httpclient.Func{}

		if reqBody == nil {
			opts = append(opts, httpclient.WithReqBody(reqBody))
		}

		requester, err := looper.NewRequester(
			strings.ToUpper(httpVerb),
			url,
			interval,
			enableMetrics,
			logLevel,
			metricsServerPort,
			opts...,
		)
		if err != nil {
			cliLogger.Fatalln(err)
		}

		// Start the requester looper.
		requester.Start()
	},
}

func init() {
	rootCmd.AddCommand(requesterCmd)

	requesterCmd.Flags().BoolVarP(&enableMetrics, "enable-metrics", "m", false, "Enable metrics")
	requesterCmd.Flags().DurationVarP(&interval, "interval", "i", 1*time.Second, "Interval between requests")
	requesterCmd.Flags().StringVarP(&httpVerb, "http-verb", "v", http.MethodGet, "HTTP verb to be used")
	requesterCmd.Flags().StringVarP(&logLevel, "log-level", "l", level.None.String(), "Log level")
	requesterCmd.Flags().StringVarP(&reqBodyFilePath, "req-body-file-path", "f", "", "File containing the request body")
	requesterCmd.Flags().StringVarP(&metricsServerPort, "metrics-server-port", "p", "8080", "Metrics server port")
	requesterCmd.Flags().StringVarP(&url, "url", "u", "", "URL to be requested")

	requesterCmd.MarkFlagRequired("http-verb")
	requesterCmd.MarkFlagRequired("url")
}
