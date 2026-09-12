algoritmo "empty call qualification"
var
marker: inteiro
callerValue: caractere
procedimento P(value: caractere)
inicio
escreval("BODY")
marker <- marker + 1
escreval("VALUE=<", value, ">")
value <- "set"
escreval("SET=<", value, ">")
fimprocedimento
inicio
marker <- 0
callerValue <- "seed"
escreval("BEFORE")
P(callerValue)
escreval("AFTER=", marker, " CALLER=<", callerValue, ">")
fimalgoritmo
