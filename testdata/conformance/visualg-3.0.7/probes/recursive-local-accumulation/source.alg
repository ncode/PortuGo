algoritmo "returns"
var
a, b, c: inteiro
funcao F(n: inteiro): inteiro
var
x: inteiro
inicio
x <- n
se n > 0 entao
x <- x + F(n-1)
fimse
retorne x
fimfuncao
inicio
a <- F(3)
escreval(a)
fimalgoritmo
