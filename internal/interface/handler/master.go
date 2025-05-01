package handler

import (
	"context"

	"connectrpc.com/connect"

	proto_master "github.com/Shakkuuu/sekai-songs-mylist/internal/gen/master"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/usecase"
	"github.com/cockroachdb/errors"
)

//go:generate gotests -w -all $GOFILE

type MasterHandler struct {
    masterUsecase usecase.MasterUsecase
}

func NewMasterHandler(masterUsecase usecase.MasterUsecase) *MasterHandler {
    return &MasterHandler{masterUsecase: masterUsecase}
}

func (h *MasterHandler) GetArtists(ctx context.Context, req *connect.Request[proto_master.GetArtistsRequest]) (*connect.Response[proto_master.GetArtistsResponse], error) {
    artists, err := h.masterUsecase.ListArtists(ctx)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    var protoArtists []*proto_master.Artist
    for _, artist := range artists {
        protoArtists = append(protoArtists, &proto_master.Artist{
            Id:   artist.ID,
            Name: artist.Name,
            Kana: artist.Kana,
        })
    }

    return connect.NewResponse(&proto_master.GetArtistsResponse{
        Artists: protoArtists,
    }), nil
}

func (h *MasterHandler) GetArtist(ctx context.Context, req *connect.Request[proto_master.GetArtistRequest]) (*connect.Response[proto_master.GetArtistResponse], error) {
    artist, err := h.masterUsecase.GetArtistByID(ctx, req.Msg.GetId())
    if err != nil {
        return nil, errors.WithStack(err)
    }

    return connect.NewResponse(&proto_master.GetArtistResponse{
        Artists: &proto_master.Artist{
            Id:   int32(artist.ID),
            Name: artist.Name,
            Kana: artist.Kana,
        },
    }), nil
}

func (h *MasterHandler) CreateArtist(ctx context.Context, req *connect.Request[proto_master.CreateArtistRequest]) (*connect.Response[proto_master.CreateArtistResponse], error) {
    err := h.masterUsecase.CreateArtist(ctx, req.Msg.GetName(), req.Msg.GetKana())
    if err != nil {
        return nil, errors.WithStack(err)
    }

    return connect.NewResponse(&proto_master.CreateArtistResponse{}), nil
}

func (h *MasterHandler) GetSongs(ctx context.Context, req *connect.Request[proto_master.GetSongsRequest]) (*connect.Response[proto_master.GetSongsResponse], error) {
    return nil, connect.NewError(connect.CodeUnimplemented, errors.New("GetSong not implemented"))
}

func (h *MasterHandler) GetSong(ctx context.Context, req *connect.Request[proto_master.GetSongRequest]) (*connect.Response[proto_master.GetSongResponse], error) {
    return nil, connect.NewError(connect.CodeUnimplemented, errors.New("GetSong not implemented"))
}
