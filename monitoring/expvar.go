package monitoring

import (
	"expvar"
	"time"
)

var (
	StartTime         = expvar.NewString("Application start time")
	Uptime            = expvar.NewString("Application uptime")
	TotalUsers        = expvar.NewInt("Total users")
	TotalRequests     = expvar.NewInt("Total Requests")
	RequestsPerSecond = expvar.NewFloat("Requests per second")
	FailedRequests    = expvar.NewInt("Total failed requests")
	FailedPerSecond   = expvar.NewFloat("Failed requests per second")
)

func TimeMonitoring(startTime time.Time) {
	StartTime.Set(startTime.String())
	for {
		uptime := time.Since(startTime)
		uptimeString := uptime.String()

		Uptime.Set(uptimeString)
		// Requests per second
		requests := float64(TotalRequests.Value())
		if requests != 0 {
			RequestsPerSecond.Set(requests / float64(uptime.Seconds()))
		}
		//Failed per second
		failed := float64(FailedRequests.Value())
		if failed != 0 {
			FailedPerSecond.Set(failed / float64(uptime.Seconds()))
		}
	}
}

func FailedRequest() {
	TotalRequests.Add(1)
	FailedRequests.Add(1)
}
