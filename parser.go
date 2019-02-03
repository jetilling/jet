package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type TokenType int

const (
	NativeFunction      TokenType = 0
	Function            TokenType = 1
	Parameter           TokenType = 2
	Attribute           TokenType = 4
	Number              TokenType = 5
	Word                TokenType = 6
	OpenBracket         TokenType = 7
	CloseBracket        TokenType = 8
	OpenParenthesis     TokenType = 9
	CloseParenthesis    TokenType = 10
	EmptyParenthesis    TokenType = 11
	OpenCurly           TokenType = 12
	CloseCurly          TokenType = 13
	VariableDeclaration TokenType = 14
	FunctionDeclaration TokenType = 15
	CallExpression      TokenType = 16
	RootExpression      TokenType = 17
	Value               TokenType = 18
)

var standardLibraryContents = [13]string{"Eq", "Func", "Var", "Add", "Run", "Update", "Jetlog", "Div", "Button", "Span", "Click", "Text"}

type Tokens struct {
	State map[int]*Body
	Do    map[int]*Body
	View  map[int]*Body
}

// type RootToken struct {
// 	Name      string
// 	Type      TokenType
// 	Body      map[int]*Body
// 	BodyCount int
// 	completed bool
// }

// type DoToken struct {
// 	Name      string
// 	Type      TokenType
// 	Body      map[int]*Func
// 	BodyCount int
// 	completed bool
// }

// type ViewToken struct {
// 	Type      TokenType
// 	Body      map[int]Body
// 	BodyCount int
// 	completed bool
// }

// SubStructs

// type Func struct {
// 	Name           string
// 	NamePos        string
// 	Type           TokenType
// 	Kind           string
// 	Parameters     map[int]*Params
// 	Arguments      map[int]*Body
// 	ArgumentCount  int
// 	ParameterCount int
// }

type Body struct {
	Name           string
	NamePos        string
	Type           TokenType
	Kind           string
	Parameters     map[int]*Params
	Arguments      map[int]*Body
	ArgumentCount  int
	ParameterCount int
}

type Params struct {
	Name    string
	NamePos string
}

// Lexical Structs

type LexicalItems struct {
	Items []Lexical
}

type Lexical struct {
	Type   TokenType
	Value  string
	IsRoot bool
	Group  string
}

// /*
// *
// *
//  */
// func (tokens *Tokens) AddStateToken(token StateToken) []StateToken {
// 	tokens.State = append(tokens.State, token)
// 	return tokens.State
// }

/*
*
*
 */
func (lexicalItems *LexicalItems) AddLexicalItem(lexical Lexical) []Lexical {
	lexicalItems.Items = append(lexicalItems.Items, lexical)
	return lexicalItems.Items
}

/*
*
*
 */
func ReadInFile(fileName string) LexicalItems {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	return lexer(scanner)

}

/*
*
*
 */
