package extedit

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// Content represents the input and output of an extedit session.
type Content struct {
	c      []string
	reader io.Reader
}

// Read satisfies the reader io.Reader interface.
func (c Content) Read(b []byte) (int, error) {
	return c.reader.Read(b)
}

// String satisfies the fmt.Stringer interface.
func (c Content) String() string {
	return strings.Join(c.c, "\n") // <- FIXME(kioopi): join string needs to come from splitfunc somehow
}

// Length returns the number of lines in Content.
func (c Content) Length() int {
	return len(c.c)
}

// contentFromReader creates a new Content object by scanning an io.Reader using a bufio.SplitFunc.
func contentFromReader(content io.Reader, split bufio.SplitFunc) (Content, error) {
	c := Content{}
	scanner := bufio.NewScanner(content)
	scanner.Split(split)

	for scanner.Scan() {
		c.c = append(c.c, scanner.Text())
	}
	c.reader = strings.NewReader(c.String())

	return c, scanner.Err()
}

func contentFromFile(filename string, split bufio.SplitFunc) (Content, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Content{}, err
	}
	defer file.Close()

	return contentFromReader(file, split)
}

func contentFromString(content string, split bufio.SplitFunc) (Content, error) {
	return contentFromReader(strings.NewReader(content), split)
}
