package model

import (
	"fmt"
	"time"
)

type Meeting struct {
	Id        int64  `db:"id"`
	TimeStamp int64  `db:"ts"`
	Author    string `db:"user_name"`
	Data      string `db:"data"`
}

type MeetingsListItem struct {
	Id        int64  `db:"id"`
	TimeStamp int64  `db:"ts"`
	Author    string `db:"user_name"`
}

func (mli MeetingsListItem) String() string {
	return fmt.Sprintf("Встреча: %d Время: %s Автор: %s", mli.Id, time.Unix(mli.TimeStamp, 0).String(), mli.Author)
}
