algoritmo "call control"
var
s,t:caractere
funcao Prefixo(QTDE:inteiro):caractere
inicio
retorne "0000"
fimfuncao
funcao Converter(STRING:caractere):caractere
var
n:inteiro
inicio
n<-1 mod 8
retorne STRING
fimfuncao
funcao Direita(STRING:caractere;QTDE:inteiro):caractere
inicio
retorne copia(STRING,compr(STRING)-QTDE+1,QTDE)
fimfuncao
inicio
s<-Direita(Converter("1")+Prefixo(4),4)
escreval(s)
fimalgoritmo
