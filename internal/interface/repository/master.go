package repository

import (
	"context"
	"database/sql"
	"slices"
	"time"

	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/entity"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/repository"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/gen/enums"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/infrastructure/db/sqlcgen"
	"github.com/cockroachdb/errors"
)

//go:generate gotests -w -all $GOFILE

type masterRepository struct {
	queries *sqlcgen.Queries
}

func NewMasterRepository(queries *sqlcgen.Queries) repository.MasterRepository {
	return &masterRepository{queries: queries}
}

// Artists
func (r *masterRepository) ListArtists(ctx context.Context) ([]*entity.Artist, error) {
	artists, err := r.queries.ListArtists(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	artistPointers := make([]*entity.Artist, len(artists))
	for i := range artists {
		artistPointers[i] = sqlToDomainArtist(&artists[i])
	}

	return artistPointers, nil
}

func (r *masterRepository) GetArtistByID(ctx context.Context, id int32) (*entity.Artist, error) {
	artist, err := r.queries.GetArtistByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.WithStack(repository.ErrNotFound)
		}
		return nil, errors.WithStack(err)
	}

	return sqlToDomainArtist(&artist), nil
}

func (r *masterRepository) CreateArtist(ctx context.Context, name, kana string) (*entity.Artist, error) {
	sqlArtist := sqlcgen.InsertArtistParams{
		Name: name,
		Kana: kana,
	}
	a, err := r.queries.InsertArtist(ctx, sqlArtist)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return sqlToDomainArtist(&a), nil
}

func (r *masterRepository) ExistsArtist(ctx context.Context, id int32) (bool, error) {
	exist, err := r.queries.ExistsArtist(ctx, id)
	if err != nil {
		return false, errors.WithStack(err)
	}

	return exist, nil
}

func sqlToDomainArtist(sqlArtist *sqlcgen.Artist) *entity.Artist {
	return &entity.Artist{
		ID:   sqlArtist.ID,
		Name: sqlArtist.Name,
		Kana: sqlArtist.Kana,
	}
}

// Singer
func (r *masterRepository) ListSingers(ctx context.Context) ([]*entity.Singer, error) {
	singers, err := r.queries.ListSingers(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	singerPointers := make([]*entity.Singer, len(singers))
	for i := range singers {
		singerPointers[i] = sqlToDomainSinger(&singers[i])
	}

	return singerPointers, nil
}

func (r *masterRepository) GetSingerByID(ctx context.Context, id int32) (*entity.Singer, error) {
	singer, err := r.queries.GetSingerByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.WithStack(repository.ErrNotFound)
		}
		return nil, errors.WithStack(err)
	}

	return sqlToDomainSinger(&singer), nil
}

func (r *masterRepository) CreateSinger(ctx context.Context, name string) (*entity.Singer, error) {
	s, err := r.queries.InsertSinger(ctx, name)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return sqlToDomainSinger(&s), nil
}

func (r *masterRepository) ExistsSinger(ctx context.Context, id int32) (bool, error) {
	exist, err := r.queries.ExistsSinger(ctx, id)
	if err != nil {
		return false, errors.WithStack(err)
	}

	return exist, nil
}

func sqlToDomainSinger(sqlSinger *sqlcgen.Singer) *entity.Singer {
	return &entity.Singer{
		ID:   sqlSinger.ID,
		Name: sqlSinger.Name,
	}
}

