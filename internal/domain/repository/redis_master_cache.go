package repository

import (
	"context"

	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/entity"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/$GOFILE -package=mock_$GOPACKAGE

const (
	ARTIST_REDIS_KEY = "artist"
	SINGER_REDIS_KEY = "singer"
	UNIT_REDIS_KEY   = "unit"
	SONG_REDIS_KEY   = "song"
	CHART_REDIS_KEY  = "chart"
)

type RedisMasterCacheRepository interface {
	// Artist
	SetArtist(ctx context.Context, id int32, data *entity.Artist) error
	SetArtists(ctx context.Context, data []*entity.Artist) error
	GetArtistByID(ctx context.Context, id int32) (*entity.Artist, error)
	GetArtists(ctx context.Context) ([]*entity.Artist, error)

	// Singer
	SetSinger(ctx context.Context, id int32, data *entity.Singer) error
	SetSingers(ctx context.Context, data []*entity.Singer) error
	GetSingerByID(ctx context.Context, id int32) (*entity.Singer, error)
	GetSingers(ctx context.Context) ([]*entity.Singer, error)

	// Unit
	SetUnit(ctx context.Context, id int32, data *entity.Unit) error
	SetUnits(ctx context.Context, data []*entity.Unit) error
	GetUnitByID(ctx context.Context, id int32) (*entity.Unit, error)
	GetUnits(ctx context.Context) ([]*entity.Unit, error)

	// Song
	SetSong(ctx context.Context, id int32, data *entity.Song) error
	SetSongs(ctx context.Context, data []*entity.Song) error
	GetSongByID(ctx context.Context, id int32) (*entity.Song, error)
	GetSongs(ctx context.Context) ([]*entity.Song, error)

	// Chart
	SetChart(ctx context.Context, id int32, data *entity.Chart) error
	SetCharts(ctx context.Context, data []*entity.Chart) error
	GetChartByID(ctx context.Context, id int32) (*entity.Chart, error)
	GetCharts(ctx context.Context) ([]*entity.Chart, error)
}
