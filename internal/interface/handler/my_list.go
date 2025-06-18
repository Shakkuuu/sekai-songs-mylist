package handler

import (
	"context"
	"fmt"
	"log"
	"sort"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/repository"
	proto_master "github.com/Shakkuuu/sekai-songs-mylist/internal/gen/master"
	proto_my_list "github.com/Shakkuuu/sekai-songs-mylist/internal/gen/mylist/v1"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/infrastructure/auth"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/usecase"
	"github.com/cockroachdb/errors"
)

//go:generate gotests -w -all $GOFILE

type MyListHandler struct {
	myListUsecase usecase.MyListUsecase
}

func NewMyListHandler(myListUsecase usecase.MyListUsecase) *MyListHandler {
	return &MyListHandler{myListUsecase: myListUsecase}
}

// MyList
func (h *MyListHandler) GetMyListsByUserID(ctx context.Context, req *connect.Request[proto_my_list.GetMyListsByUserIDRequest]) (*connect.Response[proto_my_list.GetMyListsByUserIDResponse], error) {
	id, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("user id not found in context"))
	}

	myLists, err := h.myListUsecase.GetMyListsByUserID(ctx, id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}

	protoMyLists := make([]*proto_my_list.MyList, len(myLists))
	for i, myList := range myLists {
		protoMyLists[i] = &proto_my_list.MyList{
			Id:        myList.ID,
			UserId:    myList.UserID,
			Name:      myList.Name,
			Position:  myList.Position,
			CreatedAt: timestamppb.New(myList.CreatedAt),
			UpdatedAt: timestamppb.New(myList.UpdatedAt),
		}
	}

	return connect.NewResponse(&proto_my_list.GetMyListsByUserIDResponse{
		MyLists: protoMyLists,
	}), nil
}

func (h *MyListHandler) CreateMyList(ctx context.Context, req *connect.Request[proto_my_list.CreateMyListRequest]) (*connect.Response[proto_my_list.CreateMyListResponse], error) {
	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.WithStack(err))
	}

	id, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("user id not found in context"))
	}

	if err := h.myListUsecase.CreateMyList(ctx, id, req.Msg.GetName(), req.Msg.GetPosition()); err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}

	return connect.NewResponse(&proto_my_list.CreateMyListResponse{}), nil
}

func (h *MyListHandler) ChangeMyListName(ctx context.Context, req *connect.Request[proto_my_list.ChangeMyListNameRequest]) (*connect.Response[proto_my_list.ChangeMyListNameResponse], error) {
	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.WithStack(err))
	}

	err := h.myListUsecase.ChangeMyListName(ctx, req.Msg.GetId(), req.Msg.GetName())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}

	return connect.NewResponse(&proto_my_list.ChangeMyListNameResponse{}), nil
}

func (h *MyListHandler) ChangeMyListPosition(ctx context.Context, req *connect.Request[proto_my_list.ChangeMyListPositionRequest]) (*connect.Response[proto_my_list.ChangeMyListPositionResponse], error) {
	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.WithStack(err))
	}

	err := h.myListUsecase.ChangeMyListPosition(ctx, req.Msg.GetId(), req.Msg.GetPosition())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}

	return connect.NewResponse(&proto_my_list.ChangeMyListPositionResponse{}), nil
}

func (h *MyListHandler) DeleteMyList(ctx context.Context, req *connect.Request[proto_my_list.DeleteMyListRequest]) (*connect.Response[proto_my_list.DeleteMyListResponse], error) {
	if err := h.myListUsecase.DeleteMyList(ctx, req.Msg.GetId()); err != nil {
		if errors.Is(err, usecase.ErrMyListNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, errors.WithStack(err))
		}
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}
	return connect.NewResponse(&proto_my_list.DeleteMyListResponse{}), nil
}

