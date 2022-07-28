<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# args

```go
import "klogs/args"
```

## Index

- [type Args](<#type-args>)
  - [func Parse(argv []string) *Args](<#func-parse>)
  - [func (a *Args) Usage() string](<#func-args-usage>)


## type Args

Args encapsulates all the various flags/options for klogs

```go
type Args struct {
    Help       bool
    Version    bool
    Query      []string `short:"" long:""`
    All        bool
    Label      []string
    LimitBytes string `short:"" long:"limit-bytes"`
    Since      string
    SinceTime  string `short:"" long:"since-time"`
    Tail       string `short:"" long:"tail"`
    Follow     bool
    Timestamps bool `short:"" long:"timestamps"`
    Previous   bool `short:""`
    KubeConfig string
    Context    string `short:"C"`
    Container  string
    Namespace  string
    Prefix     bool
    JSON       bool
    Theme      string
}
```

### func Parse

```go
func Parse(argv []string) *Args
```

Parse takes an array of string args and returns the parsed Args struct

### func \(\*Args\) Usage

```go
func (a *Args) Usage() string
```

Usage returns the documentation string for the command



Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)