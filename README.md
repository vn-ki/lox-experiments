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
