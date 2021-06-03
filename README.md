# Uploads API

Faz uploads de arquivos para algum backend (S3, GCS) e depois, você pode consulta-los.

# Projeto

No geral segui alguns padrões da comunidade, o mesmo que projetos como Docker e K8S vem usando.
[Aqui](https://github.com/golang-standards/project-layout) você pode encontrar uma referência sobre eles.

# Arquitetura

Tentei seguir uma arquitetura baseada no Clean Arch. Fiz algumas adaptações pela simplicidade do projeto, mas no geral  
segui os padrões descritos [aqui](https://github.com/bxcodec/go-clean-arch)

# Testes unitários

Na raiz do projeto, digite no console:

```
make tests
```

Esses comando irá gerar um container e executar os testes dentro dele.
Não tentei atingir 100% de code coverage, tentei cobrir aquilo que acreditava ser parte "do negócio".

# Executando a API

Na raiz do projeto, digite no console:

```
make start
```

Esse comando irá subir os containers do S3, GCS e da API. Quando a API estiver rodando, estará tudo pronto para executar as requisições.

# Configuração

Todas as configurações estão em variáveis de ambiente. Para trocar o backend (S3/GCS) você pode os valores da variável `STORAGE_BACKEND` no arquivo `.env`.
Os possíveis valores são:

- s3
- gcs

Qualquer coisa diferente disso irá lançar um erro.

# Postman

Para facilitar os testes, criei uma collection no Postman que pode ser encontrada aqui:  
[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/6fdabf390c205b524b49?action=collection%2Fimport)
