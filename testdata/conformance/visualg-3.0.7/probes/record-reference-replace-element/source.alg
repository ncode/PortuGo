algoritmo "record reference location"
tipo
D = registro
x: inteiro
fimregistro
var
v: vetor[1..2] de D
procedimento P(var n: inteiro)
inicio
v[1] <- v[2]
n <- 7
fimprocedimento
inicio
v[1].x <- 1
v[2].x <- 2
P(v[1].x)
escreval(v[1].x)
fimalgoritmo
