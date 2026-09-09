algoritmo "Texto e codigos"
var
  texto: caractere
inicio
  texto <- "ação"
  escreval(maiusc(texto))
  escreval(copia(texto, 2.9, 2.5))
  escreval("Tamanho: " + numpcarac(compr(texto)))
  escreval("Posicao: " + numpcarac(pos("ç", texto)))
  escreval("Codigo: " + numpcarac(asc("€")))
  escreval(carac(128))
fimalgoritmo
