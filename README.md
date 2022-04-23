# iomock
Mock Go io interfaces.

## Usage
### Install
Add iomock to your project.

```sh
go get github.com/ktnyt/iomock
```

### Mock a read error
```go
package main

import (
	"fmt"
	"io"

	"github.com/ktnyt/iomock"
)

func main() {
	r := iomock.NewReader(func(p []byte) (int, error) {
		return 0, io.ErrClosedPipe
	})

	p := make([]byte, 3)

	// Output: 0 io: read/write on closed pipe
	fmt.Println(r.Read(p))
}
```

### Mock a write error
```go
package main

import (
	"fmt"
	"io"

	"github.com/ktnyt/iomock"
)

func main() {
	w := iomock.NewWriter(func(p []byte) (int, error) {
		return 0, io.ErrClosedPipe
	})

	// Output: 0 io: read/write on closed pipe
	fmt.Println(w.Write([]byte("foo")))
}
```
