algoritmo "empty call qualification"
var
marker: inteiro
callerValue: logico
procedimento P(value: logico)
inicio
escreval("BODY")
marker <- marker + 1
escreval("VALUE=<", value, ">")
value <- falso
escreval("SET=<", value, ">")
fimprocedimento
inicio
marker <- 0
callerValue <- verdadeiro
escreval("BEFORE")
P(callerValue)
escreval("AFTER=", marker, " CALLER=<", callerValue, ">")
fimalgoritmo
