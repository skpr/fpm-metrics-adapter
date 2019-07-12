package fpm

// Status of the FPM pool.
type Status struct {
	Processes StatusProcesses `json:"processes"`
}

// StatusProcesses for listing process metrics.
type StatusProcesses struct {
	Active int64 `json:"active"`
	Total  int64 `json:"total"`
}
