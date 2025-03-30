//go:build !cuda

package runtime

func CheckCUDA() bool {
	return false
}
