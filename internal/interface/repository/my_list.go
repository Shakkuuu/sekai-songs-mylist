package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/entity"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/repository"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/gen/enums"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/infrastructure/db/sqlcgen"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

//go:generate gotests -w -all $GOFILE

type myListRepository struct {
	queries *sqlcgen.Queries
}

func NewMyListRepository(queries *sqlcgen.Queries) repository.MyListRepository {
	return &myListRepository{queries: queries}
}

// MyList
func (r *myListRepository) GetMyListByID(ctx context.Context, id int32) (*entity.MyList, error) {
	myList, err := r.queries.GetMyListByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.WithStack(repository.ErrNotFound)
		}
		return nil, errors.WithStack(err)
	}

	return sqlToDomainMyList(&myList), nil
}

func (r *myListRepository) ListMyListsByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.MyList, error) {
	myLists, err := r.queries.ListMyListsByUserID(ctx, uuid.NullUUID{UUID: userID, Valid: true})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	myListPointers := make([]*entity.MyList, len(myLists))
	for i := range myLists {
		myListPointers[i] = sqlToDomainMyList(&myLists[i])
	}

	return myListPointers, nil
}

func (r *myListRepository) CreateMyList(ctx context.Context, userID uuid.UUID, name string, position int32, createdAt, updatedAt time.Time) (*sqlcgen.MyList, error) {
	sqlMyList := sqlcgen.InsertMyListParams{
		UserID:    uuid.NullUUID{UUID: userID, Valid: true},
		Name:      name,
		Position:  position,
		CreatedAt: sql.NullTime{Time: createdAt, Valid: true},
		UpdatedAt: sql.NullTime{Time: updatedAt, Valid: true},
	}

	l, err := r.queries.InsertMyList(ctx, sqlMyList)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &l, nil
}

func (r *myListRepository) ExistsMyList(ctx context.Context, id int32) (bool, error) {
	exist, err := r.queries.ExistsMyList(ctx, id)
	if err != nil {
		return false, errors.WithStack(err)
	}

	return exist, nil
}

func (r *myListRepository) UpdateMyListName(ctx context.Context, id int32, name string, updatedAt time.Time) error {
	arg := sqlcgen.UpdateMyListNameParams{
		Name:      name,
		UpdatedAt: sql.NullTime{Time: updatedAt, Valid: true},
		ID:        id,
	}

	if err := r.queries.UpdateMyListName(ctx, arg); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *myListRepository) UpdateMyListPosition(ctx context.Context, id int32, position int32, updatedAt time.Time) error {
	arg := sqlcgen.UpdateMyListPositionParams{
		Position:  position,
		UpdatedAt: sql.NullTime{Time: updatedAt, Valid: true},
		ID:        id,
	}

	if err := r.queries.UpdateMyListPosition(ctx, arg); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *myListRepository) DeleteMyList(ctx context.Context, id int32) error {
	if err := r.queries.DeleteMyList(ctx, id); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func sqlToDomainMyList(sqlMyList *sqlcgen.MyList) *entity.MyList {
	return &entity.MyList{
		ID:        sqlMyList.ID,
		UserID:    sqlMyList.UserID.UUID.String(),
		Name:      sqlMyList.Name,
		Position:  sqlMyList.Position,
		CreatedAt: sqlMyList.CreatedAt.Time,
		UpdatedAt: sqlMyList.UpdatedAt.Time,
	}
}

// MyListChart
func (r *myListRepository) GetMyListChartByID(ctx context.Context, id int32) (*entity.MyListChart, error) {
	myListChart, err := r.queries.GetMyListChartByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.WithStack(repository.ErrNotFound)
		}
		return nil, errors.WithStack(err)
	}

	return sqlToDomainMyListChart(&myListChart), nil
}

func (r *myListRepository) ListMyListChartsByMyListID(ctx context.Context, myListID int32) ([]*entity.MyListChart, error) {
	myListCharts, err := r.queries.ListMyListChartsByMyListID(ctx, sql.NullInt32{Int32: myListID, Valid: true})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	myListChartPointers := make([]*entity.MyListChart, len(myListCharts))
	for i := range myListCharts {
		myListChartPointers[i] = sqlToDomainMyListChart(&myListCharts[i])
	}

	return myListChartPointers, nil
}

func (r *myListRepository) CreateMyListChart(ctx context.Context, myListID, chartID int32, clearType enums.ClearType, memo string, createdAt, updatedAt time.Time) (*sqlcgen.MyListChart, error) {
	sqlMyListChart := sqlcgen.InsertMyListChartParams{
		MyListID:  sql.NullInt32{Int32: myListID, Valid: true},
		ChartID:   sql.NullInt32{Int32: chartID, Valid: true},
		ClearType: sql.NullInt32{Int32: int32(clearType), Valid: true},
		Memo:      sql.NullString{String: memo, Valid: true},
		CreatedAt: sql.NullTime{Time: createdAt, Valid: true},
		UpdatedAt: sql.NullTime{Time: updatedAt, Valid: true},
	}

	l, err := r.queries.InsertMyListChart(ctx, sqlMyListChart)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &l, nil
}

func (r *myListRepository) ExistsMyListChart(ctx context.Context, id int32) (bool, error) {
	exist, err := r.queries.ExistsMyListChart(ctx, id)
	if err != nil {
		return false, errors.WithStack(err)
	}

	return exist, nil
}

