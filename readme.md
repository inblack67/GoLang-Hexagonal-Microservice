### Hexagonal Architecture

- Each piece of software maintains it's seperation of concerns
- Extremely modular
- Domain logic should be independent of frameworks
- **repo** <- **service** -> serializer (json | msgpack) -> http
