package model

import (
	"fmt"

	"github.com/gabriel-vasile/mimetype"
	"go.senan.xyz/taglib"
)

type AudioFile struct {
	LocalFilePath string
	Encoding      string
	SampleRate    int
	MIME          string
}

func MakeAudioFile(filePath string) (*AudioFile, error) {
	props, errProps := taglib.ReadProperties(filePath)
	if errProps != nil {
		return nil, errProps
	}

	mtype, errMIME := mimetype.DetectFile(filePath)
	if errMIME != nil {
		return nil, errMIME
	}

	enc, errEnc := getEncoding(mtype)
	if errEnc != nil {
		return nil, errEnc
	}

	return &AudioFile{LocalFilePath: filePath, MIME: mtype.String(), SampleRate: int(props.SampleRate), Encoding: enc}, nil
}

func getEncoding(m *mimetype.MIME) (string, error) {
	// PCM_S16LE, OPUS, MP3, FLAC, ALAW, MULAW
	const (
		ogg = "audio/ogg"
		mp3 = "audio/mpeg"
	)

	switch {
	case m.Is(ogg):
		return "OPUS", nil
	case m.Is(mp3):
		return "MP3", nil
	default:
		return "", fmt.Errorf("файлы %q в настоящий момент не поддерживаются.\nДопустимы %q и %q.", m.String(), ogg, mp3)
	}
}
