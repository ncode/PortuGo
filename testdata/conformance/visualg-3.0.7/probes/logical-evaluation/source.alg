algoritmo "logical-evaluation"
var
n: inteiro
funcao right(): logico
inicio
n <- n + 1
retorne verdadeiro
fimfuncao
inicio
escreval("BEGIN:logical-evaluation")
n <- 0
se falso e right() entao
escreval("UNEXPECTED")
fimse
escreval("E=", n)
n <- 0
se verdadeiro ou right() entao
escreval("OU=", n)
fimse
escreval("END:logical-evaluation")
fimalgoritmo
