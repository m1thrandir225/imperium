package programs

type Database interface {
	SaveProgram(program *Program) error
	GetPrograms() ([]*Program, error)
	GetProgramByID(id string) (*Program, error)
	GetProgramByPath(path string) (*Program, error)
	CleanupNonExistentPrograms() error
}
