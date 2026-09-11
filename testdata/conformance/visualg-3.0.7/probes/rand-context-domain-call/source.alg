algoritmo "random contexts"
var
x: real
i: inteiro
ok: logico
inicio
ok <- verdadeiro
para i de 1 ate 64 faca
x <- rand()
ok <- ok e (x >= 0) e (x < 1)
fimpara
escreval(ok)
fimalgoritmo
