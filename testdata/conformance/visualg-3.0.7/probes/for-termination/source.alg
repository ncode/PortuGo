algoritmo "for-termination"
var
  i: inteiro
inicio
  escreval("BEGIN:for-termination")
  para i de 5 ate 1 passo -2 faca
    escreva("D", i)
  fimpara
  escreval(";FINAL=", i)
  para i de 1 ate 6 passo 2 faca
    escreva("S", i)
  fimpara
  escreval(";FINAL=", i)
  para i de 1 ate 7 passo 3 faca
    escreva("T", i)
  fimpara
  escreval(";FINAL=", i)
  para i de 3 ate 1 passo -1 faca
    escreva("U", i)
  fimpara
  escreval(";FINAL=", i)
  para i de 1 ate 3 faca
    escreva("B", i)
    interrompa
  fimpara
  escreval(";FINAL=", i)
  escreval("END:for-termination")
fimalgoritmo
