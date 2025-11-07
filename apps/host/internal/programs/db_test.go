package programs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	InMemoryDbShared = ":memory:?cache=shared"
	InMemoryDb       = ":memory:"
)

func TestNewDatabase(t *testing.T) {
	testCases := []struct {
		name        string
		dbPath      string
		expectError bool
		expectedErr error
	}{
		{
			name:        "valid path - should create database",
			dbPath:      InMemoryDb,
			expectError: false,
		},
		{
			name:        "valid file path - should create database and directories",
			dbPath:      filepath.Join(t.TempDir(), "test", "database.db"),
			expectError: false,
		},
		{
			name:        "invalid path - should fail",
			dbPath:      "/root/restricted/database.db",
			expectError: true,
			expectedErr: ErrDatabaseDirectoryCreationFailed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, err := NewDatabase(tc.dbPath)
			if tc.expectError {
				require.Error(t, err)
				require.Nil(t, db)
				require.Equal(t, err, ErrDatabaseDirectoryCreationFailed)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, db)
				if sqliteDB, ok := db.(*sqliteDB); ok {
					sqliteDB.Close()
				}
			}
		})
	}
}

func TestSqliteDB_InitTables(t *testing.T) {
	db, err := NewDatabase(InMemoryDbShared)
	require.NoError(t, err)
	require.NotNil(t, db)

	defer func() {
		if sqliteDB, ok := db.(*sqliteDB); ok {
			_ = sqliteDB.Close()
		}
	}()

	program := &Program{
		Name:        "Test Program",
		Path:        "/test/path",
		Description: "Test Description",
	}

	err = db.SaveProgram(program)
	require.NoError(t, err)
}

func TestSqliteDB_SaveProgram(t *testing.T) {
	db, err := NewDatabase(InMemoryDbShared)
	require.NoError(t, err)
	require.NotNil(t, db)

	defer func() {
		if sqliteDB, ok := db.(*sqliteDB); ok {
			_ = sqliteDB.Close()
		}
	}()

	testCases := []struct {
		name        string
		program     *Program
		errExpected bool
	}{
		{
			name: "vaid program - should save successfully",
			program: &Program{
				Name:        "Test Program",
				Path:        "/test/path",
				Description: "Test Description",
			},
			errExpected: false,
		},
		{
			name: "invalid program - should fail",
			program: &Program{
				Name:        "",
				Path:        "",
				Description: "",
			},
			errExpected: true,
		},
		{
			name:        "empty program - should fail",
			program:     &Program{},
			errExpected: true,
		},
		{
			name:        "nil program - should fail",
			program:     nil,
			errExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := db.SaveProgram(tc.program)
			if tc.errExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSqliteDB_SaveProgram_Update(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

		})
	}
}

func TestSqliteDB_GetPrograms(t *testing.T) {
	programs := []Program{
		{
			Name:        "Test Program",
			Path:        "/test/path",
			Description: "Test Description",
		},
		{
			Name:        "Test Program 2",
			Path:        "/test/path/2",
			Description: "Test Description 2",
		},
		{
			Name:        "Test Program 3",
			Path:        "/test/path/3",
			Description: "Test Description 3",
		},
	}

	testCases := []struct {
		name           string
		build          func(db Database) ([]*Program, error)
		errExpected    bool
		programsLength int
	}{
		{
			name: "valid -  insert one return one",
			build: func(db Database) ([]*Program, error) {
				err := db.SaveProgram(&programs[0])
				require.NoError(t, err)

				return db.GetPrograms()
			},
			errExpected:    false,
			programsLength: 1,
		},
		{
			name: "valid - insert all return all",
			build: func(db Database) ([]*Program, error) {
				for _, program := range programs {
					err := db.SaveProgram(&program)
					require.NoError(t, err)
				}
				return db.GetPrograms()
			},
			errExpected:    false,
			programsLength: 3,
		},
		{
			name: "valid - no programs",
			build: func(db Database) ([]*Program, error) {
				return db.GetPrograms()
			},
			errExpected:    false,
			programsLength: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, err := NewDatabase(InMemoryDb)
			require.NoError(t, err)
			require.NotNil(t, db)

			defer func() {
				if sqliteDB, ok := db.(*sqliteDB); ok {
					_ = sqliteDB.Close()
				}
			}()

			dbPrograms, err := tc.build(db)
			if tc.errExpected {
				require.Error(t, err)
				require.Nil(t, programs)
			} else {
				require.NoError(t, err)
				if tc.programsLength != 0 {
					require.NotEmpty(t, dbPrograms)
					require.NotNil(t, dbPrograms)
				}
				require.Equal(t, tc.programsLength, len(dbPrograms))
			}
		})
	}
}

