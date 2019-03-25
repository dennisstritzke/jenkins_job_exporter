FROM golang:1.12.1-alpine as build

RUN apk add -U --no-cache ca-certificates

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY binaries/jenkins_job_exporter /jenkins_job_exporter

EXPOSE 3000
VOLUME ["/cache"]
ENTRYPOINT ["/jenkins_job_exporter"]