# Especificação consolidada da linguagem VisuAlg 3.x (dialeto Portugol)

**Propósito:** referência para implementação de lexer/parser/interpretador/compilador.
**Data de compilação:** 2026-06-10.

## 0. Metodologia e confiabilidade

Não existe especificação formal oficial do VisuAlg. O código-fonte Delphi nunca foi liberado publicamente. Este documento consolida:

| Fonte | Sigla | Confiabilidade |
|---|---|---|
| Help oficial do VisuAlg 3.0.5.2 "Madeira" (mirror documentation.help, páginas linguagem 1-5 + funções) | `[OFICIAL]` | Alta. É o texto do autor (Cláudio Morgado de Souza / Apoio Informática) |
| Manual PDF "A Linguagem de Programação do VisuAlg" (mirrors CEFET-MG, UFU, UFSC) | `[OFICIAL]` | Alta. Mesmo conteúdo do help, consolidado |
| Notas de release do VisuAlg 3.0.6/3.0.7 (Antonio Carlos Nicolodi, SourceForge) | `[3.0.6+]` | Média |
| Implementação TypeScript @designliquido/visualg (engenharia reversa de terceiros) | `[DL]` | Média/baixa. Útil como referência, diverge do real em pontos conhecidos (ver §13) |
| Comportamento amplamente reportado pela comunidade mas não documentado | `[INFERIDO]` | Baixa. Validar contra VISUALG.EXE antes de congelar comportamento |

Itens marcados `[VERIFICAR]` precisam de teste empírico contra o interpretador real (ver §12, checklist de conformidade).

---

## 1. Estrutura léxica

### 1.1 Regras gerais `[OFICIAL]`
- **Um comando por linha.** A quebra de linha é o terminador de comando. Não há separador `;` entre comandos nem blocos `begin/end`.
- **Case-insensitive** para palavras-chave e identificadores (`Escreva` = `ESCREVA` = `escreva`; `Nota` e `nota` são a mesma variável).
- Palavras-chave canônicas **sem acentos** (`logico`, `entao`, `senao`, `faca`, `funcao`).
- `[3.0.6+]` As formas acentuadas também são aceitas: `FAÇA`, `ATÉ`, `ENTÃO`, `SENÃO`, `FUNÇÃO`, `FIMFUNÇÃO`, `NÃO`. As formas sem acento continuam válidas (compatibilidade). Um lexer conforme deve aceitar ambas e normalizar.
- Não existe `goto`.
- Todo texto após a linha `fimalgoritmo` é ignorado pelo interpretador.

### 1.2 Comentários `[OFICIAL]`
- `//` até o fim da linha. Não há comentário de bloco. Comentários multi-linha exigem `//` em cada linha.
- `[VERIFICAR]` Algumas versões aceitam `{ ... }` estilo Pascal? Não documentado; assumir que não.

### 1.3 Identificadores `[OFICIAL]`
- Começam com letra; depois letras, dígitos ou `_`.
- Máximo 30 caracteres.
- Case-insensitive.
- Não podem colidir com palavras-chave.
- `[VERIFICAR]` Se letras acentuadas são aceitas em identificadores (provavelmente sim no 3.x via charset ANSI/Windows-1252, mas não documentado).

### 1.4 Literais (constantes) `[OFICIAL]`
- **Numéricos:** inteiros (`42`, `-3`) e reais (`1.5`). Separador decimal é **ponto**, sempre, independente do locale. Sem separador de milhares. Sem notação científica documentada `[VERIFICAR]`.
- **Caractere (string):** delimitado por aspas duplas `"..."`. Não há documentação de escape de aspas dentro da string `[VERIFICAR: provavelmente não existe escape; aspas internas são impossíveis]`. Aspas simples não são válidas.
- **Lógicos:** `VERDADEIRO` e `FALSO` (case-insensitive).

### 1.5 Tokens de operadores e pontuação
```
<-          atribuição (também ":" seguido de "-"? não; apenas "<-")
:           separador de tipo em declarações; formatação em escreva
..          intervalo em declaração de vetor
,           separador de listas
( ) [ ]     agrupamento / índices
+ - * / \ % ^   aritméticos (\ = divisão inteira, % = mod)
= <> < > <= >=  relacionais
;           separador de grupos de parâmetros em declarações de subprogramas
.           acesso a campo (apenas se implementar registro, ver §11)
```
Nota sobre `;`: no corpo do programa é desnecessário/ignorado, mas é **sintaticamente significativo** na lista de declaração de parâmetros de `procedimento`/`funcao` (separa grupos de parâmetros de tipos diferentes). `[DL]` trata `;` como whitespace em todo lugar; isso é uma simplificação que funciona por acaso.

