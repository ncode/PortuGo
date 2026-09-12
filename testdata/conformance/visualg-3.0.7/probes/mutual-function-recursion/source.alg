algoritmo "recursion"
funcao F(n: inteiro): inteiro
inicio
se n = 0 entao
retorne 0
senao
retorne G(n-1) + 1
fimse
fimfuncao
funcao G(n: inteiro): inteiro
inicio
se n = 0 entao
retorne 0
senao
retorne F(n-1) + 1
fimse
fimfuncao
inicio
escreval(F(3))
fimalgoritmo
