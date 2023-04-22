# GO-DISTRIBUTED


```mermaid
graph LR;
    A(Client) <==> B[WEB Protal Service]
    B <==> C[Registry Service];
    B --> D[Log Service];
    B --> E[Grading Service];
    C <==> D;
    C <==> E;
    E --> D;

    style A fill:#f9f,stroke:#333,stroke-width:4px,font-size:32px,color:#000,font-weight:bold,left:100px,top:100px;
    style B fill:#bbf,stroke:#333,stroke-width:4px,font-size:24px,color:#000,font-weight:bold,left:300px,top:100px;
    style C fill:#bfb,stroke:#333,stroke-width:4px,font-size:16px,color:#000,font-weight:bold,left:500px,top:100px;
    style D fill:#fbf,stroke:#333,stroke-width:4px,font-size:100px,color:#000,font-weight:bold,left:700px,top:300px
    style E fill:#bff,stroke:#333,stroke-width:4px,font-size:12px,color:#000,font-weight:bold,left:700px,top:300px
```