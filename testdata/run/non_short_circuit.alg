algoritmo "non_short_circuit"
var
  x: inteiro

funcao marca(): logico
inicio
  x <- x + 1
  retorne verdadeiro
fimfuncao

inicio
  x <- 0
  se falso e marca() entao
    x <- 99
  fimse
  escreval(x)
fimalgoritmo
