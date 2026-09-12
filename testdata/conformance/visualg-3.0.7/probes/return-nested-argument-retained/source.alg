algoritmo "returns"
var
a, b: inteiro
funcao G(n: inteiro): inteiro
inicio
retorne n
fimfuncao
funcao F(n: inteiro): inteiro
inicio
se n > 0 entao
retorne n
fimse
fimfuncao
inicio
a <- F(3)
b <- F(G(0))
escreval(a,b)
fimalgoritmo
