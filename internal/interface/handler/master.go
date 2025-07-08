package handler

import (
	"context"
	"log"
	"sort"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/repository"
	proto_master "github.com/Shakkuuu/sekai-songs-mylist/internal/gen/master"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/pkg/auth"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/usecase"
	"github.com/cockroachdb/errors"
)

//go:generate gotests -w -all $GOFILE

type MasterHandler struct {
	masterUsecase usecase.MasterUsecase
	userUsecase   usecase.UserUsecase
}

func NewMasterHandler(masterUsecase usecase.MasterUsecase, userUsecase usecase.UserUsecase) *MasterHandler {
	return &MasterHandler{
		masterUsecase: masterUsecase,
		userUsecase:   userUsecase,
	}
}

// Artist
func (h *MasterHandler) GetArtists(ctx context.Context, req *connect.Request[proto_master.GetArtistsRequest]) (*connect.Response[proto_master.GetArtistsResponse], error) {
	artists, err := h.masterUsecase.ListArtists(ctx)
	if err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}

	protoArtists := make([]*proto_master.Artist, len(artists))
	for i, artist := range artists {
		protoArtists[i] = &proto_master.Artist{
			Id:   artist.ID,
			Name: artist.Name,
			Kana: artist.Kana,
		}
	}

	return connect.NewResponse(&proto_master.GetArtistsResponse{
		Artists: protoArtists,
	}), nil
}

func (h *MasterHandler) GetArtist(ctx context.Context, req *connect.Request[proto_master.GetArtistRequest]) (*connect.Response[proto_master.GetArtistResponse], error) {
	artist, err := h.masterUsecase.GetArtistByID(ctx, req.Msg.GetId())
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			cerr := errors.WithStack(err)
			log.Printf("%+v\n", cerr)
			return nil, connect.NewError(connect.CodeNotFound, cerr)
		}
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}

	return connect.NewResponse(&proto_master.GetArtistResponse{
		Artist: &proto_master.Artist{
			Id:   int32(artist.ID),
			Name: artist.Name,
			Kana: artist.Kana,
		},
	}), nil
}

func (h *MasterHandler) CreateArtist(ctx context.Context, req *connect.Request[proto_master.CreateArtistRequest]) (*connect.Response[proto_master.CreateArtistResponse], error) {
	id, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok {
		err := errors.New("user id not found in context")
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeUnauthenticated, cerr)
	}
	isAdmin, err := h.userUsecase.IsAdmin(ctx, id)
	if err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}
	if !isAdmin {
		err := errors.New("permission denied: not admin")
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodePermissionDenied, cerr)
	}

	if err := req.Msg.Validate(); err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInvalidArgument, cerr)
	}

	if err := h.masterUsecase.CreateArtist(ctx, req.Msg.GetName(), req.Msg.GetKana()); err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}

	return connect.NewResponse(&proto_master.CreateArtistResponse{}), nil
}

// Singer
func (h *MasterHandler) GetSingers(ctx context.Context, req *connect.Request[proto_master.GetSingersRequest]) (*connect.Response[proto_master.GetSingersResponse], error) {
	singers, err := h.masterUsecase.ListSingers(ctx)
	if err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}

	protoSingers := make([]*proto_master.Singer, len(singers))
	for i, singer := range singers {
		protoSingers[i] = &proto_master.Singer{
			Id:   singer.ID,
			Name: singer.Name,
		}
	}

	return connect.NewResponse(&proto_master.GetSingersResponse{
		Singers: protoSingers,
	}), nil
}

func (h *MasterHandler) GetSinger(ctx context.Context, req *connect.Request[proto_master.GetSingerRequest]) (*connect.Response[proto_master.GetSingerResponse], error) {
	singer, err := h.masterUsecase.GetSingerByID(ctx, req.Msg.GetId())
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			cerr := errors.WithStack(err)
			log.Printf("%+v\n", cerr)
			return nil, connect.NewError(connect.CodeNotFound, cerr)
		}
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}

	return connect.NewResponse(&proto_master.GetSingerResponse{
		Singer: &proto_master.Singer{
			Id:   int32(singer.ID),
			Name: singer.Name,
		},
	}), nil
}

