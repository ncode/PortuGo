algoritmo "control forms"
var
i: inteiro
total: inteiro
procedimento stop
inicio
retorne
fimprocedimento
funcao value: inteiro
inicio
retorne 1
fimfuncao
inicio
se verdadeiro entao
escreva(1)
senao
escreva(2)
fimse
enquanto falso faca
interrompa
fimenquanto
repita
total <- total + 1
ate total = 1
escolha total
caso 1 ate 2, 4:
escreva("case")
outrocaso:
escreva("other")
fimescolha
para i de 3 ate 1 passo -1 faca
escreva(i)
fimpara
fimalgoritmo
