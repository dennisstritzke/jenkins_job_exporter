package main

import (
	"github.com/dennisstritzke/jenkins_job_exporter/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
)

func main() {
	jenkinsExporter, err := exporter.New()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
		return
	}

	err = jenkinsExporter.Init()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
		return
	}

	prometheus.MustRegister(jenkinsExporter)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
             <head><title>Jenkins Job Exporter</title></head>
             <body>
             <h1>Jenkins Job Exporter</h1>
             <p><a href='/metrics'>Metrics</a></p>
             </body>
             </html>`))
	})
	http.Handle("/metrics", promhttp.Handler())

	log.Println("Listening on 0.0.0.0:3000")
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
