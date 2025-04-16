package repository

import (
	"context"

	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/entity"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/repository"
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

func (r *masterRepository) ListArtists(ctx context.Context) ([]*sqlcgen.Artist, error) {
	artists, err := r.queries.ListArtists(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	artistPointers := make([]*sqlcgen.Artist, len(artists))
    for i := range artists {
        artistPointers[i] = &artists[i]
    }

	return artistPointers, nil
}

func (r *masterRepository) GetArtistByID(ctx context.Context, id int32) (*sqlcgen.Artist, error) {
	artist, err := r.queries.GetArtistByID(ctx, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &artist, nil
}

func (r *masterRepository) CreateArtist(ctx context.Context, artist *entity.Artist) error {
	sqlArtist := sqlcgen.InsertArtistParams{
		Name: artist.Name,
		Kana: artist.Kana,
	}
	if err := r.queries.InsertArtist(ctx, sqlArtist); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
