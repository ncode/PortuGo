algoritmo "record field type"
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
escreval(r.y)
fimalgoritmo
