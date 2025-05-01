package usecase

import (
	"context"

	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/entity"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/repository"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/infrastructure/db/sqlcgen"
	"github.com/cockroachdb/errors"
)

//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/$GOFILE -package=mock_$GOPACKAGE
//go:generate gotests -w -all $GOFILE

type MasterUsecase interface {
	ListArtists(ctx context.Context) ([]*entity.Artist, error)
	GetArtistByID(ctx context.Context, id int32) (*entity.Artist, error)
	CreateArtist(ctx context.Context, name, kana string) (error)
}

type masterUsecase struct {
	masterRepo  repository.MasterRepository
}

func NewMasterUsecase(repo repository.MasterRepository) MasterUsecase {
	return &masterUsecase{
		masterRepo:  repo,
	}
}

func (u *masterUsecase) ListArtists(ctx context.Context) ([]*entity.Artist, error) {
	artists, err := u.masterRepo.ListArtists(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if artists == nil {
		return []*entity.Artist{}, nil
	}

	domainArtists := make([]*entity.Artist, len(artists))
	for i, artist := range artists {
		domainArtists[i] = sqlToDomainArtist(artist)
	}

	return domainArtists, nil
}

func (u *masterUsecase) GetArtistByID(ctx context.Context, id int32) (*entity.Artist, error) {
	artist, err := u.masterRepo.GetArtistByID(ctx, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return sqlToDomainArtist(artist), nil
}

func (u *masterUsecase) CreateArtist(ctx context.Context, name, kana string) error {
	artist := entity.Artist{
		Name: name,
		Kana: kana,
	}
	if err := u.masterRepo.CreateArtist(ctx, &artist); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func sqlToDomainArtist(sqlArtist *sqlcgen.Artist) (*entity.Artist) {
	return &entity.Artist{
		ID: sqlArtist.ID,
		Name: sqlArtist.Name,
		Kana: sqlArtist.Kana,
	}
}
