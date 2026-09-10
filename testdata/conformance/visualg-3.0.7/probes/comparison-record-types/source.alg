algoritmo "comparison values"
tipo
D = registro
x: inteiro
fimregistro
T = registro
y: inteiro
fimregistro
var
a: D
b, c: T
inicio
a.x <- 2
b.y <- 7
c <- a = b
escreval(c.y)
c <- a <> b
escreval(c.y)
c <- a < b
escreval(c.y)
c <- a > b
escreval(c.y)
c <- a <= b
escreval(c.y)
c <- a >= b
escreval(c.y)
fimalgoritmo
