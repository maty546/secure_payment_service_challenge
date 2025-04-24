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
