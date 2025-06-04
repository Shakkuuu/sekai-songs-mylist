package entity

import (
	"time"

	"github.com/Shakkuuu/sekai-songs-mylist/internal/gen/enums"
)

type MyList struct {
	ID        int32
	UserID    string
	Name      string
	Position  int32
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MyListChart struct {
	ID        int32
	MyListID  int32
	ChartID   int32
	Chart     *Chart
	ClearType enums.ClearType
	Memo      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MyListChartAttachment struct {
	ID             int32
	MyListChartID  int32
	AttachmentType enums.AttachmentType
	FileURL        string
	Caption        string
	CreatedAt      time.Time
}