func (h *MasterHandler) CreateSinger(ctx context.Context, req *connect.Request[proto_master.CreateSingerRequest]) (*connect.Response[proto_master.CreateSingerResponse], error) {
	id, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok {
		err := errors.New("user id not found in context")
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeUnauthenticated, cerr)
	}
	isAdmin, err := h.userUsecase.IsAdmin(ctx, id)
	if err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}
	if !isAdmin {
		err := errors.New("permission denied: not admin")
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodePermissionDenied, cerr)
	}

	if err := req.Msg.Validate(); err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInvalidArgument, cerr)
	}

	if err := h.masterUsecase.CreateSinger(ctx, req.Msg.GetName()); err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}

	return connect.NewResponse(&proto_master.CreateSingerResponse{}), nil
}

// Unit
func (h *MasterHandler) GetUnits(ctx context.Context, req *connect.Request[proto_master.GetUnitsRequest]) (*connect.Response[proto_master.GetUnitsResponse], error) {
	units, err := h.masterUsecase.ListUnits(ctx)
	if err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}

	protoUnits := make([]*proto_master.Unit, len(units))
	for i, unit := range units {
		protoUnits[i] = &proto_master.Unit{
			Id:   unit.ID,
			Name: unit.Name,
		}
	}

	return connect.NewResponse(&proto_master.GetUnitsResponse{
		Units: protoUnits,
	}), nil
}

func (h *MasterHandler) GetUnit(ctx context.Context, req *connect.Request[proto_master.GetUnitRequest]) (*connect.Response[proto_master.GetUnitResponse], error) {
	unit, err := h.masterUsecase.GetUnitByID(ctx, req.Msg.GetId())
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			cerr := errors.WithStack(err)
			log.Printf("%+v\n", cerr)
			return nil, connect.NewError(connect.CodeNotFound, cerr)
		}
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}

	return connect.NewResponse(&proto_master.GetUnitResponse{
		Unit: &proto_master.Unit{
			Id:   int32(unit.ID),
			Name: unit.Name,
		},
	}), nil
}

func (h *MasterHandler) CreateUnit(ctx context.Context, req *connect.Request[proto_master.CreateUnitRequest]) (*connect.Response[proto_master.CreateUnitResponse], error) {
	id, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok {
		err := errors.New("user id not found in context")
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeUnauthenticated, cerr)
	}
	isAdmin, err := h.userUsecase.IsAdmin(ctx, id)
	if err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}
	if !isAdmin {
		err := errors.New("permission denied: not admin")
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodePermissionDenied, cerr)
	}

	if err := req.Msg.Validate(); err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInvalidArgument, cerr)
	}

	if err := h.masterUsecase.CreateUnit(ctx, req.Msg.GetName()); err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}

	return connect.NewResponse(&proto_master.CreateUnitResponse{}), nil
}

// VocalPattern
func (h *MasterHandler) CreateVocalPattern(ctx context.Context, req *connect.Request[proto_master.CreateVocalPatternRequest]) (*connect.Response[proto_master.CreateVocalPatternResponse], error) {
	id, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok {
		err := errors.New("user id not found in context")
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeUnauthenticated, cerr)
	}
	isAdmin, err := h.userUsecase.IsAdmin(ctx, id)
	if err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}
	if !isAdmin {
		err := errors.New("permission denied: not admin")
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodePermissionDenied, cerr)
	}

	if err := req.Msg.Validate(); err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInvalidArgument, cerr)
	}

	if len(req.Msg.GetSingerIds()) != len(req.Msg.GetSingerPositions()) {
		err := errors.New("singerIDsとsingerPositionの長さが一致していない")
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInvalidArgument, cerr)
	}

	if err := h.masterUsecase.CreateVocalPattern(ctx, req.Msg.GetSongId(), req.Msg.GetName(), req.Msg.GetSingerIds(), req.Msg.GetSingerPositions()); err != nil {
		if errors.Is(err, usecase.ErrInvalidArgument) {
			cerr := errors.WithStack(err)
			log.Printf("%+v\n", cerr)
			return nil, connect.NewError(connect.CodeInvalidArgument, cerr)
		}
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}

	return connect.NewResponse(&proto_master.CreateVocalPatternResponse{}), nil
}

