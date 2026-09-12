algoritmo "constant bounds"
funcao value(x: inteiro): inteiro
const
n = x
var
v: vetor[1..n] de inteiro
inicio
v[n] <- x
retorne v[n]
fimfuncao
inicio
escreval(value(3),"|",value(5))
fimalgoritmo
