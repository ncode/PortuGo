algoritmo "for-edges"
var
  i: inteiro
inicio
  escreval("BEGIN:for-edges")
  para i de 5 ate 1 passo -2 faca
    i <- 100
  fimpara
  escreval("DOWN=[", i, "]")
  para i de 1 ate 6 passo 2 faca
    i <- 100
  fimpara
  escreval("STEP=[", i, "]")
  i <- 42
  para i de 3 ate 1 faca
    escreval("BAD")
  fimpara
  escreval("EMPTY=[", i, "]")
  para i de 1 ate 3 faca
    i <- 100
    escreval("BODY=[", i, "]")
    interrompa
  fimpara
  escreval("BREAK=[", i, "]")
  escreval("END:for-edges")
fimalgoritmo
