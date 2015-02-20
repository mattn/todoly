# todoly

cli interface for todoly

![](http://go-gyazo.appspot.com/bd8080b81e9d663d.png)

## Usage

```
$ todoly
todoly - cli interface for todo.ly

Commands:

    add         add todo
    check       check the todo
    del         del the todo
    list        show todo list
    uncheck     uncheck the todo

Use "todoly help <command>" for more information about a command.
```

## Installation

```
$ go get github.com/mattn/todoly
```

## Setup

Put below into your `~/.netrc`

```
machine todo.ly login <YOUR-EMAIL> password <YOUR-PASSWORD>
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a mattn)
