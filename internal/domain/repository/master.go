package repository

import (
	"context"
	"time"

	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/entity"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/gen/enums"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/infrastructure/db/sqlcgen"
	"github.com/cockroachdb/errors"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/$GOFILE -package=mock_$GOPACKAGE

var (
	ErrNotFound = errors.New("not found")
)

type MasterRepository interface {
	// Artist
	ListArtists(ctx context.Context) ([]*entity.Artist, error)
	GetArtistByID(ctx context.Context, id int32) (*entity.Artist, error)
	CreateArtist(ctx context.Context, name, kana string) (*entity.Artist, error)
	ExistsArtist(ctx context.Context, id int32) (bool, error)
	// Singer
	ListSingers(ctx context.Context) ([]*entity.Singer, error)
	GetSingerByID(ctx context.Context, id int32) (*entity.Singer, error)
	CreateSinger(ctx context.Context, name string) (*entity.Singer, error)
	ExistsSinger(ctx context.Context, id int32) (bool, error)
	// Unit
	ListUnits(ctx context.Context) ([]*entity.Unit, error)
	GetUnitByID(ctx context.Context, id int32) (*entity.Unit, error)
	CreateUnit(ctx context.Context, name string) (*entity.Unit, error)
	ExistsUnit(ctx context.Context, id int32) (bool, error)
	// VocalPattern
	CreateVocalPattern(ctx context.Context, songID int32, name string) (*sqlcgen.VocalPattern, error)
	// VocalPatternSinger
	CreateVocalPatternSinger(ctx context.Context, vocalPatternID, singerID, position int32) (*sqlcgen.VocalPatternSinger, error)
	// SongUnit
	CreateSongUnit(ctx context.Context, songID, unitID int32) (*sqlcgen.SongUnit, error)
	// SongMusicVideoType
	CreateSongMusicVideoType(ctx context.Context, songID int32, musicVideoType enums.MusicVideoType) (*sqlcgen.SongMusicVideoType, error)
	// Song
	ListSongs(ctx context.Context) ([]*entity.Song, error)
	GetSongByID(ctx context.Context, id int32) (*entity.Song, error)
	CreateSong(ctx context.Context,
		name, kana string,
		lyrics_id, music_id, arrangement_id int32,
		thumbnail, originalVideo string,
		releaseTime time.Time, deleted bool,
	) (*sqlcgen.Song, error)
	ExistsSong(ctx context.Context, id int32) (bool, error)
	// Chart
	ListCharts(ctx context.Context) ([]*entity.Chart, error)
	GetChartByID(ctx context.Context, id int32) (*entity.Chart, error)
	CreateChart(ctx context.Context, songID, difficultyType, level int32, chartViewLink string) (*sqlcgen.Chart, error)
	ExistsChart(ctx context.Context, id int32) (bool, error)
}