// MyListChart
func (h *MyListHandler) GetMyListChartsByMyListID(ctx context.Context, req *connect.Request[proto_my_list.GetMyListChartsByMyListIDRequest]) (*connect.Response[proto_my_list.GetMyListChartsByMyListIDResponse], error) {
	myList, err := h.myListUsecase.GetMyListByID(ctx, req.Msg.GetMyListId())
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, errors.WithStack(err))
		}
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}

	myListCharts, err := h.myListUsecase.GetMyListChartsByMyListID(ctx, req.Msg.GetMyListId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}

	protoMyListCharts := make([]*proto_my_list.MyListChart, len(myListCharts))
	for i, myListChart := range myListCharts {
		var protoVocalPatterns []*proto_master.VocalPattern
		for _, vp := range myListChart.Chart.Song.VocalPatterns {
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
		for _, u := range myListChart.Chart.Song.Units {
			if u == nil {
				continue
			}
			protoUnits = append(protoUnits, &proto_master.Unit{
				Id:   u.ID,
				Name: u.Name,
			})
		}
		protoMyListCharts[i] = &proto_my_list.MyListChart{
			Id:       myListChart.ID,
			MyListId: myListChart.MyListID,
			Chart: &proto_master.Chart{
				Id: myListChart.Chart.ID,
				Song: &proto_master.Song{
					Id:   myListChart.Chart.Song.ID,
					Name: myListChart.Chart.Song.Name,
					Kana: myListChart.Chart.Song.Kana,
					Lyrics: &proto_master.Artist{
						Id:   myListChart.Chart.Song.Lyrics.ID,
						Name: myListChart.Chart.Song.Lyrics.Name,
						Kana: myListChart.Chart.Song.Lyrics.Kana,
					},
					Music: &proto_master.Artist{
						Id:   myListChart.Chart.Song.Music.ID,
						Name: myListChart.Chart.Song.Music.Name,
						Kana: myListChart.Chart.Song.Music.Kana,
					},
					Arrangement: &proto_master.Artist{
						Id:   myListChart.Chart.Song.Arrangement.ID,
						Name: myListChart.Chart.Song.Arrangement.Name,
						Kana: myListChart.Chart.Song.Arrangement.Kana,
					},
					Thumbnail:       myListChart.Chart.Song.Thumbnail,
					OriginalVideo:   myListChart.Chart.Song.OriginalVideo,
					ReleaseTime:     timestamppb.New(myListChart.Chart.Song.ReleaseTime),
					Deleted:         myListChart.Chart.Song.Deleted,
					VocalPatterns:   protoVocalPatterns,
					Units:           protoUnits,
					MusicVideoTypes: myListChart.Chart.Song.MusicVideoTypes,
				},
				DifficultyType: myListChart.Chart.DifficultyType,
				Level:          myListChart.Chart.Level,
				ChartViewLink:  myListChart.Chart.ChartViewLink,
			},
			ClearType: myListChart.ClearType,
			Memo:      myListChart.Memo,
			CreatedAt: timestamppb.New(myListChart.CreatedAt),
			UpdatedAt: timestamppb.New(myListChart.UpdatedAt),
		}
	}

	protoMyList := proto_my_list.MyList{
		Id:        myList.ID,
		UserId:    myList.UserID,
		Name:      myList.Name,
		Position:  myList.Position,
		CreatedAt: timestamppb.New(myList.CreatedAt),
		UpdatedAt: timestamppb.New(myList.UpdatedAt),
	}

	return connect.NewResponse(&proto_my_list.GetMyListChartsByMyListIDResponse{
		MyList:       &protoMyList,
		MyListCharts: protoMyListCharts,
	}), nil
}

func (h *MyListHandler) GetMyListChartByID(ctx context.Context, req *connect.Request[proto_my_list.GetMyListChartByIDRequest]) (*connect.Response[proto_my_list.GetMyListChartByIDResponse], error) {
	fmt.Println(req.Msg.GetId())
	myListChart, err := h.myListUsecase.GetMyListChartByID(ctx, req.Msg.GetId())
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			log.Println("point1")
			fmt.Printf("%v\n", myListChart)
			return nil, connect.NewError(connect.CodeNotFound, errors.WithStack(err))
		}
		log.Println("point2")
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}
	var protoVocalPatterns []*proto_master.VocalPattern
	for _, vp := range myListChart.Chart.Song.VocalPatterns {
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
	for _, u := range myListChart.Chart.Song.Units {
		if u == nil {
			continue
		}
		protoUnits = append(protoUnits, &proto_master.Unit{
			Id:   u.ID,
			Name: u.Name,
		})
	}

	protoMyListChart := &proto_my_list.MyListChart{
		Id:       myListChart.ID,
		MyListId: myListChart.MyListID,
		Chart: &proto_master.Chart{
			Id: myListChart.Chart.ID,
			Song: &proto_master.Song{
				Id:   myListChart.Chart.Song.ID,
				Name: myListChart.Chart.Song.Name,
				Kana: myListChart.Chart.Song.Kana,
				Lyrics: &proto_master.Artist{
					Id:   myListChart.Chart.Song.Lyrics.ID,
					Name: myListChart.Chart.Song.Lyrics.Name,
					Kana: myListChart.Chart.Song.Lyrics.Kana,
				},
				Music: &proto_master.Artist{
					Id:   myListChart.Chart.Song.Music.ID,
					Name: myListChart.Chart.Song.Music.Name,
					Kana: myListChart.Chart.Song.Music.Kana,
				},
				Arrangement: &proto_master.Artist{
					Id:   myListChart.Chart.Song.Arrangement.ID,
					Name: myListChart.Chart.Song.Arrangement.Name,
					Kana: myListChart.Chart.Song.Arrangement.Kana,
				},
				Thumbnail:       myListChart.Chart.Song.Thumbnail,
				OriginalVideo:   myListChart.Chart.Song.OriginalVideo,
				ReleaseTime:     timestamppb.New(myListChart.Chart.Song.ReleaseTime),
				Deleted:         myListChart.Chart.Song.Deleted,
				VocalPatterns:   protoVocalPatterns,
				Units:           protoUnits,
				MusicVideoTypes: myListChart.Chart.Song.MusicVideoTypes,
			},
			DifficultyType: myListChart.Chart.DifficultyType,
			Level:          myListChart.Chart.Level,
			ChartViewLink:  myListChart.Chart.ChartViewLink,
		},
		ClearType: myListChart.ClearType,
		Memo:      myListChart.Memo,
		CreatedAt: timestamppb.New(myListChart.CreatedAt),
		UpdatedAt: timestamppb.New(myListChart.UpdatedAt),
	}
	return connect.NewResponse(&proto_my_list.GetMyListChartByIDResponse{
		MyListChart: protoMyListChart,
	}), nil
}