func (r *myListRepository) ExistsMyListChartByMyListIDAndChartID(ctx context.Context, myListID, chartID int32) (bool, error) {
	arg := sqlcgen.ExistsMyListChartByMyListIDAndChartIDParams{
		MyListID: sql.NullInt32{Int32: myListID, Valid: true},
		ChartID:  sql.NullInt32{Int32: chartID, Valid: true},
	}
	exist, err := r.queries.ExistsMyListChartByMyListIDAndChartID(ctx, arg)
	if err != nil {
		return false, errors.WithStack(err)
	}

	return exist, nil
}

func (r *myListRepository) UpdateMyListChartClearType(ctx context.Context, id int32, clearType enums.ClearType, updatedAt time.Time) error {
	arg := sqlcgen.UpdateMyListChartClearTypeParams{
		ClearType: sql.NullInt32{Int32: int32(clearType), Valid: true},
		UpdatedAt: sql.NullTime{Time: updatedAt, Valid: true},
		ID:        id,
	}

	if err := r.queries.UpdateMyListChartClearType(ctx, arg); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *myListRepository) UpdateMyListChartMemo(ctx context.Context, id int32, memo string, updatedAt time.Time) error {
	arg := sqlcgen.UpdateMyListChartMemoParams{
		Memo:      sql.NullString{String: memo, Valid: true},
		UpdatedAt: sql.NullTime{Time: updatedAt, Valid: true},
		ID:        id,
	}

	if err := r.queries.UpdateMyListChartMemo(ctx, arg); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *myListRepository) DeleteMyListChart(ctx context.Context, id int32) error {
	if err := r.queries.DeleteMyListChart(ctx, id); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (r *myListRepository) DeleteMyListChartByMyListID(ctx context.Context, myListID int32) error {
	if err := r.queries.DeleteMyListChartByMyListID(ctx, sql.NullInt32{Int32: myListID, Valid: true}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func sqlToDomainMyListChart(sqlMyListChart *sqlcgen.MyListChart) *entity.MyListChart {
	return &entity.MyListChart{
		ID:        sqlMyListChart.ID,
		MyListID:  sqlMyListChart.MyListID.Int32,
		ChartID:   sqlMyListChart.ChartID.Int32,
		ClearType: enums.ClearType(sqlMyListChart.ClearType.Int32),
		Memo:      sqlMyListChart.Memo.String,
		CreatedAt: sqlMyListChart.CreatedAt.Time,
		UpdatedAt: sqlMyListChart.UpdatedAt.Time,
	}
}

// MyListChartAttachment
func (r *myListRepository) GetMyListChartAttachmentByID(ctx context.Context, id int32) (*entity.MyListChartAttachment, error) {
	myListChartAttachment, err := r.queries.GetMyListChartAttachmentByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.WithStack(repository.ErrNotFound)
		}
		return nil, errors.WithStack(err)
	}

	return sqlToDomainMyListChartAttachment(&myListChartAttachment), nil
}

func (r *myListRepository) ListMyListChartAttachmentsByMyListChartID(ctx context.Context, myListChartID int32) ([]*entity.MyListChartAttachment, error) {
	myListChartAttachments, err := r.queries.ListMyListChartAttachmentsByMyListChartID(ctx, sql.NullInt32{Int32: myListChartID, Valid: true})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	myListChartAttachmentPointers := make([]*entity.MyListChartAttachment, len(myListChartAttachments))
	for i := range myListChartAttachments {
		myListChartAttachmentPointers[i] = sqlToDomainMyListChartAttachment(&myListChartAttachments[i])
	}

	return myListChartAttachmentPointers, nil
}

func (r *myListRepository) CreateMyListChartAttachment(ctx context.Context, myListChartID int32, attachmentType enums.AttachmentType, fileURL, caption string, createdAt time.Time) (*sqlcgen.MyListChartAttachment, error) {
	sqlMyListChartAttachment := sqlcgen.InsertMyListChartAttachmentParams{
		MyListChartID:  sql.NullInt32{Int32: myListChartID, Valid: true},
		AttachmentType: sql.NullInt32{Int32: int32(attachmentType), Valid: true},
		FileUrl:        sql.NullString{String: fileURL, Valid: true},
		Caption:        sql.NullString{String: caption, Valid: true},
		CreatedAt:      sql.NullTime{Time: createdAt, Valid: true},
	}

	l, err := r.queries.InsertMyListChartAttachment(ctx, sqlMyListChartAttachment)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &l, nil
}

func (r *myListRepository) ExistsMyListChartAttachment(ctx context.Context, id int32) (bool, error) {
	exist, err := r.queries.ExistsMyListChartAttachment(ctx, id)
	if err != nil {
		return false, errors.WithStack(err)
	}

	return exist, nil
}

func (r *myListRepository) DeleteMyListChartAttachment(ctx context.Context, id int32) error {
	if err := r.queries.DeleteMyListChartAttachment(ctx, id); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (r *myListRepository) DeleteMyListChartAttachmentByMyListChartID(ctx context.Context, myListChartID int32) error {
	if err := r.queries.DeleteMyListChartAttachmentByMyListChartID(ctx, sql.NullInt32{Int32: myListChartID, Valid: true}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func sqlToDomainMyListChartAttachment(sqlMyListChartAttachment *sqlcgen.MyListChartAttachment) *entity.MyListChartAttachment {
	return &entity.MyListChartAttachment{
		ID:             sqlMyListChartAttachment.ID,
		MyListChartID:  sqlMyListChartAttachment.MyListChartID.Int32,
		AttachmentType: enums.AttachmentType(sqlMyListChartAttachment.AttachmentType.Int32),
		FileURL:        sqlMyListChartAttachment.FileUrl.String,
		Caption:        sqlMyListChartAttachment.Caption.String,
		CreatedAt:      sqlMyListChartAttachment.CreatedAt.Time,
	}
}
