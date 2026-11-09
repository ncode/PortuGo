algoritmo "benchmark strings"
var
texto: caractere
i, total: inteiro
inicio
texto <- "abcdefghijklmnop"
para i de 1 ate 200 faca
  texto <- maiusc(copia(texto, 2, 15) + copia(texto, 1, 1))
  total <- total + compr(texto)
fimpara
escreval(total)
fimalgoritmo
