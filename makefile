all: test

test: fixtag
	prove

fixtag: fixtag.go
	go build fixtag.go
