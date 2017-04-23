### Makefile --- 

## Author: shell@dskvmdeb.lan
## Version: $Id: Makefile,v 0.0 2015/06/27 08:36:39 shell Exp $
## Keywords: 
## X-URL: 

all: build

build:
	go build -o bin/tsfm github.com/shell909090/scheme-go/tsfm
	bin/tsfm -base code/macro.scm code/transformer.scm
# mkdir -p bin
# go build -o bin/scheme-go github.com/shell909090/scheme-go/main

clean:
	rm -rf bin

### Makefile ends here
