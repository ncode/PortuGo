algoritmo "record reference location"
tipo
D = registro
x: inteiro
fimregistro
var
r: D
procedimento P(var n: inteiro)
inicio
r <- r
n <- 7
fimprocedimento
inicio
r.x <- 1
P(r.x)
escreval(r.x)
fimalgoritmo