### 1.6 Palavras reservadas (conjunto 3.0.5.2 + 3.0.6)
```
algoritmo, fimalgoritmo, var, inicio, const?*, 
inteiro, real, caractere, caracter, logico, vetor, de,
se, entao, senao, fimse,
escolha, caso, outrocaso, fimescolha,
para, ate, passo, faca, fimpara,
enquanto, fimenquanto,
repita, fimrepita,
interrompa, retorne,
procedimento, fimprocedimento, funcao, fimfuncao,
leia, escreva, escreval,
e, ou, nao, xou, mod,
verdadeiro, falso,
aleatorio, arquivo, timer, pausa, debug, eco, cronometro, limpatela,
on, off,
dos?*
```
- `caracter` (sem o "e" final) é aceito como sinônimo de `caractere` em declarações de parâmetros e retorno de funções. O help oficial usa `caracter` nas assinaturas das funções builtin.
- `*` `const` e `dos`: aparecem em algumas listas da comunidade, mas não no help oficial 3.0.5.2. `[VERIFICAR]` antes de reservar.
- Formas acentuadas `[3.0.6+]`: `até, faça, então, senão, função, fimfunção, não` mapeiam para os mesmos tokens.

---

## 2. Tipos de dados `[OFICIAL]`

| Tipo | Palavra-chave | Domínio |
|---|---|---|
| Inteiro | `inteiro` | inteiros com sinal (Delphi: provavelmente Integer 32 bits `[INFERIDO]`) |
| Real | `real` | ponto flutuante (Delphi: provavelmente Double/Extended `[INFERIDO]`) |
| Caractere | `caractere` / `caracter` | string (não há tipo char separado; "caractere" é string) |
| Lógico | `logico` | VERDADEIRO / FALSO |

- Não há tipo char unitário, nem ponteiros, nem registros no VisuAlg oficial até 3.0.5.2 (ver §11 para `registro` em implementações alternativas).
- **Vetores:** `vetor [li..ls] de <tipo>` (1D) e `vetor [li1..ls1, li2..ls2] de <tipo>` (2D). Limites devem ser constantes inteiras, com `ls > li`. Máximo 2 dimensões.
- Limite global: **500 variáveis**, contando cada elemento de vetor individualmente. (Um interpretador novo provavelmente não quer reproduzir esse limite, mas ele existe no original.)

### 2.1 Coerções e compatibilidade de tipos
- Atribuição exige que o resultado da expressão "tenha tipo igual ao da variável" `[OFICIAL]`. Na prática:
  - `inteiro` → `real`: permitido (promoção implícita) `[INFERIDO, comportamento amplamente observado]`.
  - `real` → `inteiro`: **erro de execução** ("incompatibilidade de tipos"), inclusive para `x <- 5/2` com `x: inteiro` `[INFERIDO]`.
- Divisão `/` entre dois inteiros: `[VERIFICAR]` ponto crítico. Duas hipóteses observadas na comunidade: (a) sempre produz real; (b) produz inteiro quando a divisão é exata. Teste obrigatório: `x: inteiro; x <- 4/2` dá erro ou não?
- `+` entre strings: concatenação. `+` entre string e número: **erro** (use `numpcarac`) `[INFERIDO]`.
- Comparações relacionais exigem operandos do mesmo tipo (com a promoção inteiro/real).

---

## 3. Estrutura do programa `[OFICIAL]`

```
algoritmo "nome entre aspas"
// comentários opcionais
[arquivo "dados.txt"]          // opcional, no máximo 1, na seção de declarações
[var
   <declarações>]
[<declarações de procedimentos e funções>]   // ver §7; ficam entre var e inicio
inicio
   <comandos>
fimalgoritmo
```

- A linha `algoritmo "nome"` é obrigatória e o nome fica entre aspas duplas. O nome é usado como título de janelas de leitura.
- Declarações de variáveis: seção `var` com linhas no formato:
  - `<lista-de-nomes> : <tipo>`
  - `<lista-de-nomes> : vetor [<intervalos>] de <tipo>`
- `intervalo := <inteiro> .. <inteiro>`, separados por vírgula para 2D.
- `[VERIFICAR]` Se múltiplas seções `var` são permitidas no programa principal (provavelmente apenas uma).

---

## 4. Operadores e precedência

### 4.1 Operadores documentados `[OFICIAL]`

