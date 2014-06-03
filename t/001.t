#!/bin/sh
cd "`dirname $0`"

tags_fixed=tags.001.new
tags_expected=tags.001.expected

ctags 001.c
../fixtag tags > $tags_fixed
if ! diff $tags_expected $tags_fixed >/dev/null; then
	echo "not ok 1 - basic"
else
	echo "ok 1 - basic"
fi
echo "1..1"
