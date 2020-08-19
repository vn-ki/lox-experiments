# go-lox

Implementation of the tree walk interpreter from crafting interpreters in golang.

### Why?

- I like golang. Can't say the same about Java.
- So that I don't copypasta code from the book.


### Changes from the book

I have tried to maintain this close to the book, while taking some liberties in some places.

- I chose a branched namespace instead of flat one like in the book. Will comment how it is afterwards.
- I did not like the global `errno` way of doing error handling.
- Renamed `Scanner` to `Lexer` to not cause confusion with golang scanner.
- The author of the book did not want to create an interface for error handling because it would lead to more complexity. We use a closure for that in this implementation.
-

### Notes

- Go's lack of generics makes the `Visitor` pattern look awkward. I'm not sure what the best course of action here. (Go 2 will have generics. Yay!)
- There is a mix and match of design pattern
  - `parser` propagates errors across functions using the golang way of returning error
  - `interpreter` uses `panic` and `recover` in almost a throw-catch way. This was done so that changing the entire interface seemed like an exercise in vain.
  - The branch `no-visitor` removes the visitor pattern. I think this is a better way to go about it, but it was abandoned to stay close to the book

### Thoughts on visitor pattern vs no visitor pattern

no-visitor branch removes the visitor pattern. It uses a dummy interface

```golang
type Expr interface {
  expr()
}

type Unary struct {}
func (t Unary) expr()
```

Then it can be used like so insted of visitor pattern

```golang
type ExprVisitor interface {
  VisitExpr(Expr) interface{}
}

type AstVisitor struct {}

func VisitExpr(e Expr) interface{} {
  switch v := e.(type) {
  case Binary:
    ...
  case Unary:
    ...
  }
}
```

Notable differences from visitor pattern:

- lesser functions in the visitor interface. So any future requirements which may change the function signature may not affect all visitor implementations. (that is custom requirements could be handled inside the visitor implementation itself)
- No compile time guarentee that all types are handled. The visitor pattern guarentees this.
- (In my opinion) lesser abstraction and easier to follow and understand. Visitor pattern has dual indirection, while this pattern only uses a switch.
- Visitor pattern is considerably harder to track while debugging due to the double indirection.

