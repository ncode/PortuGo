algoritmo "benchmark vectors"
var
v: vetor[1..256] de inteiro
i, total: inteiro
inicio
para i de 1 ate 256 faca
  v[i] <- i
fimpara
para i de 1 ate 256 faca
  total <- total + v[i]
fimpara
escreval(total)
fimalgoritmo
