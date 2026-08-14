package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

// Service records runtime counters and timings for observability. It is safe
// for concurrent use.
type Service struct {
	uploadsTotal   atomic.Int64
	downloadsTotal atomic.Int64
	viewsTotal     atomic.Int64
	errorsTotal    atomic.Int64
	uploadBytes    atomic.Int64
	downloadBytes  atomic.Int64

	mu               sync.Mutex
	uploadEvents     []time.Time
	downloadEvents   []time.Time
	lastUpload       time.Time
	lastDownload     time.Time
	totalUploadDur   time.Duration
	totalDownloadDur time.Duration
}

// New creates a metrics service.
func New() *Service {
	return &Service{}
}

// RecordUpload tracks a completed upload of n bytes taking dur.
func (s *Service) RecordUpload(n int64, dur time.Duration) {
	s.uploadsTotal.Add(1)
	s.uploadBytes.Add(n)
	now := time.Now()
	s.mu.Lock()
	s.uploadEvents = append(s.uploadEvents, now)
	s.lastUpload = now
	s.totalUploadDur += dur
	s.mu.Unlock()
}

// RecordDownload tracks a served download of n bytes taking dur.
func (s *Service) RecordDownload(n int64, dur time.Duration) {
	s.downloadsTotal.Add(1)
	s.downloadBytes.Add(n)
	now := time.Now()
	s.mu.Lock()
	s.downloadEvents = append(s.downloadEvents, now)
	s.lastDownload = now
	s.totalDownloadDur += dur
	s.mu.Unlock()
}

// RecordView tracks an inline view.
func (s *Service) RecordView(n int64) {
	s.viewsTotal.Add(1)
}

// RecordError tracks a failed request/operation.
func (s *Service) RecordError() {
	s.errorsTotal.Add(1)
}

// Snapshot is a point-in-time view of the metrics.
type Snapshot struct {
	UploadsTotal     int64   `json:"uploads_total"`
	DownloadsTotal   int64   `json:"downloads_total"`
	ViewsTotal       int64   `json:"views_total"`
	ErrorsTotal      int64   `json:"errors_total"`
	UploadBytes      int64   `json:"upload_bytes"`
	DownloadBytes    int64   `json:"download_bytes"`
	UploadsPerMinute float64 `json:"uploads_per_minute"`
	DownloadsPerMin  float64 `json:"downloads_per_minute"`
	AvgUploadMs      float64 `json:"avg_upload_ms"`
	AvgDownloadMs    float64 `json:"avg_download_ms"`
	LastUpload       string  `json:"last_upload_at,omitempty"`
	LastDownload     string  `json:"last_download_at,omitempty"`
}

// Snapshot returns the current counters and rates.
func (s *Service) Snapshot() Snapshot {
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()

	snap := Snapshot{
		UploadsTotal:   s.uploadsTotal.Load(),
		DownloadsTotal: s.downloadsTotal.Load(),
		ViewsTotal:     s.viewsTotal.Load(),
		ErrorsTotal:    s.errorsTotal.Load(),
		UploadBytes:    s.uploadBytes.Load(),
		DownloadBytes:  s.downloadBytes.Load(),
	}

	cutoff := now.Add(-time.Minute)
	var upm, dpm int
	for _, t := range s.uploadEvents {
		if t.After(cutoff) {
			upm++
		}
	}
	for _, t := range s.downloadEvents {
		if t.After(cutoff) {
			dpm++
		}
	}
	snap.UploadsPerMinute = float64(upm)
	snap.DownloadsPerMin = float64(dpm)

	if upm > 0 {
		snap.AvgUploadMs = float64(s.totalUploadDur.Milliseconds()) / float64(upm)
	}
	if dpm > 0 {
		snap.AvgDownloadMs = float64(s.totalDownloadDur.Milliseconds()) / float64(dpm)
	}
	if !s.lastUpload.IsZero() {
		snap.LastUpload = s.lastUpload.Format(time.RFC3339)
	}
	if !s.lastDownload.IsZero() {
		snap.LastDownload = s.lastDownload.Format(time.RFC3339)
	}
	return snap
}

// Default is the process-wide metrics service.
var Default = New()
