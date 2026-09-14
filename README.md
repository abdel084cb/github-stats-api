# GitHub Stats API

## Requisitos

- Go, con una versión compatible con la indicada en `go.mod`.

SQLite se utiliza mediante el driver `modernc.org/sqlite`. No hace falta levantar un servidor de base de datos.

## Ejecución

Desde la raíz del repositorio:

```bash
go mod download
go run ./cmd
```
El servicio escucha en `http://localhost:8003`.

La base de datos y su tabla se crean automáticamente al iniciar la aplicación. 
El archivo `stats.db` se guarda en el directorio desde el que se ejecuta el programa.

## Endpoint

```http
GET /stats/{username}
```

Ejemplo de petición:

```bash
curl -i http://localhost:8003/stats/octocat
```

## Pruebas

Desde la raíz del repositorio:

```bash
go test -v ./...
```

## Limitaciones y mejoras pendientes

Los tests actuales consultan directamente la API de GitHub. Esto permite probar la integración real, pero introduce dos limitaciones:

- Los repositorios, las estrellas y los lenguajes de los usuarios pueden cambiar. Por tanto, los datos esperados definidos en los tests podrían necesitar actualizarse.
- No se puede provocar de forma controlada que GitHub deje de responder o aplique un rate limit. El tratamiento de estos errores está implementado, pero los casos `503` no están cubiertos actualmente por tests automáticos.

Una posible mejora sería ejecutar un servidor HTTP local que simule las respuestas de GitHub y permitir que el cliente utilice su dirección durante los tests. Esto haría posible utilizar datos estables y simular respuestas como `404`, `403`, `429` o errores de conexión.

He dejado esta simulación fuera de la entrega para mantener el alcance y la complejidad de los tests acordes con el tamaño del ejercicio. Sería la siguiente mejora que implementaría.