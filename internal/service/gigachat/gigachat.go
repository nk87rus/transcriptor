package gigachat

import (
	"context"

	"github.com/nk87rus/transcriptor/internal/service/sbertoken"
)

const scope = "GIGACHAT_API_PERS"

type Tokenizer interface {
	Get(context.Context) (string, error)
}

type GigaChatClient struct {
	token Tokenizer
}

func Init(ctx context.Context, authKey string) (*GigaChatClient, error) {
	return &GigaChatClient{token: sbertoken.NewTokenManager(authKey, scope)}, nil
}

