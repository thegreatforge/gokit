package splitify

type Error string

const ErrNoRules Error = "no rules added"
const ErrNoConditions Error = "no conditions added"

func (e Error) Error() string {
	return string(e)
}

func (e Error) String() string {
	return e.Error()
}
