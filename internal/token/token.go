package token

import (
	"strconv"
	"strings"
)

// Kind identifies a lexical token.
type Kind int

const (
	ILLEGAL Kind = iota
	EOF
	IDENT
	NUMBER
	STRING

	ADD
	SUB
	MUL
	QUO
	IDIV
	REM
	POW
	ASSIGN
	EQL
	NEQ
	LSS
	GTR
	LEQ
	GEQ
	LPAREN
	RPAREN
	LBRACK
	RBRACK
	COMMA
	COLON
	SEMI
	DOT
	DOTDOT

	ALGORITMO
	DOS
	RAND
	CONST
	TIPO
	REGISTRO
	FIMREGISTRO
	VAR
	INICIO
	FIMALGORITMO
	INTEIRO
	REAL
	CARACTERE
	LOGICO
	VETOR
	DE
	PROCEDIMENTO
	FIMPROCEDIMENTO
	FUNCAO
	FIMFUNCAO
	RETORNE
	SE
	ENTAO
	SENAO
	FIMSE
	ESCOLHA
	CASO
	OUTROCASO
	FIMESCOLHA
	ENQUANTO
	FACA
	FIMENQUANTO
	REPITA
	ATE
	PARA
	PASSO
	FIMPARA
	INTERROMPA
	LEIA
	ESCREVA
	ESCREVAL
	LIMPATELA
	MUDACOR
	VERDADEIRO
	FALSO
	E
	OU
	NAO
	XOU
	MOD
	NEWLINE
)

// Token is one item in the source stream.
type Token struct {
	Kind Kind
	Text string
	Pos  Pos
}

var keywords = map[string]Kind{
	"algoritmo":       ALGORITMO,
	"dos":             DOS,
	"rand":            RAND,
	"const":           CONST,
	"tipo":            TIPO,
	"registro":        REGISTRO,
	"fimregistro":     FIMREGISTRO,
	"var":             VAR,
	"inicio":          INICIO,
	"fimalgoritmo":    FIMALGORITMO,
	"inteiro":         INTEIRO,
	"real":            REAL,
	"caractere":       CARACTERE,
	"caracter":        CARACTERE,
	"logico":          LOGICO,
	"vetor":           VETOR,
	"de":              DE,
	"procedimento":    PROCEDIMENTO,
	"fimprocedimento": FIMPROCEDIMENTO,
	"funcao":          FUNCAO,
	"função":          FUNCAO,
	"fimfuncao":       FIMFUNCAO,
	"fimfunção":       FIMFUNCAO,
	"retorne":         RETORNE,
	"se":              SE,
	"entao":           ENTAO,
	"então":           ENTAO,
	"senao":           SENAO,
	"senão":           SENAO,
	"fimse":           FIMSE,
	"escolha":         ESCOLHA,
	"caso":            CASO,
	"outrocaso":       OUTROCASO,
	"fimescolha":      FIMESCOLHA,
	"enquanto":        ENQUANTO,
	"faca":            FACA,
	"faça":            FACA,
	"fimenquanto":     FIMENQUANTO,
	"repita":          REPITA,
	"ate":             ATE,
	"até":             ATE,
	"para":            PARA,
	"passo":           PASSO,
	"fimpara":         FIMPARA,
	"interrompa":      INTERROMPA,
	"leia":            LEIA,
	"escreva":         ESCREVA,
	"escreval":        ESCREVAL,
	"limpatela":       LIMPATELA,
	"mudacor":         MUDACOR,
	"verdadeiro":      VERDADEIRO,
	"falso":           FALSO,
	"e":               E,
	"ou":              OU,
	"nao":             NAO,
	"não":             NAO,
	"xou":             XOU,
	"mod":             MOD,
	"div":             IDIV,
}

var kindNames = map[Kind]string{
	ILLEGAL: "ILLEGAL", EOF: "EOF", IDENT: "IDENT", NUMBER: "NUMBER", STRING: "STRING",
	ADD: "+", SUB: "-", MUL: "*", QUO: "/", IDIV: "\\", REM: "%", POW: "^",
	ASSIGN: "<-", EQL: "=", NEQ: "<>", LSS: "<", GTR: ">", LEQ: "<=", GEQ: ">=",
	LPAREN: "(", RPAREN: ")", LBRACK: "[", RBRACK: "]", COMMA: ",", COLON: ":",
	SEMI: ";", DOT: ".", DOTDOT: "..",
	REGISTRO: "registro", FIMREGISTRO: "fimregistro",
	ALGORITMO: "algoritmo", CONST: "const", TIPO: "tipo", VAR: "var", INICIO: "inicio", FIMALGORITMO: "fimalgoritmo",
	INTEIRO: "inteiro", REAL: "real", CARACTERE: "caractere", LOGICO: "logico",
	VETOR: "vetor", DE: "de", PROCEDIMENTO: "procedimento", FIMPROCEDIMENTO: "fimprocedimento",
	FUNCAO: "funcao", FIMFUNCAO: "fimfuncao", RETORNE: "retorne", SE: "se",
	ENTAO: "entao", SENAO: "senao", FIMSE: "fimse", ESCOLHA: "escolha", CASO: "caso",
	OUTROCASO: "outrocaso", FIMESCOLHA: "fimescolha", ENQUANTO: "enquanto", FACA: "faca",
	FIMENQUANTO: "fimenquanto", REPITA: "repita", ATE: "ate", PARA: "para", PASSO: "passo",
	FIMPARA: "fimpara", INTERROMPA: "interrompa", LEIA: "leia", ESCREVA: "escreva",
	ESCREVAL: "escreval", VERDADEIRO: "verdadeiro", FALSO: "falso", E: "e", OU: "ou",
	LIMPATELA: "limpatela", MUDACOR: "mudacor", DOS: "dos", RAND: "rand",
	NAO: "nao", XOU: "xou", MOD: "mod",
	NEWLINE: "NEWLINE",
}

// Lookup returns the keyword kind for ident, or IDENT.
func Lookup(ident string) Kind {
	if kind, ok := keywords[strings.ToLower(ident)]; ok {
		return kind
	}
	return IDENT
}

// String returns the printable name for a token kind.
func (k Kind) String() string {
	if name, ok := kindNames[k]; ok {
		return name
	}
	return "token(" + strconv.Itoa(int(k)) + ")"
}

// BinaryPrecedence returns the binding strength of a binary operator, or -1.
func (k Kind) BinaryPrecedence() int {
	switch k {
	case EQL, NEQ, LSS, GTR, LEQ, GEQ:
		return 4
	case ADD, SUB, OU, XOU:
		return 5
	case MUL, QUO, IDIV, REM, MOD, E:
		return 6
	case POW:
		return 8
	default:
		return -1
	}
}

// UnaryPrecedence returns the binding strength of a unary operator, or -1.
func (k Kind) UnaryPrecedence() int {
	switch k {
	case NAO:
		return 7
	case ADD, SUB:
		return 9 // Unary signs bind more tightly than power.
	default:
		return -1
	}
}
