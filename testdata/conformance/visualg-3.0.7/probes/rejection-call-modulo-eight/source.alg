algoritmo "call control"
var
s,t:caractere
funcao Prefixo(QTDE:inteiro):caractere
inicio
retorne "0000"
fimfuncao
funcao Converter(STRING:caractere):inteiro
var
n:inteiro
inicio
n<-8 mod 3
retorne 2
fimfuncao
funcao Direita(STRING:caractere;QTDE:inteiro):caractere
inicio
retorne copia(STRING,compr(STRING)-QTDE+1,QTDE)
fimfuncao
inicio
escreval(10+Converter("1"))
fimalgoritmo
