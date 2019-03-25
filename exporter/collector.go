package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricBuildRunningCount = prometheus.NewDesc("jenkins_job_running_builds_count", "number of running Jenkins Job builds", nil, nil)
	metricTotal             = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "jenkins_job_builds_total", Help: "total number of finished Jenkins builds"}, []string{"job"})
	metricSuccess           = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "jenkins_job_build_success_total", Help: "total number of successfully finished Jenkins builds"}, []string{"job"})
	metricFailure           = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "jenkins_job_build_failure_total", Help: "total number of finished Jenkins builds in failure state"}, []string{"job"})
	metricUnstable          = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "jenkins_job_build_unstable_total", Help: "total number of finished Jenkins builds in unstable state"}, []string{"job"})
	metricAborted           = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "jenkins_job_build_aborted_total", Help: "total number of finished Jenkins builds in aborted state"}, []string{"job"})
)

func (j *Jenkins) Describe(ch chan<- *prometheus.Desc) {
	ch <- metricBuildRunningCount
	metricTotal.Describe(ch)
	metricSuccess.Describe(ch)
	metricFailure.Describe(ch)
	metricUnstable.Describe(ch)
	metricAborted.Describe(ch)
}

func (j *Jenkins) Collect(ch chan<- prometheus.Metric) {
	statusSlice := j.RefreshStatus()

	runningCounter := 0
	for _, status := range statusSlice {
		switch status.BuildStatus {
		case jobStatusSuccess:
			metricTotal.WithLabelValues(status.JobName).Add(1)
			metricSuccess.WithLabelValues(status.JobName).Add(1)
		case jobStatusFailure:
			metricTotal.WithLabelValues(status.JobName).Add(1)
			metricFailure.WithLabelValues(status.JobName).Add(1)
		case jobStatusUnstable:
			metricTotal.WithLabelValues(status.JobName).Add(1)
			metricUnstable.WithLabelValues(status.JobName).Add(1)
		case jobStatusUnknown:
			runningCounter += 1
		case jobStatusAborted:
			metricTotal.WithLabelValues(status.JobName).Add(1)
			metricAborted.WithLabelValues(status.JobName).Add(1)
		}
	}

	ch <- prometheus.MustNewConstMetric(metricBuildRunningCount, prometheus.GaugeValue, float64(runningCounter))

	for _, job := range j.jobs {
		ch <- metricTotal.WithLabelValues(job.GetName())
		ch <- metricSuccess.WithLabelValues(job.GetName())
		ch <- metricFailure.WithLabelValues(job.GetName())
		ch <- metricUnstable.WithLabelValues(job.GetName())
		ch <- metricAborted.WithLabelValues(job.GetName())
	}
}
