algoritmo "record reference location"
tipo
D = registro
x: inteiro
fimregistro
var
first, second: D
procedimento P(var n: inteiro)
inicio
first <- second
n <- 7
fimprocedimento
inicio
first.x <- 1
second.x <- 2
P(first.x)
escreval(first.x)
fimalgoritmo
