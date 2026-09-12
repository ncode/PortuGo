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
  value <- carac(abs())
  escreval("BEFORE",value,Texto(9))
  escreval("DONE")
fimalgoritmo
