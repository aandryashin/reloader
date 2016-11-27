# Reloader
[![Build Status](https://travis-ci.org/aandryashin/reloader.svg?branch=master)](https://travis-ci.org/aandryashin/reloader)
[![Coverage](https://codecov.io/github/aandryashin/reloader/reloader.svg)](https://codecov.io/gh/aandryashin/reloader)
[![Release](https://img.shields.io/github/release/aandryashin/reloader.svg)](https://github.com/aandryashin/reloader/releases/latest)

Library for automatic reloading configuration files based on file system events.

## Using

```go
package main

import (
	"log"
	"time"

	"github.com/aandryashin/reloader"
)

func loadConfig() {
	//...
}

func main() {
	err := reloader.Watch("config", loadConfig, 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	//...
}
```
