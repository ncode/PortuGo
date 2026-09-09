algoritmo "random domain"
var
  i, value: inteiro
  valid: logico
inicio
  valid <- verdadeiro
  para i de 1 ate 64 faca
    value <- randi(7)
    valid <- valid e (value >= 0) e (value < 7)
  fimpara
  escreval(valid)
fimalgoritmo
