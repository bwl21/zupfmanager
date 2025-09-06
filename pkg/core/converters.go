package core

import (
	"github.com/bwl21/zupfmanager/internal/ent"
)

// ProjectFromEnt converts an ent.Project to a domain Project
func ProjectFromEnt(entProject *ent.Project) *Project {
	if entProject == nil {
		return nil
	}
	
	return &Project{
		ID:                   entProject.ID,
		Title:                entProject.Title,
		ShortName:            entProject.ShortName,
		Config:               entProject.Config,
		AbcFileDirPreference: entProject.AbcFileDirPreference,
	}
}

// ProjectsFromEnt converts a slice of ent.Project to domain Projects
func ProjectsFromEnt(entProjects []*ent.Project) []*Project {
	projects := make([]*Project, len(entProjects))
	for i, entProject := range entProjects {
		projects[i] = ProjectFromEnt(entProject)
	}
	return projects
}

// SongFromEnt converts an ent.Song to a domain Song
func SongFromEnt(entSong *ent.Song) *Song {
	if entSong == nil {
		return nil
	}
	
	return &Song{
		ID:        entSong.ID,
		Title:     entSong.Title,
		Filename:  entSong.Filename,
		Genre:     entSong.Genre,
		Copyright: entSong.Copyright,
		Tocinfo:   entSong.Tocinfo,
	}
}

// SongsFromEnt converts a slice of ent.Song to domain Songs
func SongsFromEnt(entSongs []*ent.Song) []*Song {
	songs := make([]*Song, len(entSongs))
	for i, entSong := range entSongs {
		songs[i] = SongFromEnt(entSong)
	}
	return songs
}

// ProjectSongFromEnt converts an ent.ProjectSong to a domain ProjectSong
func ProjectSongFromEnt(entProjectSong *ent.ProjectSong) *ProjectSong {
	if entProjectSong == nil {
		return nil
	}
	
	projectSong := &ProjectSong{
		ID:         entProjectSong.ID,
		ProjectID:  entProjectSong.ProjectID,
		SongID:     entProjectSong.SongID,
		Difficulty: string(entProjectSong.Difficulty),
		Priority:   entProjectSong.Priority,
		Comment:    &entProjectSong.Comment,
	}
	
	// Add related entities if loaded
	if entProjectSong.Edges.Song != nil {
		projectSong.Song = SongFromEnt(entProjectSong.Edges.Song)
	}
	if entProjectSong.Edges.Project != nil {
		projectSong.Project = ProjectFromEnt(entProjectSong.Edges.Project)
	}
	
	return projectSong
}

// ProjectSongsFromEnt converts a slice of ent.ProjectSong to domain ProjectSongs
func ProjectSongsFromEnt(entProjectSongs []*ent.ProjectSong) []*ProjectSong {
	projectSongs := make([]*ProjectSong, len(entProjectSongs))
	for i, entProjectSong := range entProjectSongs {
		projectSongs[i] = ProjectSongFromEnt(entProjectSong)
	}
	return projectSongs
}
