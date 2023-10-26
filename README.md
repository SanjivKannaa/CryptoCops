# CryptoCops
A crawler that uses TOR proxy to crawl over the dark web in a Breadth-First fashion.

## Requirements
- Docker

## Setup
- Clone the repo
- Create `url_queue.json` and add seed urls to it.
- Create `.env` file from .env.example 
- Add the ml model to src/services/ML/fine_tuned_model_roberta/
- Run 
```bash
 $ docker-compose -f docker-compose.test.yml --env-file=src/env/.env up
 ```
- Visited website URLs will be added to `visited.json` file in the root directory