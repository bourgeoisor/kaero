package client

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"kaero/utils"
	"log"
	"sort"
	"strings"
)

type Server struct {
	host     string
	port     int
	nick     string
	username string
	realName string

	name                  string
	version               string
	availableServerModes  string
	availableChannelModes string
	iSupport              *ISupport

	conn           *tls.Conn
	logs           *utils.Logger
	listener       chan int
	channelsJoined map[string]*Channel

	bufferMotd  []string
	bufferHelp  []string
	bufferList  []string
	bufferLinks []string
	bufferInfo  []string
	BufferWho   []string
	BufferStats []string
}

func New(listener chan int, host string, port int, nick string) *Server {
	return &Server{
		host:     host,
		port:     port,
		nick:     nick,
		username: nick,
		realName: nick,

		iSupport: newISupport(),

		logs:           utils.NewLogger(),
		listener:       listener,
		channelsJoined: make(map[string]*Channel),

		bufferMotd:  make([]string, 0),
		bufferHelp:  make([]string, 0),
		bufferList:  make([]string, 0),
		bufferLinks: make([]string, 0),
		bufferInfo:  make([]string, 0),
		BufferWho:   make([]string, 0),
		BufferStats: make([]string, 0),
	}
}

func (s *Server) GetLogger() *utils.Logger {
	return s.logs
}

func (s *Server) Connect() (err error) {
	address := fmt.Sprintf("%s:%d", s.host, s.port)
	s.logs.Append("System", utils.LogStatus, fmt.Sprintf("Dialing %s...", address))

	s.conn, err = tls.Dial("tcp", address, nil)
	if err != nil {
		return err
	}

	s.logs.Append("System", utils.LogStatus, fmt.Sprintf("Connected to %s", address))

	go s.listenToMessages()

	s.SendMessage(&utils.Message{Command: "NICK", Parameters: []string{s.nick}})
	s.SendMessage(&utils.Message{Command: "USER", Parameters: []string{s.username, "0", "*", s.realName}})

	return nil
}

func (s *Server) ChannelNames() []string {
	names := make([]string, 0)
	for _, channel := range s.channelsJoined {
		names = append(names, channel.Name)
	}
	sort.Strings(names)
	return names
}

func (s *Server) GetChannel(name string) (*Channel, error) {
	if channel, ok := s.channelsJoined[name]; ok {
		return channel, nil
	}
	return nil, fmt.Errorf("channel %s not found", name)
}

func (s *Server) SendMessage(message *utils.Message) {
	data := utils.MarshalMessage(message)

	_, err := s.conn.Write([]byte(data + "\r\n"))
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (s *Server) HandleUserInput(input string, channel string) {
	message := &utils.Message{}
	if input[0] == '/' {
		message = s.handleCommand(input, channel)
		if message == nil {
			return
		}
	} else {
		if channel != "" {
			message.Command = "PRIVMSG"
			message.Parameters = []string{channel, input}
			s.channelsJoined[channel].Logs.Append(s.nick, utils.LogPrivMsg, input)
		} else {
			s.SendMessage(utils.UnmarshalMessage(string(input)))
			return
		}
	}

	s.SendMessage(message)
}

func (s *Server) listenToMessages() {
	defer s.conn.Close()
	reader := bufio.NewReader(s.conn)
	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf(err.Error())
		}
		data = strings.TrimRight(data, "\r\n")

		s.handleServerMessage(utils.UnmarshalMessage(data))
		s.listener <- 1
	}
}

func (s *Server) log(text string) {
	s.logs.Append(s.host, utils.LogStatus, text)
}
