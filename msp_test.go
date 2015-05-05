package msp

import (
  "testing"
  "errors"
)

type Database map[string][]byte

func (d Database) Get(name string) ([]byte, error) {
  out, ok := d[name]

  if ok {
    return out, nil
  } else {
    return []byte(""), errors.New("Not found!")
  }
}

func TestFormatted(t *testing.T) {
  query1 := Formatted{
    Min: 2,
    Conds: []Condition{
      String("Alice"), String("Bob"), String("Carl"),
    },
  }

  query2 := Formatted{
    Min: 3,
    Conds: []Condition{
      String("Alice"), String("Bob"), String("Carl"),
    },
  }

  query3 := Formatted{
    Min: 2,
    Conds: []Condition{
      Formatted{
        Min: 1,
        Conds: []Condition{
          String("Alice"), String("Bob"),
        },
      },
      String("Carl"),
    },
  }

  query4 := Formatted{
    Min: 2,
    Conds: []Condition{
      Formatted{
        Min: 1,
        Conds: []Condition{
          String("Alice"), String("Carl"),
        },
      },
      String("Bob"),
    },
  }

  db := UserDatabase(Database(map[string][]byte {
    "Alice": []byte("blah"),
    "Carl": []byte("herp"),
    }))

  if query1.Ok(&db) != true {
    t.Fatalf("Query #1 was wrong.")
  }

  if query2.Ok(&db) != false {
    t.Fatalf("Query #2 was wrong.")
  }

  if query3.Ok(&db) != true {
    t.Fatalf("Query #3 was wrong.")
  }

  if query4.Ok(&db) != false {
    t.Fatalf("Query #4 was wrong.")
  }

  query1String := "(2, Alice, Bob, Carl)"
  query3String := "(2, (1, Alice, Bob), Carl)"

  if query1.String() != query1String {
    t.Fatalf("Query #1 String was wrong; %v", query1.String())
  }

  if query3.String() != query3String {
    t.Fatalf("Query #3 String was wrong; %v", query3.String())
  }

  decQuery1, err := StringToFormatted(query1String)
  if err != nil || decQuery1.String() != query1String {
    t.Fatalf("Query #1 decoded wrong: %v %v", decQuery1.String(), err)
  }

  decQuery3, err := StringToFormatted(query3String)
  if err != nil || decQuery3.String() != query3String {
    t.Fatalf("Query #3 decoded wrong: %v %v", decQuery3.String(), err)
  }
}
