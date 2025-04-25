# Documentation Client KVLT

## Vue d'ensemble

KVLT est une base de données clé-valeur offrant un stockage de données simple et performant. Ce client facilite les interactions avec la base de données KVLT via des requêtes HTTP.

>**Remarque :** L'utilisation de ce client n'est pas obligatoire. Vous pouvez interagir avec l'API KVLT en utilisant n'importe quel client HTTP standard.

## Installation

Installez la bibliothèque requise :

```bash
npm install @duplojs/http-client
```

> **Information :** Si vous voulez plus d'informations sur le client HTTP, vous pouvez consulter la [documentation officielle](https://docs.duplojs.dev/fr/latest/resources/http-client/)

## Guide d'utilisation

### Configuration du client

```typescript
import { HttpClient, type TransformCodegenRouteToHttpClientRoute } from "@duplojs/http-client";
import type { kvltRoutes } from "./kvltClientRoutes.d.ts";

// Définition du type pour les routes du client
export type KvltClientRoute = TransformCodegenRouteToHttpClientRoute<
    kvltRoutes
>;

// Initialisation du client
const kvltClient = new HttpClient<KvltClientRoute>({
    baseUrl: "http://localhost:8080", // Remplacez par l'URL de votre serveur KVLT
});
```

### Exemples d'opérations

#### Insertion d'une valeur

```typescript
const insertExemple = await kvltClient
    .put(
        "/value",
        {
            body: {
                key: "test",
                value: {
                    firstName: "John",
                    lastName: "Doe",
                },
                duration: 1000, // Durée de vie en millisecondes
            },
        },
    ).iWantInformation("success.keySet");
```

#### Récupération d'une valeur

```typescript
const fetchExemple = await kvltClient
    .get(
        "/value/{key}",
        {
            params: {
                key: "test",
            },
        },
    ).iWantInformation("success.keyFound");
```

#### Remerciements

Merci d'avoir utilisé KVLT ! Nous espérons que ce client vous facilitera la vie lors de vos interactions avec notre base de données. Si vous avez des questions ou des suggestions, n'hésitez pas à nous contacter.

Merci à l'équipe de développement de [DuploJS](https://duplojs.dev/) pour leur excellent travail sur le client HTTP.
