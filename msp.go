package msp

import (
  "fmt"
  "strconv"
  "errors"
  "strings"
)

// A UserDatabase is an abstraction over the name -> share map returned by the
// secret splitter that allows an application to only decrypt or request shares
// when needed, rather than re-build a partial map of known data.
type UserDatabase interface {
  Get(string) ([]byte, error)
}


type Condition interface { // Represents one condition in a threshold gate.
  Ok(*UserDatabase) bool
}

type String string // Type of condition

func (s String) Ok(db *UserDatabase) bool {
  _, err := (*db).Get(string(s))
  return err == nil
}


type Formatted struct { // Represents threshold gate (also type of condition)
  Min int
  Conds []Condition
}

func StringToFormatted(f string) (Formatted, error) {
  var out Formatted
  var err error

  if f[0] != '(' || f[len(f) - 1] != ')' {
    return out, errors.New("Invalid string.  #1")
  }

  // Extract first value:  min.
  nextComma := strings.Index(f, ",")
  if nextComma < 0 {
    return out, errors.New("Invalid string.  #2")
  }

  out.Min, err = strconv.Atoi(f[1:nextComma])
  if err != nil {
    return out, err
  }

  f = strings.TrimSpace(f[nextComma + 1:])

  // Extract each condition (sometimes recursively)
  var nextParen, nextUnParen int
  for len(f) > 0 {
    nextComma = strings.Index(f, ",")
    nextParen = strings.Index(f, "(")

    if nextComma == -1 {
      out.Conds = append(out.Conds, String(f[0:len(f) - 1]))
      f = ""
    } else if nextParen == -1 || (nextParen != -1 && nextComma < nextParen) {
      out.Conds = append(out.Conds, String(f[0:nextComma]))
      f = strings.TrimSpace(f[nextComma + 1:])
    } else {
      nextUnParen = strings.Index(f, ")")
      if nextUnParen == -1 || nextUnParen < nextParen {
        return out, errors.New("Invalid string.  #3")
      }

      nxt, err := StringToFormatted(f[nextParen:nextUnParen + 1])
      if err != nil {
        return out, err
      }

      out.Conds = append(out.Conds, nxt)
      f = strings.TrimSpace(f[nextUnParen + 2:])
    }
  }

  return out, nil
}

func (f Formatted) String() string {
  out := fmt.Sprintf("(%v", f.Min)

  for _, cond := range f.Conds {
    switch cond.(type) {
      case String: out += fmt.Sprintf(", %v", cond)
      case Formatted: out += fmt.Sprintf(", %v", (cond.(Formatted)).String())
    }
  }

  return out + ")"
}

func (f Formatted) Ok(db *UserDatabase) bool {
  // Goes through the smallest number of conditions possible to check if the
  // threshold gate returns true.  Sometimes requires recursing down to check
  // nested threshold gates.
  rest := f.Min

  for _, cond := range f.Conds {
    if cond.Ok(db) {
      rest--
    }

    if rest == 0 {
      return true
    }
  }

  return false
}


/*
// FormatRaw takes a boolean predicate in the form of a string and parses it
// into a formatted boolean query.
func FormatRaw()


func FormattedToString(f Formatted)
*/
