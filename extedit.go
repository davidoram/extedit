// Package extedit provides functionality to open an editor with to let the user edit
// some content and get the changes the user made to process them as
// part of the user iterface of a command line programm.
package extedit

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const defaultEditor = "vim"

// Session represents an extedit session.
type Session struct {
	input     Content
	result    Content
	SplitFunc bufio.SplitFunc
}

// Invoke starts a text-editor with the contents of content.
// After the user has closed the editor Invoke returns an
// io.Reader with the edited content.
func (s *Session) Invoke(content io.Reader) (Diff, error) {
	d := Diff{}

	input, err := contentFromReader(content, s.SplitFunc)
	if err != nil {
		return d, err
	}

	fileName, err := writeTmpFile(input)
	if err != nil {
		return d, err
	}
	defer os.Remove(fileName)

	cmd := editorCmd(fileName)
	err = cmd.Run()
	if err != nil {
		return d, err
	}

	result, err := contentFromFile(fileName, s.SplitFunc)
	if err != nil {
		return d, err
	}

	return NewDiff(input, result), nil
}

// NewSession constructs a new Session with a default SplitFunc value.
func NewSession() *Session {
	return &Session{SplitFunc: bufio.ScanLines}
}

// Invoke is a shortcut to the Invoke method of a default session.
func Invoke(content io.Reader) (Diff, error) {
	s := NewSession()
	return s.Invoke(content)
}

// writeTmpFile writes content to a temporary file and returns
// the path to the file.
func writeTmpFile(content io.Reader) (string, error) {
	f, err := ioutil.TempFile("", "")

	if err != nil {
		return "", err
	}

	io.Copy(f, content)
	f.Close()
	return f.Name(), nil
}

// editorCmd creates a os/exec.Cmd to open.
// filename in an editor ready to be run().
func editorCmd(filename string) *exec.Cmd {
	editorEnvar := os.Getenv("EDITOR")
	if editorEnvar == "" {
		editorEnvar = defaultEditor
	}
	// editorEnvar might contain editor arguments which
	// need to be split out eg: "code --wait --new-window"
	editorArgs := strings.Split(editorEnvar, " ")
	cmd, editorArgs := editorArgs[0], editorArgs[1:]
	editorArgs = append(editorArgs, filename)
	fmt.Printf("%s %s", cmd, editorArgs)
	editor := exec.Command(cmd, editorArgs...)

	editor.Stdin = os.Stdin
	editor.Stdout = os.Stdout
	editor.Stderr = os.Stderr

	return editor
}