// Song
func (h *MasterHandler) GetSongs(ctx context.Context, req *connect.Request[proto_master.GetSongsRequest]) (*connect.Response[proto_master.GetSongsResponse], error) {
	songs, err := h.masterUsecase.ListSongs(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			cerr := errors.WithStack(err)
			log.Printf("%+v\n", cerr)
			return nil, connect.NewError(connect.CodeNotFound, cerr)
		}
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}

	var protoSongs []*proto_master.Song
	for _, song := range songs {
		var protoVocalPatterns []*proto_master.VocalPattern
		for _, vp := range song.VocalPatterns {
			if vp == nil {
				continue
			}
			sort.Slice(vp.Singers, func(i, j int) bool {
				return vp.Singers[i].Position < vp.Singers[j].Position
			})
			var protoSingers []*proto_master.Singer
			for _, s := range vp.Singers {
				if s == nil {
					continue
				}
				protoSingers = append(protoSingers, &proto_master.Singer{
					Id:   s.ID,
					Name: s.Name,
				})
			}
			protoVocalPatterns = append(protoVocalPatterns, &proto_master.VocalPattern{
				Id:      vp.ID,
				Name:    vp.Name,
				Singers: protoSingers,
			})
		}
		var protoUnits []*proto_master.Unit
		for _, u := range song.Units {
			if u == nil {
				continue
			}
			protoUnits = append(protoUnits, &proto_master.Unit{
				Id:   u.ID,
				Name: u.Name,
			})
		}

		protoSongs = append(protoSongs, &proto_master.Song{
			Id:   song.ID,
			Name: song.Name,
			Kana: song.Kana,
			Lyrics: &proto_master.Artist{
				Id:   song.Lyrics.ID,
				Name: song.Lyrics.Name,
				Kana: song.Lyrics.Kana,
			},
			Music: &proto_master.Artist{
				Id:   song.Music.ID,
				Name: song.Music.Name,
				Kana: song.Music.Kana,
			},
			Arrangement: &proto_master.Artist{
				Id:   song.Arrangement.ID,
				Name: song.Arrangement.Name,
				Kana: song.Arrangement.Kana,
			},
			Thumbnail:       song.Thumbnail,
			OriginalVideo:   song.OriginalVideo,
			ReleaseTime:     timestamppb.New(song.ReleaseTime),
			Deleted:         song.Deleted,
			VocalPatterns:   protoVocalPatterns,
			Units:           protoUnits,
			MusicVideoTypes: song.MusicVideoTypes,
		})
	}

	return connect.NewResponse(&proto_master.GetSongsResponse{
		Songs: protoSongs,
	}), nil
}

func (h *MasterHandler) GetSong(ctx context.Context, req *connect.Request[proto_master.GetSongRequest]) (*connect.Response[proto_master.GetSongResponse], error) {
	song, err := h.masterUsecase.GetSongByID(ctx, req.Msg.GetId())
	if err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}

	var protoSong *proto_master.Song
	var protoVocalPatterns []*proto_master.VocalPattern
	for _, vp := range song.VocalPatterns {
		if vp == nil {
			continue
		}
		sort.Slice(vp.Singers, func(i, j int) bool {
			return vp.Singers[i].Position < vp.Singers[j].Position
		})
		var protoSingers []*proto_master.Singer
		for _, s := range vp.Singers {
			if s == nil {
				continue
			}
			protoSingers = append(protoSingers, &proto_master.Singer{
				Id:   s.ID,
				Name: s.Name,
			})
		}
		protoVocalPatterns = append(protoVocalPatterns, &proto_master.VocalPattern{
			Id:      vp.ID,
			Name:    vp.Name,
			Singers: protoSingers,
		})
	}
	var protoUnits []*proto_master.Unit
	for _, u := range song.Units {
		if u == nil {
			continue
		}
		protoUnits = append(protoUnits, &proto_master.Unit{
			Id:   u.ID,
			Name: u.Name,
		})
	}
	protoSong = &proto_master.Song{
		Id:   song.ID,
		Name: song.Name,
		Kana: song.Kana,
		Lyrics: &proto_master.Artist{
			Id:   song.Lyrics.ID,
			Name: song.Lyrics.Name,
			Kana: song.Lyrics.Kana,
		},
		Music: &proto_master.Artist{
			Id:   song.Music.ID,
			Name: song.Music.Name,
			Kana: song.Music.Kana,
		},
		Arrangement: &proto_master.Artist{
			Id:   song.Arrangement.ID,
			Name: song.Arrangement.Name,
			Kana: song.Arrangement.Kana,
		},
		Thumbnail:       song.Thumbnail,
		OriginalVideo:   song.OriginalVideo,
		ReleaseTime:     timestamppb.New(song.ReleaseTime),
		Deleted:         song.Deleted,
		VocalPatterns:   protoVocalPatterns,
		Units:           protoUnits,
		MusicVideoTypes: song.MusicVideoTypes,
	}

	return connect.NewResponse(&proto_master.GetSongResponse{
		Song: protoSong,
	}), nil
}

