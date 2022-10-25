package cdn

const (
	DomainStatusRunning   = "RUNNING"
	DomainStatusOperating = "OPERATING"
	DomainStatusStopped   = "STOPPED"

	DomainStatusAll = "ALL"
)

const (
	DomainFormDefault  = "default"
	DomainFormImage    = "image"
	DomainFormDownload = "download"
	DomainFormMedia    = "media"
	DomainFormDynamic  = "dynamic"
)

func DomainFormValues() []string {
	return []string{
		DomainFormDefault,
		DomainFormImage,
		DomainFormDownload,
		DomainFormMedia,
		DomainFormDynamic,
	}
}

func validErrorPageStatusCodes() []int {
	return []int{
		401, 403, 404, 405, 414, 429,
		500, 501, 502, 503, 504,
	}
}

func validDragMode(forMp4 bool) []string {
	if forMp4 {
		return []string{"second"}
	}
	return []string{"byteAV", "byte"}
}
