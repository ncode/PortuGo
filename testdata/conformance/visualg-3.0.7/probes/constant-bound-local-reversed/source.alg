algoritmo "constant bounds"
funcao value(x: inteiro): inteiro
const
n = x
var
v: vetor[0..n] de inteiro
inicio
v[n] <- 7
retorne v[n]
fimfuncao
inicio
escreval(value(1))
escreval(value(-1))
fimalgoritmo
