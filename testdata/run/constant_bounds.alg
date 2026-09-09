algoritmo "Limites constantes"
const
  primeiro = -2
var
funcao ultimo(tamanho: inteiro): inteiro
const
  limite = tamanho - 1
var
  valores: vetor[primeiro..limite] de inteiro
inicio
  valores[limite] <- tamanho
  retorne valores[limite]
fimfuncao
inicio
  escreval(ultimo(3), "|", ultimo(5))
fimalgoritmo