func (h *MasterHandler) CreateSong(ctx context.Context, req *connect.Request[proto_master.CreateSongRequest]) (*connect.Response[proto_master.CreateSongResponse], error) {
	id, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok {
		err := errors.New("user id not found in context")
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeUnauthenticated, cerr)
	}
	isAdmin, err := h.userUsecase.IsAdmin(ctx, id)
	if err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}
	if !isAdmin {
		err := errors.New("permission denied: not admin")
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodePermissionDenied, cerr)
	}

	if err := req.Msg.Validate(); err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInvalidArgument, cerr)
	}

	if req.Msg.GetReleaseTime() == nil {
		err := errors.New("ReleaseTime can not nil")
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInvalidArgument, cerr)
	}

	if err := h.masterUsecase.CreateSong(
		ctx, req.Msg.GetName(), req.Msg.GetKana(),
		req.Msg.GetLyricsId(), req.Msg.GetMusicId(), req.Msg.GetArrangementId(),
		req.Msg.GetThumbnail(), req.Msg.GetOriginalVideo(), req.Msg.GetReleaseTime().AsTime(), req.Msg.GetDeleted(),
		req.Msg.GetUnitIds(),
		req.Msg.GetMusicVideoTypes(),
	); err != nil {
		if errors.Is(err, usecase.ErrInvalidArgument) {
			cerr := errors.WithStack(err)
			log.Printf("%+v\n", cerr)
			return nil, connect.NewError(connect.CodeInvalidArgument, cerr)
		}
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}

	return connect.NewResponse(&proto_master.CreateSongResponse{}), nil
}

