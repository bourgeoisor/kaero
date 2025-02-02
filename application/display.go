package application

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"kaero/utils"
)

func (a *Application) draw() {
	a.screen.Clear()
	a.screen.SetStyle(style(tcell.ColorReset, tcell.ColorReset))

	a.drawLogs()
	a.drawNickList()

	a.drawTopBar()
	a.drawBottomBar()
	a.drawInput()

	//a.drawSandbox()

	a.screen.Show()
}

func (a *Application) drawTopBar() {
	for col := range a.width {
		a.screen.SetContent(col, 0, ' ', nil, style(tcell.ColorWhite, tcell.ColorBlack))
		col++
	}

	channel := a.currentChannel()
	text := fmt.Sprintf("%s v%s", appName, a.version)
	if channel != nil {
		text += fmt.Sprintf(" / %s [%d users]", a.channelTab, len(channel.Nicks))
		if channel.Topic != "" {
			text += fmt.Sprintf(" - %s", channel.Topic)
		}
	}
	a.drawString(0, 0, text, style(tcell.ColorWhite, tcell.ColorBlack))
}

func (a *Application) drawBottomBar() {
	for col := range a.width {
		a.screen.SetContent(col, a.height-2, ' ', nil, style(tcell.ColorWhite, tcell.ColorBlack))
		col++
	}

	text := "0:status"
	for i, channel := range a.server.ChannelNames() {
		text += fmt.Sprintf(" %d:%s", i+1, channel)
	}

	a.drawString(0, a.height-2, text, style(tcell.ColorWhite, tcell.ColorBlack))
}

func (a *Application) drawLogs() {
	channel := a.currentChannel()

	var logs []utils.Log
	if channel == nil {
		logs = a.server.GetLogger().GetNLogs(a.height-3, a.logsOffset)
	} else {
		logs = channel.Logs.GetNLogs(a.height-3, a.logsOffset)
	}

	if channel != nil {
		for i := 1; i < a.height-2; i++ {
			a.drawString(a.MaxNickLength+1, i, "│", style(tcell.ColorReset, tcell.ColorGray))
		}
	}

	row := a.height - 3
	for i := len(logs) - 1; i >= 0; i-- {
		height := a.drawLog(row, logs[i])
		row -= height
	}
}

func (a *Application) drawLog(row int, log utils.Log) int {
	baseStyle := style(tcell.ColorReset, tcell.ColorReset)
	height := 1
	delimIndex := a.MaxNickLength + 1

	maxWidth := a.width - 2*(delimIndex) - 3

	switch log.Kind {
	case utils.LogPrivMsg:
		height = a.drawStringWrap(delimIndex+2, row, maxWidth, log.Text, baseStyle)
		a.drawString(delimIndex-len(log.Source)-1, row-height+1, log.Source, baseStyle)
		for i := row - height + 1; i <= row; i++ {
			a.drawString(delimIndex, i, "│", baseStyle)
		}
	case utils.LogSystem:
		height = a.drawStringWrap(delimIndex+2, row, maxWidth, log.Text, baseStyle.Foreground(tcell.ColorBlue))
		for i := row - height + 1; i <= row; i++ {
			a.drawString(delimIndex, i, "│", baseStyle.Foreground(tcell.ColorBlue))
		}
	case utils.LogJoined:
		a.drawString(delimIndex, row, fmt.Sprintf("│ %s %s", log.Source, log.Text), baseStyle.Foreground(tcell.ColorGreen))
	case utils.LogLeft:
		a.drawString(delimIndex, row, fmt.Sprintf("│ %s %s", log.Source, log.Text), baseStyle.Foreground(tcell.ColorRed))
	case utils.LogError:
		height = a.drawStringWrap(len(log.Source)+2, row, a.width-len(log.Source)-2, log.Text, baseStyle.Foreground(tcell.ColorRed))
		a.drawString(0, row-height+1, fmt.Sprintf("%s:", log.Source), baseStyle.Foreground(tcell.ColorRed))
	case utils.LogStatus:
		height = a.drawStringWrap(len(log.Source)+2, row, a.width-len(log.Source)-2, log.Text, baseStyle)
		a.drawString(0, row-height+1, fmt.Sprintf("%s:", log.Source), baseStyle)
	}

	return height
}

func (a *Application) drawNickList() {
	if a.channelTab == "" {
		return
	}

	delimIndex := a.width - a.MaxNickLength - 3

	nicks := a.currentChannel().NicksListByMode()

	for i := 1; i < a.height-2; i++ {
		a.drawString(delimIndex, i, "│", style(tcell.ColorReset, tcell.ColorGray))
	}

	row := 1
	for _, nick := range nicks {
		if row >= a.height-2 {
			break
		}
		a.drawString(delimIndex+2, row, nick, style(tcell.ColorReset, tcell.ColorReset))
		row++
	}
}

func (a *Application) drawInput() {
	tabName := "[status]"
	if a.channelTab != "" {
		tabName = fmt.Sprintf("[%s]", a.channelTab)
	}
	a.drawString(0, a.height-1, tabName, style(tcell.ColorReset, tcell.ColorGray))

	a.drawString(len(tabName)+1, a.height-1, string(a.inputText), style(tcell.ColorReset, tcell.ColorReset))

	if a.inputActive {
		a.screen.ShowCursor(len(tabName)+1+a.inputCursor, a.height-1)
	} else {
		a.screen.HideCursor()
	}
}

func (a *Application) drawString(x int, y int, text string, style tcell.Style) {
	row := y
	col := x
	for _, r := range []rune(text) {
		a.screen.SetContent(col, row, r, nil, style)
		_, _, _, width := a.screen.GetContent(col, row)
		col += width
	}
}

func (a *Application) drawStringWrap(x int, y int, maxWidth int, text string, style tcell.Style) int {
	chunks := make([]string, 0)
	col := x
	start := 0
	for i, r := range []rune(text) {
		a.screen.SetContent(col, 0, r, nil, style)
		_, _, _, width := a.screen.GetContent(col, 0)
		col += width
		if col >= x+maxWidth {
			col = x
			chunks = append(chunks, text[start:i])
			start = i
		}
	}
	chunks = append(chunks, text[start:])

	for i, chunk := range chunks {
		a.drawString(x, y+i-len(chunks)+1, chunk, style)
	}

	return len(chunks)
}

func (a *Application) drawSandbox() {
	colors := []tcell.Color{
		tcell.ColorBlack,
		tcell.ColorMaroon,
		tcell.ColorGreen,
		tcell.ColorOlive,
		tcell.ColorNavy,
		tcell.ColorPurple,
		tcell.ColorTeal,
		tcell.ColorSilver,
		tcell.ColorGray,
		tcell.ColorRed,
		tcell.ColorLime,
		tcell.ColorYellow,
		tcell.ColorBlue,
		tcell.ColorFuchsia,
		tcell.ColorAqua,
		tcell.ColorWhite,
	}

	for i := range 16 {
		a.screen.SetContent(i, 1, 'X', nil, style(tcell.ColorReset, colors[i]))
	}
}

func style(bg tcell.Color, fg tcell.Color) tcell.Style {
	return tcell.StyleDefault.Background(bg).Foreground(fg)
}
