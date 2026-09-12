algoritmo "empty call qualification"
var
marker: inteiro
procedimento P
inicio
escreval("BODY")
marker <- marker + 1
fimprocedimento
inicio
marker <- 0
escreval("BEFORE")
P()
escreval("AFTER=", marker)
fimalgoritmo
