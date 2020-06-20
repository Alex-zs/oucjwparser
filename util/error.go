package util

type JwError struct {
	Msg string
}

func (err *JwError) Error() string {
	return err.Msg
}
