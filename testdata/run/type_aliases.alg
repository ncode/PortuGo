algoritmo "aliases escalares"
tipo
  Numero = inteiro
  Contagem = Numero
var
  valores: vetor[1..2] de Contagem
procedimento incrementar(var n: inteiro)
inicio
  n <- n + 1
fimprocedimento
funcao calcular: inteiro
tipo
  Local = Contagem
var
  n: Local
inicio
  n <- 7
  incrementar(n)
  retorne n
fimfuncao
inicio
  valores[1] <- calcular
  escreval(valores[1], valores[2])
fimalgoritmo
