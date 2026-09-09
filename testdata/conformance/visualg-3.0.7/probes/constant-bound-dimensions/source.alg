algoritmo "constant bounds"
const
lo = -2
hi = 0
cols = 3
var
v: vetor[lo..hi,1..cols] de inteiro
inicio
v[-2,1] <- 7
v[0,3] <- 9
escreval(v[-2,1],"|",v[0,3])
fimalgoritmo
