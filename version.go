package wasmvision

var (
	version = "0.3.1"
	sha     string
)

func Version() string {
	if sha != "" {
		return version + "-" + sha
	}
	return version
}
