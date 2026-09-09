algoritmo "random builtins"
var
x, i: inteiro
ok: logico
inicio
ok <- verdadeiro
para i de 1 ate 64 faca
x <- randi(7)
ok <- ok e x >= 0 e x < 7
fimpara
escreval(ok)
fimalgoritmo
