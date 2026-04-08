package app

import (
	"context"
	"os"

	"github.com/nk87rus/transcriptor/internal/config"
	"github.com/nk87rus/transcriptor/internal/handler"

	"github.com/nk87rus/transcriptor/internal/logger"
	"github.com/nk87rus/transcriptor/internal/repository/psql"
	"github.com/nk87rus/transcriptor/internal/service/gigachat"
	"github.com/nk87rus/transcriptor/internal/service/salutespeech"
	"github.com/nk87rus/transcriptor/internal/service/telegram"
	"golang.org/x/sync/errgroup"
)

type TaskProvider interface {
	Run(context.Context) error
}

type App struct {
	tp TaskProvider
}

func Init(ctx context.Context) (*App, error) {
	logger.Init()

	cfg, err := config.InitConfig(os.Args)
	if err != nil {
		return nil, err
	}

	repo, errRepo := psql.NewPSQLRepo(ctx, cfg.DBDSN)
	if errRepo != nil {
		return nil, errRepo
	}

	salutSpeech, errSR := salutespeech.Init(ctx, cfg.SpeechRecKey)
	if errSR != nil {
		return nil, errSR
	}

	gChat, errGCH := gigachat.Init(ctx, cfg.ChatKey)
	if errGCH != nil {
		return nil, errGCH
	}

	teleBot, errTR := telegram.InitBot(ctx, cfg.TaskProvToken, handler.Init(repo, salutSpeech, gChat))
	if errTR != nil {
		return nil, errTR
	}

	return &App{tp: teleBot}, nil
}

func (a *App) Run(ctx context.Context) error {
	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		a.tp.Run(egCtx)
		return nil
	})

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
