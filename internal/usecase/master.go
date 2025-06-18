package usecase

import (
	"context"
	"time"

	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/entity"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/repository"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/gen/enums"
	"github.com/cockroachdb/errors"
	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/$GOFILE -package=mock_$GOPACKAGE
//go:generate gotests -w -all $GOFILE

var (
	ErrInvalidArgument = errors.New("invalid argument")
)

type MasterUsecase interface {
	// Artist
	ListArtists(ctx context.Context) ([]*entity.Artist, error)
	GetArtistByID(ctx context.Context, id int32) (*entity.Artist, error)
	CreateArtist(ctx context.Context, name, kana string) error
	// Singer
	ListSingers(ctx context.Context) ([]*entity.Singer, error)
	GetSingerByID(ctx context.Context, id int32) (*entity.Singer, error)
	CreateSinger(ctx context.Context, name string) error
	// Unit
	ListUnits(ctx context.Context) ([]*entity.Unit, error)
	GetUnitByID(ctx context.Context, id int32) (*entity.Unit, error)
	CreateUnit(ctx context.Context, name string) error
	// VocalPattern
	CreateVocalPattern(ctx context.Context, songID int32, name string, singerIDs, singerPositions []int32) error
	// Song
	ListSongs(ctx context.Context) ([]*entity.Song, error)
	GetSongByID(ctx context.Context, id int32) (*entity.Song, error)
	CreateSong(
		ctx context.Context,
		name, kana string,
		lyrics_id, music_id, arrangement_id int32,
		thumbnail, originalVideo string,
		releaseTime time.Time, deleted bool,
		unitIDs []int32,
		musicVideoTypes []enums.MusicVideoType,
	) error
	// Chart
	ListCharts(ctx context.Context) ([]*entity.Chart, error)
	GetChartByID(ctx context.Context, id int32) (*entity.Chart, error)
	CreateChart(
		ctx context.Context,
		songID, difficultyType, level int32,
		chartViewLink string,
	) error
}

type masterUsecase struct {
	masterRepo           repository.MasterRepository
	redisMasterCacheRepo repository.RedisMasterCacheRepository
}

func NewMasterUsecase(repo repository.MasterRepository, redisMasterCacheRepo repository.RedisMasterCacheRepository) MasterUsecase {
	return &masterUsecase{
		masterRepo:           repo,
		redisMasterCacheRepo: redisMasterCacheRepo,
	}
}

// Artist
func (u *masterUsecase) ListArtists(ctx context.Context) ([]*entity.Artist, error) {
	artists, err := u.redisMasterCacheRepo.GetArtists(ctx)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			artists, err := u.masterRepo.ListArtists(ctx)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			if artists == nil {
				return []*entity.Artist{}, nil
			}
			if err := u.redisMasterCacheRepo.SetArtists(ctx, artists); err != nil {
				return nil, errors.WithStack(err)
			}
			return artists, nil
		}
		return nil, errors.WithStack(err)
	}

	return artists, nil
}

func (u *masterUsecase) GetArtistByID(ctx context.Context, id int32) (*entity.Artist, error) {
	artist, err := u.redisMasterCacheRepo.GetArtistByID(ctx, id)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			artist, err := u.masterRepo.GetArtistByID(ctx, id)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			if err := u.redisMasterCacheRepo.SetArtist(ctx, id, artist); err != nil {
				return nil, errors.WithStack(err)
			}
			return artist, nil
		}
		return nil, errors.WithStack(err)
	}

	return artist, nil
}

func (u *masterUsecase) CreateArtist(ctx context.Context, name, kana string) error {
	_, err := u.masterRepo.CreateArtist(ctx, name, kana)
	if err != nil {
		return errors.WithStack(err)
	}

	artists, err := u.masterRepo.ListArtists(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	if artists == nil {
		return nil
	}
	if err := u.redisMasterCacheRepo.SetArtists(ctx, artists); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Singer
func (u *masterUsecase) ListSingers(ctx context.Context) ([]*entity.Singer, error) {
	singers, err := u.redisMasterCacheRepo.GetSingers(ctx)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			singers, err := u.masterRepo.ListSingers(ctx)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			if singers == nil {
				return []*entity.Singer{}, nil
			}
			if err := u.redisMasterCacheRepo.SetSingers(ctx, singers); err != nil {
				return nil, errors.WithStack(err)
			}
			return singers, nil
		}
		return nil, errors.WithStack(err)
	}

	return singers, nil
}

