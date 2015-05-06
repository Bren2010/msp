package msp

import (
  "fmt"
	"testing"
)

func TestRaw(t *testing.T) {
	alice := Condition(String("Alice"))
	bob := Condition(String("Bob"))
	carl := Condition(String("Carl"))

	query1 := Raw{
		NodeType: NodeAnd,
		Left:     &alice,
		Right:    &bob,
	}

	aliceOrBob := Condition(Raw{
		NodeType: NodeOr,
		Left:     &alice,
		Right:    &bob,
	})

	query2 := Raw{
		NodeType: NodeAnd,
		Left:     &aliceOrBob,
		Right:    &carl,
	}

	db := UserDatabase(Database(map[string][]byte{
		"Alice": []byte("blah"),
		"Carl":  []byte("herp"),
	}))

	if query1.Ok(&db) != false {
		t.Fatalf("Query #1 was wrong.")
	}

	if query2.Ok(&db) != true {
		t.Fatalf("Query #2 was wrong.")
	}

	query1String := "Alice & Bob"
	query2String := "(Alice | Bob) & Carl"

	if query1.String() != query1String {
		t.Fatalf("Query #1 String was wrong; %v", query1.String())
	}

	if query2.String() != query2String {
		t.Fatalf("Query #2 String was wrong; %v", query2.String())
	}

  fmt.Println(StringToRaw("Alice | Bob"))
}
