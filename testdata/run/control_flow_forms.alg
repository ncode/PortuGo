algoritmo "control coverage"
var
  i: inteiro
  total: inteiro
inicio
  total <- 0
  repita
    total <- total + 1
  ate total = 3
  escreval("R=", total)
  escolha total
    caso 1 ate 3:
      escreval("C=first")
    caso 3 ate 5:
      escreval("C=second")
    outrocaso:
      escreval("C=none")
  fimescolha
  para i de 3 ate 1 passo -1 faca
    escreva(i)
  fimpara
  escreval("")
  escreval("F=", i)
fimalgoritmo