**Aritméticos:**
| Op | Significado | Nota oficial |
|---|---|---|
| `+x`, `-x` | unários | "maior precedência entre os aritméticos" |
| `^` | potenciação | "maior precedência entre os aritméticos **binários**" |
| `*`, `/` | multiplicação, divisão | precedem `+` e `-` |
| `\` | divisão inteira (`5 \ 2 = 2`) | "mesma precedência da divisão tradicional" |
| `mod`, `%` | resto da divisão inteira (`8 mod 3 = 2`) | "mesma precedência da divisão tradicional" |
| `+`, `-` | adição, subtração | |

**Caracteres:** `+` concatena strings.

**Relacionais:** `=`, `<`, `>`, `<=`, `>=`, `<>`. Operandos do mesmo tipo, resultado lógico.
- Comparação de strings é **case-insensitive**: `"ABC" = "abc"` é VERDADEIRO `[OFICIAL]`. Detalhe importantíssimo e fácil de errar numa reimplementação.
- `[VERIFICAR]` Se a comparação case-insensitive usa apenas A-Z ou também acentuados (provável: tabela ANSI do Delphi, não Unicode).
- Ordem de lógicos: `FALSO < VERDADEIRO` `[OFICIAL]`.

**Lógicos:** `nao` (unário, maior precedência entre lógicos), `e`, `ou`, `xou` (xor). `[OFICIAL]`

### 4.2 Tabela de precedência consolidada `[INFERIDO + parcialmente OFICIAL]`

O manual só afirma parcialmente a hierarquia. A tabela abaixo é a reconstrução consistente com o manual e com o comportamento esperado (maior para menor):

```
1. ( )  e chamadas de função
2. unários: +, -, nao
3. ^                       (assoc. à direita? [VERIFICAR] — em Pascal-like geralmente à esquerda)
4. *, /, \, mod, %
5. +, - (binários), + (concatenação)
6. =, <>, <, >, <=, >=
7. e
8. xou
9. ou
```

Pontos a validar empiricamente (`[VERIFICAR]`, todos):
1. Se `a > 1 e a < 5` parseia sem parênteses ou exige `(a > 1) e (a < 5)`. Em Pascal exigiria parênteses (and tem precedência de multiplicação); o VisuAlg aparenta aceitar sem, mas o material didático usa sempre parênteses, o que mascara o comportamento real.
2. Precedência relativa de `e`/`ou`/`xou` entre si (a tabela acima assume e > xou > ou; alternativa: e > ou = xou).
3. Associatividade de `^`: `2^3^2` = 64 (esquerda) ou 512 (direita)?
4. `nao` aplicado a comparação: `nao a = b` parseia como `nao (a = b)` ou `(nao a) = b`?
5. Resultado de `mod` com operandos negativos (`-7 mod 3`): Delphi retorna -1 (sinal do dividendo).
6. `\` e `mod` com operandos reais: erro ou truncamento?

### 4.3 Atribuição `[OFICIAL]`
```
<variável> <- <expressão>
vet[i] <- <expressão>
matriz[i, j] <- <expressão>
```
- `<-` não é expressão (não há atribuição encadeada).
- `[3.0.6+]` `[VERIFICAR]` Algumas notas citam aceitação de `:=` como alternativa. Não confirmado em fonte primária.

---

## 5. Entrada e saída

### 5.1 `escreva (<lista-de-expressões>)` `[OFICIAL]`
- Avalia e imprime as expressões na ordem, separadas por vírgula na lista.
- **Formatação Pascal-like:**
  - `expr:N` imprime em campo de N posições, alinhado à direita.
  - `expr:N:M` (para reais) campo de N posições com M casas decimais.
- A formatação `:N[:M]` se aplica a qualquer expressão na lista, não só variáveis (ex.: `escreva(y+3:4)`).
- **Regra de espaçamento automático:** valores numéricos e lógicos são impressos com **um espaço à esquerda**; valores caractere são impressos sem espaço (para permitir concatenação visual). Exemplo oficial:
  ```
  x <- 2.5; y <- 6; a <- "teste"
  escreval ("x", x:4:1, y+3:4)   // imprime: x 2.5    9
  escreval (a, "ok")             // imprime: testeok
  escreva (VERDADEIRO)           // imprime:  VERDADEIRO
  ```
  `[VERIFICAR]` Interação exata do espaço automático com `:N` (o espaço entra dentro ou fora do campo? o exemplo oficial "x 2.5" com x:4:1 sugere que " 2.5" ocupa o campo de 4 e o espaço é o padding do próprio campo). Teste com valores que enchem o campo.
- `[VERIFICAR]` Formato default de impressão de reais sem `:N:M` (quantas casas? remove zeros à direita? `2.0` imprime como "2"?). Ponto crítico para compatibilidade de saída.
- Lógicos imprimem como ` VERDADEIRO` / ` FALSO`.

### 5.2 `escreval (...)` `[OFICIAL]`
Igual a `escreva`, pulando linha ao final. `escreval` sem argumentos pula linha `[INFERIDO, uso universal]`.

### 5.3 `leia (<lista-de-variáveis>)` `[OFICIAL]`
- Lê valores do dispositivo de entrada, atribuindo na ordem da lista.
- Aceita elementos de vetor: `leia(vet[i])`.
- GUI: abre prompt "Entre com o valor de <nome>". Esc/Cancelar **aborta o programa imediatamente**.
- Validação de tipo na entrada: valor digitado incompatível com o tipo da variável gera nova solicitação ou erro `[VERIFICAR]`.
- Interação com `arquivo` e `aleatorio`: ver §8.

---

## 6. Estruturas de controle `[OFICIAL]`

### 6.1 `se`
```
se <expr-lógica> entao
   <comandos>
