algoritmo "random details"
var
x, i: inteiro
ok: logico
inicio
ok <- verdadeiro
para i de 1 ate 64 faca
x <- randi(-7)
ok <- ok e ((x <= -8) ou (x >= 0))
fimpara
escreval(ok)
fimalgoritmo
