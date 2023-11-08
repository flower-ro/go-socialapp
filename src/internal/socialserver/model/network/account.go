package network

type IsOnWhatAppRes struct {
	Total   int    ` json:"total,omitempty"`
	Valid   int    ` json:"valid,omitempty"`
	Members string ` json:"members,omitempty"`
}

type IsOnWhatAppReq struct {
	Members []string ` json:"members,omitempty"`
}
