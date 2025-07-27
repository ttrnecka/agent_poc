package core

import "io"

type Runner interface {
	Connect() (io.Closer, error)
	Run(string) ([]byte, *ExitCodeError)
}