func (u *masterUsecase) GetSingerByID(ctx context.Context, id int32) (*entity.Singer, error) {
	singer, err := u.redisMasterCacheRepo.GetSingerByID(ctx, id)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			singer, err := u.masterRepo.GetSingerByID(ctx, id)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			if err := u.redisMasterCacheRepo.SetSinger(ctx, id, singer); err != nil {
				return nil, errors.WithStack(err)
			}
			return singer, nil
		}
		return nil, errors.WithStack(err)
	}

	return singer, nil
}

func (u *masterUsecase) CreateSinger(ctx context.Context, name string) error {
	_, err := u.masterRepo.CreateSinger(ctx, name)
	if err != nil {
		return errors.WithStack(err)
	}

	singers, err := u.masterRepo.ListSingers(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	if singers == nil {
		return nil
	}
	if err := u.redisMasterCacheRepo.SetSingers(ctx, singers); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Unit
func (u *masterUsecase) ListUnits(ctx context.Context) ([]*entity.Unit, error) {
	units, err := u.redisMasterCacheRepo.GetUnits(ctx)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			units, err := u.masterRepo.ListUnits(ctx)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			if units == nil {
				return []*entity.Unit{}, nil
			}
			if err := u.redisMasterCacheRepo.SetUnits(ctx, units); err != nil {
				return nil, errors.WithStack(err)
			}
			return units, nil
		}
		return nil, errors.WithStack(err)
	}

	return units, nil
}

func (u *masterUsecase) GetUnitByID(ctx context.Context, id int32) (*entity.Unit, error) {
	unit, err := u.redisMasterCacheRepo.GetUnitByID(ctx, id)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			unit, err := u.masterRepo.GetUnitByID(ctx, id)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			if err := u.redisMasterCacheRepo.SetUnit(ctx, id, unit); err != nil {
				return nil, errors.WithStack(err)
			}
			return unit, nil
		}
		return nil, errors.WithStack(err)
	}

	return unit, nil
}

