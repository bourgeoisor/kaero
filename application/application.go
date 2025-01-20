package application

import (
	"github.com/gdamore/tcell/v2"
	"kaero/client"
	"unicode"
)

const (
	appName = "Kaero"
)

type Application struct {
	version  string
	screen   tcell.Screen
	width    int
	height   int
	listener chan int

	server     *client.Server
	channelTab string
	logsOffset int

	inputActive bool
	inputCursor int
	inputText   []rune
}

func New(version string) (*Application, error) {
	// @todo: temporary
	listener := make(chan int)
	server := client.New(listener, "irc.libera.chat", 6697, "kaero-client")
	err := server.Connect()
	if err != nil {
		return nil, err
	}

	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	return &Application{
		version:  version,
		screen:   screen,
		listener: listener,
		server:   server,
	}, nil
}

func (a *Application) Run() error {
	err := a.screen.Init()
	if err != nil {
		return err
	}

	a.screen.EnableMouse()

	go a.listenToChannel()

	for {
		ev := a.screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			a.width, a.height = ev.Size()
			a.screen.Sync()
		case *tcell.EventMouse:
			a.handleMouseEvent(ev)
		case *tcell.EventKey:
			a.handleKeyEvent(ev)
		}

		a.draw()
	}
}

func (a *Application) Stop() {
	a.screen.Fini()
}

func (a *Application) listenToChannel() {
	for {
		<-a.listener
		a.draw()
	}
}

func (a *Application) handleMouseEvent(ev *tcell.EventMouse) {
	if ev.Buttons()&tcell.WheelUp > 0 {
		a.logsOffsetUp()
	}

	if ev.Buttons()&tcell.WheelDown > 0 {
		a.logsOffsetDown()
	}
}

func (a *Application) handleKeyEvent(ev *tcell.EventKey) {
	if ev.Modifiers() == tcell.ModAlt {
		indexes := map[rune]int{'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9}
		channels := a.server.ChannelNames()
		tab := indexes[ev.Rune()]
		if tab == 0 {
			a.channelTab = ""
		} else if tab <= len(channels) {
			a.channelTab = channels[tab-1]
		}
		a.logsOffset = 0
		return
	}

	switch ev.Key() {
	case tcell.KeyCtrlC:
		a.Stop()
		// @todo: end properly
	case tcell.KeyEnter:
		if len(a.inputText) == 0 {
			a.inputActive = !a.inputActive
		} else {
			a.server.HandleUserInput(string(a.inputText), a.channelTab)
			a.inputText = make([]rune, 0)
			a.inputCursor = 0
		}
	case tcell.KeyPgUp:
		a.logsOffsetUp()
	case tcell.KeyPgDn:
		a.logsOffsetDown()
	default:
	}

	if a.inputActive {
		switch ev.Key() {
		case tcell.KeyBackspace:
			if a.inputCursor > 0 {
				a.inputText = append(a.inputText[:a.inputCursor-1], a.inputText[a.inputCursor:]...)
				a.inputCursor--
			}
		case tcell.KeyLeft:
			a.inputCursor--
			if a.inputCursor < 0 {
				a.inputCursor = 0
			}
		case tcell.KeyRight:
			a.inputCursor++
			if a.inputCursor > len(a.inputText) {
				a.inputCursor = len(a.inputText)
			}
		default:
			if unicode.IsPrint(ev.Rune()) {
				a.inputText = append(a.inputText, ' ')
				copy(a.inputText[a.inputCursor+1:], a.inputText[a.inputCursor:])
				a.inputText[a.inputCursor] = ev.Rune()
				a.inputCursor++
			}
		}
	}
}

func (a *Application) currentChannel() *client.Channel {
	channel, err := a.server.GetChannel(a.channelTab)
	if err != nil {
		a.channelTab = ""
	}
	return channel
}

func (a *Application) logsOffsetUp() {
	a.logsOffset += 3
}

func (a *Application) logsOffsetDown() {
	a.logsOffset -= 3
	if a.logsOffset < 0 {
		a.logsOffset = 0
	}
}
