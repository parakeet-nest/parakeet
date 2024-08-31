# Comment rendre un SLM plus malin

Dans ce tutoriel, nous allons voir comment rendre un SLM plus malin en utilisant des données de contexte.

> J'utiliserais mon projet [Parakeet](https://github.com/parakeet-nest/parakeet) pour illustrer ce tutoriel. Mais vous pouvez facilement adapter les concepts avec d'autres frameworks comme LangChain.

## SLM ?

Je vous donne ma définition de SLM ou Small Language Model> C'est donc un LLM, qui est "petit" (voire même très petit) par rapport à un modèle de langage complet, qui est capable de générer du texte. En voici quelques uns :

- Tinyllama (637 MB)
- Tinydolphin (636 MB)
- Gemma:2b (1.7 GB) et Gemma2:2b (1.6 GB)
- Phi3:mini (2.2 GB) et Phi3.5 (2.2 GB)
- Qwen:0.5b (394 MB), Qwen2:0.5b (352 MB) et Qwen2:1.5b (934 MB)
- ...

Ma préférence va à des modèles capable de s'exécuter confortablement sur un **Raspberry Pi 5** avec 8Go de RAM, donc sans GPU. J'aurais tendance à dire que les modèles de moins de 1Go sont les plus adaptés (et ont ma préférence).

Je maintiens une liste de modèles que j'ai testé et qui fonctionnent bien sur un Raspberry Pi 5 avec 8Go de RAM. Vous pouvez la consulter ici : [Awesome SLMs](https://github.com/parakeet-nest/awesome-slms).

> J'utilise le projet [Ollama]() pour charger et exécuter les modèles.

## Mon objectif

Aujjourd'hui, mon objectif est de tenter de transformer le modèle **`Qwen:0.5b`** en un spécialiste des fougères. Pour cela, je vais lui fournir des données de contexte sur les fougères, et je vais lui poser diverses questions.

## Prérequis:

### Vérifier que le modèle ou les modèles que je vais utiliser n'y connaissent rien en fougères

Pour vérifier que le modèle **`Qwen2:0.5b`** ne sait rien sur les fougères, je vais lui poser une question sur les fougères. Pour cela, j'utilise le script suivant:

```bash
curl http://localhost:11434/api/chat \
-H "Content-Type: application/json" \
-d '
{
  "model": "qwen2:0.5b",
  "messages": [
    {
      "role": "system",
      "content": "You are an expert in botanics and ferns."
    },
    {
      "role": "user",
      "content": "Give me a list of ferns of the Dryopteridaceae variety"
    }
  ],
  "stream": false
}' | jq '.message.content'
```

✋ **Tout au long de mes expérimentation** j'ai fait des tests avec plusieurs modèles pour voir lequel est le plus adapté à mes exigences. Notamment celle de pouvoir fonctionner confortablement sur un Pi 5. Mon choix définitif s'est porté sur le modèle **`Qwen2:0.5b`**. Donc je ne vous présenterai que les résultats obtenus avec ce modèle. Mais rien ne vous empêche de tester avec d'autres modèles. Je vous encourage même à le faire. Je proposerais aussi quelaues hypothèses et conclusions pour expliquer les choix et les résultats obtenus. Cela devrat ensuite vous permettre de faire vos propres expérimentations et les adapter à vos besoins.

### La préparation des données

Pour commencer, je dois générer un fichier de données sur les fougères qui serq mon unique source de vérité. Pour cela, j'ai créé un fichier `ferns.json` en demandant de l'aide à **ChatGPT4o**. Pour cela j'ai utilisé le prompt suivant:

