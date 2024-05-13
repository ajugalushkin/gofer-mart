Гофермарт
===

## Table of Contents

[TOC]

## ER-диаграмма
```mermaid
erDiagram
    USER ||--o{ ORDER : is
    USER {
        string userId PK
        string firstName
        string lastName
        string middleName
        int age
    }
    ORDER {
        string orderId PK
        string name
        string status
        int accrual
        string uploaded_at
    }
    LOGIN |o--|| USER : has
    LOGIN {
        string Login PK
        string Password
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
