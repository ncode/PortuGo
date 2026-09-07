algoritmo "for-break-mutation"
var
  i: inteiro
inicio
  escreval("BEGIN:for-break-mutation")
  para i de 1 ate 3 faca
    i <- -100
    interrompa
  fimpara
  escreval("LOW=[", i, "]")
  para i de 5 ate 1 passo -2 faca
    i <- 100
    interrompa
  fimpara
  escreval("DOWNHIGH=[", i, "]")
  para i de 5 ate 1 passo -2 faca
    i <- -100
    interrompa
  fimpara
  escreval("DOWNLOW=[", i, "]")
  i <- 42
  para i de 1 ate 3 passo -1 faca
    escreval("BAD")
  fimpara
  escreval("EMPTYDOWN=[", i, "]")
  escreval("END:for-break-mutation")
fimalgoritmo
