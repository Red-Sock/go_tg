package model

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

type LocalizedMenu struct {
	menus   map[string]Menu
	command string
}

func (l *LocalizedMenu) GetCallCommand() string {
	return l.command
}

func (l LocalizedMenu) GetTransitionCommand() string {
	return MenuCall + " " + l.command
}

const LocaleContextKey = "locale_context_key_go_tg"

func (l *LocalizedMenu) GetName(ctx context.Context) string {
	locale, ok := ctx.Value(LocaleContextKey).(string)
	if !ok {
		logrus.Error(fmt.Sprintf("Couldn't find  locale to extract from context in LocalizedMenu"))
		return ""
	}
	menu := l.menus[locale]
	if menu == nil {
		logrus.Error(fmt.Sprintf("Couldn't find name of menu %s with %s locale", l.command, locale))
		return ""
	}
	return menu.GetName(ctx)
}

func (l *LocalizedMenu) GetPage(ctx context.Context, page int) InlineKeyboard {
	locale, ok := ctx.Value(LocaleContextKey).(string)
	if !ok {
		return InlineKeyboard{}
	}
	menu := l.menus[locale]
	if menu == nil {
		return InlineKeyboard{}
	}
	return menu.GetPage(ctx, page)
}

// LocalizedMenu complex menu structure
// Allows to make many versions of same menu with different labels
type LocalizedMenuPatterns struct {
	command          string
	localeToPatterns map[string]SimpleMenuPattern
}

func NewLocalizedMenuPattern(command string) *LocalizedMenuPatterns {
	return &LocalizedMenuPatterns{
		command:          command,
		localeToPatterns: make(map[string]SimpleMenuPattern),
	}
}

func AssembleMenu(command string, localeToMenus map[string]MenuPattern) {

}

func (l *LocalizedMenuPatterns) AddMenu(locale, name string) {
	l.localeToPatterns[locale] = NewMenuPattern(name)
}

func (l *LocalizedMenuPatterns) AddMenus(localeToName map[string]string) {
	for locale, name := range localeToName {
		l.AddMenu(locale, name)
	}
}

func (l *LocalizedMenuPatterns) AddLocalizedMenuButton(localeToName map[string]string, command string) {
	for locale, name := range localeToName {
		if pattern, ok := l.localeToPatterns[locale]; ok {
			pattern.AddMenuButton(name, command)
			l.localeToPatterns[locale] = pattern
		} else {
			logrus.Error(errors.New(fmt.Sprintf("Error while creating multilanguage menu. No menu was created for locale %s.", locale)))
		}
	}
}

func (l *LocalizedMenuPatterns) AddMenuButton(locale, name string, command string) {
	if pattern, ok := l.localeToPatterns[locale]; ok {
		pattern.AddMenuButton(name, command)
		l.localeToPatterns[locale] = pattern
	} else {
		logrus.Error(errors.New(fmt.Sprintf("Error while creating multilanguage menu. No menu was created for locale %s.", locale)))
	}
}

func (l *LocalizedMenuPatterns) AddEntryPoint(command string) {
	if len(command) <= 0 {
		return
	}
	l.command = command
}

func (l *LocalizedMenuPatterns) GetCallCommand() string {
	return l.command
}

func (l *LocalizedMenuPatterns) GetTransitionCommand() string {
	return MenuCall + " " + l.command
}

func (l *LocalizedMenuPatterns) compile() Menu {
	lm := LocalizedMenu{
		menus:   map[string]Menu{},
		command: l.command,
	}
	for locale, menu := range l.localeToPatterns {
		lm.menus[locale] = menu.compile()
	}
	return &lm
}
