# Backend

The backend is the interface between the frontend and database servers. It is written in **Golang**. It's responsible for:

- parse the request from the frontend
- create relevant database server queries
- parse response from database servers and handle errors
- run causal discovery algorithms
- build response for frontend
