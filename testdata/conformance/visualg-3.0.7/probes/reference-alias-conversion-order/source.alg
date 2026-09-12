algoritmo "aliasing"
var
n: inteiro
procedimento P(var x: real; var y: inteiro)
inicio
x <- 2.5
escreval(x, y, n)
y <- 7
escreval(x, y, n)
fimprocedimento
inicio
n <- 3
P(n, n)
escreval(n)
fimalgoritmo
