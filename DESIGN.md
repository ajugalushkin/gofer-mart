Гофермарт
===

[TOC]

## ER-диаграмма
```mermaid
erDiagram
    USER ||--o{ ORDER : has
    USER {
        string userId PK
        string firstName
        string lastName
        string middleName
        int age
    }
    ORDER {
        string orderId PK
        string userId FK
        string status FK
        string name
        string uploaded_at
    }
    ORDER }o--|| STATUS : contains
    STATUS{
        string code PK
        string text
    }
    USER ||--|| LOGIN : has
    LOGIN {
        string login PK
        string userId FK
        string password
    }
    ACCRUAL }o--|| ORDER : contains
    ACCRUAL }o--|| USER : contains
    ACCRUAL{
        string accrualId PK
        string userId FK
        string orderId FK
        int accrual
        string processed_at
    }
    WITHDRAWAL }o--|| ORDER : contains
    WITHDRAWAL }o--|| USER : contains
    WITHDRAWAL{
        string withdrawalId PK
        string userId FK
        string orderId FK
        int sum
        string processed_at
    }
```

## Sequence-диаграмма
```sequence
Alice->Bob: Hello Bob, how are you?
Note right of Bob: Bob thinks
Bob-->Alice: I am good thanks!
Note left of Alice: Alice responds
Alice->Bob: Where have you been?
```
