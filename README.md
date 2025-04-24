# Secure Payments Service

## Descripcion

Challenge tecnico para INI.live

## Suposiciones

- Las transferencias siempre tienen que pasar por un servicio externo, por mas que el servicio de secure payments tenga acceso a los balances y las cuentas, no existen transferencias 100% internas al sistema.
- No existen cuentas externas al sistema. Todas las cuentas involucradas deben figurar en la tabla accounts.
- Para mas simpleza no se distinguen entre IDs externos e internos de transaccion. Los mismos que se guardan en DB son los que llegan por callback

## Problemas observados

- Puede pasar que ocurra un Lost Update si llegan 2 pedidos de transferencia al mismo tiempo, que de llegar secuencialmente el segundo no seria permitido porque al usuario no le alcanza el balance.
- El endpoint que usa el worker no usa el middleware de autenticacion. Esto idealmente no seria asi pero lo decidi para simplificar.
- Hay varios valores, como urls, passwords y demas, que deberian ser configuraciones o secrets, que por falta de tiempo quedaron sin ordenar.

## Como levantar el proyecto

Primero es necesario correr un servidor de redis para el worker que se encarga de los tasks asincronos, corriendo la siguiente linea desde el directorio de "asyncServer" (exponiendo el puerto que se indica):

```
docker run -d --name redis-asynq -p 6379:6379 redis
```

Se puede pingear con este comando para chequear que esta levantado

```
docker exec -it redis redis-cli ping
```

Luego correr el worker con

```
make run
```

Una vez corriendo el worker, corremos la API de secure payments con el mismo comando pero en su respectivo directorio

```
make run
```
