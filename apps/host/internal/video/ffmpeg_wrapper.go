package video

import "os/exec"

type FFMPEGWrapper struct {
	path string
}

func NewFFMPEGWrapper(path string) *FFMPEGWrapper {
	return &FFMPEGWrapper{
		path: path,
	}
}

func (w *FFMPEGWrapper) Execute(args ...string) error {
	cmd := exec.Command(w.path, args...)
	return cmd.Run()
}

func (w *FFMPEGWrapper) ExecuteWithoutOutput(args ...string) ([]byte, error) {
	cmd := exec.Command(w.path, args...)
	return cmd.CombinedOutput()
}

func (w *FFMPEGWrapper) Version() ([]byte, error) {
	return w.ExecuteWithoutOutput("-version")
}
