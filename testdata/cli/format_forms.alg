algoritmo "formas"
dos
arquivo "data.txt"
const
  limite = 2
tipo
  item = registro
    valor: inteiro
  fimregistro
  numero = inteiro
var
  dados: vetor[0..limite, 1..2] de item
  x, i: numero
  r: real
procedimento mostrar(var n: inteiro; texto: caractere)
const
  local = 1
var
  y: inteiro
inicio
  escreval(texto, n)
fimprocedimento
funcao resultado(n: inteiro): inteiro
inicio
  retorne
  retorne n
fimfuncao
inicio
  dados[0, 1].valor <- 2
  x <- resultado(dados[0, 1].valor)
  r <- -2 ^ 2 + abs(1.0)
  mostrar(x, "texto")
  mostrar
  pi
  rand
  se (x = 2) e nao falso entao
    escreva("ok", x:3, r:8:2)
  senao
    escreval()
  fimse
  escolha x
    caso 0 ate 2, 4:
      pausa
    outrocaso:
      debug verdadeiro
  fimescolha
  enquanto falso faca
    interrompa
  fimenquanto
  repita
    x <- x + 1
  ate x > 2
  para i de 1 ate 2 passo 1 faca
    leia(x)
  fimpara
  para i de 2 ate 1 faca
  fimpara
  eco
  eco off
  eco on
  cronometro
  cronometro off
  timer 0
  aleatorio on
  aleatorio 7, 7, 2
  aleatorio off
  limpatela
  mudacor("azul", "frente")
  escreval(dos, pausa, timer, debug, aleatorio, cronometro, limpatela, mudacor)
fimalgoritmo