[senao
   <comandos>]
fimse
```
Aninhamento permitido. Não há `senao se` encadeado como token único (usa-se aninhar `se` dentro do `senao`).

### 6.2 `escolha`
```
escolha <expressão-de-seleção>
   caso <exp1>, <exp2>, ..., <expN>
      <comandos>
   caso ...
      <comandos>
   [outrocaso
      <comandos>]
fimescolha
```
- Cada `caso` aceita lista de valores separados por vírgula.
- Sem fall-through: executa o bloco do primeiro caso compatível e salta para depois de `fimescolha` `[INFERIDO, semântica Pascal case]`.
- Comparação de seleção com strings: case-insensitive (segue a regra geral de comparação) `[INFERIDO]`.
- Parênteses em volta da expressão de seleção são opcionais (`escolha x` e `escolha (x)`) `[INFERIDO via DL + exemplos da comunidade]`.
- `[VERIFICAR]` Suporte a intervalos `caso 1 ate 5`. Aparece em material da comunidade, mas NÃO no help oficial 3.0.5.2. Os valores de caso podem ser expressões ou apenas literais? O help usa apenas literais.

### 6.3 `para`
```
para <variável> de <valor-inicial> ate <valor-limite> [passo <incremento>] faca
   <comandos>
fimpara
```
Semântica oficial, precisa e completa:
- `<variável>` deve ser **inteiro**, assim como todas as expressões do comando.
- `<valor-inicial>`, `<valor-limite>` e `<incremento>` são avaliados **uma única vez**, antes da primeira iteração, e **não mudam** durante o laço mesmo que variáveis usadas neles mudem.
- `passo` default = 1. Passo negativo permitido. **Passo zero (avaliado) = erro de execução** com mensagem.
- Teste de continuação: ao chegar em `fimpara`, soma o incremento à variável e compara com o limite (`<=` para passo positivo, `>=` para negativo).
- Se já na entrada `inicial > limite` (passo positivo) ou `inicial < limite` (passo negativo), o corpo executa **zero vezes**.
- `[VERIFICAR]` Valor da variável de controle após o término do laço (no Pascal clássico é indefinido; no VisuAlg, pelo modelo de execução descrito, fica `limite + passo` ou além).
- `[VERIFICAR]` Modificar a variável de controle dentro do corpo: permitido? afeta o laço?

### 6.4 `enquanto`
```
enquanto <expr-lógica> faca
   <comandos>
fimenquanto
```
Testa antes; executa 0+ vezes.

### 6.5 `repita`
Duas formas:
```
repita
   <comandos>
ate <expr-lógica>
```
Testa depois; executa 1+ vezes; repete enquanto a expressão for FALSO.

Forma alternativa (laço infinito explícito):
```
repita
   <comandos>
fimrepita
```
Só termina via `interrompa`.

### 6.6 `interrompa`
Sai imediatamente do laço mais interno (`para`, `enquanto`, `repita`, `repita..fimrepita`). `[VERIFICAR]` comportamento fora de laço (erro de compilação ou de execução?).

---

## 7. Subprogramas `[OFICIAL]`

### 7.1 Regras gerais
- Declarados **após a seção `var` global e antes do `inicio` do programa principal**.
- **Aninhamento não permitido** (subprograma dentro de subprograma).
- **Recursão permitida.**
- Subprogramas enxergam as variáveis globais; podem declarar variáveis locais (seção `var` própria) que existem só durante a chamada.
- Quantidade, ordem e tipos dos argumentos na chamada devem bater com a declaração.

### 7.2 Procedimento
```
procedimento <nome> [(<declarações-de-parâmetros>)]
[var <declarações locais>]
inicio
   <comandos>
