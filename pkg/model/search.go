package model

// Generated with https://mholt.github.io/json-to-go/

type ErrorResponse struct {
	Error *Error `json:"error,omitempty"`
}

type Error struct {
	Type   string `json:"type,omitempty"`
	Reason string `json:"reason,omitempty"`
}

type SearchResult struct {
	Took     int        `json:"took"`
	TimedOut bool       `json:"timed_out"`
	Hits     SearchHits `json:"hits"`
}
type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}
type ApacheAccess struct {
	RemoteHost                     string `json:"remoteHost"`
	CountryCode                    string `json:"countryCode"`
	Username                       string `json:"username"`
	Timestamp                      string `json:"timestamp"`
	RequestLine                    string `json:"requestLine"`
	Status                         int    `json:"status"`
	ResponseBodySize               int    `json:"responseBodySize"`
	Referer                        string `json:"referer"`
	UserAgent                      string `json:"userAgent"`
	ServerName                     string `json:"serverName"`
	ServerIP                       string `json:"serverIP"`
	ServerPort                     int    `json:"serverPort"`
	Handler                        string `json:"handler"`
	WorkerRoute                    string `json:"workerRoute"`
	TCPStatus                      string `json:"tcpStatus"`
	Cookie                         string `json:"cookie"`
	UniqueID                       string `json:"uniqueID"`
	RequestBytes                   int    `json:"requestBytes"`
	ResponseBytes                  int    `json:"responseBytes"`
	CompressionRatio               string `json:"compressionRatio"`
	RequestDuration                int    `json:"requestDuration"`
	ModsecTimeIn                   int    `json:"modsecTimeIn"`
	ApplicationTime                int    `json:"applicationTime"`
	ModsecTimeOut                  int    `json:"modsecTimeOut"`
	ModsecAnomalyScoreIn           int    `json:"modsecAnomalyScoreIn"`
	ModsecAnomalyScoreThresholdIn  int    `json:"modsecAnomalyScoreThresholdIn"`
	ModsecAnomalyScoreOut          int    `json:"modsecAnomalyScoreOut"`
	ModsecAnomalyScoreThresholdOut int    `json:"modsecAnomalyScoreThresholdOut"`
	ModsecParanoiaLevel            int    `json:"modsecParanoiaLevel"`
}
type ModsecAlert struct {
	Description  string   `json:"description"`
	ID           int      `json:"id"`
	Client       string   `json:"client"`
	Hostname     string   `json:"hostname"`
	URI          string   `json:"uri"`
	UniqueID     string   `json:"uniqueID"`
	Msg          string   `json:"msg"`
	Data         string   `json:"data"`
	Severity     string   `json:"severity"`
	Tags         []string `json:"tags"`
	File         string   `json:"file"`
	Line         int      `json:"line"`
	Ver          string   `json:"ver"`
	RuleTemplate string   `json:"rule_template"`
}

type Source struct {
	ApacheAccess *ApacheAccess `json:"apache-access,omitempty"`
	ModsecAlert  *ModsecAlert  `json:"modsec-alert,omitempty"`
}
type Hits struct {
	Index  string  `json:"_index"`
	Type   string  `json:"_type"`
	ID     string  `json:"_id"`
	Score  float64 `json:"_score"`
	Source Source  `json:"_source"`
}
type SearchHits struct {
	Total    int     `json:"total"`
	MaxScore float64 `json:"max_score"`
	Hits     []Hits  `json:"hits"`
}
type Fields struct {
	Timestamp                           []int64 `json:"@timestamp"`
	PipelineMetadataCollectorReceivedAt []int64 `json:"pipeline_metadata.collector.received_at"`
}
type Highlight struct {
	ModsecAlertUniqueID []string `json:"modsec-alert.uniqueID"`
	Message             []string `json:"message"`
}
