package video

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

type FFMPEGWrapper struct {
	path    string
	cmd     *exec.Cmd
	running bool
	stdin   io.WriteCloser
}

func NewFFMPEGWrapper(path string) *FFMPEGWrapper {
	return &FFMPEGWrapper{
		path: path,
	}
}

func (w *FFMPEGWrapper) Execute(args ...string) error {

	w.cmd = exec.Command(w.path, args...)

	var err error
	w.stdin, err = w.cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %v", err)
	}

	w.cmd.Stdout = os.Stdout
	w.cmd.Stderr = os.Stderr

	w.running = true
	log.Printf("Executing command: %s", w.cmd.String())
	return w.cmd.Start()
}

func (w *FFMPEGWrapper) ExecuteWithoutOutput(args ...string) ([]byte, error) {
	w.cmd = exec.Command(w.path, args...)
	w.running = true
	return w.cmd.CombinedOutput()
}

func (w *FFMPEGWrapper) Version() ([]byte, error) {
	return w.ExecuteWithoutOutput("-version")
}

func (w *FFMPEGWrapper) Stop() error {
	if w.cmd != nil && w.cmd.Process != nil && w.running {
		w.running = false

		// Send 'q' command to FFmpeg through stdin
		if w.stdin != nil {
			_, err := w.stdin.Write([]byte("q"))
			if err != nil {
				log.Printf("Warning: Failed to send quit command: %v", err)
			}
			w.stdin.Close()
		}

		return w.cmd.Wait()
	}
	return fmt.Errorf("no running FFmpeg process to stop")

}