func TestSqliteDB_GetProgramByID(t *testing.T) {
	db, err := NewDatabase(InMemoryDbShared)
	require.NoError(t, err)
	require.NotNil(t, db)

	defer func() {
		if sqliteDB, ok := db.(*sqliteDB); ok {
			_ = sqliteDB.Close()
		}
	}()

	program := Program{
		Name:        "Test Program",
		Path:        "/test/path",
		Description: "Test Description",
	}

	testCases := []struct {
		name        string
		build       func() (*Program, error)
		errExpected bool
	}{
		{
			name: "valid",
			build: func() (*Program, error) {
				err = db.SaveProgram(&program)
				require.NoError(t, err)

				return db.GetProgramByID("1") // first ID is 1
			},
			errExpected: false,
		},
		{
			name: "non existent",
			build: func() (*Program, error) {
				return db.GetProgramByID("2")
			},
			errExpected: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dbProgram, err := tc.build()
			if tc.errExpected {
				require.Error(t, err)
				require.Nil(t, dbProgram)
			} else {
				require.NoError(t, err)
				require.NotNil(t, dbProgram)
				require.Equal(t, program.Name, dbProgram.Name)
				require.Equal(t, program.Path, dbProgram.Path)
			}
		})
	}
}

func TestSqliteDB_GetProgramByPath(t *testing.T) {
	db, err := NewDatabase(InMemoryDbShared)
	require.NoError(t, err)
	require.NotNil(t, db)

	defer func() {
		if sqliteDB, ok := db.(*sqliteDB); ok {
			_ = sqliteDB.Close()
		}
	}()

	program := Program{
		Name:        "Test Program",
		Path:        "/test/path",
		Description: "Test Description",
	}

	testCases := []struct {
		name        string
		build       func() (*Program, error)
		errExpected bool
	}{
		{
			name: "valid",
			build: func() (*Program, error) {
				err = db.SaveProgram(&program)
				require.NoError(t, err)
				return db.GetProgramByPath(program.Path)
			},
			errExpected: false,
		},
		{
			name: "non existent",
			build: func() (*Program, error) {
				return db.GetProgramByPath("/test/path/2")
			},
			errExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(t.Name(), func(t *testing.T) {
			dbProgram, err := tc.build()
			if tc.errExpected {
				require.Error(t, err)
				require.Nil(t, dbProgram)
			} else {
				require.NoError(t, err)
				require.NotNil(t, dbProgram)
				require.Equal(t, program.Name, dbProgram.Name)
				require.Equal(t, program.Path, dbProgram.Path)
				require.Equal(t, program.Description, dbProgram.Description)
			}
		})
	}
}

func TestSqliteDB_CleanupNonExistentPrograms(t *testing.T) {
	tempDir := t.TempDir()
	existingFile := filepath.Join(tempDir, "existing.txt")
	err := os.WriteFile(existingFile, []byte("test"), 0644)
	require.NoError(t, err)

	testCases := []struct {
		name                string
		build               func(db Database) error
		errExpected         bool
		expectedProgramsLen int
	}{
		{
			name: "cleanup non existent programs - removes non-existent",
			build: func(db Database) error {
				err = db.SaveProgram(&Program{
					Name:        "Non-Existent Program",
					Path:        "/test/path/does/not/exist",
					Description: "Test Description",
				})
				require.NoError(t, err)

				err = db.SaveProgram(&Program{
					Name:        "Existing Program",
					Path:        existingFile,
					Description: "Existing file",
				})
				require.NoError(t, err)

				return db.CleanupNonExistentPrograms()
			},
			errExpected:         false,
			expectedProgramsLen: 1,
		},
		{
			name: "cleanup with all existing programs - keeps all",
			build: func(db Database) error {
				err = db.SaveProgram(&Program{
					Name:        "Existing Program",
					Path:        existingFile,
					Description: "Existing file",
				})
				require.NoError(t, err)

				return db.CleanupNonExistentPrograms()
			},
			errExpected:         false,
			expectedProgramsLen: 1,
		},
		{
			name: "cleanup with empty database - no error",
			build: func(db Database) error {
				return db.CleanupNonExistentPrograms()
			},
			errExpected:         false,
			expectedProgramsLen: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, err := NewDatabase(InMemoryDbShared)
			require.NoError(t, err)
			require.NotNil(t, db)

			defer func() {
				if sqliteDB, ok := db.(*sqliteDB); ok {
					_ = sqliteDB.Close()
				}
			}()

			err = tc.build(db)
			if tc.errExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				programs, err := db.GetPrograms()
				require.NoError(t, err)
				require.Equal(t, tc.expectedProgramsLen, len(programs))
			}
		})
	}
}
