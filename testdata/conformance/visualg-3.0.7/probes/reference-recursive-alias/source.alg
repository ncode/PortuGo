algoritmo "aliasing"
var
n: inteiro
procedimento P(var x: inteiro; depth: inteiro)
inicio
x <- x+1
se depth > 0 entao
P(x, depth-1)
fimse
escreval(n, x)
fimprocedimento
inicio
P(n, 2)
escreval(n)
fimalgoritmo
