package vmtranslator

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

// VMCommand is a string that represents an instruction. Example: "push constant 7", "add", "label LOOP"
type VMCommand string

/*
VMCommandType is an enum that represents the type of an instruction.
C_ARITHMETIC: add, sub, neg, eq, gt, lt, and, or, not
C_PUSH: push segment i
C_POP: pop segment i
C_LABEL: label label
C_GOTO: goto label
C_IF: if-goto label
C_FUNCTION: function function n
C_RETURN: return
C_CALL: call function n
*/
type VMCommandType int

const (
	C_ARITHMETIC VMCommandType = iota
	C_PUSH
	C_POP
	C_LABEL
	C_GOTO
	C_IF
	C_FUNCTION
	C_RETURN
	C_CALL
)

// TODO: integrate this code into the hack assembler project
// CodeScanner is a struct that reads a file line by line and skips empty lines and comments. It provides a method to get the current line of the scanner without leading and trailing spaces and comments.
// TODO: support multiple comment prefixes. Example: "//" and "#"
// TODO: support multiple comment styles. Example: "//" and "/* */"
type CodeScanner struct {
	scanner       *bufio.Scanner
	commentPrefix string // the prefix that indicates a comment. Example: "//"
}

// NewCodeScanner creates a new CodeScanner with the given reader and comment prefix
func NewCodeScanner(r io.Reader, commentPrefix string) CodeScanner {
	return CodeScanner{scanner: bufio.NewScanner(r), commentPrefix: commentPrefix}
}

// Parser is a struct that reads VM commands from a file and provides methods to get the command type and arguments
type Parser struct {
	scanner        CodeScanner
	currentCommand VMCommand
	currentType    VMCommandType
}

// NewParser creates a new Parser with the given reader and comment prefix. It uses a [CodeScanner] to read the file. commentPrefix is the prefix that indicates a comment. Example: "//"
func NewParser(r io.Reader, commentPrefix string) Parser {
	cs := NewCodeScanner(r, commentPrefix)
	return Parser{scanner: cs}
}

// isEmptyLine returns true if the line is empty
func (cs CodeScanner) isEmptyLine(line string) bool {
	return len(line) == 0
}

// isCommentLine returns true if the line is a comment line.
func (cs CodeScanner) isCommentLine(line string) bool {
	line = strings.TrimSpace(line)
	if len(line) < 2 {
		return false
	}
	// check if the line starts with the comment prefix
	p := cs.commentPrefix
	return line[0:len(p)] == p
}

// Text returns the current line of the scanner without leading and trailing spaces and comments. It replaces multiple spaces with a single space and removes comments at the end of the line.
func (cs CodeScanner) Text() string {
	text := cs.scanner.Text()
	//Remove comment at the end of the line
	text = strings.Split(text, cs.commentPrefix)[0]
	// remove spaces at the beginning and end of the line
	// replace multiple spaces with a single space
	return strings.Join(strings.Fields(text), " ")
}

// Scan reads the next line from the scanner and skips empty lines and comments. It returns false if there are no more lines.
func (cs CodeScanner) Scan() bool {
	ok := cs.scanner.Scan()
	if !ok {
		return false
	}
	line := cs.Text()
	// skip empty or comment line
	if cs.isEmptyLine(line) || cs.isCommentLine(line) {
		return cs.Scan()
	}
	return true
}

// getCommandType returns the type of the given command.
func getCommandType(command VMCommand) VMCommandType {
	// split the command into words
	words := strings.Fields(string(command))
	switch words[0] {
	case "add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not":
		return C_ARITHMETIC
	case "push":
		return C_PUSH
	case "pop":
		return C_POP
	case "label":
		return C_LABEL
	case "goto":
		return C_GOTO
	case "if-goto":
		return C_IF
	case "function":
		return C_FUNCTION
	case "return":
		return C_RETURN
	case "call":
		return C_CALL
	default:
		panic("unknown command")
	}
}

// advance reads the next instruction from the input and makes it the current instruction. It returns false if there are no more instructions. advance ignores empty lines and comments.
func (p *Parser) advance() bool {
	ok := p.scanner.Scan()
	if !ok {
		return false
	}
	line := p.scanner.Text()
	command := VMCommand(line)
	p.currentCommand = command
	p.currentType = getCommandType(command)
	return true
}

// arg1 returns the first argument of the current instruction. It returns the command itself if it is an arithmetic command. It panics if the command is a return command.
func arg1(command VMCommand) string {
	words := strings.Fields(string(command))
	switch ctype := getCommandType(command); ctype {

	case C_RETURN:
		panic("return command has no arguments")
	case C_ARITHMETIC:
		return words[0]
	case C_PUSH, C_POP, C_LABEL, C_GOTO, C_IF, C_FUNCTION, C_CALL:
		return words[1]
	default:
		panic("command has no first argument")
	}

}

// arg2 returns the second argument of the current instruction. This is valid only for push, pop, function, and call commands.
func arg2(command VMCommand) int {
	words := strings.Fields(string(command))
	switch ctype := getCommandType(command); ctype {
	case C_PUSH, C_POP, C_FUNCTION, C_CALL:
		n, err := strconv.Atoi(words[2])
		if err != nil {
			panic(err)
		}
		return n
	default:
		panic("command has no second argument")
	}
}