// Unit
func (r *masterRepository) ListUnits(ctx context.Context) ([]*entity.Unit, error) {
	units, err := r.queries.ListUnits(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	unitPointers := make([]*entity.Unit, len(units))
	for i := range units {
		unitPointers[i] = sqlToDomainUnit(&units[i])
	}

	return unitPointers, nil
}

func (r *masterRepository) GetUnitByID(ctx context.Context, id int32) (*entity.Unit, error) {
	unit, err := r.queries.GetUnitByID(ctx, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return sqlToDomainUnit(&unit), nil
}

func (r *masterRepository) CreateUnit(ctx context.Context, name string) (*entity.Unit, error) {
	u, err := r.queries.InsertUnit(ctx, name)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return sqlToDomainUnit(&u), nil
}

func (r *masterRepository) ExistsUnit(ctx context.Context, id int32) (bool, error) {
	exist, err := r.queries.ExistsUnit(ctx, id)
	if err != nil {
		return false, errors.WithStack(err)
	}

	return exist, nil
}

func sqlToDomainUnit(sqlUnit *sqlcgen.Unit) *entity.Unit {
	return &entity.Unit{
		ID:   sqlUnit.ID,
		Name: sqlUnit.Name,
	}
}

// VocalPattern
func (r *masterRepository) CreateVocalPattern(ctx context.Context, songID int32, name string) (*sqlcgen.VocalPattern, error) {
	sqlVocalPattern := sqlcgen.InsertVocalPatternParams{
		SongID: sql.NullInt32{Int32: songID, Valid: true},
		Name:   sql.NullString{String: name, Valid: true},
	}
	vp, err := r.queries.InsertVocalPattern(ctx, sqlVocalPattern)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &vp, nil
}

// VocalPatternSinger
func (r *masterRepository) CreateVocalPatternSinger(ctx context.Context, vocalPatternID, singerID, position int32) (*sqlcgen.VocalPatternSinger, error) {
	sqlVocalPatternSinger := sqlcgen.InsertVocalPatternSingerParams{
		VocalPatternID: sql.NullInt32{Int32: vocalPatternID, Valid: true},
		SingerID:       sql.NullInt32{Int32: singerID, Valid: true},
		Position:       sql.NullInt32{Int32: position, Valid: true},
	}
	vps, err := r.queries.InsertVocalPatternSinger(ctx, sqlVocalPatternSinger)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &vps, nil
}

// SongUnit
func (r *masterRepository) CreateSongUnit(ctx context.Context, songID, unitID int32) (*sqlcgen.SongUnit, error) {
	sqlSongUnit := sqlcgen.InsertSongUnitParams{
		SongID: sql.NullInt32{Int32: songID, Valid: true},
		UnitID: sql.NullInt32{Int32: unitID, Valid: true},
	}
	su, err := r.queries.InsertSongUnit(ctx, sqlSongUnit)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &su, nil
}

// SongMusicVideoType
func (r *masterRepository) CreateSongMusicVideoType(ctx context.Context, songID int32, musicVideoType enums.MusicVideoType) (*sqlcgen.SongMusicVideoType, error) {
	sqlSongMusicVideoType := sqlcgen.InsertSongMusicVideoTypeParams{
		SongID:         sql.NullInt32{Int32: songID, Valid: true},
		MusicVideoType: sql.NullInt32{Int32: int32(musicVideoType), Valid: true},
	}
	smvt, err := r.queries.InsertSongMusicVideoType(ctx, sqlSongMusicVideoType)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &smvt, nil
}

// Song
func (r *masterRepository) ListSongs(ctx context.Context) ([]*entity.Song, error) {
	songs, err := r.queries.ListSongWithArtists(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return sqlToDomainListSong(songs), nil
}

func (r *masterRepository) GetSongByID(ctx context.Context, id int32) (*entity.Song, error) {
	sqlSongs, err := r.queries.GetSongDetailsByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || len(sqlSongs) <= 0 {
			return nil, errors.WithStack(repository.ErrNotFound)
		}
		return nil, errors.WithStack(err)
	}
	songs := sqlToDomainGetSong(sqlSongs)
	if len(songs) == 0 {
		return nil, errors.WithStack(repository.ErrNotFound)
	}

	return songs[0], nil
}

func (r *masterRepository) CreateSong(
	ctx context.Context,
	name, kana string,
	lyrics_id, music_id, arrangement_id int32,
	thumbnail, originalVideo string,
	releaseTime time.Time, deleted bool,
) (*sqlcgen.Song, error) {
	sqlSong := sqlcgen.InsertSongParams{
		Name:          name,
		Kana:          kana,
		LyricsID:      sql.NullInt32{Int32: lyrics_id, Valid: true},
		MusicID:       sql.NullInt32{Int32: music_id, Valid: true},
		ArrangementID: sql.NullInt32{Int32: arrangement_id, Valid: true},
		Thumbnail:     sql.NullString{String: thumbnail, Valid: true},
		OriginalVideo: sql.NullString{String: originalVideo, Valid: true},
		ReleaseTime:   sql.NullTime{Time: releaseTime, Valid: true},
		Deleted:       sql.NullBool{Bool: deleted, Valid: true},
	}
	s, err := r.queries.InsertSong(ctx, sqlSong)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &s, nil
}

func (r *masterRepository) ExistsSong(ctx context.Context, id int32) (bool, error) {
	exist, err := r.queries.ExistsSong(ctx, id)
	if err != nil {
		return false, errors.WithStack(err)
	}

	return exist, nil
}

func sqlToDomainListSong(sqlSongs []sqlcgen.ListSongWithArtistsRow) []*entity.Song {
	var songs []*entity.Song
	for _, v := range sqlSongs {
		songIndex := slices.IndexFunc(songs, func(song *entity.Song) bool {
			return song.ID == v.ID
		})
		if songIndex == -1 {
			songs = append(songs, &entity.Song{
				ID:   v.ID,
				Name: v.Name,
				Kana: v.Kana,
				Lyrics: entity.Artist{
					ID:   v.LyricsArtistID.Int32,
					Name: v.LyricsArtistName.String,
					Kana: v.LyricsArtistKana.String,
				},
				Music: entity.Artist{
					ID:   v.MusicArtistID.Int32,
					Name: v.MusicArtistName.String,
					Kana: v.MusicArtistKana.String,
				},
				Arrangement: entity.Artist{
					ID:   v.ArrangementArtistID.Int32,
					Name: v.ArrangementArtistName.String,
					Kana: v.ArrangementArtistKana.String,
				},
				Thumbnail:     v.Thumbnail.String,
				OriginalVideo: v.OriginalVideo.String,
				ReleaseTime:   v.ReleaseTime.Time,
				Deleted:       v.Deleted.Bool,
				VocalPatterns: []*entity.VocalPattern{
					{
						ID:   v.VocalPatternID.Int32,
						Name: v.VocalPatternName.String,
						Singers: []*entity.Singer{
							{
								ID:       v.SingerID.Int32,
								Name:     v.SingerName.String,
								Position: v.SingerOrder.Int32,
							},
						},
					},
				},
				Units: []*entity.Unit{
					{
						ID:   v.UnitID.Int32,
						Name: v.UnitName.String,
					},
				},
				MusicVideoTypes: []enums.MusicVideoType{
					enums.MusicVideoType(v.MusicVideoTypeID.Int32),
				},
			})
		} else {
			vpIndex := slices.IndexFunc(songs[songIndex].VocalPatterns, func(vp *entity.VocalPattern) bool {
				return vp.ID == v.VocalPatternID.Int32
			})
			if vpIndex != -1 {
				sgFind := slices.ContainsFunc(songs[songIndex].VocalPatterns[vpIndex].Singers, func(sg *entity.Singer) bool {
					return sg.ID == v.SingerID.Int32
				})
				if !sgFind {
					songs[songIndex].VocalPatterns[vpIndex].Singers = append(songs[songIndex].VocalPatterns[vpIndex].Singers, &entity.Singer{
						ID:       v.SingerID.Int32,
						Name:     v.SingerName.String,
						Position: v.SingerOrder.Int32,
					})
				}
			} else {
				songs[songIndex].VocalPatterns = append(songs[songIndex].VocalPatterns, &entity.VocalPattern{
					ID:      v.VocalPatternID.Int32,
					Name:    v.VocalPatternName.String,
					Singers: []*entity.Singer{{ID: v.SingerID.Int32, Name: v.SingerName.String, Position: v.SingerOrder.Int32}},
				})
			}

			utFind := slices.ContainsFunc(songs[songIndex].Units, func(ut *entity.Unit) bool {
				return ut.ID == v.UnitID.Int32
			})
			if !utFind {
				songs[songIndex].Units = append(songs[songIndex].Units, &entity.Unit{
					ID:   v.UnitID.Int32,
					Name: v.UnitName.String,
				})
			}

			mvtFind := slices.ContainsFunc(songs[songIndex].MusicVideoTypes, func(mvt enums.MusicVideoType) bool {
				return mvt == enums.MusicVideoType(v.MusicVideoTypeID.Int32)
			})
			if !mvtFind {
				songs[songIndex].MusicVideoTypes = append(songs[songIndex].MusicVideoTypes, enums.MusicVideoType(v.MusicVideoTypeID.Int32))
			}
		}
	}
	return songs
}
func sqlToDomainGetSong(sqlSongs []sqlcgen.GetSongDetailsByIDRow) []*entity.Song {
	var songs []*entity.Song
	for _, v := range sqlSongs {
		songIndex := slices.IndexFunc(songs, func(song *entity.Song) bool {
			return song.ID == v.ID
		})
		if songIndex == -1 {
			songs = append(songs, &entity.Song{
				ID:   v.ID,
				Name: v.Name,
				Kana: v.Kana,
				Lyrics: entity.Artist{
					ID:   v.LyricsArtistID.Int32,
					Name: v.LyricsArtistName.String,
					Kana: v.LyricsArtistKana.String,
				},
				Music: entity.Artist{
					ID:   v.MusicArtistID.Int32,
					Name: v.MusicArtistName.String,
					Kana: v.MusicArtistKana.String,
				},
				Arrangement: entity.Artist{
					ID:   v.ArrangementArtistID.Int32,
					Name: v.ArrangementArtistName.String,
					Kana: v.ArrangementArtistKana.String,
				},
				Thumbnail:     v.Thumbnail.String,
				OriginalVideo: v.OriginalVideo.String,
				ReleaseTime:   v.ReleaseTime.Time,
				Deleted:       v.Deleted.Bool,
				VocalPatterns: []*entity.VocalPattern{
					{
						ID:   v.VocalPatternID.Int32,
						Name: v.VocalPatternName.String,
						Singers: []*entity.Singer{
							{
								ID:       v.SingerID.Int32,
								Name:     v.SingerName.String,
								Position: v.SingerOrder.Int32,
							},
						},
					},
				},
				Units: []*entity.Unit{
					{
						ID:   v.UnitID.Int32,
						Name: v.UnitName.String,
					},
				},
				MusicVideoTypes: []enums.MusicVideoType{
					enums.MusicVideoType(v.MusicVideoTypeID.Int32),
				},
			})
		} else {
			vpIndex := slices.IndexFunc(songs[songIndex].VocalPatterns, func(vp *entity.VocalPattern) bool {
				return vp.ID == v.VocalPatternID.Int32
			})
			if vpIndex != -1 {
				sgFind := slices.ContainsFunc(songs[songIndex].VocalPatterns[vpIndex].Singers, func(sg *entity.Singer) bool {
					return sg.ID == v.SingerID.Int32
				})
				if !sgFind {
					songs[songIndex].VocalPatterns[vpIndex].Singers = append(songs[songIndex].VocalPatterns[vpIndex].Singers, &entity.Singer{
						ID:       v.SingerID.Int32,
						Name:     v.SingerName.String,
						Position: v.SingerOrder.Int32,
					})
				}
			} else {
				songs[songIndex].VocalPatterns = append(songs[songIndex].VocalPatterns, &entity.VocalPattern{
					ID:      v.VocalPatternID.Int32,
					Name:    v.VocalPatternName.String,
					Singers: []*entity.Singer{{ID: v.SingerID.Int32, Name: v.SingerName.String, Position: v.SingerOrder.Int32}},
				})
			}

			utFind := slices.ContainsFunc(songs[songIndex].Units, func(ut *entity.Unit) bool {
				return ut.ID == v.UnitID.Int32
			})
			if !utFind {
				songs[songIndex].Units = append(songs[songIndex].Units, &entity.Unit{
					ID:   v.UnitID.Int32,
					Name: v.UnitName.String,
				})
			}

			mvtFind := slices.ContainsFunc(songs[songIndex].MusicVideoTypes, func(mvt enums.MusicVideoType) bool {
				return mvt == enums.MusicVideoType(v.MusicVideoTypeID.Int32)
			})
			if !mvtFind {
				songs[songIndex].MusicVideoTypes = append(songs[songIndex].MusicVideoTypes, enums.MusicVideoType(v.MusicVideoTypeID.Int32))
			}
		}
	}
	return songs
}

// Chart
func (r *masterRepository) ListCharts(ctx context.Context) ([]*entity.Chart, error) {
	charts, err := r.queries.ListChartWithSongWithArtists(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return sqlToDomainListChart(charts), nil
}

func (r *masterRepository) GetChartByID(ctx context.Context, id int32) (*entity.Chart, error) {
	sqlCharts, err := r.queries.GetChartWithSongWithArtistsByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || len(sqlCharts) <= 0 {
			return nil, errors.WithStack(repository.ErrNotFound)
		}
		return nil, errors.WithStack(err)
	}
	charts := sqlToDomainGetChart(sqlCharts)
	if len(charts) == 0 {
		return nil, errors.WithStack(repository.ErrNotFound)
	}

	return charts[0], nil
}

func (r *masterRepository) CreateChart(ctx context.Context, songID, difficultyType, level int32, chartViewLink string) (*sqlcgen.Chart, error) {
	sqlChart := sqlcgen.InsertChartParams{
		SongID:         sql.NullInt32{Int32: songID, Valid: true},
		DifficultyType: sql.NullInt32{Int32: difficultyType, Valid: true},
		Level:          sql.NullInt32{Int32: level, Valid: true},
		ChartViewLink:  sql.NullString{String: chartViewLink, Valid: true},
	}
	c, err := r.queries.InsertChart(ctx, sqlChart)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &c, nil
}

func (r *masterRepository) ExistsChart(ctx context.Context, id int32) (bool, error) {
	exist, err := r.queries.ExistsChart(ctx, id)
	if err != nil {
		return false, errors.WithStack(err)
	}

	return exist, nil
}

func sqlToDomainListChart(sqlCharts []sqlcgen.ListChartWithSongWithArtistsRow) []*entity.Chart {
	var charts []*entity.Chart
	for _, v := range sqlCharts {
		chartIndex := slices.IndexFunc(charts, func(chart *entity.Chart) bool {
			return chart.ID == v.ID
		})
		if chartIndex == -1 {
			charts = append(charts, &entity.Chart{
				ID: v.ID,
				Song: entity.Song{
					ID:   v.SongID.Int32,
					Name: v.SongName.String,
					Kana: v.SongKana.String,
					Lyrics: entity.Artist{
						ID:   v.LyricsArtistID.Int32,
						Name: v.LyricsArtistName.String,
						Kana: v.LyricsArtistKana.String,
					},
					Music: entity.Artist{
						ID:   v.MusicArtistID.Int32,
						Name: v.MusicArtistName.String,
						Kana: v.MusicArtistKana.String,
					},
					Arrangement: entity.Artist{
						ID:   v.ArrangementArtistID.Int32,
						Name: v.ArrangementArtistName.String,
						Kana: v.ArrangementArtistKana.String,
					},
					Thumbnail:     v.Thumbnail.String,
					OriginalVideo: v.OriginalVideo.String,
					ReleaseTime:   v.ReleaseTime.Time,
					Deleted:       v.Deleted.Bool,
					VocalPatterns: []*entity.VocalPattern{
						{
							ID:   v.VocalPatternID.Int32,
							Name: v.VocalPatternName.String,
							Singers: []*entity.Singer{
								{
									ID:       v.SingerID.Int32,
									Name:     v.SingerName.String,
									Position: v.SingerOrder.Int32,
								},
							},
						},
					},
					Units: []*entity.Unit{
						{
							ID:   v.UnitID.Int32,
							Name: v.UnitName.String,
						},
					},
					MusicVideoTypes: []enums.MusicVideoType{
						enums.MusicVideoType(v.MusicVideoTypeID.Int32),
					},
				},
				DifficultyType: enums.DifficultyType(v.DifficultyType.Int32),
				Level:          v.Level.Int32,
				ChartViewLink:  v.ChartViewLink.String,
			})
		} else {
			vpIndex := slices.IndexFunc(charts[chartIndex].Song.VocalPatterns, func(vp *entity.VocalPattern) bool {
				return vp.ID == v.VocalPatternID.Int32
			})
			if vpIndex != -1 {
				sgFind := slices.ContainsFunc(charts[chartIndex].Song.VocalPatterns[vpIndex].Singers, func(sg *entity.Singer) bool {
					return sg.ID == v.SingerID.Int32
				})
				if !sgFind {
					charts[chartIndex].Song.VocalPatterns[vpIndex].Singers = append(charts[chartIndex].Song.VocalPatterns[vpIndex].Singers, &entity.Singer{
						ID:       v.SingerID.Int32,
						Name:     v.SingerName.String,
						Position: v.SingerOrder.Int32,
					})
				}
			} else {
				charts[chartIndex].Song.VocalPatterns = append(charts[chartIndex].Song.VocalPatterns, &entity.VocalPattern{
					ID:      v.VocalPatternID.Int32,
					Name:    v.VocalPatternName.String,
					Singers: []*entity.Singer{{ID: v.SingerID.Int32, Name: v.SingerName.String, Position: v.SingerOrder.Int32}},
				})
			}

			utFind := slices.ContainsFunc(charts[chartIndex].Song.Units, func(ut *entity.Unit) bool {
				return ut.ID == v.UnitID.Int32
			})
			if !utFind {
				charts[chartIndex].Song.Units = append(charts[chartIndex].Song.Units, &entity.Unit{
					ID:   v.UnitID.Int32,
					Name: v.UnitName.String,
				})
			}

			mvtFind := slices.ContainsFunc(charts[chartIndex].Song.MusicVideoTypes, func(mvt enums.MusicVideoType) bool {
				return mvt == enums.MusicVideoType(v.MusicVideoTypeID.Int32)
			})
			if !mvtFind {
				charts[chartIndex].Song.MusicVideoTypes = append(charts[chartIndex].Song.MusicVideoTypes, enums.MusicVideoType(v.MusicArtistID.Int32))
			}
		}
	}
	return charts
}
func sqlToDomainGetChart(sqlCharts []sqlcgen.GetChartWithSongWithArtistsByIDRow) []*entity.Chart {
	var charts []*entity.Chart
	for _, v := range sqlCharts {
		chartIndex := slices.IndexFunc(charts, func(chart *entity.Chart) bool {
			return chart.ID == v.ID
		})
		if chartIndex == -1 {
			charts = append(charts, &entity.Chart{
				ID: v.ID,
				Song: entity.Song{
					ID:   v.SongID.Int32,
					Name: v.SongName.String,
					Kana: v.SongKana.String,
					Lyrics: entity.Artist{
						ID:   v.LyricsArtistID.Int32,
						Name: v.LyricsArtistName.String,
						Kana: v.LyricsArtistKana.String,
					},
					Music: entity.Artist{
						ID:   v.MusicArtistID.Int32,
						Name: v.MusicArtistName.String,
						Kana: v.MusicArtistKana.String,
					},
					Arrangement: entity.Artist{
						ID:   v.ArrangementArtistID.Int32,
						Name: v.ArrangementArtistName.String,
						Kana: v.ArrangementArtistKana.String,
					},
					Thumbnail:     v.Thumbnail.String,
					OriginalVideo: v.OriginalVideo.String,
					ReleaseTime:   v.ReleaseTime.Time,
					Deleted:       v.Deleted.Bool,
					VocalPatterns: []*entity.VocalPattern{
						{
							ID:   v.VocalPatternID.Int32,
							Name: v.VocalPatternName.String,
							Singers: []*entity.Singer{
								{
									ID:       v.SingerID.Int32,
									Name:     v.SingerName.String,
									Position: v.SingerOrder.Int32,
								},
							},
						},
					},
					Units: []*entity.Unit{
						{
							ID:   v.UnitID.Int32,
							Name: v.UnitName.String,
						},
					},
					MusicVideoTypes: []enums.MusicVideoType{
						enums.MusicVideoType(v.MusicVideoTypeID.Int32),
					},
				},
				DifficultyType: enums.DifficultyType(v.DifficultyType.Int32),
				Level:          v.Level.Int32,
				ChartViewLink:  v.ChartViewLink.String,
			})
		} else {
			vpIndex := slices.IndexFunc(charts[chartIndex].Song.VocalPatterns, func(vp *entity.VocalPattern) bool {
				return vp.ID == v.VocalPatternID.Int32
			})
			if vpIndex != -1 {
				sgFind := slices.ContainsFunc(charts[chartIndex].Song.VocalPatterns[vpIndex].Singers, func(sg *entity.Singer) bool {
					return sg.ID == v.SingerID.Int32
				})
				if !sgFind {
					charts[chartIndex].Song.VocalPatterns[vpIndex].Singers = append(charts[chartIndex].Song.VocalPatterns[vpIndex].Singers, &entity.Singer{
						ID:       v.SingerID.Int32,
						Name:     v.SingerName.String,
						Position: v.SingerOrder.Int32,
					})
				}
			} else {
				charts[chartIndex].Song.VocalPatterns = append(charts[chartIndex].Song.VocalPatterns, &entity.VocalPattern{
					ID:      v.VocalPatternID.Int32,
					Name:    v.VocalPatternName.String,
					Singers: []*entity.Singer{{ID: v.SingerID.Int32, Name: v.SingerName.String, Position: v.SingerOrder.Int32}},
				})
			}

			utFind := slices.ContainsFunc(charts[chartIndex].Song.Units, func(ut *entity.Unit) bool {
				return ut.ID == v.UnitID.Int32
			})
			if !utFind {
				charts[chartIndex].Song.Units = append(charts[chartIndex].Song.Units, &entity.Unit{
					ID:   v.UnitID.Int32,
					Name: v.UnitName.String,
				})
			}

			mvtFind := slices.ContainsFunc(charts[chartIndex].Song.MusicVideoTypes, func(mvt enums.MusicVideoType) bool {
				return mvt == enums.MusicVideoType(v.MusicVideoTypeID.Int32)
			})
			if !mvtFind {
				charts[chartIndex].Song.MusicVideoTypes = append(charts[chartIndex].Song.MusicVideoTypes, enums.MusicVideoType(v.MusicVideoTypeID.Int32))
			}
		}
	}
	return charts
}
