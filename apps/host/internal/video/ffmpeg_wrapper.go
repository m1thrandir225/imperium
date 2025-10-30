package video

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/m1thrandir225/imperium/apps/host/internal/util"
)

type FFMPEGWrapper struct {
	path    string
	cmd     *exec.Cmd
	running bool
	stdin   io.WriteCloser
}

func NewFFMPEGWrapper(path string) (*FFMPEGWrapper, error) {

	if !util.IsValidPath(path) {
		return nil, ErrInvalidPath
	}

	return &FFMPEGWrapper{
		path: path,
	}, nil
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

func (w *FFMPEGWrapper) ExecuteWithStdout(args ...string) (io.ReadCloser, error) {
	w.cmd = exec.Command(w.path, args...)

	var err error
	w.stdin, err = w.cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %v", err)
	}

	stdout, err := w.cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	w.cmd.Stderr = os.Stderr

	w.running = true
	log.Printf("Executing command with stdout: %s", w.cmd.String())

	if err := w.cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %v", err)
	}

	return &ffmpegStream{
		rc:   stdout,
		stop: w.Stop,
	}, nil
}

// ffmpegStream Implements io.ReadCloser
type ffmpegStream struct {
	rc   io.ReadCloser
	stop func() error
}

func (s *ffmpegStream) Read(p []byte) (int, error) {
	return s.rc.Read(p)
}

func (s *ffmpegStream) Close() error {
	if s.stop != nil {
		_ = s.stop()
	}
	return s.rc.Close()
}
