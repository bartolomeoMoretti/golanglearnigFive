package telegram

import (
	"errors"
	"fmt"
	"golanglearningFive/lib/e"
	"golanglearningFive/storage"
	"log"
	"strings"
)

const (
	RndCmd   = "/get"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	// tm := time.Unix(datesaved,0)
	//strconv.Itoa(timeT)
	//fmt.Printf(timeT)
	//fmt.Sprint(timeT)

	//tNow := time.Now()
	//tUnix := tNow.Unix()
	//timeT := time.Unix(tUnix, 0)
	//timeT := time.Unix(int64(datesaved), 0)

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
		//DateCreated: fmt.Sprint(datesaved),
		//DateSaved: fmt.Sprint(timeT),
	}

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, fmt.Sprint(msgSaved, msgSaved2, pageURL)); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send random", err) }()

	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil { //to-do: haven't username - who will get message in Public-group?
		return err
	} //can be given more than once

	/*
	   if err := p.tg.SendMessage(ctx, chatID, page.URL); err != nil { //to-do: haven't username - who will get message in Public-group?
	   return err
	   }
	*/

	//return p.storage.Remove(ctx, page)
	return
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	if text == StartCmd {
		return false
	} else if text == HelpCmd {
		return false
	} else if text == RndCmd {
		return false
	}
	return true
}

/*func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}*/
