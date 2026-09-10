algoritmo "type alias callable"
tipo
T = inteiro
var
x: T
procedimento P(var n: inteiro)
inicio
n <- n + 1
fimprocedimento
funcao F(n: inteiro): inteiro
inicio
retorne n + 1
fimfuncao
inicio
x <- 7
P(x)
escreval(F(x))
fimalgoritmo
