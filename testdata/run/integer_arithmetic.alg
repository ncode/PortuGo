algoritmo "integer arithmetic"
var
  value: inteiro
inicio
  value <- 2147483647 + 1
  escreval(value)
  escreval(value * 2)
  escreval(abs(value), "|", int(2147483648.0))
fimalgoritmo
