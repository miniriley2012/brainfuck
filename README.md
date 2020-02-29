# Brainfuck
Yet another Brainfuck compiler. This one will compile to [Go]("https://golang.org") or C.
```sh
go install github.com/miniriley2012/brainfuck/cmd/bf
```

## Flags
```
Usage of bf:
  -c string
        Language to compile to. Can be "go" or "c".
  -i    Interpret (not implemented yet)
```

## Caveats
This hasn't been tested extensively. If you run into any problems then open an issue. Compile-time checks are
also very basic as of now. If the inputted Brainfuck is wrong, the emitted code will most likely be too. Expect panics (Go)
or segfaults (C) if your Brainfuck has issues. 

## TODO
- [ ] Input from stdin
- [ ] Implement interpreter
- [ ] More compile-time checks
- [ ] Simplify Go compiler
- [ ] REPL?