func (h *MyListHandler) AddMyListChart(ctx context.Context, req *connect.Request[proto_my_list.AddMyListChartRequest]) (*connect.Response[proto_my_list.AddMyListChartResponse], error) {
	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.WithStack(err))
	}

	if err := h.myListUsecase.AddMyListChart(ctx, req.Msg.GetMyListId(), req.Msg.GetChartId(), req.Msg.GetClearType(), req.Msg.GetMemo()); err != nil {
		if errors.Is(err, usecase.ErrDuplicateMyListChart) {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.WithStack(err))
		}
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}

	return connect.NewResponse(&proto_my_list.AddMyListChartResponse{}), nil
}

func (h *MyListHandler) ChangeMyListChartClearType(ctx context.Context, req *connect.Request[proto_my_list.ChangeMyListChartClearTypeRequest]) (*connect.Response[proto_my_list.ChangeMyListChartClearTypeResponse], error) {
	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.WithStack(err))
	}

	err := h.myListUsecase.ChangeMyListChartClearType(ctx, req.Msg.GetId(), req.Msg.GetClearType())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}

	return connect.NewResponse(&proto_my_list.ChangeMyListChartClearTypeResponse{}), nil
}

func (h *MyListHandler) ChangeMyListChartMemo(ctx context.Context, req *connect.Request[proto_my_list.ChangeMyListChartMemoRequest]) (*connect.Response[proto_my_list.ChangeMyListChartMemoResponse], error) {
	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.WithStack(err))
	}

	err := h.myListUsecase.ChangeMyListChartMemo(ctx, req.Msg.GetId(), req.Msg.GetMemo())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}

	return connect.NewResponse(&proto_my_list.ChangeMyListChartMemoResponse{}), nil
}

func (h *MyListHandler) DeleteMyListChart(ctx context.Context, req *connect.Request[proto_my_list.DeleteMyListChartRequest]) (*connect.Response[proto_my_list.DeleteMyListChartResponse], error) {
	if err := h.myListUsecase.DeleteMyListChart(ctx, req.Msg.GetId()); err != nil {
		if errors.Is(err, usecase.ErrMyListChartNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, errors.WithStack(err))
		}
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}
	return connect.NewResponse(&proto_my_list.DeleteMyListChartResponse{}), nil
}

