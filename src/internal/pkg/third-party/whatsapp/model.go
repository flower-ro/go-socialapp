package whatsapp

type ExtractedMedia struct {
	MediaPath string `json:"media_path"`
	MimeType  string `json:"mime_type"`
	Caption   string `json:"caption"`
}

type evtMessage struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type evtReaction struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}
