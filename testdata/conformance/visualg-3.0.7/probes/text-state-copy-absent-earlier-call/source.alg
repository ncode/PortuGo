algoritmo "text call state"
var
  value: caractere
funcao Texto(idx: inteiro): caractere
inicio
  escreval(idx)
  retorne "abc"
fimfuncao
funcao Silencioso: caractere
inicio
  retorne "abc"
fimfuncao
funcao Numero(idx: inteiro): real
inicio
  escreval(idx)
  retorne idx + 0.9
fimfuncao
inicio
  escreval(copia(Texto(1),abs(),3))
  escreval("DONE")
fimalgoritmo
