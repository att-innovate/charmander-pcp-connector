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
