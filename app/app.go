package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"zupfmanager/core" // Importiere das core-Package
	"zupfmanager/models"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

const dbFile = "zupfmanager.db"
const abcDir = "/Users/beweiche/Dropbox/RuthVeehNoten" // Pfad zu den ABC-Dateien

// App struct
type App struct {
	ctx context.Context
	db  *sql.DB // Datenbankverbindung speichern
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// --- Unsere Initialisierungslogik ---
	log.Println("Initialisiere Datenbank und synchronisiere Songs...")
	db, err := core.initDB(dbFile)
	if err != nil {
		// Im Fehlerfall sollten wir das der GUI irgendwie mitteilen oder loggen.
		// Fürs Erste loggen wir es und beenden ggf. die App? Oder laufen ohne DB weiter?
		// Hier loggen wir es und speichern nil, die Methoden müssen das prüfen.
		log.Printf("FATAL: Fehler beim Initialisieren der Datenbank: %v", err)
		// runtime.Quit(a.ctx) // Beendet die App
		a.db = nil // Stellen sicher, dass db nil ist
		return     // Verhindert Sync-Versuch ohne DB
	}
	a.db = db // Datenbankverbindung speichern

	log.Println("Datenbank erfolgreich initialisiert.")

	// Songs synchronisieren
	err = core.syncSongsWithDirectory(a.db, abcDir)
	if err != nil {
		// Auch hier Fehler behandeln - Loggen reicht fürs Erste
		log.Printf("WARNUNG: Fehler beim Synchronisieren der Songs: %v", err)
	} else {
		log.Printf("Songs aus Verzeichnis '%s' synchronisiert.\n", abcDir)
	}
	// --- Ende unserer Initialisierungslogik ---
}

// shutdown wird aufgerufen, wenn die App beendet wird.
func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		log.Println("Schließe Datenbankverbindung...")
		err := a.db.Close()
		if err != nil {
			log.Printf("Fehler beim Schließen der Datenbank: %v", err)
		}
	}
}

// --- Methoden, die vom Frontend aufgerufen werden können ---

// checkDB prüft, ob die DB-Verbindung besteht. Gibt einen Fehler zurück, wenn nicht.
func (a *App) checkDB() error {
	if a.db == nil {
		return fmt.Errorf("datenbankverbindung nicht initialisiert")
	}
	// Optional: Ping zur Überprüfung der Verbindung
	err := a.db.PingContext(a.ctx)
	if err != nil {
		return fmt.Errorf("datenbankverbindung verloren: %w", err)
	}
	return nil
}

// GetAllSongs gibt alle Songs zurück.
func (a *App) GetAllSongs() ([]models.Song, error) {
	if err := a.checkDB(); err != nil {
		return nil, err
	}
	return core.GetAllSongs(a.db)
}

// GetAllProjects gibt alle Projekte zurück.
func (a *App) GetAllProjects() ([]models.Project, error) {
	if err := a.checkDB(); err != nil {
		return nil, err
	}
	return core.GetAllProjects(a.db)
}

// CreateProject erstellt ein neues Projekt.
func (a *App) CreateProject(title string) (int64, error) {
	if err := a.checkDB(); err != nil {
		return 0, err
	}
	return core.CreateProject(a.db, title)
}

// AddSongToProject fügt einen Song zu einem Projekt hinzu.
// Beachte: Die Parameter müssen mit den JSON-Typen aus dem Frontend übereinstimmen.
// Ggf. müssen wir hier Typkonvertierungen vornehmen, wenn das Frontend z.B. alles als String sendet.
func (a *App) AddSongToProject(projectID int, songID int, priority int, difficulty string, comment string) error {
	if err := a.checkDB(); err != nil {
		return err
	}
	// Stelle sicher, dass Priority im gültigen Bereich liegt (obwohl die DB es auch prüft)
	if priority < 1 || priority > 4 {
		priority = 3 // Setze auf Default, wenn ungültig
	}
	return core.AddSongToProject(a.db, projectID, songID, priority, difficulty, comment)
}

// GetSongsInProject gibt die Songs eines bestimmten Projekts zurück.
func (a *App) GetSongsInProject(projectID int) ([]models.ProjectSong, error) {
	if err := a.checkDB(); err != nil {
		return nil, err
	}
	return core.GetSongsInProject(a.db, projectID)
}

// BuildProject startet den Build-Prozess für ein Projekt.
func (a *App) BuildProject(projectID int) error {
	if err := a.checkDB(); err != nil {
		return err
	}
	return core.BuildProject(a.db, projectID)
}

// SyncSongs startet die Synchronisation manuell (optional).
func (a *App) SyncSongs() error {
	if err := a.checkDB(); err != nil {
		return err
	}
	log.Println("Manuelle Synchronisation gestartet...")
	err := core.syncSongsWithDirectory(a.db, abcDir)
	if err != nil {
		log.Printf("Fehler bei manueller Synchronisation: %v", err)
		return fmt.Errorf("fehler beim Synchronisieren: %w", err)
	}
	log.Println("Manuelle Synchronisation abgeschlossen.")
	return nil
}

// UpdateSongGenre aktualisiert das Genre eines Songs.
func (a *App) UpdateSongGenre(songID int, genre string) error {
	if err := a.checkDB(); err != nil {
		return err
	}
	// Genre validieren (optional, DB macht es auch)
	validGenres := map[string]bool{"altchr": true, "neuchr": true, "klass": true, "saek": true, "": true} // Leer erlaubt?
	if !validGenres[genre] {
		return fmt.Errorf("ungültiges Genre: %s", genre)
	}

	// SQL ausführen
	stmt := "UPDATE songs SET genre = ? WHERE id = ?"
	_, err := a.db.ExecContext(a.ctx, stmt, genre, songID)
	if err != nil {
		return fmt.Errorf("fehler beim Aktualisieren des Genres für Song ID %d: %w", songID, err)
	}
	log.Printf("Genre für Song ID %d auf '%s' aktualisiert.\n", songID, genre)
	return nil
}

// RemoveSongFromProject entfernt einen Song aus einem Projekt.
func (a *App) RemoveSongFromProject(projectID int, songID int) error {
	if err := a.checkDB(); err != nil {
		return err
	}
	stmt := "DELETE FROM project_songs WHERE project_id = ? AND song_id = ?"
	res, err := a.db.ExecContext(a.ctx, stmt, projectID, songID)
	if err != nil {
		return fmt.Errorf("fehler beim Entfernen von Song ID %d aus Projekt ID %d: %w", songID, projectID, err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected > 0 {
		log.Printf("Song ID %d aus Projekt ID %d entfernt.\n", songID, projectID)
	} else {
		log.Printf("Keine Zuordnung gefunden für Song ID %d in Projekt ID %d.\n", songID, projectID)
	}
	return nil
}

// DeleteProject löscht ein Projekt und alle zugehörigen Song-Zuordnungen (via CASCADE).
func (a *App) DeleteProject(projectID int) error {
	if err := a.checkDB(); err != nil {
		return err
	}
	stmt := "DELETE FROM projects WHERE project_id = ?"
	res, err := a.db.ExecContext(a.ctx, stmt, projectID)
	if err != nil {
		return fmt.Errorf("fehler beim Löschen von Projekt ID %d: %w", projectID, err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected > 0 {
		log.Printf("Projekt ID %d gelöscht.\n", projectID)
	} else {
		log.Printf("Projekt ID %d nicht gefunden.\n", projectID)
	}
	return nil
}
