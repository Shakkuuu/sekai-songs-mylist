package repository

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/entity"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/repository"
	"github.com/cockroachdb/errors"
	"github.com/redis/go-redis/v9"
)

type redisMasterCacheRepository struct {
	rc *redis.Client
}

func NewRedisMasterCacheRepository(rc *redis.Client) repository.RedisMasterCacheRepository {
	return &redisMasterCacheRepository{
		rc: rc,
	}
}

// Artist
func (r *redisMasterCacheRepository) SetArtist(ctx context.Context, id int32, data *entity.Artist) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return errors.WithStack(err)
	}
	r.rc.Set(ctx, repository.ARTIST_REDIS_KEY+":"+strconv.Itoa(int(id)), jsonBytes, 0)
	return nil
}
func (r *redisMasterCacheRepository) SetArtists(ctx context.Context, data []*entity.Artist) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return errors.WithStack(err)
	}
	r.rc.Set(ctx, repository.ARTIST_REDIS_KEY+":all", jsonBytes, 0)
	return nil
}
func (r *redisMasterCacheRepository) GetArtistByID(ctx context.Context, id int32) (*entity.Artist, error) {
	data, err := r.rc.Get(ctx, repository.ARTIST_REDIS_KEY+":"+strconv.Itoa(int(id))).Bytes()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var artist *entity.Artist
	err = json.Unmarshal(data, &artist)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return artist, nil
}
func (r *redisMasterCacheRepository) GetArtists(ctx context.Context) ([]*entity.Artist, error) {
	data, err := r.rc.Get(ctx, repository.ARTIST_REDIS_KEY+":all").Bytes()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var artists []*entity.Artist
	err = json.Unmarshal(data, &artists)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return artists, nil
}

// Singer
func (r *redisMasterCacheRepository) SetSinger(ctx context.Context, id int32, data *entity.Singer) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return errors.WithStack(err)
	}
	r.rc.Set(ctx, repository.SINGER_REDIS_KEY+":"+strconv.Itoa(int(id)), jsonBytes, 0)
	return nil
}
func (r *redisMasterCacheRepository) SetSingers(ctx context.Context, data []*entity.Singer) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return errors.WithStack(err)
	}
	r.rc.Set(ctx, repository.SINGER_REDIS_KEY+":all", jsonBytes, 0)
	return nil
}
func (r *redisMasterCacheRepository) GetSingerByID(ctx context.Context, id int32) (*entity.Singer, error) {
	data, err := r.rc.Get(ctx, repository.SINGER_REDIS_KEY+":"+strconv.Itoa(int(id))).Bytes()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var singer *entity.Singer
	err = json.Unmarshal(data, &singer)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return singer, nil
}
func (r *redisMasterCacheRepository) GetSingers(ctx context.Context) ([]*entity.Singer, error) {
	data, err := r.rc.Get(ctx, repository.SINGER_REDIS_KEY+":all").Bytes()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var singers []*entity.Singer
	err = json.Unmarshal(data, &singers)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return singers, nil
}

// Unit
func (r *redisMasterCacheRepository) SetUnit(ctx context.Context, id int32, data *entity.Unit) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return errors.WithStack(err)
	}
	r.rc.Set(ctx, repository.UNIT_REDIS_KEY+":"+strconv.Itoa(int(id)), jsonBytes, 0)
	return nil
}
func (r *redisMasterCacheRepository) SetUnits(ctx context.Context, data []*entity.Unit) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return errors.WithStack(err)
	}
	r.rc.Set(ctx, repository.UNIT_REDIS_KEY+":all", jsonBytes, 0)
	return nil
}
func (r *redisMasterCacheRepository) GetUnitByID(ctx context.Context, id int32) (*entity.Unit, error) {
	data, err := r.rc.Get(ctx, repository.UNIT_REDIS_KEY+":"+strconv.Itoa(int(id))).Bytes()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var unit *entity.Unit
	err = json.Unmarshal(data, &unit)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return unit, nil
}
func (r *redisMasterCacheRepository) GetUnits(ctx context.Context) ([]*entity.Unit, error) {
	data, err := r.rc.Get(ctx, repository.UNIT_REDIS_KEY+":all").Bytes()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var units []*entity.Unit
	err = json.Unmarshal(data, &units)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return units, nil
}

// Song
func (r *redisMasterCacheRepository) SetSong(ctx context.Context, id int32, data *entity.Song) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return errors.WithStack(err)
	}
	r.rc.Set(ctx, repository.SONG_REDIS_KEY+":"+strconv.Itoa(int(id)), jsonBytes, 0)
	return nil
}
func (r *redisMasterCacheRepository) SetSongs(ctx context.Context, data []*entity.Song) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return errors.WithStack(err)
	}
	r.rc.Set(ctx, repository.SONG_REDIS_KEY+":all", jsonBytes, 0)
	return nil
}
func (r *redisMasterCacheRepository) GetSongByID(ctx context.Context, id int32) (*entity.Song, error) {
	data, err := r.rc.Get(ctx, repository.SONG_REDIS_KEY+":"+strconv.Itoa(int(id))).Bytes()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var song *entity.Song
	err = json.Unmarshal(data, &song)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return song, nil
}
func (r *redisMasterCacheRepository) GetSongs(ctx context.Context) ([]*entity.Song, error) {
	data, err := r.rc.Get(ctx, repository.SONG_REDIS_KEY+":all").Bytes()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var songs []*entity.Song
	err = json.Unmarshal(data, &songs)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return songs, nil
}

// Chart
func (r *redisMasterCacheRepository) SetChart(ctx context.Context, id int32, data *entity.Chart) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return errors.WithStack(err)
	}
	r.rc.Set(ctx, repository.CHART_REDIS_KEY+":"+strconv.Itoa(int(id)), jsonBytes, 0)
	return nil
}
func (r *redisMasterCacheRepository) SetCharts(ctx context.Context, data []*entity.Chart) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return errors.WithStack(err)
	}
	r.rc.Set(ctx, repository.CHART_REDIS_KEY+":all", jsonBytes, 0)
	return nil
}
func (r *redisMasterCacheRepository) GetChartByID(ctx context.Context, id int32) (*entity.Chart, error) {
	data, err := r.rc.Get(ctx, repository.CHART_REDIS_KEY+":"+strconv.Itoa(int(id))).Bytes()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var chart *entity.Chart
	err = json.Unmarshal(data, &chart)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return chart, nil
}
func (r *redisMasterCacheRepository) GetCharts(ctx context.Context) ([]*entity.Chart, error) {
	data, err := r.rc.Get(ctx, repository.CHART_REDIS_KEY+":all").Bytes()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var charts []*entity.Chart
	err = json.Unmarshal(data, &charts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return charts, nil
}
