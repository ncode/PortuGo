algoritmo "empty call qualification"
var
marker: inteiro
callerValue: real
procedimento P(value: real)
inicio
escreval("BODY")
marker <- marker + 1
escreval("VALUE=<", value, ">")
value <- 9.5
escreval("SET=<", value, ">")
fimprocedimento
inicio
marker <- 0
callerValue <- 2.5
escreval("BEFORE")
P()
escreval("AFTER=", marker, " CALLER=<", callerValue, ">")
fimalgoritmo
