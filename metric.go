package pcp

type Metric struct {
	Name        string `json:"name"`
	ID          uint32 `json:"pmID"`
	Indom       uint32 `json:"indom"`
	Type        string `json:"type"`
	Sem         string `json:"instant"`
	Units       string `json:"units"`
	TextOneline string `json:"text-oneline"`
	TextHelp    string `json:"text-help"`
}

type Instance struct {
	Instance int32       `json:"instance"`
	Value    interface{} `json:"value"`
}

type MetricValue struct {
	Name      string `json:"name"`
	Pmid      uint32 `json:"pmid"`
	Instances []Instance
}
