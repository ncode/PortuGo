algoritmo "nested halt"
procedimento stop(x: inteiro)
inicio
escreval("missing:", x)
fimprocedimento
funcao f(): inteiro
inicio
stop()
retorne 7
fimfuncao
inicio
escreval("BEFORE")
escreval("prefix", f(), 1 \ 0)
escreval("AFTER")
fimalgoritmo
