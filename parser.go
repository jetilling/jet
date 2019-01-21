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
	State RootToken
	Do    RootToken
	View  RootToken
}

type RootToken struct {
	Name      string
	Type      TokenType
	StateBody map[int]*Body
	DoBody    map[int]*Func
	BodyCount int
	completed bool
}

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

type Func struct {
	Name           string
	NamePos        string
	Type           TokenType
	Kind           string
	Parameters     map[int]*Params
	Arguments      map[int]*Body
	ArgumentCount  int
	ParameterCount int
}

type Body struct {
	Name          string
	NamePos       string
	Type          TokenType
	Kind          string
	Arguments     map[int]*Body
	ArgumentCount int
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
	currentStep := RootToken{}
	index := 0
	pushAsArgument := false
	pushAsParameter := false

	tokenList = buildAST(lexicalList, index, tokenList, currentStep, pushAsArgument, pushAsParameter)

	// fmt.Println("STUFF: ", tokenList)
	// fmt.Println(" ")
	// fmt.Println("FINAL: ", tokenList)
	// fmt.Println(" ")
	// fmt.Println("STATE: ", tokenList.State)
	// fmt.Println("STATE BODY: ", tokenList.State.StateBody)
	// fmt.Println("STATE BODY 0: ", tokenList.State.StateBody[1])
	// fmt.Println("STATE BODY 0 ARGUMENTS: ", tokenList.State.StateBody[1].Arguments)
	// fmt.Println("STATE BODY 0 ARGUMENTS 0: ", tokenList.State.StateBody[1].Arguments[0])
	// fmt.Println("STATE BODY 0 ARGUMENTS 1: ", tokenList.State.StateBody[1].Arguments[1])
	// fmt.Println("----------------------------------------")
	fmt.Println("DO: ", tokenList.Do)
	fmt.Println("DO BODY: ", tokenList.Do.DoBody)
	fmt.Println("DO BODY 0: ", tokenList.Do.DoBody[0])
	fmt.Println("DO BODY 0 Parameters: ", tokenList.Do.DoBody[0].ParameterCount)
	fmt.Println("DO BODY 0 ARGUMENTS: ", tokenList.Do.DoBody[0].Arguments[0])
	fmt.Println("DO BODY 0 ARGUMENTS 0: ", tokenList.Do.DoBody[0].Arguments[1])
	fmt.Println("DO BODY 0 ARGUMENTS 1: ", tokenList.Do.DoBody[0].Arguments[2])
	// fmt.Println("----------------------------------------")
	// fmt.Println("VIEW: ", tokenList.View)
	// fmt.Println("VIEW BODY: ", tokenList.View.Body)
	// fmt.Println("----------------------------------------")
}

/*
*
*
*
 */
