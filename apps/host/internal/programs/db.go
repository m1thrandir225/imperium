package programs

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
)

type Database interface {
	SaveProgram(program *Program) error
	GetPrograms() ([]*Program, error)
	GetProgramByID(id string) (*Program, error)
	GetProgramByPath(path string) (*Program, error)
	CleanupNonExistentPrograms() error
}

type sqliteDB struct {
	db *sql.DB
}

// NewDatabase returns a new instance of a local database
func NewDatabase(dbPath string) (Database, error) {
	return newSqliteDB(dbPath)
}

func newSqliteDB(dbPath string) (*sqliteDB, error) {
	dir := filepath.Dir(dbPath)

	// dont create a directory if the path is in-memory
	if dbPath != ":memory" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("failed to create directory: %v", err)
			return nil, ErrDatabaseDirectoryCreationFailed
		}
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("failed to open database: %v", err)
		return nil, ErrDatabaseOpenFailed
	}

	if err := db.Ping(); err != nil {
		log.Printf("failed to ping database: %v", err)
		return nil, ErrDatabasePingFailed
	}

	programDB := &sqliteDB{
		db: db,
	}
	if err := programDB.initTables(); err != nil {
		log.Printf("failed to initialize tables: %v", err)
		return nil, ErrDatabaseInitializeTablesFailed
	}
	return programDB, nil
}

func (pdb *sqliteDB) initTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS programs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		path TEXT UNIQUE NOT NULL,
		description TEXT,
		last_modified DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_programs_path ON programs(path);
	CREATE INDEX IF NOT EXISTS idx_programs_name ON programs(name);
	`

	_, err := pdb.db.Exec(query)
	return err
}

func (pdb *sqliteDB) Close() error {
	return pdb.db.Close()
}

func (pdb *sqliteDB) SaveProgram(program *Program) error {
	query := `
	INSERT OR REPLACE INTO programs (name, path, description, last_modified, updated_at)
	VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
	`

	fileInfo, err := os.Stat(program.Path)
	lastModified := time.Time{}
	if err == nil {
		lastModified = fileInfo.ModTime()
	}

	result, err := pdb.db.Exec(query, program.Name, program.Path, program.Description, lastModified)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	program.ID = strconv.FormatInt(id, 10)
	return err
}

func (pdb *sqliteDB) GetPrograms() ([]*Program, error) {
	query := `SELECT id, name, path, description FROM programs ORDER BY name`

	rows, err := pdb.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get programs: %w", err)
	}

	var programs []*Program
	for rows.Next() {
		program := &Program{}
		err := rows.Scan(&program.ID, &program.Name, &program.Path, &program.Description)
		if err != nil {
			log.Printf("Error scanning program: %v", err)
			continue
		}
		programs = append(programs, program)
	}
	return programs, nil
}

func (pdb *sqliteDB) GetProgramByID(id string) (*Program, error) {
	query := `SELECT id, name, path, description FROM programs WHERE id = ?`
	program := &Program{}
	err := pdb.db.QueryRow(query, id).Scan(&program.ID, &program.Name, &program.Path, &program.Description)
	if err != nil {
		return nil, err
	}
	return program, nil
}

func (pdb *sqliteDB) GetProgramByPath(path string) (*Program, error) {
	query := `SELECT id, name, path, description FROM programs WHERE path = ?`

	program := &Program{}
	err := pdb.db.QueryRow(query, path).Scan(&program.ID, &program.Name, &program.Path, &program.Description)
	if err != nil {
		return nil, err
	}

	return program, nil
}

func (pdb *sqliteDB) CleanupNonExistentPrograms() error {
	query := `SELECT path FROM programs`
	rows, err := pdb.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			continue
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			// Program no longer exists, remove from database
			deleteQuery := `DELETE FROM programs WHERE path = ?`
			if _, err := pdb.db.Exec(deleteQuery, path); err != nil {
				log.Printf("Failed to delete non-existent program %s: %v", path, err)
			}
		}
	}

	return nil
}
