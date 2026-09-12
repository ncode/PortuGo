algoritmo "text details"
var
i: inteiro
s: caractere
inicio
s <- " !#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_`abcdefghijklmnopqrstuvwxyz{|}~€‚ƒ„…†‡ˆ‰Š‹Œ‘’“”•–—˜™š›œŸ ¡¢£¤¥¦§¨©ª«¬­®¯°±²³´µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏĞÑÒÓÔÕÖ×ØÙÚÛÜİŞßàáâãäåæçèéêëìíîïğñòóôõö÷øùúûüışÿ"
para i de 1 ate compr(s) faca
escreval(asc(copia(s,i,1)),"|",asc(maiusc(copia(s,i,1))),"|",asc(minusc(copia(s,i,1))))
fimpara
fimalgoritmo