```
You are an expert in botanics and your main topic is about the ferns.

I want a list of 10 varieties of ferns with their characteristics and 5 ferns per variety.

the output format is in JSON:
```json
[
    {
        "variety": "name of the variety of ferns",
        "description": "Descriptionof the variety of ferns",
        "ferns": [
            {
                "name": "scientifict name of the fern",
                "common_ame": "common name of the fern",
                "description": "description and characteristics of the fern"
            },
        ]
    },
]
```

Voici un extrait de ce fichier:

```json
[
    {
        "variety": "Polypodiaceae",
        "description": "A family of ferns known for their leathery fronds and widespread distribution in tropical and subtropical regions.",
        "ferns": [
            {
                "name": "Polypodium vulgare",
                "common_name": "Common Polypody",
                "description": "A hardy fern with leathery, evergreen fronds that thrive in rocky and shaded areas."
            },
            {
                "name": "Polypodium glycyrrhiza",
                "common_name": "Licorice Fern",
                "description": "Known for its sweet-tasting roots, this fern grows on moist, shaded rocks and tree trunks."
            },
```

📝 vous pouvez trouver le fichier complet ici: [ferns.json](https://github.com/parakeet-nest/parakeet/tree/main/examples/42-make-a-slm-smarter/data/ferns.json)

### Mise en forme des données

Ensuite, j'ai transformé le fichier JSON en deux fichiers markdown avec les même données mais avec des structures différentes. Voici les structures de ces fichiers:

#### `ferns.1.md`

```markdown
# Title of the report

## Variety: name of the variety of ferns
*Description:*
Description of the variety of ferns

### Ferns:
- **name:** **scientific name of the fern**
  **common name:** common name of the fern
  **description:** description and characteristics of the fern
```

#### `ferns.2.md`

```markdown
# Title of the report

## Variety: name of the variety of ferns
*Description:*
Description of the variety of ferns

### Ferns:

#### Name of the fern

**name:** **scientific name of the fern**
**common name:** common name of the fern
**description:** description and characteristics of the fern
```

> Le second fichier est un peu plus structuré que le premier grâce aux titres de paragraphes. Cela me permettra de tester si la structure des données a un impact sur les résultats obtenus.

📝 Vous pouvez trouver les fichiers complets ici: [ferns.1.md](https://github.com/parakeet-nest/parakeet/tree/main/examples/42-make-a-slm-smarter/data/ferns.1.md) et [ferns.2.md](https://github.com/parakeet-nest/parakeet/tree/main/examples/42-make-a-slm-smarter/data/ferns.2.md)

## 1ères expérimentations

J'ai décidé d'aller au plus simple est de créer un prompt qui sera composé des éléments suivants:

- Les instructions pour le modèle.
- Le contenu**complet** du fichier`ferns.1.md` ou`ferns.2.md`, ce que j'appelle le**contexte**.
- La question que je vais poser au modèle.

En Go avec **Parakeet**, cela donne quelque chose comme ça:

```golang
Messages: []llm.Message{
    {Role: "system", Content: systemContent},
    {Role: "system", Content: contextContext},
    {Role: "user", Content: question},
},
```

Et les instructions pour le modèle:

```golang
systemContent := `**Instruction:**
You are an expert in botanics.
Please use only the content provided below to answer the question.
Do not add any external knowledge or assumptions.`
```

📝 Vous pouvez trouver le code complet ici: [01-context](https://github.com/parakeet-nest/parakeet/tree/main/examples/42-make-a-slm-smarter/01-context)

### Questions

Pour faire mes tests, j'ai utilisé les questions suivantes:

- **Question 1:** Give me a list of ferns of the Dryopteridaceae variety
- **Question 2:** What is the common name Dryopteris cristata?

Pour cette première série de tests, j'ai utilisé trois modèles **`Qwen:0.5b`**, **`Qwen2:0.5b`** et **`CognitiveComputations/dolphin-gemma2`** (1.6 GB).

Les résultats obtenus sont les suivants:

#### Avec `ferns.1.md`

| LLM + ferns.1.md | Question 1 | Question 2 |
| ---------------- | ---------- | ---------- |
| qwen:0.5b        | 😡         | 😡         |
| qwen2:0.5b       | 😡         | 😡         |
| dolphin-gemma2   | 🙂         | 🙂         |

> - 😡: résultat faux ou incomplet
> - 🙂: résultat satisfaisant

#### Avec `ferns.2.md`

| LLM + ferns.2.md | Question 1 | Question 2 |
| ---------------- | ---------- | ---------- |
| qwen:0.5b        | 😡         | 😡         |
| qwen2:0.5b       | 😡         | 😡         |
| dolphin-gemma2   | 😡         | 🙂         |

> - 😡: résultat faux ou incomplet
> - 🙂: résultat satisfaisant

### Hypothèses et Observations

- **dolphin-gemma2** a déjà des connaissances sur les fougères, donc il a pu répondre correctement à la question en s'appuyant sur ses connaissances et les informations fournis (même si théoriquement les instructions étaient de n'utiliser que le contexte fourni).
- **dolphin-gemma2** semble être capable de s'en sortir avec un contexte de grande taille.
- *qwen:0.5b** et**qwen2:0.5b** n'ont pas pu répondre correctement aux questiona. Cela peut être dû à sa capacité à traiter des contextes plus grands.
- La structure des données semble avoir un impact sur les résultats obtenus avec**dolphin-gemma2** - mais pourtant j'aurais pensé que la seconde structure de document serait plus facilement exploitable. Nous verrons si cela se confirme par la suite.

Ma première hypothèse ou constatation serait que **qwen:0.5b** et **qwen2:0.5b** (donc des trés petits LLMs) ne sont pas capable d'exploiter un très grand contexte, même si il est structurée.

Il faut donc aider **qwen:0.5b** et **qwen2:0.5b** à "se concentrer" sur des données précises et "plus proches" des questions. J'ai donc reproduit les mêmes expérimentations mais en réduisant la taille du contexte.

## 2ème série d'expérimentations: réduire la taille du contexte

Pour cette seconde série d'expérimentations, j'ai réduit la taille du contexte en ne gardant que les informations sur une seule variété de fougères. J'ai donc créé un fichier `ferns.1.extract.md` et `ferns.2.extract.md` qui contiennent les informations sur la variété **Dryopteridaceae**.

📝 J'utilise le même code pour exécuter les exemples : [01-context](https://github.com/parakeet-nest/parakeet/tree/main/examples/42-make-a-slm-smarter/01-context)

**Pour `ferns.1.extract.md`**:

```markdown
# Fern Varieties Report

## Variety: Dryopteridaceae
*Description:*
A large family of ferns with robust and often leathery fronds, commonly found in woodlands.

### Ferns:
- **name:** **Dryopteris filix-mas**
  **common name:** Male Fern
  **description:** A sturdy fern with pinnate fronds, commonly found in temperate forests.

- **name:** **Dryopteris marginalis**
  **common name:** Marginal Wood Fern
  **description:** Known for its evergreen fronds, this fern thrives in rocky, shaded environments.

- **name:** **Dryopteris erythrosora**
  **common name:** Autumn Fern
  **description:** This fern features striking copper-red fronds that mature to green.

- **name:** **Dryopteris cristata**
  **common name:** Crested Wood Fern
  **description:** A fern with uniquely crested fronds, typically found in wetland areas.

- **name:** **Dryopteris affinis**
  **common name:** Golden Male Fern
  **description:** A robust fern with yellowish fronds and a preference for moist, shaded habitats.
```

**Pour `ferns.2.extract.md`**:

```markdown
# Fern Varieties Report

## Variety: Dryopteridaceae
*Description:*
A large family of ferns with robust and often leathery fronds, commonly found in woodlands.

### Ferns:

#### Dryopteris filix-mas

**name:** **Dryopteris filix-mas**  
**common name:** Male Fern  
**description:** A sturdy fern with pinnate fronds, commonly found in temperate forests.

#### Dryopteris marginalis

**name:** **Dryopteris marginalis**  
**common name:** Marginal Wood Fern  
**description:** Known for its evergreen fronds, this fern thrives in rocky, shaded environments.

#### Dryopteris erythrosora

**name:** **Dryopteris erythrosora**  
**common name:** Autumn Fern  
**description:** This fern features striking copper-red fronds that mature to green.

#### Dryopteris cristata

**name:** **Dryopteris cristata**  
**common name:** Crested Wood Fern  
**description:** A fern with uniquely crested fronds, typically found in wetland areas.

#### Dryopteris affinis

**name:** **Dryopteris affinis**  
**common name:** Golden Male Fern  
**description:** A robust fern with yellowish fronds and a preference for moist, shaded habitats.
```

📝 Vous pouvez trouver les fichiers complets ici: [ferns.1.extract.md](https://github.com/parakeet-nest/parakeet/tree/main/examples/42-make-a-slm-smarter/data/ferns.1.extract.md) et [ferns.2.extract.md](https://github.com/parakeet-nest/parakeet/tree/main/examples/42-make-a-slm-smarter/data/ferns.2.extract.md)

J'ai ensuite répété les mêmes expérimentations que précédemment (Questions identiques).

### Résultats

| LLM + ferns.1.extract.md | Question 1 | Question 2 |
| ------------------------ | ---------- | ---------- |
| qwen:0.5b                | 😡         | 😡         |
| qwen2:0.5b               | 🙂         | 🙂         |
| dolphin-gemma2           | 🙂         | 😡         |

| LLM + ferns.2.extract.md | Question 1 | Question 2 |
| ------------------------ | ---------- | ---------- |
| qwen:0.5b                | 🙂         | 😡         |
| qwen2:0.5b               | 🙂         | 🙂         |
| dolphin-gemma2           | 🙂         | 🙂         |

Clairement, le fait de réduire la taille du contexte a permis à **qwen2:0.5b** de mieux exploiter les informations et de répondre correctement aux questions. **dolphin-gemma2** a eu des résultats mitigés, mais il a quand même été capable de répondre correctement à certaines questions.

De plus, la structure des données semble avoir un impact sur les résultats obtenus avec les trois modèles. La seconde structure de document semble être plus facilement exploitable.

### Conclusion

✋ **Je réalise que mes hypothèses et conclusions sont basées sur un nombre limité de tests. Il serait intéressant de répéter ces expérimentations avec un plus grand nombre de modèles et de questions.**

Néanmoins, pour la suite de mes expérimentations, je vais continuer à utiliser **qwen2:0.5b** et la seconde structure de document.

Bien sûr, je voudrais pouvoir interroger mon expert en fougeres sur d'autres variétés de fougères. Pour cela, je vais devoir mettre en place un système qui sera capable d'extraire uniquement les informations pertinentes pour une question donnée. Nous allons donc passer à la phase suivante de nos expérimentations et faire du **RAG** (
Retrieval-augmented generation) pour extraire les informations pertinentes.

## 3ème série d'expérimentations: recherche de similarité pour fournir un contexte plus pertinent mais plus petit

📝 Cette fois-ci le code pour exécuter les exemples est ici : [02-rag](https://github.com/parakeet-nest/parakeet/tree/main/examples/42-make-a-slm-smarter/02-rag)

Ce programme va faire plusieurs choses:

- Il va charger le document`ferns.2.split.md`
- Le splitter en plusieurs parties (une par variété de fougères)
- Calculer les vecteurs (embeddings) de chaque partie (ou chunk) avec l'aide d'un LLM approprié pour faire de la génération d'embeddings. Je ferais des essais avec**all-minilm:33m**,**nomic-embed-text** et**mxbai-embed-large**.
- Attendre une question utilisateur
- Calculer le vecteur de la question
- Calculer la ou les similaritées entre le vecteur de la question et les vecteurs des parties du document.
- Créer un contexte avec la ou les parties du document qui ont la plus grande similarité avec la question.
- Poser la question au modèle**qwen2:0.5b**  avec le contexte généré.

> Je ferqis aussi des essais avec**qwen2:1.5b**.

✋ `ferns.2.split.md` est un fichier markdown qui contient les mêmes informations que `ferns.2.md` (Il est aussi structuré de la même manière), mais dans lequels j'ai ajouté un **marqueur** `<!-- SPLIT -->` à chaque fin de section d'une variété de fougère, pour indiquer les différentes parties du document. Cela me permettra de découper le document en plusieurs parties et de calculer les embeddings de chaque partie.

📝 Vous pouvez trouver le fichier complet ici: [ferns.2.split.md](https://github.com/parakeet-nest/parakeet/tree/main/examples/42-make-a-slm-smarter/data/ferns.2.split.md)

✋ Pour cette expérimentation, j'ai utilisé un "in memory vectore store" pour stocker les vecteurs des parties du document, et la distance "cossine" pour calculer la similarité entre les vecteurs des parties du document et le vecteur de la question. Ces fonctionnalités sont disponibles dans le projet [Parakeet](https://github.com/parakeet-nest/parakeet).

> Parakeet fournit d'autres fonctionnalités pour faire du RAG, notament avec **Elasticsearch**. Vous pouvez consulter la documentation et les exemples pour plus d'informations.

### Recherche de similarités

Pour faire les recherche de similarités, j'ai utilisé la fonction `SearchTopNSimilarities` de [Parakeet](https://github.com/parakeet-nest/parakeet). Voici la signature de cette fonction:

```golang
func (mvs *embeddings.MemoryVectorStore) SearchTopNSimilarities(embeddingFromQuestion llm.VectorRecord, limit float64, max int) ([]llm.VectorRecord, error)
```

```text
SearchTopNSimilarities searches for the top N similar vector records based on the given embedding from a question. It returns a slice of vector records and an error if any. The limit parameter specifies the minimum similarity score for a record to be considered similar. The max parameter specifies the maximum number of vector records to return.
```

Donc,par exemple si j'utilise:

```golang
similarities, err := store.SearchTopNSimilarities(embeddingFromQuestion, 0.5, 1)
```

La fonction va chercher les vecteurs (dans le vector store) qui ont une distance cosine supérieure ou égale à 0.5 avec le vecteur de la question. Et elle va retourner le vecteur avec le meilleure score.

Et si je veux chercher les 3 meilleurs vecteurs:

```golang
similarities, err := store.SearchTopNSimilarities(embeddingFromQuestion, 0.5, 3)
```

### Questions

Pour faire cette nouvelle expérimentation, j'ai utilisé les mêmes questions que pour l'experimentation précédente:

- **Question 1:** Give me a list of ferns of the Dryopteridaceae variety
- **Question 2:** What is the common name Dryopteris cristata?

### Résultats

Les résultats obtenus sont les suivants:

|   | LLM + ferns.2.split.md | Question 1 | Question 2 | TopNSimilarities |
| - | ---------------------- | ---------- | ---------- | ---------------- |
| 1 | qwen2:0.5b             | 😡         | 🙂😡       | 3                |
| 2 | qwen2:0.5b             | 🙂         | 🙂😡       | 1                |

1. **Retourne jusqu'à 3 similarités**: Dans le 1ers cas**qwen2:0.5b** ne peut pas répondre correctement à la question, il boucle. Concernant la question 2, il répond correctement si je n'ai qu'une seule similarité. Mais il ne peut pas répondre correctement si j'ai plus d'une similarité. Et certaines fois, il ne trouve aucune similarité même si l'information existe.
2. **Retourne 1 similarité**: Dans le 2ème cas,**qwen2:0.5b** répond correctement à la question 1 et 2. Mais parfois il ne trouve aucune similarité même si l'information existe.

#### 1ères Hypothèses et Observations

- Il ne faut remonter qu'une seule similarité pour que**qwen2:0.5b** puisse répondre correctement à la question (pour être focus).
- Il ne faut pas remonter 0 similarité pour que**qwen2:0.5b** puisse répondre correctement à la question (pour avoir l'information).
- Je dois donc trouver un moyen d'améliorer la recherche de similarité pour être sûr de retrouver l'information.

Je vais donc utiliser un autre modèle pour faire la génération d'embeddings: ****nomic-embed-text****.

### Utilisation de nomic-embed-text

|   | LLM + ferns.2.split.md | Question 1 | Question 2 | TopNSimilarities |
| - | ---------------------- | ---------- | ---------- | ---------------- |
| 1 | qwen2:0.5b             | 😡         | 🙂         | 3                |
| 2 | qwen2:0.5b             | 🙂         | 😡         | 1                |

Même si il y a une amélioration, je me retrouve encore avec des recherches de similarité qui ne sont pas satisfaisantes (0 similarité même si l'information existe). Je vais donc essayer avec **mxbai-embed-large**.

### Utilisation de mxbai-embed-large

|   | LLM + ferns.2.split.md | Question 1 | Question 2 | TopNSimilarities |
| - | ---------------------- | ---------- | ---------- | ---------------- |
| 1 | qwen2:0.5b             | 😡         | 🙂         | 3                |
| 2 | qwen2:0.5b             | 🙂         | 🙂😡       | 1                |

Clairement je dois limiter à une seule similarité pour que **qwen2:0.5b** puisse répondre le plus correctement à la question.

L'utilisation de **mxbai-embed-large** a apporté une amélioration significative. Je n'ai plus de resultat de recherche de similarité qui ne sont pas satisfaisants (l'information est bien retrouvée).

Cependant **qwen2:0.5b** ne répond pas toujours correctement à la question 2 et utilise l'information d'une autre fougère de la même variété.

Je vais donc faire le même test avec d'autres modéles pour voir si je peux obtenir de meilleurs résultats.

## 4ème série d'expérimentations: essais avec d'autres modèles

Je conserve **mxbai-embed-large** pour la génération d'embeddings et je vais faire des essais avec d'autres modèles pour la complétion.

|   | LLM + ferns.2.split.md | Question 1 | Question 2 | TopNSimilarities |
| - | ---------------------- | ---------- | ---------- | ---------------- |
| 1 | qwen2:1.5b             | 🙂         | 🙂         | 3                |
| 2 | qwen2:1.5b             | 🙂         | 🙂         | 1                |
| 3 | tinydolphin            | 😡         | 😡         | 3                |
| 4 | tinydolphin            | 😡🙂       | 🙂         | 1                |

Les tailles de paramètres des modèles sont les suivantes:

- **qwen2:0.5b** (352 MB)`0.5b`
- **qwen2:1.5b** (934 MB)`1.5b`
- **tinydolphin** (636 MB)`1.1b`

Malheureusement, **qwen2:0.5b** ne me donne pas entièrement satisfaction pour mon cas d'usage d'expert en fougère.

**qwen2:1.5b** est beaucoup mieux que **qwen2:0.5b**. Il répond correctement à la question 1 et 2. Il est capable de répondre correctement même si je remonte 3 similarités.

**tinydolphin** est aussi capable de répondre correctement à la question 1 et 2. Mais pour la question 1, il arrive qu'il me retourne des résultats en doubles. Par contre il est nécessaire de limiter à un seul résultat la recherche de similarité pour avoir un résultat satisfaisant.

Je me demande si je ne pourrais pas aider un peu **tinydolphin** en lui fournissant un contexte plus précis et structuré pour lui permettre de répondre correctement à la question 1 (obtenir la liste des fougères d'une variété donnée).

## 5ème série d'expérimentations: essais avec un contexte plus précis et structuré

J'ai donc créé un nouveau fichier `ferns.2.split.list.md` qui contient les informations les mêmes informations que `ferns.2.split.md` mais á la fin de chaue section de la variété de fougère, j'ai ajouté une liste des noms des fougères de la variété. Comme ceci par exemple:

```markdown
## List of Ferns of the variety: Dryopteridaceae

- **Dryopteris filix-mas** (Male Fern)
- **Dryopteris marginalis** (Marginal Wood Fern)
- **Dryopteris erythrosora** (Autumn Fern)
- **Dryopteris cristata** (Crested Wood Fern)
- **Dryopteris affinis** (Golden Male Fern)
```

En fait, je présente la même information plusieurs fois dans le document. Mais sous une forme différente.

📝 Vous pouvez trouver le fichier complet ici: [ferns.2.split.list.md](https://github.com/parakeet-nest/parakeet/tree/main/examples/42-make-a-slm-smarter/data/ferns.2.split.list.md)

Je conserve **mxbai-embed-large** pour la génération d'embeddings et je vais faire les mêmes texte que dans l'expérimentation précédente mais qvec le fichier `ferns.2.split.list.md`.

📝 Cette fois-ci le code pour exécuter les exemples est ici : [03-rag-list](https://github.com/parakeet-nest/parakeet/tree/main/examples/42-make-a-slm-smarter/03-rag-list)

### Résultats

|   | LLM + ferns.2.split.list.md | Question 1 | Question 2 | TopNSimilarities |
| - | --------------------------- | ---------- | ---------- | ---------------- |
| 1 | qwen2:1.5b                  | 🙂         | 🙂         | 3                |
| 2 | qwen2:1.5b                  | 🙂         | 🙂         | 1                |
| 3 | tinydolphin                 | 😡         | 😡🙂       | 3                |
| 4 | tinydolphin                 | 🙂         | 🙂         | 1                |

Cette fois ci en ajoutant des informations supplémentaires dans le contexte, j'ai pu obtenir de meilleurs résultats avec **tinydolphin**. Il répond correctement à la question 1 et 2 si je conserve une seule similarité.

## Conclusion

Ce type d'expérimentation est très intéressant mais peut durer indéfiniment. Il est important de bien définir les objectifs et les contraintes de l'expérimentation pour ne pas s'égarer.

Pour mon cas d'usage, mes conclusions sont les suivantes:

Pour permettre à un SLM d'être un expert en fougères, il est nécessaire de lui fournir des données de contexte sur les fougères. Il est important de structurer ces données de manière à ce qu'elles soient facilement exploitables par le modèle. Il est également important de limiter la taille du contexte pour permettre au modèle de traiter les informations de manière efficace.

Et je retiendrais les candidats suivants pour créer un bon expert en fougères:

|   | LLM + ferns.2.split.list.md | ferns.2.list.md | ferns.2.split.list.md | TopNSimilarities |
| - | --------------------------- | --------------- | --------------------- | ---------------- |
| 1 | qwen2:1.5b                  | ✅              | ✅                    | 3                |
| 2 | qwen2:1.5b                  | ✅              | ✅                    | 1                |
| 4 | tinydolphin                 |                 | ✅                    | 1                |


Je vous encourage à faire vos propres expérimentations et à adapter les concepts présentés ici à vos besoins.
