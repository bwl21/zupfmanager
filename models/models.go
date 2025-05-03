package models

// Song repräsentiert ein Notenblatt in der Datenbank.
type Song struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Filename string `json:"filename"`
	Genre    string `json:"genre,omitempty"` // Kann leer sein
}

// Project repräsentiert ein Projekt in der Datenbank.
type Project struct {
	ProjectID    int    `json:"project_id"`
	ProjectTitle string `json:"project_title"`
}

// ProjectSong repräsentiert die Zuordnung eines Songs zu einem Projekt
// mit zusätzlichen Attributen.
type ProjectSong struct {
	ProjectID  int    `json:"project_id"`
	SongID     int    `json:"song_id"`
	Priority   int    `json:"priority"`   // 1-4
	Difficulty string `json:"difficulty"` // z.B. "leicht", "mittel", "schwer"
	Comment    string `json:"comment"`
	// Optional: Felder aus dem Song-Struct für Anzeigezwecke hinzufügen
	SongTitle    string `json:"song_title,omitempty"`
	SongFilename string `json:"song_filename,omitempty"`
}
