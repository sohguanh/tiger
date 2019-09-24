// clientUtil is a thin wrapper over http.Client and net.* to provide different timeout,number of retries,time to wait between each retry for each http and net request.
package clientUtil

// RequestInfo is the struct to store the information for doing retry.
type RequestInfo struct {
	TimeoutSec         int
	RetryTimes         int
	WaitBeforeRetrySec int
}
