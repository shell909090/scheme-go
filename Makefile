### Makefile --- 

## Author: shell@dskvmdeb.lan
## Version: $Id: Makefile,v 0.0 2015/06/27 08:36:39 shell Exp $
## Keywords: 
## X-URL: 

all: build

build:
	mkdir -p bin
	go build -o bin/scheme-go bitbucket.org/shell909090/scheme-go/main

### Makefile ends here
