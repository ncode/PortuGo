algoritmo "record layout"
tipo
Child = registro
x: inteiro
fimregistro
Parent = registro
y: Child
fimregistro
var
r: Parent
inicio
r.y.x <- 7
escreval(r.y.x)
fimalgoritmo
