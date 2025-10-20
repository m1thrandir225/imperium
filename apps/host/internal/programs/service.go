package programs

import "os/exec"

// Service represents the main functionality of the programs part
// Handles Program Discovery, Saving and Launching
// Rawg API Integration
type Service interface {
	DiscoverAndSavePrograms(paths []string) error
	GetLocalPrograms() ([]*Program, error)
	GetLocalProgramByPath(path string) (*Program, error)
	GetLocalProgramByID(id string) (*Program, error)
	DiscoverPrograms() ([]Program, error)
	DiscoverProgramsIn(paths []string) ([]Program, error) //??? should return pointer no?
	SaveProgram(req CreateProgramRequest) (*Program, error)
	LaunchProgram(path string) (*exec.Cmd, error)
	GetWindowTitleByProcessID(pid uint32) (string, error)
	RawgSearch(program Program) Program
}
