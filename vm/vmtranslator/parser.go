package vmtranslator

import (
	"bufio"
	"strconv"
	"strings"
)

// VMCommand is a string that represents an instruction
// Example: "push constant 7", "add", "label LOOP"
type VMCommand string

// VMCommandType is an enum that represents the type of an instruction.
// C_ARITHMETIC: add, sub, neg, eq, gt, lt, and, or, not
// C_PUSH: push segment i
// C_POP: pop segment i
// C_LABEL: label label
// C_GOTO: goto label
// C_IF: if-goto label
// C_FUNCTION: function function n
// C_RETURN: return
// C_CALL: call function n
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

type CodeScanner struct {
	scanner       *bufio.Scanner
	commentPrefix string
}

type Parser struct {
	scanner        CodeScanner
	currentCommand VMCommand
	currentType    VMCommandType
}

// advance reads the next instruction from the input and makes it the current instruction
// It returns false if there are no more instructions
// advance ignores empty lines and comments

func (cs CodeScanner) isEmptyLine(line string) bool {
	return len(line) == 0
}

func (cs CodeScanner) isCommentLine(line string) bool {
	return line[0:2] == cs.commentPrefix
}

func (cd CodeScanner) text() string {
	text := cd.scanner.Text()
	// remove spaces at the beginning and end of the line
	// replace multiple spaces with a single space
	return strings.Join(strings.Fields(text), " ")
}

func (cs CodeScanner) scan() bool {
	ok := cs.scanner.Scan()
	if !ok {
		return false
	}
	line := cs.text()
	// skip empty or comment line
	if cs.isEmptyLine(line) || cs.isCommentLine(line) {
		return cs.scan()
	}
	return true
}

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

func (p *Parser) advance() bool {
	ok := p.scanner.scan()
	if !ok {
		return false
	}
	line := p.scanner.text()
	command := VMCommand(line)
	p.currentCommand = command
	p.currentType = getCommandType(command)
	return true
}

// arg1 returns the first argument of the current instruction
// It returns the command itself if it is an arithmetic command
// It panics if the command is a return command
func (p *Parser) arg1() string {
	words := strings.Fields(string(p.currentCommand))
	if p.currentType == C_RETURN {
		panic("return command has no arguments")
	} else if p.currentType == C_ARITHMETIC {
		return words[0]
	}
	return words[1]
}

// arg2 returns the second argument of the current instruction
// It panics if the command has no second argument
func (p *Parser) arg2() int {
	words := strings.Fields(string(p.currentCommand))
	switch p.currentType {
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