fimprocedimento
```
- `<declarações-de-parâmetros>` = grupos `[var] <nomes separados por vírgula>: <tipo>` separados por **ponto e vírgula**.
- `var` no grupo = passagem **por referência**; sem `var` = **por valor**.
- Chamada como comando: `soma` ou `soma(n, m)`.

### 7.3 Função
```
funcao <nome> [(<declarações-de-parâmetros>)]: <tipo-de-retorno>
[var <declarações locais>]
inicio
   <comandos>
fimfuncao
```
- Retorno via comando `retorne <expressão>`.
- Chamada é **expressão**.
- **Detalhe crucial de parsing:** função sem parâmetros pode ser chamada **sem parênteses** (`res <- soma`). Ou seja, um identificador em posição de expressão que nomeia uma função é uma chamada. O parser precisa conhecer os subprogramas declarados (eles vêm antes do corpo principal, então uma passada de declarações antes de resolver expressões resolve isso).
- `[VERIFICAR]` `retorne` em procedimento (retorno antecipado sem valor): aceito?
- `[VERIFICAR]` Função que termina sem executar `retorne`: erro ou valor default?
- `[VERIFICAR]` Vetores como parâmetros: o help oficial não documenta. Comunidade reporta suporte parcial/bugado no 3.x. Testar `procedimento p(v: vetor[1..10] de inteiro)`.
- `[VERIFICAR]` `retorne` dentro do programa principal: erro?

### 7.4 Passagem por referência
- Parâmetro `var` recebe o endereço; alterações afetam a variável do chamador.
- `[VERIFICAR]` Passar expressão/literal para parâmetro `var`: erro de compilação?
- `[VERIFICAR]` Passar elemento de vetor (`v[i]`) por referência: suportado?

---

## 8. Comandos especiais de ambiente `[OFICIAL]`

### 8.1 `aleatorio`
Substitui digitação em `leia` por geração aleatória.
```
aleatorio [on]              // ativa; faixa default 0..100 inclusive
aleatorio <v1> [, <v2>]     // ativa com faixa: 0..v1, ou v1..v2 (troca se v2 < v1)
aleatorio off               // desativa (off obrigatório)
```
- `v1`/`v2` devem ser **constantes numéricas**, não expressões.
- Não afeta leitura de variáveis lógicas.
- Variáveis caractere: gera strings de **5 letras maiúsculas**.
- Fica na seção de comandos `[INFERIDO: o help mostra na seção de comandos]`.
- `[VERIFICAR]` Para variável inteira com faixa, gera inteiro uniforme; para real, gera real? Granularidade?

### 8.2 `arquivo "<nome>"`
Redireciona `leia` para arquivo-texto.
- Vai na **seção de declarações** (antes de `var`? junto? `[VERIFICAR posição exata]`), no máximo **um** por programa.
- Se o arquivo não existe: lê por digitação e **grava** os valores lidos no arquivo, na ordem.
- Se existe: consome valores do arquivo até o fim; depois volta à digitação.
- Sem caminho: resolve relativo à pasta de trabalho corrente (onde está o VISUALG.EXE). Sem extensão default.
- `[VERIFICAR]` Formato do arquivo: um valor por linha? separado por espaço?

### 8.3 `timer`
```
timer on | timer <ms> | timer off
```
- Atraso antes de cada linha; default 500 ms; faixa 0..10000 (valores fora são clampados).
- Vários `timer` permitidos, todos na seção de comandos.
- Para um interpretador headless: implementar como no-op configurável.

### 8.4 `pausa`
Breakpoint incondicional (interativo no IDE). Headless: no-op ou hook de debugger.

### 8.5 `debug <expr-lógica>`
Breakpoint condicional: interrompe se a expressão for VERDADEIRO.

### 8.6 `eco on | off`
Ativa/desativa a impressão dos dados de entrada na saída padrão.

### 8.7 `cronometro on | off`
- `on`: imprime "Cronômetro iniciado." e inicia contagem em ms.
- `off`: imprime "Cronômetro terminado. Tempo decorrido: xx segundo(s) e xx ms".

### 8.8 `limpatela`
Limpa a "tela DOS" (simulação de console). Não afeta a área de saída padrão do IDE.

---

## 9. Biblioteca de funções predefinidas `[OFICIAL]`

Todas case-insensitive. Usáveis em qualquer posição de expressão (nunca no lado esquerdo de `<-`).

### 9.1 Numéricas, algébricas e trigonométricas
| Função | Assinatura | Semântica |
|---|---|---|
| `abs(x)` | inteiro/real → mesmo tipo | valor absoluto |
| `arccos(x)` | real → real | arco cosseno, radianos |
| `arcsen(x)` | real → real | arco seno, radianos |
| `arctan(x)` | real → real | arco tangente, radianos |
| `cos(x)` | real → real | cosseno (radianos) |
| `cotan(x)` | real → real | cotangente (radianos) |
| `exp(base, expoente)` | real, real → real | base^expoente (NÃO é e^x!) |
| `grauprad(x)` | real → real | graus → radianos |
| `int(x)` | real → inteiro | parte inteira. `[VERIFICAR]` trunca em direção a zero ou floor? Para negativos: `int(-1.8)` = -1 (trunc, comportamento Delphi Int()) ou -2 (floor, comportamento DL)? **DL usa Math.floor, o que provavelmente diverge do original Delphi (Trunc).** |
| `log(x)` | real → real | log base 10 |
| `logn(x)` | real → real | log natural (base e) |
| `pi` | → real | saída padrão registrada: `3.14159265358979` |
| `quad(x)` | num → num | x*x |
| `radpgrau(x)` | real → real | radianos → graus |
| `raizq(x)` | real → real | raiz quadrada |
| `rand` | → real | aleatório em [0, 1) |
| `randi(limite)` | inteiro → inteiro | aleatório em [0, limite) |
| `sen(x)` | real → real | seno (radianos) |
| `tan(x)` | real → real | tangente (radianos) |

`pi` é usado sem parênteses; a forma `pi()` é rejeitada com diagnóstico de sintaxe nas gravações. As formas de chamada de `rand` continuam em verificação.

### 9.2 Manipulação de strings
| Função | Assinatura | Semântica |
|---|---|---|
| `asc(s)` | caracter → inteiro | código ASCII do **primeiro** caractere |
| `carac(c)` | inteiro → caracter | caractere do código ASCII |
| `caracpnum(c)` | caracter → inteiro ou real | parse numérico (StrToInt/StrToFloat). `[VERIFICAR]` comportamento com string inválida: erro de execução? |
| `compr(c)` | caracter → inteiro | comprimento |
| `copia(c, p, n)` | caracter, inteiro, inteiro → caracter | substring a partir da posição p (base 1), com n caracteres |
| `maiusc(c)` | caracter → caracter | uppercase. `[VERIFICAR]` se acentuados são convertidos (Delphi UpperCase só converte A-Z) |
| `minusc(c)` | caracter → caracter | lowercase. Mesma ressalva |
| `numpcarac(n)` | inteiro/real → caracter | número → string. `[VERIFICAR]` formato de real |
| `pos(subc, c)` | caracter, caracter → inteiro | posição (base 1) de subc em c; 0 se ausente. Ordem dos args: agulha, palheiro. `[VERIFICAR]` se a busca é case-sensitive (provável que sim, ao contrário das comparações!) |

### 9.3 Observações para o implementador
- Os nomes das builtins NÃO são impedidos de colidir com subprogramas do usuário? `[VERIFICAR]` (provavelmente são reservados).
- `exp` é a maior pegadinha da biblioteca: tem 2 argumentos e é potenciação, não exponencial.

---

## 10. Modelo de execução e runtime

- **Dois displays** `[OFICIAL]`: a "saída padrão" (painel do IDE) e a "tela DOS" (simulação de console, afetada por `limpatela`). Em implementação headless, trate como um único stdout.
- Escopo: global (programa principal) + local por chamada de subprograma. Sem blocos léxicos internos. Shadowing de global por local/parâmetro: permitido `[INFERIDO]`.
- Pilha de ativação visível no IDE (Ctrl-F3); profundidade máxima de recursão não documentada `[VERIFICAR: estouro de pilha vs limite artificial]`.
- Valores iniciais de variáveis: **não documentado**. Comunidade reporta inicialização com zero/""/FALSO `[VERIFICAR — decisão importante: ler variável não inicializada é erro ou retorna zero-value?]`.
- Erros de execução param o programa com mensagem apontando a linha.

---

## 11. Diferenças entre versões

| Versão | Mudanças relevantes para a linguagem |
|---|---|
| 2.0/2.5 (Cláudio Morgado de Souza / Apoio Informática) | Base da linguagem como descrita aqui. O help "Funções do Visualg Versão 2.0" é o mesmo conjunto de builtins do 3.x |
| 3.0–3.0.5.2 "Madeira" | Reescrita do IDE; linguagem essencialmente idêntica à 2.5 |
| 3.0.6 / 3.0.7 (Antonio Carlos Nicolodi, 2016+) | Aceita palavras-chave acentuadas (`PARA...FAÇA`, `SE..ENTÃO..SENÃO`, operador `NÃO`), mantendo as antigas. Mudanças de IDE (skins etc.) sem impacto na linguagem `[3.0.6+]` |
| @designliquido/visualg (TS) | Implementa extensões que NÃO existem no VisuAlg oficial: `tipo ... registro ... fimregistro` e acesso a campos com `.`. Se seu objetivo for compatibilidade com o VisuAlg real, **não** implemente registro |

Se a meta é "rodar qualquer .ALG de sala de aula", o alvo é o conjunto 3.0.5.2 + keywords acentuadas do 3.0.6.

---

## 12. Checklist de conformidade (testes empíricos recomendados)

Antes de congelar a semântica do seu interpretador Go, rode estes casos no VISUALG.EXE 3.0.7 (Wine roda bem) e registre a saída como golden files:

1. `x: inteiro; x <- 5/2` → erro? E `x <- 4/2`?
2. `escreva(2.0)`, `escreva(2.5)`, `escreva(1/3)` → formato default de reais.
3. `escreva(123:2)` → campo menor que o valor: trunca ou expande?
4. `escreva("a", 1, "b")` → posicionamento exato dos espaços automáticos.
5. `se a > 1 e a < 5 entao` sem parênteses → compila?
6. `2^3^2` → 64 ou 512?
7. `nao a = b` → parse.
8. `-7 mod 3`, `-7 \ 3` → sinais.
9. `7.5 mod 2`, `7.5 \ 2` → erro ou trunca?
10. `int(-1.8)` → -1 ou -2?
11. `pos("a", "BANANA")` → 0 (case-sensitive) ou 2 (insensitive)?
12. `maiusc("ção")` → "ÇÃO" ou "çãO"...?
13. `escolha x` com `caso 1 ate 5` → aceita?
14. Variável não inicializada em expressão → erro ou zero?
15. Valor da variável de controle após `fimpara`.
16. Modificar variável de controle dentro do `para`.
17. Função sem `retorne` executado.
18. `interrompa` fora de laço.
19. Vetor como parâmetro de subprograma.
20. `v[i]` como argumento de parâmetro `var`.
21. Leitura com `arquivo`: formato do arquivo gerado.
22. String com 31+ caracteres no identificador → erro em que fase?
23. Literais: `1.` , `.5`, `1e3` → aceitos?
24. `escreval(verdadeiro e falso ou verdadeiro)` → precedência e/ou.
25. Comparação `"a" < "B"` → ordem com case-insensitivity.

## 13. Divergências conhecidas da implementação Design Líquido (não copiar)

1. `escreva(x:5)`: DL prefixa 5 espaços; o real alinha à direita em campo de largura 5 (Pascal).
2. `int()`: DL usa floor; o real (Delphi `Int`/`Trunc`) trunca em direção a zero `[a confirmar no teste 10]`.
3. DL implementa `registro`/`tipo`, inexistentes no VisuAlg real.
4. DL trata `;` como whitespace universal.
5. DL não implementa `caso` com intervalo nem vários comandos especiais (`eco`, `cronometro`, `timer`, `arquivo`, `pausa`, `debug` parcial).
6. DL usa parser herdado do Delégua para precedência abaixo de `ou`/`e`; não é evidência do comportamento do VisuAlg real.

## 14. Gramática EBNF reconstruída `[INFERIDO]`

Reconstrução minha a partir das fontes acima. Não é oficial. NEWLINE é token significativo (terminador de comando). Tokens case-insensitive; formas acentuadas omitidas por brevidade (normalizar no lexer).

```ebnf
programa        = "algoritmo" STRING NEWLINE
                  { NEWLINE }
                  [ "arquivo" STRING NEWLINE ]
                  [ secao_var ]
                  { subprograma }
                  "inicio" NEWLINE
                  { comando }
                  "fimalgoritmo" (* resto do arquivo ignorado *) ;

