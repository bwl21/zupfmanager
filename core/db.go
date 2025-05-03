package core

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3" // SQLite driver needed for sql.Open
	"zupfmanager/models"
)

// initDB initialisiert die SQLite-Datenbank und erstellt die Tabellen, falls sie nicht existieren.
func initDB(filepath string) (*sql.DB, error) {
	// Erstellt die DB-Datei, falls sie nicht existiert
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		file, err := os.Create(filepath)
		if err != nil {
			return nil, fmt.Errorf("konnte DB-Datei nicht erstellen: %w", err)
		}
		file.Close()
	}

	db, err := sql.Open("sqlite3", filepath+"?_foreign_keys=on") // Foreign keys aktivieren
	if err != nil {
		return nil, fmt.Errorf("konnte Datenbank nicht öffnen: %w", err)
	}

	// SQL-Statements zum Erstellen der Tabellen
	createSongsTableSQL := `CREATE TABLE IF NOT EXISTS songs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		filename TEXT NOT NULL UNIQUE,
		genre TEXT CHECK(genre IN ('altchr', 'neuchr', 'klass', 'saek'))
	);`

	createProjectsTableSQL := `CREATE TABLE IF NOT EXISTS projects (
		project_id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_title TEXT NOT NULL UNIQUE
	);`

	// Verbindungstabelle für die m:n-Beziehung zwischen Songs und Projects
	createProjectSongsTableSQL := `CREATE TABLE IF NOT EXISTS project_songs (
		project_id INTEGER NOT NULL,
		song_id INTEGER NOT NULL,
		priority INTEGER DEFAULT 3 CHECK(priority BETWEEN 1 AND 4),
		difficulty TEXT,
		comment TEXT,
		PRIMARY KEY (project_id, song_id),
		FOREIGN KEY (project_id) REFERENCES projects(project_id) ON DELETE CASCADE,
		FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE CASCADE
	);`

	// Tabellen erstellen
	statements := []string{createSongsTableSQL, createProjectsTableSQL, createProjectSongsTableSQL}
	for _, stmt := range statements {
		_, err = db.Exec(stmt)
		if err != nil {
			db.Close() // Wichtig: DB schließen bei Fehler
			return nil, fmt.Errorf("konnte Tabelle nicht erstellen: %w\nStatement: %s", err, stmt)
		}
	}

	return db, nil
}

// syncSongsWithDirectory durchsucht ein Verzeichnis nach .abc Dateien und fügt sie zur DB hinzu.
func syncSongsWithDirectory(db *sql.DB, dirPath string) error {
	log.Printf("Starte Synchronisation für Verzeichnis: %s\n", dirPath)
	foundFiles := 0
	addedSongs := 0

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Fehler beim Zugriff auf Pfad %q: %v\n", path, err)
			// Bei bestimmten Fehlern (z.B. Permission denied) kann es sinnvoll sein, hier abzubrechen
			// return err
			return filepath.SkipDir // Überspringe dieses Verzeichnis bei Fehler
		}
		if info.IsDir() {
			// Ignoriere versteckte Verzeichnisse (optional)
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil // Gehe in Unterverzeichnisse
		}

		if strings.HasSuffix(strings.ToLower(info.Name()), ".abc") {
			foundFiles++
			filename := info.Name()
			// log.Printf("Gefunden: %s\n", filename) // Weniger verbose

			// Prüfen, ob Song bereits existiert
			var existingID int
			errCheck := db.QueryRow("SELECT id FROM songs WHERE filename = ?", filename).Scan(&existingID)
			if errCheck == nil {
				// Song existiert bereits, überspringen
				return nil
			} else if errCheck != sql.ErrNoRows {
				// Anderer Fehler bei der DB-Abfrage
				log.Printf("Fehler bei der Abfrage für Song '%s': %v\n", filename, errCheck)
				return fmt.Errorf("datenbankfehler bei Abfrage für %s: %w", filename, errCheck)
			}

			// Song existiert nicht, Titel extrahieren und einfügen
			title, errExtract := extractTitleFromABC(path)
			if errExtract != nil {
				log.Printf("Konnte Titel für '%s' nicht extrahieren: %v. Überspringe Datei.\n", filename, errExtract)
				return nil // Fehler beim Extrahieren ignorieren und weitermachen
			}
			if title == "" {
				log.Printf("Kein T: Feld in '%s' gefunden. Verwende Dateinamen als Titel.\n", filename)
				title = strings.TrimSuffix(filename, filepath.Ext(filename)) // Fallback: Dateiname ohne Endung
			}

			// Song in DB einfügen
			insertSQL := "INSERT INTO songs(title, filename) VALUES (?, ?)"
			_, errInsert := db.Exec(insertSQL, title, filename)
			if errInsert != nil {
				// Prüfen auf UNIQUE constraint violation (kann passieren, wenn gleichzeitig hinzugefügt wird oder Logikfehler)
				if strings.Contains(errInsert.Error(), "UNIQUE constraint failed") {
					log.Printf("Song '%s' existiert bereits (UNIQUE constraint), übersprungen.\n", filename)
					return nil
				}
				log.Printf("Fehler beim Einfügen von Song '%s' (Titel: '%s'): %v\n", filename, title, errInsert)
				// Bei anderen Fehlern weitermachen oder abbrechen? Hier: weitermachen
				return nil
			}
			addedSongs++
			log.Printf("Song hinzugefügt: '%s' (Titel: '%s')\n", filename, title)
		}
		return nil // Weiter zum nächsten Eintrag
	})

	if err != nil {
		// Fehler, die von filepath.Walk zurückgegeben werden (z.B. initiale Zugriffsfehler)
		return fmt.Errorf("fehler beim Durchlaufen des Verzeichnisses %s: %w", dirPath, err)
	}

	log.Printf("Synchronisation abgeschlossen. %d .abc Dateien gefunden, %d neue Songs hinzugefügt.\n", foundFiles, addedSongs)
	return nil
}

