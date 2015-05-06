Monotone Span Programs
======================

A *Monotone Span Program* (or *MSP*) is a cryptographic technique for splitting
a secret into several *shares* that are then distributed to *parties* or
*users*.  (Have you heard of [Shamir's Secret Sharing](http://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing)?  It's like that.)

Unlike Sharmir's Secret Sharing, MSPs allow *arbitrary monotone access
structures*.  An access structure is just a boolean predicate on a set of users
that tells us whether or not that set is allowed to recover the secret.  A
monotone access structure is the same thing, but with the invariant that adding
a user to a set will never turn the predicate's output from `true` to
`false`--negations or boolean `nots` are disallowed.

**Example:**  `(Alice or Bob) and Carl` is good, but `(Alice or Bob) and !Carl`
is not because excluding people is rude.


#### Types of Predicates

An MSP itself is a type of predicate and the reader is probably familiar with
raw boolean predicates like in the example above, but another important type is
a *formatted boolean predicate*.

Formatted boolean predicates are isomorphic to all MSPs and therefore all
monotone raw boolean predicates.  They're built by nesting threshold gates.

**Example:**  Let `(2, Alice, Bob, Carl)` denote that at least 2 members of the
set `{Alice, Bob, Carl}` must be present to recover the secret.  Then,
`(2, (1, Alice, Bob), Carl)` is the formatted version of
`(Alice or Bob) and Carl`.

It is possible to convert between different types of predicates (and its one of
the fundamental operations of splitting secrets with an MSP), but circuit
minimization is a non-trivial and computationally complex problem.  Therefore,
its best that the 