secao_var       = "var" NEWLINE { declaracao NEWLINE } ;
declaracao      = lista_ids ":" tipo ;
lista_ids       = ID { "," ID } ;
tipo            = tipo_simples
                | "vetor" "[" intervalo [ "," intervalo ] "]" "de" tipo_simples ;
tipo_simples    = "inteiro" | "real" | "caractere" | "caracter" | "logico" ;
intervalo       = INT ".." INT ;

subprograma     = procedimento | funcao ;
procedimento    = "procedimento" ID [ "(" params ")" ] NEWLINE
                  [ secao_var ] "inicio" NEWLINE { comando }
                  "fimprocedimento" NEWLINE ;
funcao          = "funcao" ID [ "(" params ")" ] ":" tipo_simples NEWLINE
                  [ secao_var ] "inicio" NEWLINE { comando }
                  "fimfuncao" NEWLINE ;
params          = grupo_param { ";" grupo_param } ;
grupo_param     = [ "var" ] lista_ids ":" tipo ;   (* vetor como tipo de param: [VERIFICAR] *)

comando         = ( atribuicao | chamada_proc | cmd_leia | cmd_escreva
                  | cmd_se | cmd_escolha | cmd_para | cmd_enquanto
                  | cmd_repita | "interrompa" | "retorne" [ expr ]
                  | cmd_especial | (* linha vazia *) ) NEWLINE ;

