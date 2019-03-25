# Jenkins Job Exporter
Prometheus exporter for Jenkins Job metrics, written in Go.

## Metrics
- `jenkins_job_running_builds_count`
- `jenkins_job_builds_total`
- `jenkins_job_build_success_total`
- `jenkins_job_build_failure_total`
- `jenkins_job_build_unstable_total`
- `jenkins_job_build_aborted_total`

## Configuration
| Value | Description | Default | Example |
|-------|-------------|---------|---------|
| `JENKINS_URL` | URL of the Jenkins instance. | `none` | `https://jenkins.example.com/` |
| `JENKINS_USER` | User that is permitted to access the Jenkins instance. | `none` | `jdoe` |
| `JENKINS_API_TOKEN` | API token of the Jenkins user. | `none` | `93012748faa8f49a45137976181d4e3e18` | 
| `CACHE_FILE_LOCATION` | Location of the file used to cache Jenkins Job status. | `/cache/buildStatusCache.json` | `/any/path/you/like.json` |
| `JENKINS_VIEW` | Jenkins View in which the Jobs should be scraped. | `none` | `Staging` | 
| `JENKINS_JOB_FILTER_REGEX` | Golang Regex to filter the Jobs of the configured view. | `.*` | `[0-9]+_.*_important` | 