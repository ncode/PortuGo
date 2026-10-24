algoritmo "record observations"
tipo Person = registro
name: caractere
age: inteiro
fimregistro
Empty = registro
fimregistro
Repeated = registro
code: inteiro
code: caractere
fimregistro
var person: Person
people: vetor[1..2] de Person
empty: Empty
repeated: Repeated
inicio
person.name <- "Ana"
person.age <- 21
people[1] <- person
person.name <- "Bia"
people[2] <- person
repeated.code <- 7
fimalgoritmo
