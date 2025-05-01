package fpm

// QueryResponse provided by the FPM status request with query string "json&full".
// This is a temporay struct and is marshalled into our Status struct.
// https://www.php.net/manual/en/fpm.status.php
type QueryResponse struct {
	ProcessManager     string `json:"process manager"`
	ListenQueue        int64  `json:"listen queue"`
	ListenQueueLen     int64  `json:"listen queue len"`
	IdleProcesses      int64  `json:"idle processes"`
	ActiveProcesses    int64  `json:"active processes"`
	TotalProcesses     int64  `json:"total processes"`
	MaxActiveProcesses int64  `json:"max active processes"`
}

// Status of the FPM pool.
type Status struct {
	// The process manager type - static, dynamic or ondemand.
	ProcessManager string `json:"phpfpm_process_manager"`
	// The number of requests (backlog) currently waiting for a free process.
	ListenQueue int64 `json:"phpfpm_listen_queue"`
	// The maximum allowed size of the listen queue.
	ListenQueueLen int64 `json:"phpfpm_listen_queue_len"`
	// The number of processes that are currently idle (waiting for requests).
	IdleProcesses int64 `json:"phpfpm_idle_processes"`
	// The number of processes that are currently processing requests.
	ActiveProcesses int64 `json:"phpfpm_active_processes"`
	// The current total number of processes.
	TotalProcesses int64 `json:"phpfpm_total_processes"`
	// The maximum number of concurrently active processes.
	MaxActiveProcesses int64 `json:"phpfpm_max_active_processes"`
}