atribuicao      = lvalue "<-" expr ;
lvalue          = ID [ "[" expr [ "," expr ] "]" ] ;
chamada_proc    = ID [ "(" [ expr { "," expr } ] ")" ] ;
cmd_leia        = "leia" "(" lvalue { "," lvalue } ")" ;
cmd_escreva     = ( "escreva" | "escreval" ) [ "(" [ item_escrita { "," item_escrita } ] ")" ] ;
item_escrita    = expr [ ":" expr [ ":" expr ] ] ;

cmd_se          = "se" expr "entao" NEWLINE { comando }
                  [ "senao" NEWLINE { comando } ] "fimse" ;
cmd_escolha     = "escolha" [ "(" ] expr [ ")" ] NEWLINE
                  { "caso" valores_caso NEWLINE { comando } }
                  [ "outrocaso" NEWLINE { comando } ]
                  "fimescolha" ;
valores_caso    = valor_caso { "," valor_caso } ;
valor_caso      = literal [ "ate" literal ]      (* intervalo: [VERIFICAR] *) ;
cmd_para        = "para" ID "de" expr "ate" expr [ "passo" expr ] "faca" NEWLINE
                  { comando } "fimpara" ;
cmd_enquanto    = "enquanto" expr "faca" NEWLINE { comando } "fimenquanto" ;
cmd_repita      = "repita" NEWLINE { comando } ( "ate" expr | "fimrepita" ) ;

