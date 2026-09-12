algoritmo "frames"
funcao F(n: inteiro): inteiro
var
x: inteiro
inicio
x <- n
se n > 0 entao
escreva(F(n-1))
fimse
retorne x
fimfuncao
inicio
escreval(F(3))
fimalgoritmo
