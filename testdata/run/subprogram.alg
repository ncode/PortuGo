algoritmo "subprogram"
var
  x: inteiro

procedimento inc(var n: inteiro)
inicio
  n <- n + 1
fimprocedimento

funcao dobro(n: inteiro): inteiro
inicio
  retorne n * 2
fimfuncao

inicio
  x <- 10
  inc(x)
  escreval(dobro(x))
fimalgoritmo
