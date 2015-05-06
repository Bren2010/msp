package msp

import (
  "strings"
  "container/list"
  "errors"
)

type NodeType int // Types of node in the binary expression tree.

const (
    NodeAnd NodeType = iota
    NodeOr
)

func (t NodeType) Type() NodeType {
  return t
}

type Raw struct { // Represents one node in the tree.
  NodeType

  Left *Condition
  Right *Condition
}

func StringToRaw(r string) (out Raw, err error) {
  // Automaton.  Modification of Dijkstra's Two-Stack Algorithm for parsing
  // infix notation.  Reads one long unbroken expression (several operators and
  // operands with no parentheses) at a time and parses it into a binary
  // expression tree (giving AND operators precedence).
  //
  // Steps to the next (un)parenthesis.
  //     (     -> Push new queue onto staging stack
  //     value -> Push onto back of queue at top of staging stack.
  //     )     -> Pop queue off top of staging stack, build BET, and push tree
  //              onto the back of the top queue.
  //
  // Staging stack is empty on initialization and should have exactly 1 node
  // (the root node) at the end of the string.}
  min := func(a, b, c int) int { // Return smallest non-negative argument.
    if a > b { a, b = b, a } // Sort {a, b, c}
    if b > c { b, c = c, b }
    if a > b { a, b = b, a }

    if a != -1 {
      return a
    } else if b != -1 {
      return b
    } else {
      return c
    }
  }

  getNext := func(r string) (string, string) { // r -> (next, rest)
    r = strings.TrimSpace(r)

    if r[0] == '(' || r[0] == ')' || r[0] == '&' || r[0] == '|' {
      return r[0:1], r[1:]
    }

    nextOper := min(
      strings.Index(r, "&"),
      strings.Index(r, "|"),
      strings.Index(r, ")"),
    )

    if nextOper == -1 {
      return r, ""
    }
    return strings.TrimSpace(r[0:nextOper]), r[nextOper:]
  }

  staging := list.New() // Stack of (Condition list, operator list)

  var nxt string
  for len(r) > 0 {
    nxt, r = getNext(r)

    switch nxt {
      case "(":
        staging.PushFront([2]*list.List{list.New(), list.New()})
      case ")":
        top := staging.Remove(staging.Front()).([2]*list.List)
        if top[0].Len() != (top[1].Len() + 1) {
          return out, errors.New("Stacks are invalid size.")
        }

        for typ := NodeAnd; typ <= NodeOr; typ++ {
          leftOperand := top[0].Front().Next()

          for oper := top[1].Front(); oper != nil; oper = oper.Next() {
            if oper.Value.(NodeType) == typ {
              left := leftOperand.Value.(Condition)
              right := leftOperand.Next().Value.(Condition)

              leftOperand.Value = Raw{
                NodeType: typ,
                Left: &left,
                Right: &right,
              }

              top[0].Remove(leftOperand.Next())
              top[1].Remove(oper)
            }

            leftOperand = leftOperand.Next()
          }
        }

        if top[0].Len() != 1 || top[1].Len() != 0 {
          return out, errors.New("Invalid expression--couldn't evaluate.")
        }

        if staging.Len() == 0 {
          if len(r) == 0 {
            return top[0].Front().Value.(Raw), nil
          }
          return out, errors.New("Invalid string--terminated early.")
        }
        staging.Front().Value.([2]*list.List)[0].PushBack(top[0].Front().Value)

      case "&":
        staging.Front().Value.([2]*list.List)[1].PushBack(NodeAnd)
      case "|":
        staging.Front().Value.([2]*list.List)[1].PushBack(NodeOr)
      default:
        staging.Front().Value.([2]*list.List)[0].PushBack(nxt)
    }
  }

  return out, errors.New("Invalid string--never terminated.")
}

func (r Raw) String() string {
  out := ""

  switch (*r.Left).(type) {
  case String: out += string((*r.Left).(String))
  default:     out += "(" + (*r.Left).(Raw).String() + ")"
  }

  if r.Type() == NodeAnd {
    out += " & "
  } else {
    out += " | "
  }

  switch (*r.Right).(type) {
  case String: out += string((*r.Right).(String))
  default:     out += "(" + (*r.Right).(Raw).String() + ")"
  }

  return out
}

func (r Raw) Ok(db *UserDatabase) bool {
  if r.Type() == NodeAnd {
    return (*r.Left).Ok(db) && (*r.Right).Ok(db)
  } else {
    return (*r.Left).Ok(db) || (*r.Right).Ok(db)
  }
}
