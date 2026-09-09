algoritmo "Faixas em escolha"
var
  valor, menor, maior: real
inicio
  valor <- 2.5
  menor <- 1
  maior <- 3
  escolha valor
    caso menor ate maior, 5:
      escreval("faixa")
    caso 2 ate 4:
      escreval("segunda faixa")
    outrocaso:
      escreval("fora")
  fimescolha
  escreval("fim")
fimalgoritmo
