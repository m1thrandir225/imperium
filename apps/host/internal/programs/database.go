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

type ProgramDB struct {
	db *sql.DB
}

func NewProgramDB(dbPath string) (*ProgramDB, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	programDB := &ProgramDB{
		db: db,
	}
	if err := programDB.initTables(); err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}
	return programDB, nil
}

func (pdb *ProgramDB) initTables() error {
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

func (pdb *ProgramDB) Close() error {
	return pdb.db.Close()
}

func (pdb *ProgramDB) SaveProgram(program *Program) error {
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

func (pdb *ProgramDB) GetPrograms() ([]*Program, error) {
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

func (pdb *ProgramDB) GetProgramByID(id string) (*Program, error) {
	query := `SELECT id, name, path, description FROM programs WHERE id = ?`
	program := &Program{}
	err := pdb.db.QueryRow(query, id).Scan(&program.ID, &program.Name, &program.Path, &program.Description)
	if err != nil {
		return nil, err
	}
	return program, nil
}

func (pdb *ProgramDB) GetProgramByPath(path string) (*Program, error) {
	query := `SELECT id, name, path, description FROM programs WHERE path = ?`

	program := &Program{}
	err := pdb.db.QueryRow(query, path).Scan(&program.ID, &program.Name, &program.Path, &program.Description)
	if err != nil {
		return nil, err
	}

	return program, nil
}

func (pdb *ProgramDB) CleanupNonExistentPrograms() error {
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
