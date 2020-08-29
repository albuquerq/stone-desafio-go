# Stone Desafio Go


Implementa o [desafio de Go da stone](https://gist.github.com/guilhermebr/fb0d5896d76634703d385a4c68b730d8). Foi utilizado o PostgreSQL como banco de dados.

# API

Implementa uma REST API JSON com rotas autenticadas por meio de JWT passadas pelo cabeçalho `Authorization: Bearer <jwt-token>`.

Método | Rota | JWT | Função
--:|--|:--:|--
POST | /api/v1/login | - | Autentica utilizando CPF e Secret da conta, retornando JWT para acesso de rotas privadas.
GET  | /api/v1/accounts | - | Retorna todas contas cadastradas.
POST | /api/v1/accounts | - | Cria uma nova conta.
GET  | /api/v1/accounts/{account_id}/balance | - | Retorna o saldo da conta.
GET  | /api/v1/transfers | Sim | Retorna a lista de transferências em que a conta serviu como origem da operação.
POST | /api/v1/transfers | Sim | Faz uma transferência da conta autenticada para uma conta informada.


## Tipos de recursos

Valores retornados são embutidos no tipo Response:

```
Response {
    value <Object|Array>
    error ErrorObject // Omitido quando não ocorre.
}

ErrorObject {
    code int
    message string
    go_err <string|Object> // Campo temporário para fins depuração.
}
```

Os objetos de manipulação de contas, seguem os formatos dependendo da operação.

Rota *POST /api/v1/login*

```
// Request body
Credential {
    cpf string
    secret string
}

// Reponse body
Response {
    value: {
        token string // JWT token usado no Authorization header.
    }
}
```

Rota *GET /api/v1/accounts*

```
// Response body
Respose {
    value: [
        Account {
            id uuid
            name string
            cpf string
            balance integer // Representa um valor monetário em centavos de real brasileiro (BRL)
            created_at Date
        }
    ]  
}
```

Rota *POST /api/v1/accounts*
```
// Request body
Account {
    name string
    cpf string
    balance integer // Saldo em centavos de real brasileiro (BRL)
    secret string
}

// Response body
Response {
    value: Account {
        id uuid
        name string
        cpf string
        balance integer // Representa um valor monetário em centavos de real brasileiro (BRL)
        created_at Date
    }
}
```

Rota *GET /api/v1/accounts/{account_id}/balance*
```
Balance {
    balance integer // Saldo em centavos de real brasileiro (BRL)
}
```

Rota *GET /api/v1/transfers*

Header: `Authorization: Bearer <token>`

```
Response {
    value: [
        Transfer {
            id uuid
            account_origin_id uuid
            account_destination_id uuid
            amount integer // Em centavos de real brasileiro (BRL)
            created_at Date
        }
    ]
}
```

Rota *POST /api/v1/transfers*

Header: `Authorization: Bearer <token>`
```
// Request body
Transfer {
    account_destination_id uuid
    amount integer // Saldo em centavos de ral brasileiro (BRL)
}

// Reponse body
Response {
    value: Transfer {
        id uuid
        account_origin_id uuid
        account_destination_id uuid
        amount integer
        created_at Date
    }
}
```


# Como executar

Esse projeto depende da ferramenta [migrate](https://github.com/golang-migrate/migrate) para execução das migrações do banco de dados. Um ambiente docker de execução simples foi elaborado com docker-compose, para executá-lo basta utilizar o arquivo de composição localizado na pasta `deployments`. Assim, clone o repositório e execute os comandos:


```bash
> cd ./deployments
> docker-compose up
 ```

        Considerações, docker compose não sincroniza os serviços de forma que esses realizem suas operações na sequência orquestrada para manter suas integridades. 
        Desse forma, o serviços de migração e o de api pode não serem executados no tempo correto, a migração antes da utilização do banco de dados, para o correto funcionamento. Para diminuir essa possibilidade o serviço de API utiliza um tempo de espera, para que o serviço de migração do banco de dados possa executar anteriormente.