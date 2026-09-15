# PortuGo

Português (Brasil) | [English](README.md)

**Aprenda e pratique Portugol no Linux e no macOS, direto do terminal.**

O PortuGo é para quem usa Linux ou macOS e quer acompanhar cursos com
VisuAlg, resolver exercícios ou ensinar programação em português.
Ele executa programas `.alg` no **dialeto VisuAlg 3.x**, com o editor que
você já gosta de usar. Você não precisa instalar Windows, máquina virtual
ou Wine.

Está escrevendo seu primeiro `algoritmo` ou voltando a um exercício das
aulas? Seja bem-vindo! Os passos abaixo mostram como compilar o projeto
e colocar seu primeiro programa para rodar.

## Primeiros passos

### Baixar uma versão pronta

As versões publicadas estão no [GitHub Releases](https://github.com/ncode/PortuGo/releases).
Quando uma versão é publicada, a seção **Assets** inclui o executável
`portugo` para Linux, macOS e Windows. Você não precisa instalar o Go para
usar um executável pronto.

Escolha o arquivo correspondente ao seu computador:

| Sistema | Final do nome do arquivo |
| --- | --- |
| Linux, Intel/AMD de 64 bits | `linux-amd64.tar.gz` |
| Linux, ARM64 | `linux-arm64.tar.gz` |
| macOS, Intel | `darwin-amd64.tar.gz` |
| macOS, Apple Silicon | `darwin-arm64.tar.gz` |
| Windows, Intel/AMD de 64 bits | `windows-amd64.zip` |
| Windows, ARM64 | `windows-arm64.zip` |

Descompacte o arquivo e abra um terminal na pasta extraída. Execute
`./portugo --version` para conferir a versão instalada e
`./portugo run your-program.alg` para executar um exercício, substituindo
`your-program.alg` pelo caminho do seu arquivo. No PowerShell do Windows,
use `.\portugo.exe` no lugar de `./portugo`. Cada versão também inclui o
arquivo `SHA256SUMS` para verificar a integridade do arquivo baixado.

### Compilar a partir do código-fonte

Você vai precisar do **[Git](https://git-scm.com/install/)** e do **Go 1.27 ou
mais recente**. O [guia oficial de instalação do Go](https://go.dev/doc/install),
em inglês, tem instruções para Linux e macOS. Você não precisa saber
programar em Go para usar o PortuGo; ele cuida da compilação da ferramenta
de linha de comando para você.

Abra um terminal e confira se as duas ferramentas estão disponíveis:

```sh
git --version
go version
```

Se algum comando não for encontrado, instale a ferramenta que falta e abra
um novo terminal antes de continuar.

Depois, baixe o projeto, compile e execute o primeiro exemplo:

```sh
git clone https://github.com/ncode/PortuGo.git
cd PortuGo
go build -o portugo ./cmd/portugo
./portugo run examples/hello.alg
```

Você deve ver:

```text
Ola, Portugol!
```

Seu primeiro programa já está rodando! `portugo` é o executável que você
acabou de compilar. O `./` indica ao terminal que ele está na pasta atual.
Os comandos abaixo consideram que você continua na pasta `PortuGo`.

## Escreva seu próprio programa

Abra seu editor de texto favorito e salve o código abaixo como
**`boas-vindas.alg`** na pasta do projeto. Use texto simples com
codificação UTF-8:

```portugol
algoritmo "boas-vindas"
var
  nome: caractere
inicio
  escreva("Como voce se chama? ")
  leia(nome)
  escreval("Ola, ", nome, "! Vamos programar?")
fimalgoritmo
```

Para executar:

```sh
./portugo run boas-vindas.alg
```

Digite seu nome e pressione Enter. O programa vai cumprimentar você pelo
nome. `leia` lê sua resposta, `escreva` mostra a pergunta e `escreval` mostra
a saudação, com uma quebra de linha no final. O PortuGo também exibe a
entrada após a leitura, então é normal ver sua resposta repetida.

Experimente mudar a saudação e executar o arquivo de novo. Cada arquivo
`.alg` contém um programa completo, que começa com `algoritmo` e termina
com `fimalgoritmo`.

## Comandos do dia a dia

| O que você quer fazer | Comando |
| --- | --- |
| Executar um programa | `./portugo run boas-vindas.alg` |
| Verificar a sintaxe e os tipos sem executar | `./portugo check boas-vindas.alg` |
| Ver o código formatado no terminal | `./portugo fmt boas-vindas.alg` |
| Formatar e salvar o arquivo | `./portugo fmt -w boas-vindas.alg` |
| Verificar se o arquivo já está formatado | `./portugo fmt --check boas-vindas.alg` |
| Iniciar uma sessão interativa | `./portugo repl` |
| Ver a ajuda de um comando | `./portugo run -h` |
| Mostrar a versão do executável | `./portugo --version` |

Coloque opções como `-w` e `--check` **antes** do nome do arquivo.
O comando `check` não exibe nada quando está tudo certo. Se houver algum
problema, as mensagens indicam o arquivo, a linha e a coluna para você
localizá-lo no editor. `fmt --check` termina com código de saída 1 quando
o arquivo precisa de formatação.

Na sessão interativa, digite um **programa completo**. A linha com
`fimalgoritmo` inicia a execução. Se o programa usar `leia`, digite os
valores solicitados em seguida. Para encerrar a sessão, digite `:sair`
em uma linha separada.

Está experimentando laços de repetição? Você pode limitar o número de
passos da execução:

```sh
./portugo run --max-steps 10000 boas-vindas.alg
```

O programa para e exibe uma mensagem se atingir esse limite. Você também
pode pressionar Ctrl+C para interromper um comando em execução.

## Continue explorando

A [pasta de exemplos](examples/) é um bom lugar para escolher seu próximo
exercício:

| Exemplo | O que explorar |
| --- | --- |
| [Olá, Portugol](examples/hello.alg) | Um primeiro programa completo |
| [Fatorial](examples/fatorial.alg) | Entrada de dados, variáveis e um laço `para`; experimente digitar `5` |
| [Entrada pelo terminal](examples/console_input.alg) | Leitura de nome e idade, uma linha de entrada por vez |
| [Condições lógicas](examples/logical_conditions.alg) | Comparações e expressões lógicas |
| [Vetores](examples/vector_bounds.alg) | Declaração e acesso aos elementos de vetores |

Execute qualquer exemplo da mesma forma:

```sh
./portugo run examples/fatorial.alg
```

Já tem exercícios salvos no VisuAlg? Experimente executá-los com
`./portugo run path/to/exercise.alg`. Substitua `path/to/exercise.alg` pelo
caminho do seu arquivo. O PortuGo lê arquivos-fonte em UTF-8 e Windows-1252,
incluindo a codificação normalmente usada pelo VisuAlg.

## Compatibilidade e estado do projeto

O PortuGo usa o dialeto VisuAlg 3.x, com comportamento verificado por
comparação com o VisuAlg 3.0.7. Ele aceita variáveis, condições, laços de
repetição, procedimentos, funções, vetores e funções predefinidas para
números e texto.

O [relatório de conformidade](docs/quality-acceptance-2026-09-14.md), em inglês,
registra a aprovação nas verificações de todos os programas aceitos do
conjunto de referência. Alguns casos de programas rejeitados e comportamentos
fora desse conjunto ainda não foram verificados. Por isso, esses resultados
não garantem compatibilidade com todo programa escrito para VisuAlg.

O PortuGo é um interpretador de linha de comando. Ele não reproduz a
interface do VisuAlg nem seu depurador gráfico. O Portugol Studio e outros
dialetos de Portugol usam sintaxes diferentes e estão fora do escopo atual.

A [referência da linguagem](docs/language.md), em inglês, descreve a sintaxe
aceita, o comportamento e as limitações conhecidas. O
[histórico de alterações](CHANGELOG.md), também em inglês, registra as mudanças.

## Dúvidas, ideias e contribuições

Encontrou uma instrução confusa ou um exercício que não funciona como
esperava? [Abra uma issue](https://github.com/ncode/PortuGo/issues) com um
pequeno exemplo `.alg`, o resultado esperado e o que aconteceu. Dúvidas e
melhorias na documentação também são bem-vindas.

Para contribuir com o código, comece pelo
[guia de desenvolvimento](docs/development.md) e pelas
[instruções para contribuir](AGENTS.md), ambos em inglês.
Você pode executar os testes na pasta do projeto com:

```sh
go test ./...
```

Bons estudos!
