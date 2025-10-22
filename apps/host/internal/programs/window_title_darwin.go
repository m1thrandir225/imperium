//go:build darwin
// +build darwin

package programs

func (s *programService) GetWindowTitleByProcessID(pid uint32) (string, error) {
	return "", ErrOSNotSupported
}
