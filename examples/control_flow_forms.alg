algoritmo "control flow"
var
  i: inteiro
  total: inteiro
inicio
  total <- 0
  repita
    total <- total + 1
  ate total = 3
  escolha total
    caso 1 ate 3:
      escreval("first range")
    caso 3 ate 5:
      escreval("second range")
    outrocaso:
      escreval("no range")
  fimescolha
  para i de 3 ate 1 passo -1 faca
    escreva(i)
  fimpara
  escreval("")
fimalgoritmo
