# TP Go - API météo

## Lancer le serveur

se mettre dans le dossier server et faire :

```
go run .
```

le serveur tourne sur le port 8080

## Les routes

- GET /stations → retourne toutes les stations (200)
- GET /stations/{id} → retourne une station (200) ou 404 si elle existe pas
- POST /stations → crée une station (201), 409 si doublon, 400 si JSON invalide
- PUT /stations/{id} → remplace une station (200) ou la crée (201)
- DELETE /stations/{id} → supprime une station (204) ou 404
- GET /stations/{id}/observations → retourne les observations d'une station (200)

## Postman

la collection utilisée est EFREI Golang J3 — API REST météo.postman_collection.json
