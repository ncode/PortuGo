algoritmo "record layout"
tipo
D = registro
i: inteiro
r: real
s: caractere
b: logico
fimregistro
var
x: D
inicio
escreval(x.i,"|",x.r,"|",x.s,"|",x.b)
x.i <- 7
x.r <- 1.5
x.s <- "ok"
x.b <- verdadeiro
escreval(x.i,"|",x.r,"|",x.s,"|",x.b)
fimalgoritmo
