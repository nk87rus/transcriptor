package salutespeech

import (
	"context"

	"github.com/nk87rus/transcriptor/internal/service/sbertoken"
)

const scope = "SALUTE_SPEECH_PERS"

type Tokenizer interface {
	Get(context.Context) (string, error)
}

type SaluteSpeechClient struct {
	token Tokenizer
}

func Init(ctx context.Context, authKey string) (*SaluteSpeechClient, error) {
	return &SaluteSpeechClient{token: sbertoken.NewTokenManager(authKey, scope)}, nil
}
