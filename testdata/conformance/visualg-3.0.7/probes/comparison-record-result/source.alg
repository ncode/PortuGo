algoritmo "comparison values"
tipo
D = registro
x: inteiro
fimregistro
var
a, b, c: D
inicio
a.x <- 2
b.x <- 7
c <- a = b
escreval(c.x)
c <- a <> b
escreval(c.x)
c <- a < b
escreval(c.x)
c <- a > b
escreval(c.x)
c <- a <= b
escreval(c.x)
c <- a >= b
escreval(c.x)
fimalgoritmo
