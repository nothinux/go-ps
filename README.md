# go-ps
go-ps is process library for find, list, and get information from process. go-ps read information about process from `/proc` file.

## Documentation
see [pkg.go.dev](https://pkg.go.dev/github.com/nothinux/go-ps)


## Installation

```
$ go get github.com/nothinux/go-ps
```

### Getting Started
#### Get All Running Process Name
``` go
package main

import (
    "log"
    "fmt"
    "github.com/nothinux/go-ps"
)

func main() {
    process, err := ps.GetProcess()
    if err != nil {
        log.Fatal(err)
    }

    for _, p := range process {
        fmt.Println(p.Comm)
    }
}
```

#### Find Pid from Process Name
``` go
package main

import (
    "log"
    "github.com/nothinux/go-ps"
    "fmt"
)

func main() {
    pid, err := ps.FindPid("nginx")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(pid)
}
```

#### Find Process then get information from that process
``` go
package main


import (
    "log"
    "github.com/nothinux/go-ps"
    "fmt"
)

func main() {
    p, err := ps.FindProcessName("nginx")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("process id:", p.Pid)
    fmt.Println("process name:", p.Comm)
    fmt.Println("process cmd:", p.CmdLine)
    fmt.PrintLn("process state:", p.State)

}
`` `

[more](https://pkg.go.dev/github.com/nothinux/go-ps)


## LICENSE
[MIT](https://github.com/nothinux/go-ps/blob/master/LICENSE)
