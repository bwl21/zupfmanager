package core

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"zupfmanager/models" // Expliziter Modul-Pfad
)

const exportsDir = "exports" // Unterverzeichnis für Exporte

// ProjectExportData definiert die Struktur für die JSON-Exportdatei.
type ProjectExportData struct {
	ProjectTitle string        `json:"project_title"`
	ProjectID    int           `json:"project_id"`
	Songs        []ProjectSong `json:"songs"` // Verwendet ProjectSong aus models.go
}

// buildProject führt den Build-Prozess für ein gegebenes Projekt aus.
func buildProject(db *sql.DB, projectID int) error {
	log.Printf("Starte Build für Projekt ID: %d\n", projectID)

	// 1. Projektdaten abrufen
	var projectTitle string
	err := db.QueryRow("SELECT project_title FROM projects WHERE project_id = ?", projectID).Scan(&projectTitle)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("projekt mit ID %d nicht gefunden", projectID)
		}
		return fmt.Errorf("fehler beim Abrufen des Projekttitels für ID %d: %w", projectID, err)
	}

	songs, err := getSongsInProject(db, projectID)
	if err != nil {
		return fmt.Errorf("fehler beim Abrufen der Songs für Projekt ID %d: %w", projectID, err)
	}

	if len(songs) == 0 {
		log.Printf("Projekt '%s' (ID: %d) enthält keine Songs. Build übersprungen.\n", projectTitle, projectID)
		return nil // Kein Fehler, aber nichts zu tun
	}

	exportData := ProjectExportData{
		ProjectTitle: projectTitle,
		ProjectID:    projectID,
		Songs:        songs,
	}

	// 2. Als JSON exportieren
	jsonFilePath, err := exportProjectToJson(exportData)
	if err != nil {
		return fmt.Errorf("fehler beim Exportieren des Projekts %d nach JSON: %w", projectID, err)
	}
	log.Printf("Projekt erfolgreich nach '%s' exportiert.\n", jsonFilePath)

	// 3. Externes Programm aufrufen
	err = runZupfgenerator(jsonFilePath)
	if err != nil {
		return fmt.Errorf("fehler beim Ausführen von zupfgenerator für Projekt %d: %w", projectID, err)
	}
	log.Printf("zupfgenerator erfolgreich für '%s' ausgeführt.\n", jsonFilePath)

	log.Printf("Build für Projekt '%s' (ID: %d) erfolgreich abgeschlossen.\n", projectTitle, projectID)
	return nil
}

// exportProjectToJson schreibt die Projektdaten in eine JSON-Datei.
func exportProjectToJson(data ProjectExportData) (string, error) {
	// Sicherstellen, dass das Exportverzeichnis existiert
	if err := os.MkdirAll(exportsDir, 0755); err != nil {
		return "", fmt.Errorf("konnte Exportverzeichnis '%s' nicht erstellen: %w", exportsDir, err)
	}

	// Dateinamen generieren (z.B. "Projekt Titel.json")
	// Ersetze ungültige Zeichen für Dateinamen
	safeTitle := strings.ReplaceAll(data.ProjectTitle, "/", "-")
	safeTitle = strings.ReplaceAll(safeTitle, "\\", "-")
	// Weitere Ersetzungen könnten nötig sein, je nach OS
	jsonFilename := fmt.Sprintf("%s.json", safeTitle)
	jsonFilePath := filepath.Join(exportsDir, jsonFilename)

	// Daten nach JSON marshallen
	jsonData, err := json.MarshalIndent(data, "", "  ") // Mit Einrückung für Lesbarkeit
	if err != nil {
		return "", fmt.Errorf("fehler beim Marshalling der Projektdaten nach JSON: %w", err)
	}

	// JSON-Daten in Datei schreiben
	err = os.WriteFile(jsonFilePath, jsonData, 0644)
	if err != nil {
		return "", fmt.Errorf("fehler beim Schreiben der JSON-Datei '%s': %w", jsonFilePath, err)
	}

	return jsonFilePath, nil
}

// runZupfgenerator ruft das externe Programm 'zupfgenerator' auf.
func runZupfgenerator(jsonFilePath string) error {
	cmd := exec.Command("zupfgenerator", jsonFilePath) // Annahme: zupfgenerator nimmt JSON-Pfad als Argument

	log.Printf("Führe Befehl aus: %s\n", cmd.String())

	// Standard Output und Standard Error des Befehls abfangen
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("Fehler beim Ausführen von zupfgenerator:\n--- Output ---\n%s\n--------------\n", string(output))
		return fmt.Errorf("zupfgenerator fehlgeschlagen: %w", err)
	}

	log.Printf("zupfgenerator Output:\n--- Output ---\n%s\n--------------\n", string(output))
	return nil
}
