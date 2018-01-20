# gh-ranking [![Build Status](https://travis-ci.org/fern4lvarez/gh-ranking.png)](https://travis-ci.org/fern4lvarez/gh-ranking) 
[Documentation online](http://godoc.org/github.com/fern4lvarez/gh-ranking)

Languages ranking on GitHub

**gh-ranking** is a simple tool to get the position of a programming language in the GitHub popularity ranking.

## Install

* Step 1: Get the package. Then you will be able to use ´gh-ranking´ as a executable program. 

```
go get github.com/fern4lvarez/gh-ranking
```

* Step 2 (Optional): Run tests

```
$ go test -v ./...
```

##Usage

You can use this package via CLI or API

### CLI
As easy as doing this:

```
$ gh-ranking go
go is on position #23 on GitHub.
For more information, visit https://github.com/languages/go
```

If the language has multiple words, use quotes:

```
$ gh-ranking "Visual Basic"
Visual Basic is on position #31 on GitHub.
For more information, visit https://github.com/languages/Visual%20Basic
```

What if that language doesn't exist yet?
```
$ gh-ranking Goffeeejure
Something went wrong! Try again with some other real language, Mr. @HipsterHacker.
```

### API

```
package main

import (
    "fmt"
    rank "github.com/fern4lvarez/gh-ranking"
)

func main() {
    pos, err := rank.Position("go")
    if err != nil {
        fmt.Println("Error")
    }
    fmt.Println(pos)
    // Output: 23
}
```

##Contribute!
Yes, contribute! Spam me with issues, fork it, improve it, do it!

##TODO (aka Nice To Have)
* Multiple languages selection
* Show full list

##License
gh-ranking is MIT licensed, see [here](https://github.com/fern4lvarez/gh-ranking/blob/master/LICENSE)
