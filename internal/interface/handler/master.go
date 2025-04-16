package handler

import (
	"context"

	proto_master "github.com/Shakkuuu/sekai-songs-mylist/internal/gen/master"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/usecase"
	"github.com/cockroachdb/errors"
)

//go:generate gotests -w -all $GOFILE

type MasterHandler struct {
	proto_master.UnimplementedMasterServiceServer
    masterUsecase usecase.MasterUsecase
}

func NewMasterHandler(masterUsecase usecase.MasterUsecase) *MasterHandler {
    return &MasterHandler{masterUsecase: masterUsecase}
}

func (h *MasterHandler) GetArtists(ctx context.Context, req *proto_master.GetArtistsRequest) (*proto_master.GetArtistsResponse, error) {
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

    return &proto_master.GetArtistsResponse{
        Artists: protoArtists,
    }, nil
}

func (h *MasterHandler) GetArtist(ctx context.Context, req *proto_master.GetArtistRequest) (*proto_master.GetArtistResponse, error) {
    artist, err := h.masterUsecase.GetArtistByID(ctx, req.Id)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    return &proto_master.GetArtistResponse{
        Artists: &proto_master.Artist{
            Id:   int32(artist.ID),
            Name: artist.Name,
            Kana: artist.Kana,
        },
    }, nil
}

func (h *MasterHandler) CreateArtist(ctx context.Context, req *proto_master.CreateArtistRequest) (*proto_master.CreateArtistResponse, error) {
    err := h.masterUsecase.CreateArtist(ctx, req.Name, req.Kana)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    return &proto_master.CreateArtistResponse{}, nil
}
