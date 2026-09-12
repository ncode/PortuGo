algoritmo "empty call qualification"
var
marker: inteiro
callerValue: inteiro
procedimento P(value: inteiro)
inicio
escreval("BODY")
marker <- marker + 1
escreval("VALUE=<", value, ">")
value <- 9
escreval("SET=<", value, ">")
fimprocedimento
inicio
marker <- 0
callerValue <- 7
escreval("BEFORE")
P(callerValue)
escreval("AFTER=", marker, " CALLER=<", callerValue, ">")
fimalgoritmo
