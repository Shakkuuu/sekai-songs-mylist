package repository

import (
	"context"
	"time"

	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/entity"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/gen/enums"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/infrastructure/db/sqlcgen"
	"github.com/google/uuid"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/$GOFILE -package=mock_$GOPACKAGE

type MyListRepository interface {
	// MyList
	GetMyListByID(ctx context.Context, id int32) (*entity.MyList, error)
	ListMyListsByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.MyList, error)
	CreateMyList(ctx context.Context, userID uuid.UUID, name string, position int32, createdAt, updatedAt time.Time) (*sqlcgen.MyList, error)
	ExistsMyList(ctx context.Context, id int32) (bool, error)
	UpdateMyListName(ctx context.Context, id int32, name string, updatedAt time.Time) error
	UpdateMyListPosition(ctx context.Context, id int32, position int32, updatedAt time.Time) error
	DeleteMyList(ctx context.Context, id int32) error

	// MyListChart
	GetMyListChartByID(ctx context.Context, id int32) (*entity.MyListChart, error)
	ListMyListChartsByMyListID(ctx context.Context, myListID int32) ([]*entity.MyListChart, error)
	CreateMyListChart(ctx context.Context, myListID, chartID int32, clearType enums.ClearType, memo string, createdAt, updatedAt time.Time) (*sqlcgen.MyListChart, error)
	ExistsMyListChart(ctx context.Context, id int32) (bool, error)
	ExistsMyListChartByMyListIDAndChartID(ctx context.Context, myListID, chartID int32) (bool, error)
	UpdateMyListChartClearType(ctx context.Context, id int32, clearType enums.ClearType, updatedAt time.Time) error
	UpdateMyListChartMemo(ctx context.Context, id int32, memo string, updatedAt time.Time) error
	DeleteMyListChart(ctx context.Context, id int32) error
	DeleteMyListChartByMyListID(ctx context.Context, myListID int32) error

	// MyListChartAttachment
	GetMyListChartAttachmentByID(ctx context.Context, id int32) (*entity.MyListChartAttachment, error)
	ListMyListChartAttachmentsByMyListChartID(ctx context.Context, myListChartID int32) ([]*entity.MyListChartAttachment, error)
	CreateMyListChartAttachment(ctx context.Context, myListChartID int32, attachmentType enums.AttachmentType, fileURL, caption string, createdAt time.Time) (*sqlcgen.MyListChartAttachment, error)
	ExistsMyListChartAttachment(ctx context.Context, id int32) (bool, error)
	DeleteMyListChartAttachment(ctx context.Context, id int32) error
	DeleteMyListChartAttachmentByMyListChartID(ctx context.Context, myListChartID int32) error
}
