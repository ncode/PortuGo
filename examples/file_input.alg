// Execute from the repository root; the supplied data file remains unchanged.
algoritmo "Entrada de arquivo"
arquivo "testdata/run/file_input.txt"
var
  primeiro, segundo: inteiro
inicio
  leia(primeiro, segundo)
  escreval("Soma:", primeiro + segundo)
fimalgoritmo
