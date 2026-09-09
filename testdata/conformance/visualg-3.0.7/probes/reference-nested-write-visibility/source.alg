algoritmo "aliasing"
var
n: inteiro
procedimento Q(var y: inteiro)
inicio
y <- 7
escreval(n, y)
fimprocedimento
procedimento P(var x: inteiro)
inicio
x <- 5
Q(x)
escreval(n, x)
fimprocedimento
inicio
n <- 3
P(n)
escreval(n)
fimalgoritmo
