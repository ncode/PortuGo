algoritmo "record boundary"
tipo
D = registro
x: inteiro
fimregistro
var
r: D
procedimento P(var n: inteiro)
inicio
n <- n + 1
fimprocedimento
inicio
r.x <- 7
P(r.x)
escreval(r.x)
fimalgoritmo
