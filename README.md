# Golang Telegram Bot `FRAMEWORK`

## Project is in _deeeeep_ developing and far from complete

## Allows Telegram Bot developer to
1. Create smart menus
2. Supports multiple language menus
3. Allows to create static and custom dynamic menus
 
### Example
   ![example](https://s8.gifyu.com/images/ezgif.com-gif-maker91dac04b6bc3438d.gif)

## What is ahead ?
1. ...

## Api
### Usage
1. (Can be copied and pasted). Example of simple echo bot. Can do two commands: 
   1. /echo "text" - returns "text" to user.
   2. /menu - returns menu with one button "/echo 123"
```go
    package main

import (
   "os"
   "os/signal"
   "strings"

   "github.com/Red-Sock/go_tg/client"
   tg "github.com/Red-Sock/go_tg/interfaces"
   "github.com/Red-Sock/go_tg/model"
   "github.com/Red-Sock/go_tg/model/response/menu"
)

func main() {
   // Create bot instance
   bot := client.NewBot("1816053505:AAG9VfAEJtRFZlqp58WgESyAByjlBhmR0Ik")
   // Fill bot with commands/menus
   bot.AddCommandHandler(&EchoHandler{}, "/echo")

   m := menu.NewSimple("/menu", "MainMenu")
   m.AddButton("echo 123", "/echo 123")
   bot.AddMenu(m)

   bot.Start()

   done := make(chan os.Signal, 1)
   signal.Notify(done, os.Interrupt)
   <-done
   bot.Stop()
}

type EchoHandler struct{}

// Handle - Simple Echo handler
func (s *EchoHandler) Handle(in *model.MessageIn, out tg.Sender) {
   if len(in.Args) == 0 {
      out.Send(model.NewMessage("empty input!"))
      return
   }
   out.Send(model.NewMessage(strings.Join(in.Args, "-")))
}

// Dump - is unnecessary here
func (s *EchoHandler) Dump(_ int64) {}
```

2. Work with context.
If you are using localized menus or need to work with user's metadata - recommended to use interfaces.ExternalContext
```go
    type ContextEnricher struct {}
    func (c *ContextEnricher) GetContext(message *tgmodel.MessageIn) (context.Context, error){
    	// can make DB calls to extract users metadata and put it in context
        ctx := context.WithValue(context.Background(), "someKey", "some value")
    	return ctx, nil
    }

    bt.ExternalContext = &ContextEnricher{}
 ```


Then can be extracted in handler


```go
      func (s *Start) Handle(in *tgmodel.MessageIn, out tg.Sender) {
        someVal, ok := in.Ctx.Value("someKey").(string)
        if !ok {
           out.Send("Unexpected context error.")
		   return
        }
	    out.Send(tgmodel.NewMessage(someVal))
      }
```
4. Add simple non localized menu:
```go
    // "/menu" - command that opens menu 
	// "MainMenu" - name of menu (will be displayed as text message above buttons)
    m := menu.NewSimple("/menu", "MainMenu") // creates menu 
	
    m.AddButton("echo 123", "/echo 123") // adds simple button
    bot.AddMenu(m) // registers menu in bot framework
```
4. Add localized menu
```go
     // "/openLM" - command that opens localized menu 
     localizedMenu := tgmodel.NewLocalizedMenu("/openLM") // create localized menu

	 // "/openLM" - command
	 // "Локализированное меню" - displayed name
     mRU := tgmodel.NewSimple("/openLM", "Локализированное меню")
     mRU.AddButton("эхо 123", "/echo 123")
     localizedMenu.AddMenu("ru", mRU)
   
	 mEN := tgmodel.NewSimple("/openLM", "Localized menu")
	 mEN.AddButton("echo 123", "/echo 123")
	 localizedMenu.AddMenu("en", mEN)
```
## Examples of usage
1. [Gitlab webhook to Telegram Notificator](https://github.com/AlexSkilled/GitM8)
