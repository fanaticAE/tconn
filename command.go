package tconn

import "bufio"

type Runner interface {
	Run([]string, bufio.ReadWriter) error
}

type DefaultCommandNotFoundRunner struct {
}

func (DefaultCommandNotFoundRunner) Run(args []string, rw bufio.ReadWriter) error {
	_, err := rw.Writer.WriteString("This command could not be found\n")
	if err != nil {
		return err
	}
	return rw.Writer.Flush()

}
