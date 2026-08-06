# Desafio Fullstack - Mini Kanban Veritas

Olá! Bem-vindo(a) ao meu repositório. Este é o projeto que desenvolvi para o desafio técnico da Veritas Consultoria Empresarial. 

A ideia aqui foi construir um Kanban simples, mas robusto, com três colunas fixas ("A Fazer", "Em Progresso" e "Concluídas"). Foquei bastante em entregar não só o que foi pedido, mas também uma experiência fluida para o usuário final e uma arquitetura limpa para quem for ler o código.

---

## Como rodar o projeto na sua máquina

Para facilitar a avaliação, eu coloquei toda a aplicação em containers usando o **Docker**. Isso significa que você não precisa instalar o Go, o Node ou configurar portas manualmente.

**Passo a passo:**
1. Clone este repositório para o seu computador.

2. Abra o terminal na pasta raiz do projeto.

3. Execute o comando abaixo para construir e subir as imagens:

   ```bash
   docker compose up --build

4. Aguarde a mensagem de sucesso no terminal e acesse no seu navegador:

Interface do Kanban: http://localhost:3000

API do Backend: http://localhost:8080/tasks

## Minhas Decisões Técnicas
Durante o desenvolvimento, precisei tomar algumas decisões de arquitetura e design. Aqui explico o "porquê" de cada uma:

Arquitetura em Containers Separados: Decidi criar um Dockerfile para o Go e outro para o React (usando Nginx), orquestrando tudo com o docker-compose. Fiz isso para mostrar um ambiente mais próximo do real, mantendo as responsabilidades totalmente separadas.

Persistência em JSON com Segurança: O desafio pedia armazenamento em memória com bônus para persistência. Usei o arquivo tasks.json. Para evitar dor de cabeça com leitura/escrita simultânea, implementei um sync.Mutex no Go, bloqueando o arquivo enquanto uma ação acontece.

Drag and Drop sem Bugs (Otimismo na UI): No React, usei a biblioteca @hello-pangea/dnd. A mágica aqui é a "Optimistic UI": quando você arrasta um card, a interface atualiza instantaneamente para você, enquanto o Axios avisa o backend de forma silenciosa por trás.

Modo Escuro e Flexbox: Pensei no avaliador que vai testar isso à noite! Criei um toggle de Modo Escuro/Claro e usei Flexbox para garantir que o board ocupe a tela toda no desktop e role horizontalmente no celular (estilo Trello).

Segurança Básica: Criei um middleware no backend não apenas para fazer logs bonitos no terminal, mas para injetar cabeçalhos defensivos na resposta (X-Frame-Options e X-Content-Type-Options).

## Documentação Visual (UML)
Abaixo estão os diagramas representando o comportamento do sistema. O GitHub renderiza esses diagramas automaticamente!

1. User Flow (Como o usuário usa o sistema)

stateDiagram-v2
    [*] --> Visualizar_Board: Acessa a aplicação
    Visualizar_Board --> Adicionar_Tarefa: Clica em "+ Adicionar"
    Adicionar_Tarefa --> Visualizar_Board: Preenche título e salva
    
    Visualizar_Board --> Editar_Tarefa: Clica no ícone de Lápis
    Editar_Tarefa --> Visualizar_Board: Modifica e salva
    
    Visualizar_Board --> Mover_Tarefa: Segura e arrasta o card (Drag & Drop)
    Mover_Tarefa --> Visualizar_Board: Solta na nova coluna
    
    Visualizar_Board --> Excluir_Tarefa: Clica no ícone de Lixeira
    Excluir_Tarefa --> Visualizar_Board: Confirma a exclusão

2. Data Flow (Como os dados viajam por baixo dos panos) - Bônus

sequenceDiagram
    participant U as Usuário
    participant F as Frontend (React)
    participant B as Backend (Go API)
    participant J as tasks.json

    U->>F: Interage com a tela (ex: mover tarefa)
    F->>B: Requisição HTTP (PUT /tasks/{id})
    activate B
    B->>B: Lock no Mutex (Segurança de concorrência)
    B->>J: Atualiza o arquivo
    J-->>B: Confirmação de gravação
    B->>B: Unlock no Mutex
    B-->>F: Retorna JSON (Status 200 OK)
    deactivate B
    F-->>U: Atualiza a interface visual