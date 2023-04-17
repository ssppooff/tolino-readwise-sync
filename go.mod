module github.com/ssppooff/tolino-readwise-sync

go 1.20

replace (
	github.com/ssppooff/tolino-readwise-sync/readwise => ./readwise
	github.com/ssppooff/tolino-readwise-sync/tolino => ./tolino
	github.com/ssppooff/tolino-readwise-sync/utils => ./utils
)

require golang.org/x/exp v0.0.0-20230321023759-10a507213a29
