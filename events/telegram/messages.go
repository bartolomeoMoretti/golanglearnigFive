package telegram

const msgHelp = `I can save and keep you notes. Also I can offer you them to read.

In order to save the page, just send me a text to it.

In order to get a random note from your list, send me command /get.
Caution! After that, this page won't be removed from your list, keep calm!`

const msgHello = "Hi there! š¾\n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command š¤"
	msgNoSavedPages   = "You have no saved pages š"
	msgSaved          = ", your this note has been saved:\n"
	msgSaved2         = "- - - - - - -\n"
	msgSaved22        = "\n- - - - - - -"
	msgSaved3         = "\nš£ The one was send from Telegram at "
	msgSaved4         = "š Dear "
	msgAlreadyExists  = "You have already have this page in your list š¤"
)
