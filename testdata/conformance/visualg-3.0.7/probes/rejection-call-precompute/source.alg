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
t<-Converter("1")
s<-Direita(Prefixo(4)+t,4)
escreval(s)
fimalgoritmo