// Chart
func (h *MasterHandler) GetCharts(ctx context.Context, req *connect.Request[proto_master.GetChartsRequest]) (*connect.Response[proto_master.GetChartsResponse], error) {
	charts, err := h.masterUsecase.ListCharts(ctx)
	if err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}

	var protoCharts []*proto_master.Chart
	for _, chart := range charts {
		var protoVocalPatterns []*proto_master.VocalPattern
		for _, vp := range chart.Song.VocalPatterns {
			if vp == nil {
				continue
			}
			sort.Slice(vp.Singers, func(i, j int) bool {
				return vp.Singers[i].Position < vp.Singers[j].Position
			})
			var protoSingers []*proto_master.Singer
			for _, s := range vp.Singers {
				if s == nil {
					continue
				}
				protoSingers = append(protoSingers, &proto_master.Singer{
					Id:   s.ID,
					Name: s.Name,
				})
			}
			protoVocalPatterns = append(protoVocalPatterns, &proto_master.VocalPattern{
				Id:      vp.ID,
				Name:    vp.Name,
				Singers: protoSingers,
			})
		}

		var protoUnits []*proto_master.Unit
		for _, u := range chart.Song.Units {
			if u == nil {
				continue
			}
			protoUnits = append(protoUnits, &proto_master.Unit{
				Id:   u.ID,
				Name: u.Name,
			})
		}
		protoSong := proto_master.Song{
			Id:   chart.Song.ID,
			Name: chart.Song.Name,
			Kana: chart.Song.Kana,
			Lyrics: &proto_master.Artist{
				Id:   chart.Song.Lyrics.ID,
				Name: chart.Song.Lyrics.Name,
				Kana: chart.Song.Lyrics.Kana,
			},
			Music: &proto_master.Artist{
				Id:   chart.Song.Music.ID,
				Name: chart.Song.Music.Name,
				Kana: chart.Song.Music.Kana,
			},
			Arrangement: &proto_master.Artist{
				Id:   chart.Song.Arrangement.ID,
				Name: chart.Song.Arrangement.Name,
				Kana: chart.Song.Arrangement.Kana,
			},
			Thumbnail:       chart.Song.Thumbnail,
			OriginalVideo:   chart.Song.OriginalVideo,
			ReleaseTime:     timestamppb.New(chart.Song.ReleaseTime),
			Deleted:         chart.Song.Deleted,
			VocalPatterns:   protoVocalPatterns,
			Units:           protoUnits,
			MusicVideoTypes: chart.Song.MusicVideoTypes,
		}

		protoCharts = append(protoCharts, &proto_master.Chart{
			Id:             chart.ID,
			Song:           &protoSong,
			DifficultyType: chart.DifficultyType,
			Level:          chart.Level,
			ChartViewLink:  chart.ChartViewLink,
		})
	}

	return connect.NewResponse(&proto_master.GetChartsResponse{
		Charts: protoCharts,
	}), nil
}
func (h *MasterHandler) GetChart(ctx context.Context, req *connect.Request[proto_master.GetChartRequest]) (*connect.Response[proto_master.GetChartResponse], error) {
	chart, err := h.masterUsecase.GetChartByID(ctx, req.Msg.GetId())
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			cerr := errors.WithStack(err)
			log.Printf("%+v\n", cerr)
			return nil, connect.NewError(connect.CodeNotFound, cerr)
		}
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}

	var protoChart *proto_master.Chart
	var protoVocalPatterns []*proto_master.VocalPattern
	for _, vp := range chart.Song.VocalPatterns {
		if vp == nil {
			continue
		}
		sort.Slice(vp.Singers, func(i, j int) bool {
			return vp.Singers[i].Position < vp.Singers[j].Position
		})
		var protoSingers []*proto_master.Singer
		for _, s := range vp.Singers {
			if s == nil {
				continue
			}
			protoSingers = append(protoSingers, &proto_master.Singer{
				Id:   s.ID,
				Name: s.Name,
			})
		}
		protoVocalPatterns = append(protoVocalPatterns, &proto_master.VocalPattern{
			Id:      vp.ID,
			Name:    vp.Name,
			Singers: protoSingers,
		})
	}

	var protoUnits []*proto_master.Unit
	for _, u := range chart.Song.Units {
		if u == nil {
			continue
		}
		protoUnits = append(protoUnits, &proto_master.Unit{
			Id:   u.ID,
			Name: u.Name,
		})
	}
	protoSong := proto_master.Song{
		Id:   chart.Song.ID,
		Name: chart.Song.Name,
		Kana: chart.Song.Kana,
		Lyrics: &proto_master.Artist{
			Id:   chart.Song.Lyrics.ID,
			Name: chart.Song.Lyrics.Name,
			Kana: chart.Song.Lyrics.Kana,
		},
		Music: &proto_master.Artist{
			Id:   chart.Song.Music.ID,
			Name: chart.Song.Music.Name,
			Kana: chart.Song.Music.Kana,
		},
		Arrangement: &proto_master.Artist{
			Id:   chart.Song.Arrangement.ID,
			Name: chart.Song.Arrangement.Name,
			Kana: chart.Song.Arrangement.Kana,
		},
		Thumbnail:       chart.Song.Thumbnail,
		OriginalVideo:   chart.Song.OriginalVideo,
		ReleaseTime:     timestamppb.New(chart.Song.ReleaseTime),
		Deleted:         chart.Song.Deleted,
		VocalPatterns:   protoVocalPatterns,
		Units:           protoUnits,
		MusicVideoTypes: chart.Song.MusicVideoTypes,
	}

	protoChart = &proto_master.Chart{
		Id:             chart.ID,
		Song:           &protoSong,
		DifficultyType: chart.DifficultyType,
		Level:          chart.Level,
		ChartViewLink:  chart.ChartViewLink,
	}

	return connect.NewResponse(&proto_master.GetChartResponse{
		Chart: protoChart,
	}), nil
}
func (h *MasterHandler) CreateChart(ctx context.Context, req *connect.Request[proto_master.CreateChartRequest]) (*connect.Response[proto_master.CreateChartResponse], error) {
	id, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok {
		err := errors.New("user id not found in context")
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeUnauthenticated, cerr)
	}
	isAdmin, err := h.userUsecase.IsAdmin(ctx, id)
	if err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}
	if !isAdmin {
		err := errors.New("permission denied: not admin")
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodePermissionDenied, cerr)
	}

	if err := req.Msg.Validate(); err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInvalidArgument, cerr)
	}

	if err := h.masterUsecase.CreateChart(
		ctx, req.Msg.GetSongId(), int32(req.Msg.GetDifficultyType()), req.Msg.GetLevel(), req.Msg.GetChartViewLink(),
	); err != nil {
		if errors.Is(err, usecase.ErrInvalidArgument) {
			cerr := errors.WithStack(err)
			log.Printf("%+v\n", cerr)
			return nil, connect.NewError(connect.CodeInvalidArgument, cerr)
		}
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}

	return connect.NewResponse(&proto_master.CreateChartResponse{}), nil
}
