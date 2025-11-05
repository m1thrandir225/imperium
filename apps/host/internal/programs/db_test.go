package programs

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
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
			dbPath:      ":memory:",
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
	db, err := NewDatabase(":memory")
	require.NoError(t, err)
	require.NotNil(t, db)

	defer func() {
		if sqliteDB, ok := db.(*sqliteDB); ok {
			sqliteDB.Close()
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

func TestSqliteDB_SaveProgram(t *test)