func (h *MyListHandler) GetMyListChartAttachmentsByMyListChartID(ctx context.Context, req *connect.Request[proto_my_list.GetMyListChartAttachmentsByMyListChartIDRequest]) (*connect.Response[proto_my_list.GetMyListChartAttachmentsByMyListChartIDResponse], error) {
	myListChart, err := h.myListUsecase.GetMyListChartByID(ctx, req.Msg.GetMyListChartId())
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, errors.WithStack(err))
		}
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}
	var protoVocalPatterns []*proto_master.VocalPattern
	for _, vp := range myListChart.Chart.Song.VocalPatterns {
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
	for _, u := range myListChart.Chart.Song.Units {
		if u == nil {
			continue
		}
		protoUnits = append(protoUnits, &proto_master.Unit{
			Id:   u.ID,
			Name: u.Name,
		})
	}

	protoMyListChart := &proto_my_list.MyListChart{
		Id:       myListChart.ID,
		MyListId: myListChart.MyListID,
		Chart: &proto_master.Chart{
			Id: myListChart.Chart.ID,
			Song: &proto_master.Song{
				Id:   myListChart.Chart.Song.ID,
				Name: myListChart.Chart.Song.Name,
				Kana: myListChart.Chart.Song.Kana,
				Lyrics: &proto_master.Artist{
					Id:   myListChart.Chart.Song.Lyrics.ID,
					Name: myListChart.Chart.Song.Lyrics.Name,
					Kana: myListChart.Chart.Song.Lyrics.Kana,
				},
				Music: &proto_master.Artist{
					Id:   myListChart.Chart.Song.Music.ID,
					Name: myListChart.Chart.Song.Music.Name,
					Kana: myListChart.Chart.Song.Music.Kana,
				},
				Arrangement: &proto_master.Artist{
					Id:   myListChart.Chart.Song.Arrangement.ID,
					Name: myListChart.Chart.Song.Arrangement.Name,
					Kana: myListChart.Chart.Song.Arrangement.Kana,
				},
				Thumbnail:       myListChart.Chart.Song.Thumbnail,
				OriginalVideo:   myListChart.Chart.Song.OriginalVideo,
				ReleaseTime:     timestamppb.New(myListChart.Chart.Song.ReleaseTime),
				Deleted:         myListChart.Chart.Song.Deleted,
				VocalPatterns:   protoVocalPatterns,
				Units:           protoUnits,
				MusicVideoTypes: myListChart.Chart.Song.MusicVideoTypes,
			},
			DifficultyType: myListChart.Chart.DifficultyType,
			Level:          myListChart.Chart.Level,
			ChartViewLink:  myListChart.Chart.ChartViewLink,
		},
		ClearType: myListChart.ClearType,
		Memo:      myListChart.Memo,
		CreatedAt: timestamppb.New(myListChart.CreatedAt),
		UpdatedAt: timestamppb.New(myListChart.UpdatedAt),
	}

	myListChartAttachments, err := h.myListUsecase.GetMyListChartAttachmentsByMyListChartID(ctx, req.Msg.GetMyListChartId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}

	protoMyListChartAttachments := make([]*proto_my_list.MyListChartAttachment, len(myListChartAttachments))
	for i, myListChartAttachment := range myListChartAttachments {
		protoMyListChartAttachments[i] = &proto_my_list.MyListChartAttachment{
			Id:             myListChartAttachment.ID,
			MyListChartId:  myListChartAttachment.MyListChartID,
			AttachmentType: myListChartAttachment.AttachmentType,
			FileUrl:        myListChartAttachment.FileURL,
			Caption:        myListChartAttachment.Caption,
			CreatedAt:      timestamppb.New(myListChartAttachment.CreatedAt),
		}
	}

	return connect.NewResponse(&proto_my_list.GetMyListChartAttachmentsByMyListChartIDResponse{
		MyListChart:            protoMyListChart,
		MyListChartAttachments: protoMyListChartAttachments,
	}), nil
}

func (h *MyListHandler) AddMyListChartAttachment(ctx context.Context, req *connect.Request[proto_my_list.AddMyListChartAttachmentRequest]) (*connect.Response[proto_my_list.AddMyListChartAttachmentResponse], error) {
	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.WithStack(err))
	}

	if err := h.myListUsecase.AddMyListChartAttachment(ctx, req.Msg.GetMyListChartId(), req.Msg.GetAttachmentType(), req.Msg.GetFileUrl(), req.Msg.GetCaption()); err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}

	return connect.NewResponse(&proto_my_list.AddMyListChartAttachmentResponse{}), nil
}

func (h *MyListHandler) DeleteMyListChartAttachment(ctx context.Context, req *connect.Request[proto_my_list.DeleteMyListChartAttachmentRequest]) (*connect.Response[proto_my_list.DeleteMyListChartAttachmentResponse], error) {
	if err := h.myListUsecase.DeleteMyListChartAttachment(ctx, req.Msg.GetId()); err != nil {
		if errors.Is(err, usecase.ErrMyListChartAttachmentNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, errors.WithStack(err))
		}
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}
	return connect.NewResponse(&proto_my_list.DeleteMyListChartAttachmentResponse{}), nil
}
