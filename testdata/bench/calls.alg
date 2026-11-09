algoritmo "benchmark calls"
funcao Fibonacci(n: inteiro): inteiro
inicio
se n < 2 entao
  retorne n
senao
  retorne Fibonacci(n - 1) + Fibonacci(n - 2)
fimse
fimfuncao
inicio
escreval(Fibonacci(12))
fimalgoritmo
