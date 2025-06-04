package usecase

import (
	"context"
	"time"

	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/entity"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/repository"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/gen/enums"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/$GOFILE -package=mock_$GOPACKAGE
//go:generate gotests -w -all $GOFILE

var (
	ErrMyListNotFound                = errors.New("this id mylist not found")
	ErrMyListChartNotFound           = errors.New("this id mylistchart not found")
	ErrMyListChartAttachmentNotFound = errors.New("this id mylistchartattachment not found")
	ErrDuplicateMyListChart          = errors.New("this chart already exists")
)

type MyListUsecase interface {
	GetMyListByID(ctx context.Context, id int32) (*entity.MyList, error)
	GetMyListsByUserID(ctx context.Context, userID string) ([]*entity.MyList, error)
	CreateMyList(ctx context.Context, userID string, name string, position int32) error
	ChangeMyListName(ctx context.Context, id int32, name string) error
	ChangeMyListPosition(ctx context.Context, id, position []int32) error
	DeleteMyList(ctx context.Context, id int32) error
	GetMyListChartByID(ctx context.Context, id int32) (*entity.MyListChart, error)
	GetMyListChartsByMyListID(ctx context.Context, myListID int32) ([]*entity.MyListChart, error)
	AddMyListChart(ctx context.Context, myListID, chartID int32, clearType enums.ClearType, memo string) error
	ChangeMyListChartClearType(ctx context.Context, id int32, clearType enums.ClearType) error
	ChangeMyListChartMemo(ctx context.Context, id int32, memo string) error
	DeleteMyListChart(ctx context.Context, id int32) error
	GetMyListChartAttachmentsByMyListChartID(ctx context.Context, myListChartID int32) ([]*entity.MyListChartAttachment, error)
	AddMyListChartAttachment(ctx context.Context, myListChartID int32, attachmentType enums.AttachmentType, fileURL, caption string) error
	DeleteMyListChartAttachment(ctx context.Context, id int32) error
}

type myListUsecase struct {
	myListRepo           repository.MyListRepository
	masterRepo           repository.MasterRepository
	redisMasterCacheRepo repository.RedisMasterCacheRepository
}

func NewMyListUsecase(repo repository.MyListRepository, masterRepo repository.MasterRepository, redisMasterCacheRepo repository.RedisMasterCacheRepository) MyListUsecase {
	return &myListUsecase{
		myListRepo:           repo,
		masterRepo:           masterRepo,
		redisMasterCacheRepo: redisMasterCacheRepo,
	}
}

func (u *myListUsecase) GetMyListByID(ctx context.Context, id int32) (*entity.MyList, error) {
	myList, err := u.myListRepo.GetMyListByID(ctx, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return myList, nil
}

func (u *myListUsecase) GetMyListsByUserID(ctx context.Context, userID string) ([]*entity.MyList, error) {
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	myLists, err := u.myListRepo.ListMyListsByUserID(ctx, parsedID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if myLists == nil {
		return []*entity.MyList{}, nil
	}

	return myLists, nil
}

func (u *myListUsecase) CreateMyList(ctx context.Context, userID string, name string, position int32) error {
	// TODO: リストの個数制限を設けるならここでチェック

	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return errors.WithStack(err)
	}
	now := time.Now()
	createdAt := now
	updatedAt := now
	if _, err := u.myListRepo.CreateMyList(ctx, parsedID, name, position, createdAt, updatedAt); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (u *myListUsecase) ChangeMyListName(ctx context.Context, id int32, name string) error {
	if err := u.myListRepo.UpdateMyListName(ctx, id, name, time.Now()); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (u *myListUsecase) ChangeMyListPosition(ctx context.Context, id, position []int32) error {
	for i := range len(id) {
		if err := u.myListRepo.UpdateMyListPosition(ctx, id[i], position[i], time.Now()); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (u *myListUsecase) DeleteMyList(ctx context.Context, id int32) error {
	exist, err := u.myListRepo.ExistsMyList(ctx, id)
	if err != nil {
		return errors.WithStack(err)
	}
	if !exist {
		return errors.WithStack(ErrMyListNotFound)
	}

	myListCharts, err := u.myListRepo.ListMyListChartsByMyListID(ctx, id)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, myListChart := range myListCharts {
		if err := u.myListRepo.DeleteMyListChartAttachmentByMyListChartID(ctx, myListChart.ID); err != nil {
			return errors.WithStack(err)
		}
	}

	if err := u.myListRepo.DeleteMyListChartByMyListID(ctx, id); err != nil {
		return errors.WithStack(err)
	}

	if err := u.myListRepo.DeleteMyList(ctx, id); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (u *myListUsecase) GetMyListChartByID(ctx context.Context, id int32) (*entity.MyListChart, error) {
	myListChart, err := u.myListRepo.GetMyListChartByID(ctx, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	chart, err := u.redisMasterCacheRepo.GetChartByID(ctx, myListChart.ChartID)
	if errors.Is(err, redis.Nil) {
		chart, err := u.masterRepo.GetChartByID(ctx, myListChart.ChartID)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if err := u.redisMasterCacheRepo.SetChart(ctx, myListChart.ChartID, chart); err != nil {
			return nil, errors.WithStack(err)
		}
		myListChart.Chart = chart
		return myListChart, nil
	} else if err != nil {
		return nil, errors.WithStack(err)
	}
	myListChart.Chart = chart

	return myListChart, nil
}

func (u *myListUsecase) GetMyListChartsByMyListID(ctx context.Context, myListID int32) ([]*entity.MyListChart, error) {
	myListCharts, err := u.myListRepo.ListMyListChartsByMyListID(ctx, myListID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if myListCharts == nil {
		return []*entity.MyListChart{}, nil
	}

	for i := range len(myListCharts) {
		chart, err := u.redisMasterCacheRepo.GetChartByID(ctx, myListCharts[i].ChartID)
		if errors.Is(err, redis.Nil) {
			chart, err := u.masterRepo.GetChartByID(ctx, myListCharts[i].ChartID)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			if err := u.redisMasterCacheRepo.SetChart(ctx, myListCharts[i].ChartID, chart); err != nil {
				return nil, errors.WithStack(err)
			}
			myListCharts[i].Chart = chart
			continue
		} else if err != nil {
			return nil, errors.WithStack(err)
		}
		myListCharts[i].Chart = chart
	}

	return myListCharts, nil
}

func (u *myListUsecase) AddMyListChart(ctx context.Context, myListID, chartID int32, clearType enums.ClearType, memo string) error {
	exist, err := u.myListRepo.ExistsMyListChartByMyListIDAndChartID(ctx, myListID, chartID)
	if err != nil {
		return errors.WithStack(err)
	}
	if exist {
		return errors.WithStack(ErrDuplicateMyListChart)
	}

	now := time.Now()
	createdAt := now
	updatedAt := now
	if _, err := u.myListRepo.CreateMyListChart(ctx, myListID, chartID, clearType, memo, createdAt, updatedAt); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (u *myListUsecase) ChangeMyListChartClearType(ctx context.Context, id int32, clearType enums.ClearType) error {
	if err := u.myListRepo.UpdateMyListChartClearType(ctx, id, clearType, time.Now()); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (u *myListUsecase) ChangeMyListChartMemo(ctx context.Context, id int32, memo string) error {
	if err := u.myListRepo.UpdateMyListChartMemo(ctx, id, memo, time.Now()); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (u *myListUsecase) DeleteMyListChart(ctx context.Context, id int32) error {
	exist, err := u.myListRepo.ExistsMyListChart(ctx, id)
	if err != nil {
		return errors.WithStack(err)
	}
	if !exist {
		return errors.WithStack(ErrMyListChartNotFound)
	}

	if err := u.myListRepo.DeleteMyListChartAttachmentByMyListChartID(ctx, id); err != nil {
		return errors.WithStack(err)
	}

	if err := u.myListRepo.DeleteMyListChart(ctx, id); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (u *myListUsecase) GetMyListChartAttachmentsByMyListChartID(ctx context.Context, myListChartID int32) ([]*entity.MyListChartAttachment, error) {
	myListChartAttachments, err := u.myListRepo.ListMyListChartAttachmentsByMyListChartID(ctx, myListChartID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if myListChartAttachments == nil {
		return []*entity.MyListChartAttachment{}, nil
	}

	return myListChartAttachments, nil
}

func (u *myListUsecase) AddMyListChartAttachment(ctx context.Context, myListChartID int32, attachmentType enums.AttachmentType, fileURL, caption string) error {
	now := time.Now()
	createdAt := now
	if _, err := u.myListRepo.CreateMyListChartAttachment(ctx, myListChartID, attachmentType, fileURL, caption, createdAt); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (u *myListUsecase) DeleteMyListChartAttachment(ctx context.Context, id int32) error {
	exist, err := u.myListRepo.ExistsMyListChartAttachment(ctx, id)
	if err != nil {
		return errors.WithStack(err)
	}
	if !exist {
		return errors.WithStack(ErrMyListChartAttachmentNotFound)
	}

	if err := u.myListRepo.DeleteMyListChartAttachment(ctx, id); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