func buildAST(lexicalList LexicalItems, index int, tokenList Tokens, currentStep RootToken, pushAsArgument, pushAsParameter bool) Tokens {
	if index <= len(lexicalList.Items)-1 {
		item := lexicalList.Items[index]
		switch item.Type {
		case 5, 6, 16:
			// fmt.Println("Item value: ", item.Value, ", Item type: ", item.Type, ", Item is Root: ", item.IsRoot)
			if item.IsRoot {
				rootToken := RootToken{Type: RootExpression, Name: item.Value}
				switch item.Value {
				case "State":
					tokenList.State = rootToken
					currentStep = tokenList.State
				case "Do":
					tokenList.Do = rootToken
					currentStep = tokenList.Do
				case "View":
					tokenList.View = rootToken
					currentStep = tokenList.View
				}
			} else {
				switch currentStep.Name {
				case "State":
					currentStep = handleStateBlock(item, currentStep)
				case "Do":
					currentStep = handleDoBlock(item, currentStep, pushAsArgument, pushAsParameter)
				case "View":
					currentStep = handleViewBlock(item, currentStep, pushAsArgument)
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
			switch currentStep.Name {
			case "State":
				// Completing State Section
				currentStep.completed = true
				tokenList.State = currentStep
			case "Do":
				// Completing Do Section
				currentStep.completed = true
				tokenList.Do = currentStep
			case "View":
				// Completing View Section
				currentStep.completed = true
				tokenList.View = currentStep
			}
		}
		index = index + 1
		return buildAST(lexicalList, index, tokenList, currentStep, pushAsArgument, pushAsParameter)
	} else {
		fmt.Println("COMPLETE!")
	}
	return tokenList
}

/*
*
* Input:
*
* Output:

*
 */
// func Parser(lexicalList LexicalItems) {
// 	tokenList := Tokens{}
// 	currentStep := RootToken{}
// 	pushAsArgument := false
// 	pushAsParameter := false

// 	for _, value := range lexicalList.Items {
// 		switch value.Type {
// 		// case 5:
// 		// 	fmt.Println("CASE 5 (Number): ", value.Value)
// 		// 	currentBodyCount := VerifyIndex(currentStep.BodyCount - 1)
// 		// 	argumentLength := len(currentStep.Body[currentBodyCount].Arguments)
// 		// 	if argumentLength < 1 {
// 		// 		m := make(map[int]*Body)
// 		// 		m[0] = &Body{Name: value.Value, Type: Number, Kind: "Number"}
// 		// 		currentStep.Body[currentStep.BodyCount].Arguments = m
// 		// 		currentStep.Body[currentStep.BodyCount].ArgumentCount = argumentLength
// 		// 	} else {
// 		// 		currentStep.Body[currentBodyCount].Arguments[argumentLength] = &Body{Name: value.Value, Type: Number, Kind: "Number"}
// 		// 		currentStep.Body[currentBodyCount].ArgumentCount = argumentLength + 1
// 		// 	}
// 		case 5:
// 		case 6:
// 			if value.IsMainGroup {
// 				item := RootToken{Type: RootExpression, Name: value.Value}
// 				switch value.Value {
// 				case "State":
// 					tokenList.State = item
// 					currentStep = tokenList.State
// 				case "Do":
// 					tokenList.Do = item
// 					currentStep = tokenList.Do
// 				case "View":
// 					tokenList.View = item
// 					currentStep = tokenList.View
// 				}
// 			} else {
// 				switch currentStep.Name {
// 				case "State":
// 					currentStep = handleStateBlock(value, currentStep)
// 				case "Do":
// 					currentStep = handleDoBlock(value, currentStep, pushAsArgument, pushAsParameter)
// 				case "View":
// 					currentStep = handleViewBlock(value, currentStep, pushAsArgument)
// 				}
// 				// Here we need to check if we're in State, Do, or View
// 				// Then run specific functions accordingly
// 				// This will combine case 5 and 6
// 				// if StringInSlice(value.Value, standardLibraryContents) {
// 				// 	if !pushAsArgument {
// 				// 		fmt.Println("CASE 6 && VALUE IN STANDARD LIBRARY: ", value.Value)
// 				// 		if len(currentStep.Body) < 1 {
// 				// 			m := make(map[int]*Body)
// 				// 			m[0] = &Body{Name: value.Value, Type: CallExpression, Kind: "Func"}
// 				// 			currentStep.Body = m
// 				// 			currentStep.BodyCount = len(currentStep.Body)
// 				// 		} else {
// 				// 			currentStep.Body[len(currentStep.Body)] = &Body{Name: value.Value, Type: CallExpression, Kind: "Func"}
// 				// 			currentStep.BodyCount = len(currentStep.Body)
// 				// 		}
// 				// 	} else {
// 				// 		fmt.Println("MAYBE THIS IS A SL EXPRESSION IN A FUNCTION BODY? ", value.Value)
// 				// 		fmt.Println("Current step: ", currentStep.Body[0])
// 				// 	}
// 				// } else {
// 				// 	fmt.Println("CASE 6 (Word): ", value.Value)
// 				// 	currentBodyCount := VerifyIndex(currentStep.BodyCount - 1)
// 				// 	argumentLength := len(currentStep.Body[currentBodyCount].Arguments)
// 				// 	// fmt.Println("CURRENT BODY COUNT: ", currentBodyCount)
// 				// 	// fmt.Println("ARGUMENT LENGTH: ", argumentLength)
// 				// 	if argumentLength < 1 {
// 				// 		m := make(map[int]*Body)
// 				// 		m[0] = &Body{Name: value.Value, Type: FunctionDeclaration, Kind: "Func"}
// 				// 		currentStep.Body[currentBodyCount].Arguments = m
// 				// 		currentStep.Body[currentBodyCount].ArgumentCount = argumentLength
// 				// 		// fmt.Println("CURRENT ARGUMENTS when 0: ", currentStep.Body[currentBodyCount].Arguments)
// 				// 	} else {
// 				// 		currentStep.Body[currentBodyCount].Arguments[argumentLength] = &Body{Name: value.Value, Type: VariableDeclaration, Kind: "??"}
// 				// 		currentStep.Body[currentBodyCount].ArgumentCount = argumentLength
// 				// 		// fmt.Println("CURRENT ARGUMENTS when +: ", currentStep.Body[currentBodyCount].Arguments)
// 				// 	}
// 				// }
// 			}
// 		case 7:
// 			// Open Bracket
// 			fmt.Println("CASE 7 (OpenBracket): ", value.Value)
// 			pushAsArgument = true
// 		case 8:
// 			// Close Bracket
// 			fmt.Println("CASE 8 (CloseBracket): ", value.Value)
// 			pushAsArgument = false
// 		case 9:
// 			// Open Parenthesis
// 			fmt.Println("CASE 9 (OpenParenthesis): ", value.Value)
// 			pushAsParameter = true
// 			fmt.Println("If in do, then this is a parameter else its an attribut list....hmm maybe refactor")
// 		case 10:
// 			// Close Parenthesis
// 			fmt.Println("CASE 10 (CloseParenthesis): ", value.Value)
// 			pushAsParameter = false
// 		case 11:
// 			// fmt.Println("CASE 11 (EmptyParenthesis): ", value.Value)
// 			// fmt.Println("Move to next thing")
// 		case 12:
// 			// fmt.Println("CASE 12 (OpenCurly): ", value.Value)
// 			// fmt.Println("Do we care about this? Maybe throw an error if it doesn't exist?")
// 		case 13:
// 			switch currentStep.Name {
// 			case "State":
// 				fmt.Println("COMPLETEING STATE!!")
// 				// Completing State Section
// 				currentStep.completed = true
// 				tokenList.State = currentStep
// 			case "Do":
// 				// Completing Do Section
// 				fmt.Println("COMPLETEING DO!!")
// 				currentStep.completed = true
// 				tokenList.Do = currentStep
// 			case "View":
// 				// Completing View Section
// 				fmt.Println("COMPLETEING VIEW!!")
// 				currentStep.completed = true
// 				tokenList.View = currentStep
// 			}

// 		}

// 	}

// 	// fmt.Println(" ")
// 	// fmt.Println("FINAL: ", tokenList)
// 	// fmt.Println(" ")
// 	// fmt.Println("STATE: ", tokenList.State)
// 	// fmt.Println("STATE BODY: ", tokenList.State.Body)
// 	// fmt.Println("STATE BODY 0: ", tokenList.State.Body[1])
// 	// fmt.Println("STATE BODY 0 ARGUMENTS: ", tokenList.State.Body[1].Arguments)
// 	// fmt.Println("STATE BODY 0 ARGUMENTS 0: ", tokenList.State.Body[1].Arguments[0])
// 	// fmt.Println("STATE BODY 0 ARGUMENTS 1: ", tokenList.State.Body[1].Arguments[1])
// 	// fmt.Println("----------------------------------------")
// 	fmt.Println("DO: ", tokenList.Do)
// 	fmt.Println("DO BODY: ", tokenList.Do.DoBody)
// 	fmt.Println("DO BODY 0: ", tokenList.Do.DoBody[0])
// 	fmt.Println("DO BODY 0 ARGUMENTS: ", tokenList.Do.DoBody[0].Arguments[0])
// 	fmt.Println("DO BODY 0 Parameters: ", tokenList.Do.DoBody[0].ParameterCount)
// 	fmt.Println("DO BODY 0 ARGUMENTS 0: ", tokenList.Do.DoBody[0].Arguments[1])
// 	fmt.Println("DO BODY 0 ARGUMENTS 1: ", tokenList.Do.DoBody[0].Arguments[2])
// 	// fmt.Println("----------------------------------------")
// 	// fmt.Println("VIEW: ", tokenList.View)
// 	// fmt.Println("VIEW BODY: ", tokenList.View.Body)
// 	// fmt.Println("----------------------------------------")

// }

//----------------------------- HELPERS ------------------------------//

/*
*
*
*
 */
func handleStateBlock(value Lexical, currentStep RootToken) RootToken {
	if StringInSlice(value.Value, standardLibraryContents) {
		if len(currentStep.StateBody) < 1 {
			m := make(map[int]*Body)
			m[0] = &Body{Name: value.Value, Type: CallExpression, Kind: "Func"}
			currentStep.StateBody = m
			currentStep.BodyCount = len(currentStep.StateBody)
		} else {
			currentStep.StateBody[len(currentStep.StateBody)] = &Body{Name: value.Value, Type: CallExpression, Kind: "Func"}
			currentStep.BodyCount = len(currentStep.StateBody)
		}
	} else {
		currentBodyCount := VerifyIndex(currentStep.BodyCount - 1)
		argumentLength := len(currentStep.StateBody[currentBodyCount].Arguments)
		if argumentLength < 1 {
			m := make(map[int]*Body)
			m[0] = &Body{Name: value.Value, Type: FunctionDeclaration, Kind: "Func"}
			currentStep.StateBody[currentBodyCount].Arguments = m
			currentStep.StateBody[currentBodyCount].ArgumentCount = argumentLength + 1
		} else {
			currentStep.StateBody[currentBodyCount].Arguments[argumentLength] = &Body{Name: value.Value, Type: VariableDeclaration, Kind: "??"}
			currentStep.StateBody[currentBodyCount].ArgumentCount = argumentLength
		}
	}
	return currentStep
}

/*
*
*
*
 */
func handleDoBlock(value Lexical, currentStep RootToken, pushAsArgument, pushAsParameter bool) RootToken {
	fmt.Println("IN DO: ", value)

	// Everything here should start as Func, as only Functions will be listed in Do
	if !pushAsArgument && !pushAsParameter {
		if len(currentStep.DoBody) < 1 {
			m := make(map[int]*Func)
			m[0] = &Func{Name: value.Value, Type: CallExpression, Kind: "Func"}
			currentStep.DoBody = m
			currentStep.BodyCount = len(currentStep.DoBody)
		} else {
			currentStep.DoBody[len(currentStep.DoBody)] = &Func{Name: value.Value, Type: CallExpression, Kind: "Func"}
			currentStep.BodyCount = len(currentStep.DoBody)
		}

	} else {

		currentBodyCount := VerifyIndex(currentStep.BodyCount - 1)
		argumentLength := len(currentStep.DoBody[currentBodyCount].Arguments)
		parametersLength := len(currentStep.DoBody[currentBodyCount].Parameters)
		if pushAsParameter {
			currentStep.DoBody[currentBodyCount] = addParameterToAST(parametersLength, value.Value, currentStep.DoBody[currentBodyCount])
		} else {
			//Maybe here we check if value.Type == 16
			// If so, then we need a way to denote that the next item to come through is an argument for that SL entity
			currentStep.DoBody[currentBodyCount] = addArgumentToAST(argumentLength, value.Value, currentStep.DoBody[currentBodyCount])
		}

	}
	return currentStep
}

/*
*
*
*
 */
func handleViewBlock(value Lexical, currentStep RootToken, pushAsArgument bool) RootToken {
	fmt.Println("IN VIEW: ", value)
	return currentStep
}

/*
*
*
*
 */
func addArgumentToAST(argumentLength int, value string, body *Func) *Func {
	if argumentLength < 1 {
		m := make(map[int]*Body)
		m[0] = &Body{Name: value, Type: CallExpression, Kind: "Func"}
		body.Arguments = m
		body.ArgumentCount = argumentLength + 1
	} else {
		body.Arguments[argumentLength] = &Body{Name: value, Type: VariableDeclaration, Kind: "??"}
		body.ArgumentCount = argumentLength + 1
	}
	return body
}

/*
*
*
*
 */
func addParameterToAST(parametersLength int, value string, body *Func) *Func {
	if parametersLength < 1 {
		m := make(map[int]*Params)
		m[0] = &Params{Name: value}
		body.Parameters = m
		body.ParameterCount = parametersLength + 1
		fmt.Println("PUSHING AS PARAMETER: ", value)
	} else {
		body.Parameters[parametersLength] = &Params{Name: value}
		body.ParameterCount = parametersLength + 1
	}

	return body
}
