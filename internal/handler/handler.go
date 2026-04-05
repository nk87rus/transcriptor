package handler

import (
	"context"

	"github.com/nk87rus/stenographer/internal/model"
)

type RepoMgr interface {
	RegisterUser(ctx context.Context, userID int64, userName string) error
	GetMeetingsList(ctx context.Context) ([]model.MeetingsListItem, error)
	GetMeeting(ctx context.Context, mID int64) (*model.Meeting, error)
}

type Handler struct {
	repo RepoMgr
}

func Init(storage RepoMgr) *Handler {
	return &Handler{repo: storage}
}

func (h *Handler) RegisterUser(ctx context.Context, userID int64, userName string) error {
	return h.repo.RegisterUser(ctx, userID, userName)
}

func (h *Handler) GetMeetingsList(ctx context.Context) ([]model.MeetingsListItem, error) {
	return h.repo.GetMeetingsList(ctx)
}

func (h *Handler) GetMeeting(ctx context.Context, mID int64) (*model.Meeting, error) {
	return h.repo.GetMeeting(ctx, mID)
}
