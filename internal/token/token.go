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
	DOTDOT

	ALGORITMO
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
	VERDADEIRO
	FALSO
	E
	OU
	NAO
	XOU
	MOD
)

// Token is one item in the source stream.
type Token struct {
	Kind Kind
	Text string
	Pos  Pos
}

var keywords = map[string]Kind{
	"algoritmo":       ALGORITMO,
	"var":             VAR,
	"inicio":          INICIO,
	"fimalgoritmo":    FIMALGORITMO,
	"inteiro":         INTEIRO,
	"real":            REAL,
	"caractere":       CARACTERE,
	"logico":          LOGICO,
	"vetor":           VETOR,
	"de":              DE,
	"procedimento":    PROCEDIMENTO,
	"fimprocedimento": FIMPROCEDIMENTO,
	"funcao":          FUNCAO,
	"fimfuncao":       FIMFUNCAO,
	"retorne":         RETORNE,
	"se":              SE,
	"entao":           ENTAO,
	"senao":           SENAO,
	"fimse":           FIMSE,
	"escolha":         ESCOLHA,
	"caso":            CASO,
	"outrocaso":       OUTROCASO,
	"fimescolha":      FIMESCOLHA,
	"enquanto":        ENQUANTO,
	"faca":            FACA,
	"fimenquanto":     FIMENQUANTO,
	"repita":          REPITA,
	"ate":             ATE,
	"para":            PARA,
	"passo":           PASSO,
	"fimpara":         FIMPARA,
	"interrompa":      INTERROMPA,
	"leia":            LEIA,
	"escreva":         ESCREVA,
	"escreval":        ESCREVAL,
	"verdadeiro":      VERDADEIRO,
	"falso":           FALSO,
	"e":               E,
	"ou":              OU,
	"nao":             NAO,
	"xou":             XOU,
	"mod":             MOD,
}

var kindNames = map[Kind]string{
	ILLEGAL: "ILLEGAL", EOF: "EOF", IDENT: "IDENT", NUMBER: "NUMBER", STRING: "STRING",
	ADD: "+", SUB: "-", MUL: "*", QUO: "/", IDIV: "\\", REM: "%", POW: "^",
	ASSIGN: "<-", EQL: "=", NEQ: "<>", LSS: "<", GTR: ">", LEQ: "<=", GEQ: ">=",
	LPAREN: "(", RPAREN: ")", LBRACK: "[", RBRACK: "]", COMMA: ",", COLON: ":",
	SEMI: ";", DOTDOT: "..",
	ALGORITMO: "algoritmo", VAR: "var", INICIO: "inicio", FIMALGORITMO: "fimalgoritmo",
	INTEIRO: "inteiro", REAL: "real", CARACTERE: "caractere", LOGICO: "logico",
	VETOR: "vetor", DE: "de", PROCEDIMENTO: "procedimento", FIMPROCEDIMENTO: "fimprocedimento",
	FUNCAO: "funcao", FIMFUNCAO: "fimfuncao", RETORNE: "retorne", SE: "se",
	ENTAO: "entao", SENAO: "senao", FIMSE: "fimse", ESCOLHA: "escolha", CASO: "caso",
	OUTROCASO: "outrocaso", FIMESCOLHA: "fimescolha", ENQUANTO: "enquanto", FACA: "faca",
	FIMENQUANTO: "fimenquanto", REPITA: "repita", ATE: "ate", PARA: "para", PASSO: "passo",
	FIMPARA: "fimpara", INTERROMPA: "interrompa", LEIA: "leia", ESCREVA: "escreva",
	ESCREVAL: "escreval", VERDADEIRO: "verdadeiro", FALSO: "falso", E: "e", OU: "ou",
	NAO: "nao", XOU: "xou", MOD: "mod",
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