cmd_especial    = "aleatorio" ( [ "on" ] | NUM [ "," NUM ] | "off" )
                | "timer" ( "on" | "off" | INT )
                | "pausa" | "debug" expr
                | "eco" ( "on" | "off" )
                | "cronometro" ( "on" | "off" )
                | "limpatela" ;

(* expressões: precedência conforme §4.2; pontos [VERIFICAR] anotados lá *)
expr            = expr_ou ;
expr_ou         = expr_xou { "ou" expr_xou } ;
expr_xou        = expr_e { "xou" expr_e } ;
expr_e          = expr_rel { "e" expr_rel } ;
expr_rel        = expr_adit [ ( "=" | "<>" | "<" | ">" | "<=" | ">=" ) expr_adit ] ;
expr_adit       = expr_mult { ( "+" | "-" ) expr_mult } ;
expr_mult       = expr_pot { ( "*" | "/" | "\\" | "mod" | "%" ) expr_pot } ;
expr_pot        = expr_unaria { "^" expr_unaria } ;
expr_unaria     = ( "+" | "-" | "nao" ) expr_unaria | primario ;
primario        = NUM | STRING | "verdadeiro" | "falso"
                | "(" expr ")"
                | ID [ "(" [ expr { "," expr } ] ")" | "[" expr [ "," expr ] "]" ] ;
                (* ID nu que nomeia função = chamada sem parênteses, ver §7.3 *)
```

## 15. Fontes

- Help oficial VisuAlg 3.0.5.2: https://documentation.help/VISUALG-Versao/ (páginas `linguagem.htm` a `linguagem5.htm`, `funcoes.htm`)
- Manual PDF (mirror CEFET-MG): https://www.acad.cefetmg.br/uploads/MATERIAIS_AULAS/340988-A_Linguagem_de_Programação_do_VisuAlg.pdf
- Manual VisuAlg 2.x (mirror UFSC): http://www.inf.ufsc.br/~bosco.sobral/ensino/ine5201/Visualg2_manual.pdf
- VisuAlg 3.0 (Nicolodi, binários + LEIA-ME): https://sourceforge.net/projects/visualg30/
- Implementação TS de referência cruzada: https://github.com/DesignLiquido/visualg
