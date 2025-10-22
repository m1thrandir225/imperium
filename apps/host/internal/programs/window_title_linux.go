//go:build linux
// +build linux

package programs

func (s *programService) GetWindowTitleByProcessID(pid uint32) (string, error) {
	return "", ErrOSNotSupported
}
