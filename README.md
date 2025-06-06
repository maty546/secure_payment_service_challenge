# Secure Payments Service

## Descripcion

Challenge tecnico para INI.live

## Problemas observados en la implementacion/puntos de mejora

- Falta de tests. Por tiempo no hice todos los tests unitarios, pero deje uno de ejemplo en el servicio secure_payments, para mostrar el tipo de tests que suelo hacer (aproximadamente)
- Puede pasar que ocurra un Lost Update si llegan 2 pedidos de transferencia al mismo tiempo, que de llegar secuencialmente el segundo no seria permitido porque al usuario no le alcanza el balance.
- El endpoint que usa el worker no usa el middleware de autenticacion. Esto idealmente no seria asi pero lo decidi para simplificar.
- Hay varios valores, como urls, passwords y demas, que deberian ser configuraciones o secrets, que por falta de tiempo quedaron sin ordenar.
- Hay metodos de la capa de repositorio que mezclan entidades (repo de transfer toca accounts) de una forma que hace un poco de ruido. Decidi dejarlo asi porque me parecio mas importante asegurar que las operaciones de balances y cambio de estado de transferencias sean una sola transaccion de DB. Para que quede mejor se podria refactorizar pero no hubo tiempo.
- Por falta de tiempo quedo afuera el uso de concurrencia de go. Identifique un par de casos en los que se podría agregar pero no me dio tiempo, iba a usar la concurrencia de go para hacer una especie de workqueue en memoria (como para mostrar el uso de goroutines y channels) pero al final me decidí por usar un servicio externo, que es el async worker.

## Como levantar el proyecto

### Opcion 1 - Usar la script run.sh

Despues de correr las migraciones de la carpeta "migrations" (para crear las tablas) se puede ejecutar la script run.sh. Se debe darle permiso con

```
chmod +x run.sh
```

Y luego

```
./run.sh
```

### Opcion 2 - Levantar servicios a mano

Ejecutar migraciones de la carpeta "migrations".

Es necesario correr un servidor de redis para el worker que se encarga de los tasks asincronos, corriendo la siguiente linea desde el directorio de "asyncServer" (exponiendo el puerto que se indica):

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

## Como probar el proyecto

Se puede usar el postman incluido en la carpeta tools. Cosas a tener en cuenta al probar:

- El tiempo de timeout configurado para el task asincrono que chequea estado de transacciones es de 30 segundos
- Los endpoints dentro de la carpeta "Secure Endpoints" utilizan el middleware con validacion de jwt, por lo cual es necesario tomar el token desde el endpoint de login, con el usuario "api" y la clave "123" y enviarlo como Bearer Token.
- Para que una cuenta sea considerada como externa debe contener la substring "ext"

## Misc

Dejo links de unos docs de google que use para organizar mi trabajo y escribir casos de aceptacion, como curiosidad

https://docs.google.com/spreadsheets/d/1DIH-u6JdDtWNDwCZcJFhSbZVVJ9hbeeq1RbPWzIjdck/edit?usp=sharing

https://docs.google.com/document/d/1PoqFRjYlRfd8_Ah6BPCPaJ8nw45DuL1jronyexxR5tg/edit?usp=sharing
