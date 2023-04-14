module github.com/ssppooff/tolino-readwise-sync

go 1.20

replace (
	github.com/ssppooff/tolino-readwise-sync/readwise => ./readwise
	github.com/ssppooff/tolino-readwise-sync/tolino => ./tolino
	github.com/ssppooff/tolino-readwise-sync/utils => ./utils
)
