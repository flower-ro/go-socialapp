package request

type GroupCreateReq struct {
	Name    string   ` json:"name,omitempty"`
	Member  []string ` json:"member,omitempty"`
	Creator string   ` json:"creator,omitempty"`
}
