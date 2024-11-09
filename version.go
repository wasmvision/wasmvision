package wasmvision

var (
	version = "0.1.0"
	sha     string
)

func Version() string {
	if sha != "" {
		return version + "-" + sha
	}
	return version
}
