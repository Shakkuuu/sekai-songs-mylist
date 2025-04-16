package repository

import (
	"context"

	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/entity"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/infrastructure/db/sqlcgen"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/$GOFILE -package=mock_$GOPACKAGE

type MasterRepository interface {
	// Artist
	ListArtists(ctx context.Context) ([]*sqlcgen.Artist, error)
	GetArtistByID(ctx context.Context, id int32) (*sqlcgen.Artist, error)
	CreateArtist(ctx context.Context, artist *entity.Artist) (error)
}
