package exporter

import (
	"encoding/json"
	"fmt"
	"github.com/bndr/gojenkins"
	"github.com/dennisstritzke/jenkins_job_exporter/util"
	"io/ioutil"
	"log"
	"regexp"
)

const (
	jenkinsUrlEnvVarName              = "JENKINS_URL"
	jenkinsUserEnvVarName             = "JENKINS_USER"
	jenkinsPasswordEnvVarName         = "JENKINS_API_TOKEN"
	statusMapCacheFileEnvVarName      = "CACHE_FILE_LOCATION"
	statusMapCacheFileDefaultLocation = "/cache/buildStatusCache.json"

	jobsFromViewEnvVarName   = "JENKINS_VIEW"
	jobFilterRegexEnvVarName = "JENKINS_JOB_FILTER_REGEX"
	jobFilterDefaultRegex    = ".*"
)

type Jenkins struct {
	jenkins   *gojenkins.Jenkins
	jobs      []*gojenkins.Job
	statusMap map[string]map[int64]BuildStatus
}

func New() (*Jenkins, error) {
	jenkinsUrl, err := util.EnvValue(jenkinsUrlEnvVarName)
	if err != nil {
		return nil, err
	}
	jenkinsUser, err := util.EnvValue(jenkinsUserEnvVarName)
	if err != nil {
		return nil, err
	}
	jenkinsPassword, err := util.EnvValue(jenkinsPasswordEnvVarName)
	if err != nil {
		return nil, err
	}

	jenkins := gojenkins.CreateJenkins(nil, jenkinsUrl, jenkinsUser, jenkinsPassword)
	_, err = jenkins.Init()
	if err != nil {
		return nil, err
	}

	return &Jenkins{
		jenkins:   jenkins,
		statusMap: make(map[string]map[int64]BuildStatus),
	}, nil
}

func (j *Jenkins) persistStatusMap() error {
	statusMapJson, err := json.Marshal(j.statusMap)
	if err != nil {
		return err
	}

	statusMapFileName := util.EnvValueWithDefault(statusMapCacheFileEnvVarName, statusMapCacheFileDefaultLocation)
	err = ioutil.WriteFile(statusMapFileName, statusMapJson, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (j *Jenkins) retrieveStatusMapFromFile() (map[string]map[int64]BuildStatus, error) {
	statusMapFileName := util.EnvValueWithDefault(statusMapCacheFileEnvVarName, statusMapCacheFileDefaultLocation)

	var statusMap map[string]map[int64]BuildStatus
	statusMapJson, err := ioutil.ReadFile(statusMapFileName)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(statusMapJson, &statusMap)
	if err != nil {
		return nil, err
	}

	return statusMap, nil
}

func (j *Jenkins) Init() error {
	_, err := j.initialiseJobs()
	if err != nil {
		return err
	}

	j.initialiseStatusMap()

	return nil
}

func (j *Jenkins) initialiseJobs() ([]*gojenkins.Job, error) {
	jenkinsViewName := util.EnvValueWithDefault(jobsFromViewEnvVarName, "All")
	folder, err := j.jenkins.GetView(jenkinsViewName)
	if err != nil {
		return nil, err
	}

	jobFilterRegexString := util.EnvValueWithDefault(jobFilterRegexEnvVarName, jobFilterDefaultRegex)
	jobFilterRegex := regexp.MustCompile(jobFilterRegexString)

	var jobs []*gojenkins.Job

	for _, innerJob := range folder.Raw.Jobs {
		if jobFilterRegex.Match([]byte(innerJob.Name)) {
			job, err := j.jenkins.GetJob(innerJob.Name)
			if err != nil {
				log.Println(fmt.Sprintf("Unable to retrieve Job '%s'", innerJob.Name))
				continue
			}

			jobs = append(jobs, job)
		}
	}

	j.jobs = jobs
	return jobs, nil
}

func (j *Jenkins) initialiseStatusMap() {
	log.Println("===> Initialising Jobs")

	statusMapCache, err := j.retrieveStatusMapFromFile()
	if err != nil {
		log.Println("Unable to load cache. Continuing...")
	} else {
		j.statusMap = statusMapCache
	}

	for index, job := range j.jobs {
		log.Println(fmt.Sprintf("(%d/%d) %s...", index+1, len(j.jobs), job.GetName()))
		buildMap := make(map[int64]BuildStatus)

		jobBuilds, err := job.GetAllBuildIds()
		if err != nil {
			continue
		}

		for _, build := range jobBuilds {
			if _, ok := j.statusMap[job.GetName()][build.Number]; !ok {
				buildStatus := j.getBuildStatus(job.GetName(), build.Number)
				if buildStatus == jobStatusSuccess || buildStatus == jobStatusFailure || buildStatus == jobStatusUnstable || buildStatus == jobStatusAborted {
					// The job has a result, so we don't want to report metrics on that one as they are already collected.
					buildMap[build.Number] = jobStatusIgnored
				} else {
					// The job hasn't finished yet, so we want to collect metrics.
					buildMap[build.Number] = jobStatusUnknown
				}
			} else {
				buildMap[build.Number] = j.statusMap[job.GetName()][build.Number]
			}
		}
		j.statusMap[job.GetName()] = buildMap
	}
	log.Println("Writing cache")
	err = j.persistStatusMap()
	if err != nil {
		log.Println("Error writing cache. Continuing...")
	}
	log.Println("Done")
}

func (j *Jenkins) updateStatusMap() {
	for _, job := range j.jobs {
		jobBuilds, err := job.GetAllBuildIds()
		if err != nil {
			continue
		}

		for _, build := range jobBuilds {
			if _, ok := j.statusMap[job.GetName()][build.Number]; !ok {
				j.statusMap[job.GetName()][build.Number] = jobStatusUnknown
			}
		}
	}
}

func (j *Jenkins) RefreshStatus() []StatusMetric {
	j.updateStatusMap()

	var statusSlice []StatusMetric
	for job, buildStatusMap := range j.statusMap {
		for buildNumber := range buildStatusMap {
			if j.statusMap[job][buildNumber] == jobStatusUnknown {
				status := j.refreshBuildStatus(job, buildNumber)

				statusSlice = append(statusSlice, StatusMetric{
					JobName:     job,
					BuildNumber: buildNumber,
					BuildStatus: status,
				})
			}
		}
	}

	return statusSlice
}

func (j *Jenkins) getBuildStatus(job string, buildNumber int64) BuildStatus {
	buildObject, err := j.jenkins.GetBuild(job, buildNumber)
	if err != nil {
		return jobStatusUnknown
	}

	return DetermineBuildStatus(buildObject.GetResult())
}

func (j *Jenkins) refreshBuildStatus(job string, buildNumber int64) BuildStatus {
	status := j.getBuildStatus(job, buildNumber)
	j.statusMap[job][buildNumber] = status
	return status
}
