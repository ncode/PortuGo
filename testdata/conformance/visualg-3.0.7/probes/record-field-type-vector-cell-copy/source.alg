algoritmo "record field type"
tipo
D = registro
x: inteiro
fimregistro
var
v: vetor[1..2] de D
inicio
v[1].x <- 7
v[2] <- v[1]
v[1].x <- 8
escreval(v[1].x,v[2].x)
fimalgoritmo
