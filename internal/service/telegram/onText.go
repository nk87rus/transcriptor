package telegram

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nk87rus/stenographer/internal/model"
	"github.com/rs/zerolog/log"
	tele "gopkg.in/telebot.v3"
)

var reCmd = regexp.MustCompile(`^\/(start|get|list|find|chat)`)

type TextHandler interface {
	RegisterUser(ctx context.Context, userID int64, userName string) error
	GetMeetingsList(ctx context.Context) ([]model.MeetingsListItem, error)
	GetMeeting(ctx context.Context, mID int64) (*model.Meeting, error)
}

func (tb *TeleBot) OnText(ctx tele.Context) error {
	msg := Message{
		MsgType: MsgText,
		MsgCtx:  ctx,
	}
	tb.inChan <- msg

	return nil
}

func (tb *TeleBot) ProcessText(ctx context.Context, msg tele.Context) error {
	if len(msg.Entities()) > 0 && msg.Entities()[0].Type == tele.EntityCommand {
		log.Debug().Int64("ID", msg.Sender().ID).Str("UserName", msg.Sender().Username).Msgf("обработка команды %q", msg.Text())
		defer log.Debug().Int64("ID", msg.Sender().ID).Str("UserName", msg.Sender().Username).Msgf("обработка команды %q завершена", msg.Text())

		switch reCmd.FindString(msg.Text()) {
		case "/start":
			return tb.CmdStart(msg)
		case "/list":
			return tb.CmdList(msg)
		case "/get":
			return tb.CmdGet(msg)
		case "/find":
		case "/chat":
		default:
			// tb.SendResponse(ctx)
		}
	}
	return nil
}

func (tb *TeleBot) CmdStart(ctx tele.Context) error {
	// TODO: добавить обработку ошибки при попытке повторной регистрации пользователя
	if errHdlr := tb.hdlr.RegisterUser(tb.ctx, ctx.Sender().ID, ctx.Sender().Username); errHdlr != nil {
		return errHdlr
	}

	tb.outChan <- Response{MsgCtx: ctx, Data: fmt.Sprintf("Пользователь %q успешно зарегистрирован", ctx.Sender().Username)}
	return nil
}

func (tb *TeleBot) CmdList(ctx tele.Context) error {
	data, errHdlr := tb.hdlr.GetMeetingsList(tb.ctx)
	if errHdlr != nil {
		return errHdlr
	}

	var resp = Response{MsgCtx: ctx}
	if len(data) == 0 {
		resp.Data = "не найдено ни одной сохранённой встречи"
	} else {
		var sb strings.Builder
		for i, m := range data {
			sb.WriteString("📝 ")
			sb.WriteString(m.String())
			if i != len(data) {
				sb.WriteString("\n")
			}
		}
		resp.Data = sb.String()
	}

	tb.outChan <- resp
	return nil
}
func (tb *TeleBot) CmdGet(ctx tele.Context) error {
	meetingID, errID := strconv.ParseInt(ctx.Message().Payload, 10, 64)
	if errID != nil {
		return fmt.Errorf("не корректный формат идентификатора встречи")
	}

	data, errHdlr := tb.hdlr.GetMeeting(tb.ctx, meetingID)
	if errHdlr != nil {
		return errHdlr
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📝 Встреча %d от %s\n", data.Id, time.Unix(data.TimeStamp, 0).String()))
	sb.WriteString(fmt.Sprintf("Автор: %s\n", data.Author))
	sb.WriteString(strings.Repeat("-", 15) + "\n")
	sb.WriteString(data.Data)

	tb.outChan <- Response{MsgCtx: ctx, Data: sb.String()}
	return nil
}
func HdlrCmdFind(c tele.Context) error {
	return c.Send("поиск встречи по ключевым словам.")
}
func HdlrCmdChat(c tele.Context) error {
	return c.Send("запрос к GigaChat")
}
