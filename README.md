# link-preview

A Go module gathers elements required for linkPreview.

## Usage
```go
package main

import (
	"fmt"
	"github.com/junzki/link-preview"
)

func main() {
    link := "http://custom-domain.local/case.html"
    
    result, err := link_preview.PreviewLink(link, nil)
    if err != nil {
    	panic(err)
    }
    
    fmt.Println(result.Title)
}
```

## References:
Thanks to [aakash4525]'s [py_link_Preview], this package is mostly inspired by his awesome work. 


[aakash4525]: https://github.com/aakash4525
[py_link_Preview]: https://github.com/aakash4525/py_link_preview

## License
_**BSD 3-Clause License**_
