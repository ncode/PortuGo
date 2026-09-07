algoritmo "record-field"
tipo
  dado = registro
    x: inteiro
  fimregistro
var
  r: dado
inicio
  escreval("BEGIN:record-field")
  r.x <- 3
  escreval(r.x)
  escreval("END:record-field")
fimalgoritmo
