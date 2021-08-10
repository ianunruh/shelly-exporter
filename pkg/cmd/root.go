package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/ianunruh/shelly-exporter/pkg/collect"
)

var listenAddr string

var rootCmd = &cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		log, err := zap.NewProduction()
		if err != nil {
			return err
		}

		http.Handle("/metrics", promhttp.Handler())
		http.Handle("/probe", collect.ProbeHandler(log))

		log.Info("Starting metrics server",
			zap.String("address", listenAddr))
		return http.ListenAndServe(listenAddr, nil)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&listenAddr, "listen", "l", ":9090", "Host/port to listen on")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
