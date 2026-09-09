algoritmo "returns"
var
a, b: inteiro
funcao F(n: inteiro): inteiro
var
x: inteiro
inicio
se n = 0 entao
retorne 7
senao
x <- F(n-1)
fimse
fimfuncao
inicio
a <- F(0)
b <- F(2)
escreval(a,b)
fimalgoritmo
