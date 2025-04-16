package entity

import (
	"time"

	"github.com/Shakkuuu/sekai-songs-mylist/internal/gen/enums"
)

type Artist struct {
	ID   int32
	Name string
	Kana string
}

type Chart struct {
	ID             int32
	SongID         int32
	DifficultyType enums.DifficultyType
	Level          int32
	ChartViewLink  string
}

type Singer struct {
	ID   int32
	Name string
}

type Song struct {
	ID            int32
	Name          string
	Kana          string
	LyricsID      int32
	MusicID       int32
	ArrangementID int32
	Thumbnail     string
	OriginalVideo string
	ReleaseTime   time.Time
	Deleted       bool
}

type SongMusicVideoType struct {
	ID             int32
	SongID         int32
	MusicVideoType enums.MusicVideoType
}

type Unit struct {
	ID   int32
	Name string
}

type VocalPattern struct {
	ID     int32
	SongID int32
	Name   string
}

type VocalPatternSinger struct {
	ID             int32
	VocalPatternID int32
	SingerID       int32
	Position       int32
}

type VocalPatternUnit struct {
	ID             int32
	VocalPatternID int32
	UnitID         int32
}
