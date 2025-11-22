internal/
├── domain/                     # Ядро - сущности и правила
│   ├── entities/               # Сущности с методами
│   │   ├── user.go
│   │   └── order.go
│   ├── value_objects/          # VO: Email, Money, Address
│   │   └── email.go
│   └── repositories/           # ИНТЕРФЕЙСЫ репозиториев
│       ├── user_repository.go
│       └── order_repository.go

├── services/                   # Use Cases - ВСЯ бизнес-логика
│   ├── user_service.go         # type UserService interface + impl
│   ├── order_service.go
│   └── auth_service.go

├── infrastructure/             # ВСЯ инфраструктура
│   ├── http/                   # HTTP слой
│   │   ├── handlers/
│   │   │   ├── user_handler.go
│   │   │   ├── order_handler.go
│   │   │   └── auth_handler.go
│   │   └── middleware/
│   │       ├── auth.go
│   │       └── logger.go
│   └── persistence/            # Работа с данными
│       ├── postgres/
│       │   ├── user_repository.go  # impl domain.UserRepository
│       │   └── order_repository.go
│       └── redis/
│           └── cache_repository.go

├── config/                     # Конфигурация
│   ├── config.go
│   └── database.go

└── shared/                     # Общие утилиты
    ├── logger/
    │   └── logger.go
    ├── database/
    │   └── connection.go
    └── utils/
        └── crypto.go