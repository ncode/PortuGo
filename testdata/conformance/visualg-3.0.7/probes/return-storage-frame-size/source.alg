algoritmo "returns"
var
a, b: inteiro
funcao G(n: inteiro): inteiro
inicio
se n > 0 entao
retorne n
fimse
fimfuncao
funcao F: inteiro
var
x,y,z,w: inteiro
inicio
x <- G(9)
retorne 11
fimfuncao
funcao H: inteiro
inicio
retorne G(0)
fimfuncao
inicio
a <- F()
b <- H()
escreval(a,b)
fimalgoritmo
