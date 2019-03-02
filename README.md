# mac2vendor
Provides a mac address to vendor lookup

## Usage

### CLI
```bash
./mac2vendor resolve -mac 84:38:35:77:aa:52 [-quiet]
```

### Library

```go
package main

import "fmt"
import m2v "github.com/n3integration/mac2vendor"

func main() {
  vnd, err := m2v.Lookup("84:38:35:70:aa:52")
  if err != nil {
    fmt.Println("lookup error:", err)
  } else if vnd == "" {
    fmt.Println("not found")
  } else {
    fmt.Println("found ==>", vnd)
  }
}
```

### Web Service

```bash
./mac2vendor serve [-port 9000]
```

```curl
curl -siv 127.0.0.1:9000/84:38:35:70:aa:52
```

## License

Copyright 2019 n3integration@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
