algoritmo "record layout"
tipo
D = registro
x: inteiro
fimregistro
E = D
var
a: D
b: E
inicio
a.x <- 7
b <- a
a.x <- 8
escreval(a.x,b.x)
fimalgoritmo