func lexer(scanner *bufio.Scanner) LexicalItems {
	items := []Lexical{}
	lexicalList := LexicalItems{items}

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {

			words := strings.Fields(line)

			for _, value := range words {
				fmt.Println("VALUE: ", value)
				//Maybe put a check here to see if value has a ' and if it does create a flag to
				// run a function that reconstructs the string until the next ' is found, at which point,
				// the string should be added to the lexical list

				// Checks if value is a number
				if _, err := strconv.Atoi(value); err == nil {
					item := Lexical{Type: Number, Value: value, IsRoot: false}
					lexicalList.AddLexicalItem(item)
				} else {
					switch value {
					case "()":
						item := Lexical{Type: EmptyParenthesis, Value: value, IsRoot: false}
						lexicalList.AddLexicalItem(item)
					case "(":
						item := Lexical{Type: OpenParenthesis, Value: value, IsRoot: false}
						lexicalList.AddLexicalItem(item)
					case ")":
						item := Lexical{Type: CloseParenthesis, Value: value, IsRoot: false}
						lexicalList.AddLexicalItem(item)
					case "[":
						item := Lexical{Type: OpenBracket, Value: value, IsRoot: false}
						lexicalList.AddLexicalItem(item)
					case "]":
						item := Lexical{Type: CloseBracket, Value: value, IsRoot: false}
						lexicalList.AddLexicalItem(item)
					case "{":
						item := Lexical{Type: OpenCurly, Value: value, IsRoot: false}
						lexicalList.AddLexicalItem(item)
					case "}":
						item := Lexical{Type: CloseCurly, Value: value, IsRoot: false}
						lexicalList.AddLexicalItem(item)
					default:
						switch {
						case value == "State":
							item := Lexical{Type: Word, Value: value, IsRoot: true, Group: "State"}
							lexicalList.AddLexicalItem(item)
						case value == "Do":
							item := Lexical{Type: Word, Value: value, IsRoot: true, Group: "Do"}
							lexicalList.AddLexicalItem(item)
						case value == "View":
							item := Lexical{Type: Word, Value: value, IsRoot: true, Group: "View"}
							lexicalList.AddLexicalItem(item)
						case StringInSlice(value, standardLibraryContents):
							item := Lexical{Type: CallExpression, Value: value, IsRoot: false, Group: ""}
							lexicalList.AddLexicalItem(item)
						default:
							item := Lexical{Type: Word, Value: value, IsRoot: false, Group: ""}
							lexicalList.AddLexicalItem(item)
						}

					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return lexicalList

}

/*
*
*
*
 */
func Parser(lexicalList LexicalItems) {
	tokenList := Tokens{}
	currentStep := make(map[int]*Body)
	index := 0
	pushAsArgument := false
	pushAsParameter := false
	block := ""

	tokenList = buildAST(lexicalList, index, tokenList, &currentStep, block, pushAsArgument, pushAsParameter)

	// fmt.Println("STUFF: ", tokenList)
	// fmt.Println(" ")
	// fmt.Println("FINAL: ", tokenList)
	// fmt.Println(" ")
	// fmt.Println("STATE: ", tokenList.State)
	// fmt.Println("STATE BODY: ", tokenList.State)
	// fmt.Println("STATE BODY 0: ", tokenList.State[0])
	// fmt.Println("STATE BODY 0 ARGUMENTS: ", tokenList.State[0].Arguments)
	// fmt.Println("STATE BODY 0 ARGUMENTS 0: ", tokenList.State[0].Arguments[0])
	// fmt.Println("STATE BODY 0 ARGUMENTS 1: ", tokenList.State[0].Arguments[1])
	// fmt.Println("STATE BODY 1: ", tokenList.State[1])
	// fmt.Println("STATE BODY 1 ARGUMENTS: ", tokenList.State[1].Arguments)
	// fmt.Println("STATE BODY 1 ARGUMENTS 0: ", tokenList.State[1].Arguments[0])
	// fmt.Println("STATE BODY 1 ARGUMENTS 1: ", tokenList.State[1].Arguments[1])
	// fmt.Println("----------------------------------------")
	// fmt.Println("DO: ", tokenList.Do)
	// fmt.Println("DO BODY: ", tokenList.Do)
	// fmt.Println("DO BODY 0: ", tokenList.Do[0])
	// fmt.Println("DO BODY 0 Parameters: ", tokenList.Do[0].ParameterCount)
	// fmt.Println("DO BODY 0 ARGUMENTS 0: ", tokenList.Do[0].Arguments[0])
	// fmt.Println("DO BODY 0 ARGUMENTS 1: ", tokenList.Do[0].Arguments[1])
	// fmt.Println("DO BODY 0 ARGUMENTS 2: ", tokenList.Do[0].Arguments[2])
	// fmt.Println("DO BODY 0 ARGUMENTS 3: ", tokenList.Do[0].Arguments[3])
	// fmt.Println("----------------------------------------")
	// PrettyPrint(tokenList.Do)
	// fmt.Println("VIEW: ", tokenList.View)
	// fmt.Println("VIEW BODY: ", tokenList.View.Body)
	// fmt.Println("----------------------------------------")
	PrettyPrint(tokenList.View)
}

/*
*
*
*
 */
func buildAST(lexicalList LexicalItems, index int, tokenList Tokens, currentStep *map[int]*Body, block string, pushAsArgument, pushAsParameter bool) Tokens {
	if index <= len(lexicalList.Items)-1 {
		item := lexicalList.Items[index]
		switch item.Type {
		case 5, 6, 16:
			if item.IsRoot {
				block = item.Value
			} else {
				switch block {
				case "State":
					handleStateBlock(item, currentStep)
				case "Do":
					index = handleDoBlock(lexicalList, index, currentStep, pushAsArgument, pushAsParameter)
				case "View":
					index = handleViewBlock(lexicalList, index, currentStep, pushAsArgument, pushAsParameter)
				}
			}
		case 7:
			// Open Bracket
			pushAsArgument = true
		case 8:
			// Close Bracket
			pushAsArgument = false
		case 9:
			// Open Parenthesis
			pushAsParameter = true
			// fmt.Println("If in do, then this is a parameter else its an attribut list....hmm maybe refactor")
		case 10:
			// Close Parenthesis
			pushAsParameter = false
		case 11:
			// fmt.Println("CASE 11 (EmptyParenthesis): ", value.Value)
			// fmt.Println("Move to next thing?")
		case 12:
			// fmt.Println("CASE 12 (OpenCurly): ", value.Value)
			// fmt.Println("Do we care about this? Maybe throw an error if it doesn't exist?")
		case 13:
			switch block {
			case "State":
				// Completing State Section
				tokenList.State = (*currentStep)
				newMap := make(map[int]*Body)
				currentStep = &newMap
			case "Do":
				// Completing Do Section
				tokenList.Do = (*currentStep)
				newMap := make(map[int]*Body)
				currentStep = &newMap
			case "View":
				// Completing View Section
				tokenList.View = (*currentStep)
				newMap := make(map[int]*Body)
				currentStep = &newMap
			}
		}
		index = index + 1
		// fmt.Println("before return: ", currentStep, " block: ", block)
		return buildAST(lexicalList, index, tokenList, currentStep, block, pushAsArgument, pushAsParameter)
	} else {
		fmt.Println("COMPLETE!")
	}
	return tokenList
}

//----------------------------- HELPERS ------------------------------//

/*
*
*
*
 */
func handleStateBlock(value Lexical, currentStep *map[int]*Body) {
	if StringInSlice(value.Value, standardLibraryContents) {
		if len(*currentStep) < 1 {
			m := make(map[int]*Body)
			m[0] = &Body{Name: value.Value, Type: CallExpression, Kind: "Func"}
			*currentStep = m
		} else {
			(*currentStep)[len(*currentStep)] = &Body{Name: value.Value, Type: CallExpression, Kind: "Func"}
		}
	} else {
		currentBodyCount := VerifyIndex(len(*currentStep) - 1)
		argumentLength := len((*currentStep)[currentBodyCount].Arguments)
		if argumentLength < 1 {
			m := make(map[int]*Body)
			m[0] = &Body{Name: value.Value, Type: FunctionDeclaration, Kind: "Func"}
			(*currentStep)[currentBodyCount].Arguments = m
		} else {
			(*currentStep)[currentBodyCount].Arguments[argumentLength] = &Body{Name: value.Value, Type: VariableDeclaration, Kind: "SL"}
		}
	}
	return
}

/*
*
*
*
 */
func handleDoBlock(lexicalList LexicalItems, index int, currentStep *map[int]*Body, pushAsArgument, pushAsParameter bool) int {
	value := lexicalList.Items[index]
	// Everything here should start as Func, as only Functions will be listed in Do
	if !pushAsArgument && !pushAsParameter {
		if len(*currentStep) < 1 {
			m := make(map[int]*Body)
			m[0] = &Body{Name: value.Value, Type: CallExpression, Kind: "Func"}
			*currentStep = m
		} else {
			(*currentStep)[len(*currentStep)] = &Body{Name: value.Value, Type: CallExpression, Kind: "Func"}
		}

	} else {

		currentBodyIndex := VerifyIndex(len(*currentStep) - 1)
		if pushAsParameter {
			parametersLength := len((*currentStep)[currentBodyIndex].Parameters)
			(*currentStep)[currentBodyIndex] = addParameterToAST(parametersLength, value.Value, (*currentStep)[currentBodyIndex])
		} else {
			argumentLength := len((*currentStep)[currentBodyIndex].Arguments)
			if value.Type == 16 {
				(*currentStep)[currentBodyIndex] = addArgumentToAST(argumentLength, value.Value, (*currentStep)[currentBodyIndex])
				index = index + 1
				arguments := &(*currentStep)[currentBodyIndex].Arguments
				return handleDoBlock(lexicalList, index, arguments, pushAsArgument, pushAsParameter)
			} else {
				(*currentStep)[currentBodyIndex] = addArgumentToAST(argumentLength, value.Value, (*currentStep)[currentBodyIndex])

				// Check if the next value is a word or number, if so then we want to add it to the current argument map
				if lexicalList.Items[index+1].Type == 6 || lexicalList.Items[index+1].Type == 5 {
					index = index + 1
					value = lexicalList.Items[index]
					(*currentStep)[currentBodyIndex] = addArgumentToAST(argumentLength+1, value.Value, (*currentStep)[currentBodyIndex])
					// NOTE: Currently a standard library function can only accept a maximum of 2 arguments
					// To fix this, another handleDoBlock, but we don't want to pass in the next arguments map, that's the challenge
				}
			}
		}
	}
	return index
}

/*
*
*
*
 */
func handleViewBlock(lexicalList LexicalItems, index int, currentStep *map[int]*Body, pushAsArgument, pushAsParameter bool) int {
	value := lexicalList.Items[index]
	fmt.Println("IN VIEW: ", value.Value)

	return index
}

/*
*
*
*
 */
func addArgumentToAST(argumentLength int, value string, body *Body) *Body {
	if argumentLength < 1 {
		m := make(map[int]*Body)
		m[0] = &Body{Name: value, Type: CallExpression, Kind: "Func"}
		// fmt.Println("body before: ", body)
		body.Arguments = m
		body.ArgumentCount = argumentLength + 1
		// fmt.Println("body after: ", body)
	} else {
		body.Arguments[argumentLength] = &Body{Name: value, Type: VariableDeclaration, Kind: "SL"}
		body.ArgumentCount = argumentLength + 1
	}
	return body
}

/*
*
*
*
 */
func addParameterToAST(parametersLength int, value string, body *Body) *Body {
	if parametersLength < 1 {
		m := make(map[int]*Params)
		m[0] = &Params{Name: value}
		body.Parameters = m
		body.ParameterCount = parametersLength + 1
		// fmt.Println("PUSHING AS PARAMETER: ", value)
	} else {
		body.Parameters[parametersLength] = &Params{Name: value}
		body.ParameterCount = parametersLength + 1
	}

	return body
}
