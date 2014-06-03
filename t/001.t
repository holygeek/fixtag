#!/bin/sh
cd "`dirname $0`"
ctags foo.c
../fixtag tags > tags.new
if ! diff tags.expected tags.new >/dev/null; then
	echo "not ok 1 - basic"
else
	echo "ok 1 - basic"
fi
echo "1..1"
