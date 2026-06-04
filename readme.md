# ⚔️ Arena RPG

> Uma API REST de combate de RPG, escrita em **Go** — projeto de fim de semana pra conhecer e brincar com a linguagem.

## Sobre o projeto

Este é um projeto **de aprendizado**, sem pretensão de produção. A ideia foi sair do zero em Go construindo algo divertido e "completinho": uma API onde você cria personagens de RPG, rola dados e faz dois personagens se enfrentarem numa arena, ganhando XP e subindo de nível.

O foco nunca foi o produto final, e sim **passar pelos conceitos** da linguagem na prática — sintaxe, arquitetura, ecossistema e ferramentas — fazendo cada peça funcionar antes de seguir pra próxima.

## Stack

- **Go** (1.26)
- **`net/http`** da biblioteca padrão — roteamento nativo com casamento de método + caminho (recurso do Go 1.22+), sem framework externo
- **SQLite** via [`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite) — implementação em Go puro, sem necessidade de compilador C (CGO)

## O que a API faz

| Método | Rota                  | O que faz                                            |
| ------ | --------------------- | ---------------------------------------------------- |
| GET    | `/ping`               | Health check — responde `{"status":"ok"}`            |
| POST   | `/personagens`        | Cria um personagem                                   |
| GET    | `/personagens`        | Lista todos os personagens                           |
| GET    | `/personagens/{id}`   | Busca um personagem pelo ID                          |
| POST   | `/rolar`              | Rola dados no formato `NdX` (ex: `3d6`)              |
| POST   | `/combate`            | Resolve uma luta entre dois personagens, turno a turno; o vencedor ganha XP |

## Como rodar

Pré-requisito: Go 1.22 ou superior.

```bash
go run .
```

O servidor sobe em `http://localhost:8080` e cria o arquivo `arena.db` (SQLite) automaticamente.

### Exemplos

```bash
# criar personagens
curl -X POST localhost:8080/personagens \
  -d '{"nome":"Aragorn","classe":"guerreiro","hp":100,"ataque":20,"defesa":10}'

curl -X POST localhost:8080/personagens \
  -d '{"nome":"Gimli","classe":"guerreiro","hp":100,"ataque":10,"defesa":30}'

# listar
curl localhost:8080/personagens

# buscar por id
curl localhost:8080/personagens/1

# rolar dados
curl -X POST localhost:8080/rolar -d '{"dados":"3d6"}'

# combate (o vencedor ganha XP, e a cada 100 XP sobe de nível)
curl -X POST localhost:8080/combate -d '{"atacante_id":1,"defensor_id":2}'
```

## Estrutura do projeto

```
arena_rpg/
├── go.mod
├── arena.db                  # banco SQLite (gerado em runtime, fora do git)
├── main.go                   # ponto de entrada: liga o banco e registra as rotas
└── internal/
    ├── personagem/           # tudo de personagem: modelo, persistência e handlers
    │   ├── personagem.go      #   struct + acesso ao banco (Inserir, Listar, Buscar, DarXP)
    │   └── handler.go         #   handlers HTTP
    ├── dados/                # rolagem de dados
    │   ├── dados.go           #   lógica pura (sem HTTP)
    │   └── handler.go         #   handler HTTP
    └── combate/              # combate
        ├── combate.go         #   lógica pura da luta (sem HTTP)
        └── handler.go         #   handler HTTP
```

A ideia de arquitetura que guiou tudo: **separar as camadas por responsabilidade.**

- O **handler** fala HTTP (lê a requisição, escreve a resposta).
- A **lógica pura** (`combate.go`, `dados.go`) não sabe o que é HTTP — só resolve o domínio.
- A **persistência** (`personagem.go`) cuida de guardar e buscar.

Essa separação é o que permitiu trocar o armazenamento em memória por SQLite sem reescrever a lógica de combate.

## O que aprendi (e por quê)

Esta foi a parte mais valiosa. Cada conceito tem o motivo de existir:

- **Pacotes (`package`)** — em Go, a pasta é a unidade de pacote. Aprendi que pastas diferentes são pacotes diferentes, e que pra um usar o outro é preciso **exportar** (nomes com maiúscula = públicos, minúscula = privados), **importar** pelo caminho do módulo e **chamar qualificado** (`personagem.Buscar`). _Por quê:_ é a base da organização e do encapsulamento em Go.

- **A pasta `internal/`** — pacotes ali dentro só podem ser importados pelo próprio projeto. _Por quê:_ privacidade garantida pela linguagem, evitando que o "código de dentro" vire API pública sem querer.

- **Structs e tags de JSON** (`` `json:"nome"` ``) — modelar dados tipados e mapear pro JSON. _Por quê:_ campos Go começam com maiúscula pra serem exportados, mas o JSON usa minúsculo; as tags fazem a ponte.

- **Tratamento de erro explícito** (`if err != nil`) — Go não tem `try/catch`; o erro é um valor que você checa na hora. _Por quê:_ deixa o caminho de falha visível e obrigatório, em vez de escondido.

- **`map` + o idioma "comma, ok"** (`p, ok := mapa[chave]`) — acessar um map devolve o valor e um booleano de "existe?". _Por quê:_ é o jeito idiomático de checar existência sem ser enganado pelo valor zero.

- **Valor zero útil** — um `sync.Mutex` já nasce pronto; uma `string` vazia é `""`; um slice `nil` vira `null` no JSON. _Por quê:_ Go desenha os tipos pra que o "zero" já sirva, poupando inicialização.

- **Ponteiros e _pointer receivers_** (`func (p *Personagem) GanharXP(...)`) — métodos com receiver de ponteiro alteram o original, não uma cópia. _Por quê:_ foi o que permitiu o XP "grudar" no personagem guardado (combinado com regravar no armazenamento).

- **Roteamento nativo método + caminho** (`"POST /personagens"`, `r.PathValue("id")`) — recurso do Go 1.22+. _Por quê:_ dá pra fazer uma API limpa só com a stdlib, sem framework.

- **Concorrência** — cada requisição HTTP roda numa goroutine própria; por isso o map em memória precisava de `sync.Mutex`. _Por quê:_ entender que o servidor atende em paralelo evita bugs sutis de "concurrent map writes".

- **`database/sql` + driver + blank import** (`_ "modernc.org/sqlite"`) — a stdlib define a interface, o driver implementa e se registra sozinho. Queries com `?` (placeholders) blindam contra SQL injection; `Scan(&campo)` lê colunas pra variáveis. _Por quê:_ é o jeito padrão e seguro de falar com banco em Go — e com o SQLite, o mutex manual pôde ser removido, pois o `*sql.DB` já é seguro pra concorrência.

- **Separação lógica pura × HTTP** — a regra do domínio mora em funções/pacotes que não conhecem HTTP. _Por quê:_ código testável, reutilizável e fácil de trocar a "casca" (HTTP, CLI, etc.) sem mexer no miolo.

## Próximos passos

Ideias pra evoluir o projeto no futuro, caso valha a pena retomar:

**Organização e qualidade**
- Refatorar respostas JSON num helper compartilhado (em vez de setar o header em cada handler).
- Adicionar **middlewares** para preocupações transversais: log de requisições, CORS, recuperação de panics.
- Separar as regras de validação dos handlers, devolvendo erros de domínio que o handler traduz em status HTTP.
- Escrever **testes automatizados** (pacote `testing`, table-driven tests) — especialmente da lógica de combate.

**Funcionalidades**
- `PUT` e `DELETE` de personagens.
- Inventário/armas de verdade (a lista de armas precisaria ser guardada como JSON numa coluna ou numa tabela separada).
- Persistir o **histórico de batalhas** numa tabela própria.
- Calibrar a fórmula de combate (hoje defesa alta domina demais e a luta arrasta).

**Arquitetura e infra**
- Introduzir uma **interface de repositório**, pra trocar SQLite por Postgres sem tocar na lógica.
- Configuração via variáveis de ambiente (porta, caminho do banco).
- Autenticação de usuários.
- **Dockerizar** a aplicação.

**Indo além (rumo ao MMO)**
- Um **frontend** consumindo a API.
- **WebSockets** para combate em tempo real entre jogadores.

## Notas de decisão

- **Memória antes de banco:** os primeiros cards guardavam tudo num `map` em RAM; o SQLite só entrou depois que a lógica estava redonda. Debugar regra de negócio e SQL ao mesmo tempo seria furada.
- **Funcionando antes de bonito:** a organização em pacotes e os refinamentos vieram conforme a necessidade, não de cara.
- **`modernc.org/sqlite` em vez de `mattn/go-sqlite3`:** o primeiro é Go puro e dispensa compilador C, o que evita dor de cabeça com CGO (especialmente rodando em WSL).

---

_Projeto de fim de semana feito para aprender Go. Construído card a card, do `go mod init` ao banco persistente._