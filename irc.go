package network

import (
	"github.com/thoj/go-ircevent"
)


func StreamIrc(url, channel, nickname, user string, callback func(e *irc.Event)) {
	irccon := irc.IRC(nickname, user)

	//irccon.VerboseCallbackHandler = true
	//irccon.Debug = true

	err := irccon.Connect(url)
	if err != nil { panic(err) }
	defer irccon.Quit()

	irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(channel) })
	irccon.AddCallback("PRIVMSG", callback)

	irccon.Loop()
}