package cmd

type Error struct {
	Message string
}

func (r *Error) Error() string {
	return r.Message
}
