# mac2vendor
Provides a mac address to vendor lookup

## Usage

### CLI
```bash
./cmd/cli/mac2vendor/mac2vendor -quiet -mac 84:38:35:77:aa:52
```

### Library

```go
package main

import "fmt"
import m2v "github.com/n3integration/mac2vendor"

func main() {
  mac2vnd, err := m2v.Load(m2v.Dat)
  if err != nil {
    fmt.Println("error:", err)
  }

  vnd, err := mac2vnd.Lookup("84:38:35:70:aa:52")
  if err != nil {
    fmt.Println("lookup error:", err)
  } else if vnd == "" {
    fmt.Println("not found") 
  } else {
    fmt.Println("found ==>", vnd)
  }
}
```
