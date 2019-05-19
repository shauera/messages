package persistence

//Error - a persistence error type
type Error string

//Error - implemenation of Error interface
func (e Error) Error() string {
	return string(e)
}

//ErrorNotFound - record could not be found in the repository
const ErrorNotFound = Error("Not found")
