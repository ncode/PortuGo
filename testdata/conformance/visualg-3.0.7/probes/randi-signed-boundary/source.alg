algoritmo "random details"
var
x, i: inteiro
ok: logico
inicio
ok <- verdadeiro
para i de 1 ate 64 faca
x <- randi(2147483648)
ok <- ok e (x >= 0) e (x < 2147483648)
fimpara
escreval(ok)
fimalgoritmo
