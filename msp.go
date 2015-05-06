package msp

// A UserDatabase is an abstraction over the name -> share map returned by the
// secret splitter that allows an application to only decrypt or request shares
// when needed, rather than re-build a partial map of known data.
type UserDatabase interface {
	Get(string) ([]byte, error)
}

type Condition interface { // Represents one condition in a predicate
	Ok(*UserDatabase) bool
}

type String string // Type of condition

func (s String) Ok(db *UserDatabase) bool {
	_, err := (*db).Get(string(s))
	return err == nil
}
