algoritmo "Domínios aleatórios"
var
  fracao: real
  inteiro_sorteado, tentativa: inteiro
  valido: logico
inicio
  valido <- verdadeiro
  para tentativa de 1 ate 32 faca
    fracao <- rand
    inteiro_sorteado <- randi(7)
    valido <- valido e (fracao >= 0) e (fracao < 1)
    valido <- valido e (inteiro_sorteado >= 0) e (inteiro_sorteado < 7)
  fimpara
  escreval(valido)
fimalgoritmo
