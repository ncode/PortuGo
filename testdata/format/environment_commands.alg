algoritmo "environment comments"
dos // console directive
arquivo "data.txt" // file directive
var
  value: inteiro
inicio
  eco // unspecified echo
  eco on // enable echo
  eco off // disable echo
  cronometro on // start clock
  cronometro off // stop clock
  timer 0 // no delay
  timer 1 + // first delay operand
    0 // second delay operand
  pausa // unconditional pause
  debug verdadeiro // conditional pause
  debug verdadeiro // first condition operand
    e falso // second condition operand
  aleatorio on // default bounds
  aleatorio off // console input
  aleatorio 1 // lower bound
  aleatorio 1, 3 // both bounds
  aleatorio 1 + // first bound operand
    0, 3, 2 // remaining operands
  limpatela // clear screen
  mudacor("azul", "frente") // foreground color
  rand // random statement
  leia(value) // input
  escreva(value) // output without newline
  escreval(value) // output with newline
  escreval(dos, aleatorio, timer, pausa, debug, cronometro, limpatela, mudacor) // absent values
  escreval(timer( // ignored argument
    1)) // absent call value
fimalgoritmo // terminator
// opaque trailing note
