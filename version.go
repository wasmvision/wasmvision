package wasmvision

var (
	version = "0.0.1"
	sha     string
)

func Version() string {
	if sha != "" {
		return version + "-" + sha
	}
	return version
}
