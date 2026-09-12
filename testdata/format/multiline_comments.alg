algoritmo "comments"
const
  n = 1 // first constant operand
    + 2 // second constant operand
tipo
  item = // before record keyword
    registro // record header
    value: // before field type
      inteiro // field type
  fimregistro
var
  i: inteiro
  v: vetor[1.. // lower bound
    2] de inteiro // upper bound

procedimento P(a: inteiro; // first parameter
  b: inteiro) // second parameter
inicio
fimprocedimento
inicio
  escreval(1 // first operand
    + 2 // second operand
    )
  P(1, // first argument
    2) // second argument
  se verdadeiro // condition
    entao // then
  fimse
  escolha 1
    caso 1 + // first label operand
      1: // second label operand
  fimescolha
  enquanto falso // condition
    faca // while header
  fimenquanto
  repita
  ate verdadeiro // first condition operand
    ou falso // second condition operand
  para i de 1 // first bound
    ate 0 faca // second bound
  fimpara
fimalgoritmo
