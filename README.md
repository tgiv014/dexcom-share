# dexcom-share

`dexcom-share` provides a golang client for the undocumented Dexcom Share API. This client can be used to retrieve near-realtime glucose data uploaded from a Dexcom device.

# Usage

```go
package main

import (
	"fmt"
  "log"

	dexcomshare "github.com/tgiv014/dexcom-share"
)

func main() {
	c, err := dexcomshare.NewClient("username", "password")
	if err != nil {
		log.Fatal(err)
	}

	entries, err := c.ReadGlucose(1440, 1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(entries)
}
```

