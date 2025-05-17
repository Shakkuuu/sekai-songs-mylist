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
	Song         Song
	DifficultyType enums.DifficultyType
	Level          int32
	ChartViewLink  string
}

type Singer struct {
	ID   int32
	Name string
	Position int32
}

type Song struct {
	ID            int32
	Name          string
	Kana          string
	Lyrics Artist
	Music Artist
	Arrangement Artist
	Thumbnail     string
	OriginalVideo string
	ReleaseTime   time.Time
	Deleted       bool
	VocalPatterns []*VocalPattern
	MusicVideoTypes []enums.MusicVideoType
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
	Name   string
	Singers []*Singer
	Units []*Unit
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
