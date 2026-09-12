algoritmo "cadastro"
tipo
  Pessoa = registro
    nome: caractere
    idade: inteiro
  fimregistro
var
  pessoas: vetor[1..2] de Pessoa
procedimento aniversario(var anos: inteiro)
inicio
  anos <- anos + 1
fimprocedimento
inicio
  pessoas[1].nome <- "Ana"
  pessoas[1].idade <- 20
  aniversario(pessoas[1].idade)
  pessoas[2] <- pessoas[1]
  pessoas[1].nome <- "Bia"
  escreval(pessoas[1].nome, ":", pessoas[1].idade)
  escreval(pessoas[2].nome, ":", pessoas[2].idade)
fimalgoritmo
