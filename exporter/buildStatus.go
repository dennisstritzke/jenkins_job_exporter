package exporter

const (
	jobStatusSuccess  BuildStatus = 0
	jobStatusFailure  BuildStatus = 1
	jobStatusUnstable BuildStatus = 2
	jobStatusUnknown  BuildStatus = 3
	jobStatusIgnored  BuildStatus = 4
	jobStatusAborted  BuildStatus = 5
)

type BuildStatus int

func DetermineBuildStatus(input string) BuildStatus {
	switch input {
	case "SUCCESS":
		return jobStatusSuccess
	case "FAILURE":
		return jobStatusFailure
	case "UNSTABLE":
		return jobStatusUnstable
	case "ABORTED":
		return jobStatusAborted
	default:
		return jobStatusUnknown
	}
}

type StatusMetric struct {
	JobName     string
	BuildNumber int64
	BuildStatus BuildStatus
}
