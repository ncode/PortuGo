algoritmo "Limite de texto"
var
  s: caractere
  i: inteiro
inicio
  s <- "a"
  para i de 1 ate 8 faca
    s <- s + s
  fimpara
  escreval(compr(s))
  escreval(copia(s, 250, 20))
fimalgoritmo
