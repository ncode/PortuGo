algoritmo "call control"
var
s,t:caractere
funcao Prefixo(QTDE:inteiro):caractere
inicio
retorne "0000"
fimfuncao
funcao Converter(STRING:caractere):caractere
inicio
retorne STRING
fimfuncao
funcao Direita(STRING:caractere;QTDE:inteiro):caractere
inicio
retorne copia(STRING,compr(STRING)-QTDE+1,QTDE)
fimfuncao
inicio
s<-Direita(Prefixo(4)+Converter("1"),4)
escreval(s)
fimalgoritmo
