algoritmo "global var after type subprogram"
tipo
  tdado = registro
    codigo: inteiro
  fimregistro
procedimento incrementar
inicio
  dados.codigo <- dados.codigo + 1
fimprocedimento
var
  dados: tdado
inicio
  dados.codigo <- 4
  incrementar()
  escreval(dados.codigo)
fimalgoritmo
