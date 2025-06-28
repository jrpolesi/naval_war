Padrao de eventos

prefixo "client\_"

são mesagens que o cliente envia para o servidor

prefixo "server\_"
são mensagens que o servidor envia para o cliente

Endpoint de conexão WebSocket
`ws://localhost:8080/websocket`

Fluxo de mensagens

1º

```json
{
  "event": "server_created_new_player",
  "payload": {
    "id": "9cf454bd-546a-488d-bbf4-07fe4787aa9f",
    "name": "",
    "is_ready": false
  }
}
```

2º

```json
{
  "event": "client_updated_player_info",
  "payload": {
    "name": "My Name",
    "is_ready": false
  }
}
```

3º

```json
{
  "event": "server_updated_players_info",
  "payload": {
    "teams": [
      {
        "id": "447ba043-d4a5-47ac-b703-b0e0dd99b010",
        "players": [
          {
            "id": "47911c41-7eba-4e5c-89c4-c2f3badd646f",
            "name": "My Name25",
            "is_ready": true
          }
        ]
      },
      {
        "id": "4948c34d-45ee-4b75-8635-65df2ed1317e",
        "players": [
          {
            "id": "fa57ac85-240f-4185-9b31-05f1a9d39fbf",
            "name": "My Name2",
            "is_ready": true
          }
        ]
      }
    ]
  }
}
```

4º

```json
{
  "event": "server_game_started",
  "payload": {
    "map": {
      "size": [10, 10]
    },
    "teams": [
      {
        "id": "447ba043-d4a5-47ac-b703-b0e0dd99b010",
        "players": [
          {
            "id": "47911c41-7eba-4e5c-89c4-c2f3badd646f",
            "name": "My Name25",
            "is_ready": true
          }
        ]
      },
      {
        "id": "4948c34d-45ee-4b75-8635-65df2ed1317e",
        "players": [
          {
            "id": "fa57ac85-240f-4185-9b31-05f1a9d39fbf",
            "name": "My Name2",
            "is_ready": true
          }
        ]
      }
    ],
    "ships": [
      {
        "id": "53c1bfef-a5c5-4fc5-941d-55133dd356e1",
        "position": {
          "x": 8,
          "y": 3
        },
        "is_damaged": false,
        "is_owner": true
      },
      {
        "id": "f955b0c9-f5bf-4df6-9172-528e3230cddb",
        "position": {
          "x": 8,
          "y": 6
        },
        "is_damaged": false,
        "is_owner": true
      },
      {
        "id": "a83a6efc-95dd-49b0-b10a-ff3db2b8ef2b",
        "position": {
          "x": 1,
          "y": 6
        },
        "is_damaged": false,
        "is_owner": true
      }
    ]
  }
}
```

5º

```json
{
  "event": "client_player_performed_action",
  "payload": {
    "type": "attack",
    "payload": {
      "position": {
        "x": 8,
        "y": 4
      }
    }
  }
}
```
