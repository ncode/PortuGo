algoritmo "fatorial"
var
  n, i, acc: inteiro
inicio
  leia(n)
  acc <- 1
  para i de 2 ate n faca
    acc <- acc * i
  fimpara
  escreval(acc)
fimalgoritmo
