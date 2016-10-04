package tconn

import "net"
import "log"
import "bufio"
import "errors"
import "regexp"
import "strings"

var commands map[string]Runner

var ConnectionType string
var Address string

var cmdRegexp = regexp.MustCompile(`"([^"]+)"|((?:[^ ]|)+)`)

var CommandNotFoundRunner Runner

var Log = true

func RegisterCommand(command string, runner Runner) {
	commands[command] = runner
}

func Listen() error {
	if ConnectionType == "" {
		return errors.New("Connection Type is not set")
	}
	if Address == "" {
		return errors.New("Address is not set")
	}
	listener, err := net.Listen(ConnectionType, Address)
	if err != nil {
		return err
	}
	go listen(listener)
	return nil
}

func listen(listener net.Listener) {
	defer listener.Close()
	for {
		c, err := listener.Accept()
		if err != nil {
			if Log {
				log.Println(err)
			}
			continue
		}
		go handle(c)
	}
}

func handle(c net.Conn) {
	var err error
	var str string

	rw := bufio.ReadWriter{
		Reader: bufio.NewReader(c),
		Writer: bufio.NewWriter(c),
	}

	for func() bool { str, err = rw.Reader.ReadString('\n'); return err == nil }() {
		line := strings.TrimSpace(str)
		tmp := cmdRegexp.FindAllStringSubmatch(line, -1)

		params := make([]string, 0, len(tmp))
		for _, v := range tmp {
			if len(v) < 2 {
				continue
			}
			if v[1] != "" {
				params = append(params, v[1])
			} else {
				params = append(params, v[2])
			}
		}

		if len(params) < 1 {
			continue
		}

		if command, ok := commands[params[0]]; ok {
			command.Run(params, rw)

		} else {
			CommandNotFoundRunner.Run(params, rw)
		}

		if err != nil && Log {
			log.Println(err)
		}

	}

	if Log {
		log.Println(err)
	}

}

func init() {
	commands = make(map[string]Runner)
	CommandNotFoundRunner = &DefaultCommandNotFoundRunner{}
}
