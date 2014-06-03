package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type Tag struct {
	items []string
}

type StructTag struct {
	*Tag
}

func NewStructTag(t *Tag) *StructTag {
	if !t.IsStructTag() {
		return nil
	}
	return &StructTag{t}
}

func (s *StructTag) StructName() string {
	return whiteSpaceRe.Split(s.Locator(), 4)[2]
}

type Tags struct {
	Lines  []string
	tags   []*Tag
	tagMap map[string][]*Tag
}

func NewTags() *Tags {
	t := &Tags{}
	t.tagMap = make(map[string][]*Tag)
	return t
}

func (t *Tags) Add(line string) {
	t.Lines = append(t.Lines, line)
	tag := NewTag(line)
	t.tags = append(t.tags, tag)
	if tags, exist := t.tagMap[tag.Tag()]; exist {
		tags = append(tags, tag)
	} else {
		t.tagMap[tag.Tag()] = []*Tag{tag}
	}
}

func (t *Tags) Get(tag string) []*Tag {
	return t.tagMap[tag]
}

func (t *Tags) IsUnique(tag string) bool {
	return len(t.tagMap[tag]) == 1
}

func (t *Tags) Tags() []*Tag {
	return t.tags
}

func NewTag(line string) *Tag {
	return &Tag{strings.Split(line, "\t")}
}

func (o *Tag) String() string {
	return strings.Join(o.items, "\t")
}

func (o *Tag) Tag() string {
	return o.items[0]
}

func (o *Tag) File() string {
	return o.items[1]
}

func (o *Tag) Locator() string {
	return o.items[2]
}

func (o *Tag) SetLocator(l string) {
	o.items[2] = l
}

func (o *Tag) UsesRegex() bool {
	return o.Locator()[0] == '/'
}

var structTagRe = regexp.MustCompile(`^/\^typedef\s+struct\b`)
var whiteSpaceRe = regexp.MustCompilePOSIX(`  *`)

func (o *Tag) IsStructTag() bool {
	return len(structTagRe.FindString(o.Locator())) > 0
}

func usage() {
	fmt.Println(`Usage: fixtags [tagfile]

If not given, [tagfile] defaults to "tags".
`)
}

func mustScanTagFile(tagfile string) *bufio.Scanner {
	f, err := os.Open(tagfile)
	if err != nil {
		log.Fatal(err)
	}

	return bufio.NewScanner(f)
}

func NewTagsFromFile(tagfile string) *Tags {
	scanner := mustScanTagFile(tagfile)

	var tags = NewTags()
	for scanner.Scan() {
		line := scanner.Text()
		tags.Add(line)
	}
	return tags
}

func mustGetTagFile() string {
	if len(os.Args) > 2 {
		usage()
		os.Exit(1)
	}

	tagfile := "tags"
	if len(os.Args) == 2 {
		tagfile = os.Args[1]
	}
	return tagfile
}

func main() {
	tags := NewTagsFromFile(mustGetTagFile())

	for _, tag := range tags.Tags() {
		st := NewStructTag(tag)
		if st == nil {
			fmt.Println(tag)
			continue
		}
		candidates := tags.Get(st.StructName())
		matches := []*Tag{}
		for _, tag := range candidates {
			if tag.File() == st.File() {
				matches = append(matches, tag)
			}
		}
		if len(matches) == 1 {
			st.SetLocator(matches[0].Locator())
		}
		fmt.Println(st)
	}
}
