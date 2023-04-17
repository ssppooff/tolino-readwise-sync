module github.com/ssppooff/tolino-readwise-sync

go 1.20

replace (
	github.com/ssppooff/tolino-readwise-sync/readwise => ./readwise
	github.com/ssppooff/tolino-readwise-sync/tolino => ./tolino
	github.com/ssppooff/tolino-readwise-sync/utils => ./utils
)

require (
	github.com/alexflint/go-arg v1.4.3
	golang.org/x/exp v0.0.0-20230321023759-10a507213a29
)

require github.com/alexflint/go-scalar v1.1.0 // indirect