func (u *masterUsecase) CreateUnit(ctx context.Context, name string) error {
	_, err := u.masterRepo.CreateUnit(ctx, name)
	if err != nil {
		return errors.WithStack(err)
	}

	units, err := u.masterRepo.ListUnits(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	if units == nil {
		return nil
	}
	if err := u.redisMasterCacheRepo.SetUnits(ctx, units); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// VocalPattern
func (u *masterUsecase) CreateVocalPattern(ctx context.Context, songID int32, name string, singerIDs, singerPositions []int32) error {
	exist, err := u.masterRepo.ExistsSong(ctx, songID)
	if err != nil {
		return errors.WithStack(err)
	}
	if !exist {
		return errors.WithStack(ErrInvalidArgument)
	}

	for _, id := range singerIDs {
		exist, err := u.masterRepo.ExistsSinger(ctx, id)
		if err != nil {
			return errors.WithStack(err)
		}
		if !exist {
			return errors.WithStack(ErrInvalidArgument)
		}
	}

	vp, err := u.masterRepo.CreateVocalPattern(ctx, songID, name)
	if err != nil {
		return errors.WithStack(err)
	}

	for i := range len(singerIDs) {
		_, err := u.masterRepo.CreateVocalPatternSinger(ctx, vp.ID, singerIDs[i], singerPositions[i])
		if err != nil {
			return errors.WithStack(err)
		}
	}

	songs, err := u.masterRepo.ListSongs(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	if songs == nil {
		return nil
	}
	if err := u.redisMasterCacheRepo.SetSongs(ctx, songs); err != nil {
		return errors.WithStack(err)
	}

	charts, err := u.masterRepo.ListCharts(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	if charts == nil {
		return nil
	}
	if err := u.redisMasterCacheRepo.SetCharts(ctx, charts); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Song
func (u *masterUsecase) ListSongs(ctx context.Context) ([]*entity.Song, error) {
	songs, err := u.redisMasterCacheRepo.GetSongs(ctx)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			songs, err := u.masterRepo.ListSongs(ctx)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			if songs == nil {
				return []*entity.Song{}, nil
			}
			if err := u.redisMasterCacheRepo.SetSongs(ctx, songs); err != nil {
				return nil, errors.WithStack(err)
			}
			return songs, nil
		}
		return nil, errors.WithStack(err)
	}

	return songs, nil
}

func (u *masterUsecase) GetSongByID(ctx context.Context, id int32) (*entity.Song, error) {
	song, err := u.redisMasterCacheRepo.GetSongByID(ctx, id)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			song, err := u.masterRepo.GetSongByID(ctx, id)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			if err := u.redisMasterCacheRepo.SetSong(ctx, id, song); err != nil {
				return nil, errors.WithStack(err)
			}
			return song, nil
		}
		return nil, errors.WithStack(err)
	}

	return song, nil
}

func (u *masterUsecase) CreateSong(
	ctx context.Context,
	name, kana string,
	lyrics_id, music_id, arrangement_id int32,
	thumbnail, originalVideo string,
	releaseTime time.Time, deleted bool,
	unitIDs []int32,
	musicVideoTypes []enums.MusicVideoType,
) error {
	lyricsExist, err := u.masterRepo.ExistsArtist(ctx, lyrics_id)
	if err != nil {
		return errors.WithStack(err)
	}
	if !lyricsExist {
		return errors.WithStack(ErrInvalidArgument)
	}
	musicExist, err := u.masterRepo.ExistsArtist(ctx, music_id)
	if err != nil {
		return errors.WithStack(err)
	}
	if !musicExist {
		return errors.WithStack(ErrInvalidArgument)
	}
	arrangementExist, err := u.masterRepo.ExistsArtist(ctx, arrangement_id)
	if err != nil {
		return errors.WithStack(err)
	}
	if !arrangementExist {
		return errors.WithStack(ErrInvalidArgument)
	}

	for _, id := range unitIDs {
		exist, err := u.masterRepo.ExistsUnit(ctx, id)
		if err != nil {
			return errors.WithStack(err)
		}
		if !exist {
			return errors.WithStack(ErrInvalidArgument)
		}
	}

	s, err := u.masterRepo.CreateSong(ctx, name, kana, lyrics_id, music_id, arrangement_id, thumbnail, originalVideo, releaseTime, deleted)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, unitID := range unitIDs {
		_, err := u.masterRepo.CreateSongUnit(ctx, s.ID, unitID)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	for _, musicVideoType := range musicVideoTypes {
		_, err := u.masterRepo.CreateSongMusicVideoType(ctx, s.ID, musicVideoType)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	songs, err := u.masterRepo.ListSongs(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	if songs == nil {
		return nil
	}
	if err := u.redisMasterCacheRepo.SetSongs(ctx, songs); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Chart
func (u *masterUsecase) ListCharts(ctx context.Context) ([]*entity.Chart, error) {
	charts, err := u.redisMasterCacheRepo.GetCharts(ctx)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			charts, err := u.masterRepo.ListCharts(ctx)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			if charts == nil {
				return []*entity.Chart{}, nil
			}
			if err := u.redisMasterCacheRepo.SetCharts(ctx, charts); err != nil {
				return nil, errors.WithStack(err)
			}
			return charts, nil
		}
		return nil, errors.WithStack(err)
	}

	return charts, nil
}

func (u *masterUsecase) GetChartByID(ctx context.Context, id int32) (*entity.Chart, error) {
	chart, err := u.redisMasterCacheRepo.GetChartByID(ctx, id)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			chart, err := u.masterRepo.GetChartByID(ctx, id)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			if err := u.redisMasterCacheRepo.SetChart(ctx, id, chart); err != nil {
				return nil, errors.WithStack(err)
			}
			return chart, nil
		}
		return nil, errors.WithStack(err)
	}

	return chart, nil
}

func (u *masterUsecase) CreateChart(
	ctx context.Context,
	songID, difficultyType, level int32,
	chartViewLink string,
) error {
	exist, err := u.masterRepo.ExistsSong(ctx, songID)
	if err != nil {
		return errors.WithStack(err)
	}
	if !exist {
		return errors.WithStack(ErrInvalidArgument)
	}

	if _, err := u.masterRepo.CreateChart(ctx, songID, difficultyType, level, chartViewLink); err != nil {
		return errors.WithStack(err)
	}

	charts, err := u.masterRepo.ListCharts(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	if charts == nil {
		return nil
	}
	if err := u.redisMasterCacheRepo.SetCharts(ctx, charts); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
