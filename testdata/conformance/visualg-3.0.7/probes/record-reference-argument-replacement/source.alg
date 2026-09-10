algoritmo "record reference location"
tipo
D = registro
x: inteiro
fimregistro
var
first, second: D
funcao Reset: inteiro
inicio
first <- second
retorne 0
fimfuncao
procedimento P(var n: inteiro; unused: inteiro)
inicio
n <- n + 5
fimprocedimento
inicio
first.x <- 1
second.x <- 2
P(first.x, Reset)
escreval(first.x)
fimalgoritmo
