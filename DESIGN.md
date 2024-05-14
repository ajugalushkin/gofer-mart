Гофермарт
===

[TOC]

## Описание основных сущностей и связей между ними ( ER-диаграмма)
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
## Список логических компонентов системы и схема из зависимостей друг от друга (общая архитектура решения)
![image](https://hackmd.io/_uploads/rJJ-lcxX0.png)

## Две диаграммы последовательностей (sequence-диаграммы): для операции начисления балов - отразить взаимодействие с внешней системой), для операции снятия с баланса