// extractTitleFromABC liest eine ABC-Datei und extrahiert den Titel aus dem T: Feld.
func extractTitleFromABC(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("konnte Datei nicht öffnen %s: %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// ABC standard: T: field must be at the beginning of the line
		if strings.HasPrefix(line, "T:") {
			title := strings.TrimSpace(line[2:])
			// Nimm nur die erste T: Zeile, falls mehrere vorhanden sind
			if title != "" {
				return title, nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("fehler beim Lesen der Datei %s: %w", filePath, err)
	}

	// Kein T: Feld gefunden
	return "", nil
}

// --- Projekt Management Funktionen ---

// createProject fügt ein neues Projekt hinzu und gibt dessen ID zurück.
func createProject(db *sql.DB, title string) (int64, error) {
	res, err := db.Exec("INSERT INTO projects(project_title) VALUES (?)", title)
	if err != nil {
		// Prüfen auf UNIQUE constraint
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return 0, fmt.Errorf("projekt mit Titel '%s' existiert bereits", title)
		}
		return 0, fmt.Errorf("fehler beim Erstellen des Projekts '%s': %w", title, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("fehler beim Abrufen der ID für Projekt '%s': %w", title, err)
	}
	log.Printf("Projekt '%s' (ID: %d) erstellt.\n", title, id)
	return id, nil
}

// addSongToProject fügt einen Song zu einem Projekt hinzu.
func addSongToProject(db *sql.DB, projectID, songID int, priority int, difficulty, comment string) error {
	// Optional: Prüfen, ob Projekt und Song existieren
	// ...

	sqlStmt := `INSERT INTO project_songs(project_id, song_id, priority, difficulty, comment)
				VALUES (?, ?, ?, ?, ?)
				ON CONFLICT(project_id, song_id) DO UPDATE SET
				priority=excluded.priority, difficulty=excluded.difficulty, comment=excluded.comment;`
	_, err := db.Exec(sqlStmt, projectID, songID, priority, difficulty, comment)
	if err != nil {
		return fmt.Errorf("fehler beim Hinzufügen/Aktualisieren von Song ID %d zu Projekt ID %d: %w", songID, projectID, err)
	}
	log.Printf("Song ID %d zu Projekt ID %d hinzugefügt/aktualisiert.\n", songID, projectID)
	return nil
}

// getSongsInProject ruft alle Songs ab, die zu einem bestimmten Projekt gehören.
func getSongsInProject(db *sql.DB, projectID int) ([]ProjectSong, error) {
	rows, err := db.Query(`
		SELECT ps.project_id, ps.song_id, ps.priority, ps.difficulty, ps.comment, s.title, s.filename
		FROM project_songs ps
		JOIN songs s ON ps.song_id = s.id
		WHERE ps.project_id = ?
		ORDER BY ps.priority ASC, s.title ASC;
	`, projectID)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Abrufen der Songs für Projekt ID %d: %w", projectID, err)
	}
	defer rows.Close()

	var songs []ProjectSong
	for rows.Next() {
		var ps ProjectSong
		err := rows.Scan(&ps.ProjectID, &ps.SongID, &ps.Priority, &ps.Difficulty, &ps.Comment, &ps.SongTitle, &ps.SongFilename)
		if err != nil {
			return nil, fmt.Errorf("fehler beim Scannen der Song-Daten für Projekt ID %d: %w", projectID, err)
		}
		songs = append(songs, ps)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("fehler nach dem Iterieren über die Song-Zeilen für Projekt ID %d: %w", projectID, err)
	}

	return songs, nil
}

// getAllProjects ruft alle Projekte ab.
func getAllProjects(db *sql.DB) ([]Project, error) {
	rows, err := db.Query("SELECT project_id, project_title FROM projects ORDER BY project_title ASC")
	if err != nil {
		return nil, fmt.Errorf("fehler beim Abrufen aller Projekte: %w", err)
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ProjectID, &p.ProjectTitle); err != nil {
			return nil, fmt.Errorf("fehler beim Scannen der Projekt-Daten: %w", err)
		}
		projects = append(projects, p)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("fehler nach dem Iterieren über die Projekt-Zeilen: %w", err)
	}
	return projects, nil
}

// getAllSongs ruft alle Songs aus der Datenbank ab.
func getAllSongs(db *sql.DB) ([]Song, error) {
	rows, err := db.Query("SELECT id, title, filename, genre FROM songs ORDER BY title ASC")
	if err != nil {
		return nil, fmt.Errorf("fehler beim Abrufen aller Songs: %w", err)
	}
	defer rows.Close()

	var songs []Song
	for rows.Next() {
		var s Song
		var genre sql.NullString // Behandelt NULL-Werte für Genre
		if err := rows.Scan(&s.ID, &s.Title, &s.Filename, &genre); err != nil {
			return nil, fmt.Errorf("fehler beim Scannen der Song-Daten: %w", err)
		}
		if genre.Valid {
			s.Genre = genre.String
		} else {
			s.Genre = "" // Setze auf leeren String, wenn NULL
		}
		songs = append(songs, s)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("fehler nach dem Iterieren über die Song-Zeilen: %w", err)
	}
	return songs, nil
}

// --- Weitere DB-Funktionen (z.B. Update, Delete) können hier hinzugefügt werden ---
