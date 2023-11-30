package parse

// 解析接收到的webHook参数

type Alert struct {
	Status       string      `json:"status"`
	Labels       Labels      `json:"labels"`
	Annotations  Annotations `json:"annotations"`
	StartsAt     string      `json:"startsAt"`
	EndsAt       string      `json:"endsAt"`
	GeneratorURL string      `json:"generatorURL"`
	Fingerprint  string      `json:"fingerprint"`
	SilenceURL   string      `json:"silenceURL"`
	DashboardURL string      `json:"dashboardURL"`
	PanelURL     string      `json:"panelURL"`
	Values       interface{} `json:"values"`
	ValueString  string      `json:"valueString"`
}

type Labels struct {
	Addr          string `json:"addr"`
	AlertName     string `json:"alertname"`
	Alter         string `json:"alter"`
	Device        string `json:"device"`
	Fstype        string `json:"fstype"`
	GrafanaFolder string `json:"grafana_folder"`
	Instance      string `json:"instance"`
	Job           string `json:"job"`
	MountPoint    string `json:"mountpoint"`
	Project       string `json:"project"`
	HostName      string `json:"hostname"`
}

type Annotations struct {
	Description string `json:"description"`
	Level       string `json:"level"`
}
