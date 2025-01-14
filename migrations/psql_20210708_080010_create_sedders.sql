-- +migrate Up
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (1, 'La estructura del cuerpo de la petición no es válida',
        'La estructura del cuerpo de la petición no es válida', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (2, 'Intentando loguearse', 'Intentando loguearse', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (3, 'No se pudo crear el registro', 'No se pudo crear el registro', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (4, 'Registro Creado', 'Registro Creado', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (5, 'El registro ya existe', 'El registro ya existe', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (6, 'Cantidad de campos invalidos', 'Cantidad de campos invalidos', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (7, 'Formato del registro invalido', 'Formato del registro invalido', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (8, 'Error en la sentencia', 'Error en la sentencia', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (9, 'Acceso denegado', 'Acceso denegado', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (10, 'Usuario o contraseña incorrectos', 'Usuario o contraseña incorrectos', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (11, 'Error de conexion', 'Error de conexion', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (12, 'No se encuentra el servidor remoto', 'No se encuentra el servidor remoto', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (13, 'Acceso permitido', 'Acceso permitido', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (14, 'Conexión establecida', 'Conexión establecida', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (15, 'La estructura no cumple validaciones', 'La estructura no cumple validaciones', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (16, 'No se pudo insertar el rol', 'No se pudo insertar el rol', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (17, 'El id debe ser numérico', 'El id debe ser numerico', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (18, 'No se pudo actualizar el registro', 'No se pudo actualizar el registro', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (19, 'Actualizado correctamente', 'Actualizado correctamente', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (20, 'No se pudo eliminar el registro', 'No se pudo eliminar el registro', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (21, 'No se pudo actualizar el rol', 'No se pudo actualizar el rol', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (22, 'No se pudo consultar el registro', 'No se pudo consultar el registro', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (23, 'Buscando si el usuario esta logueado', 'Buscando si el usuario esta logueado', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (24, 'El usuario ya esta logueado', 'El usuario ya esta logueado', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (25, 'No se pudo registrar el usuario', 'No se pudo registrar el usuario', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (26, 'Obteniendo los modulos del usuario', 'Obteniendo los modulos del usuario', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (27, 'Generando el token del usuario', 'Generando el token del usuario', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (28, 'Eliminado correctamente', 'Eliminado correctamente', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (29, 'Procesado Correctamente', 'Procesado Procesado Correctamente', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (30, 'Incicio de sesion existoso', 'Inicio de sesion existoso', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (31, 'No se puedo realizar el lout', 'No se puedo realizar el lout', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (32, 'Error buscando usuario autenticado', 'Error buscando usuario autenticado', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (33, 'Error Generando Token', 'Error Generando Token', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (34, 'La petición no contiene el token en el encabezado', 'La petición no contiene el token en el encabezado',
        '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (35, 'Su token ha expirado, por favor vuelva a ingresar', 'Su token ha expirado, por favor vuelva a ingresar',
        '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (36, 'Error de validación del token', 'Error de validación del token', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (37, 'Error al procesar el token', 'Error al procesar el token', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (38, 'Su token no es válido', 'Su token no es válido', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (39, 'El status debe ser un entero', 'El status deber ser un entero', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (40, 'No se pudo leer el archivo a cargar', 'No se pudo leer el archivo a cargar', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (41, 'No se pudo crear el documento', 'No se pudo crear el documento', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (42, 'No se pudo cargar el archivo', 'No se pudo cargar el archivo', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (43, 'Documento creado correctamente', 'Documento creado correctamente', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (44, 'No se pudo obtener las palabras claves', 'No se pudo obtener las palabras claves', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (45, 'Palabras claves obtenidas correctamente', 'Palabras claves obtenidas correctamente', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (46, 'Llave privada de enpoint no es la correcta', 'Llave privada de enpoint no es la correcta', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (47, 'No se pudo insertar datos', 'No se pudo insertar datos', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (48, 'No se pudo crear archivo', 'No se pudo crear archivo', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (49, 'La peticion no trae el private key', 'La peticion no trae el private key', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (50, 'Endpoint no se encuentra activo', 'Endpoint no se encuentra activo', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (51, 'Intentando preparar la transaccion', 'Intentando preparar la transaccion', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (52, 'Error ejecutando Actividades', 'ErrorEjecutando actividades', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (53, 'Obteniendo los elementos del componente', 'Obteniendo los elementos del componente', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (54, 'Usuario Bloqueado', 'Usuario Bloqueado', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (55, 'Error Interno del servidor', 'Error Interno del servidor', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (56, 'No se pudo cargar archivo a la tabla', 'No se pudo cargar archivo a la tabla', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (57, 'No se pudo registrar los password en la base de datos',
        'No se pudo registrar los password en la base de datos', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (58, 'No se pudo Bloquear / Desbloquear el usuario', 'No se pudo Bloquear / Desbloquear el usuario', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (59, 'Contraseña no Coincide', 'Contraseña no Coincide', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (60, 'La contraseña es insegura por pertenecer a listas Negras ',
        'La contraseña es insegura por pertenecer a listas Negras ', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (61, 'Esta contraseña ya se ha utilizado, por favor utilizar otra',
        'Esta contraseña ya se ha utilizado, por favor utilizar otra', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (62, 'Contraseña no cumple el tamaño Minimo o Maximo requerido',
        'Contraseña no cumple el tamaño Minimo o Maximo requerido', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (63, 'Contraseña no cumple la cantidad de caracteres alfabeticos requeridos',
        'Contraseña no cumple la cantidad de caracteres alfabeticos requeridos', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (64, 'Contraseña no cumple la cantidad de caracteres numericos requeridos',
        'Contraseña no cumple la cantidad de caracteres numericos requeridos', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (65, 'Contraseña no cumple con la cantidad de letras minusculas requerdas',
        'Contraseña no cumple con la cantidad de letras minusculas requerdas', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (66, 'Contraseña no cumple con la cantidad de letras mayusculas requerdas',
        'Contraseña no cumple con la cantidad de letras mayusculas requerdas', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (67, 'Contraseña no cumple con la cantidad de caracteres especiales requerdos',
        'Contraseña no cumple con la cantidad de caracteres especiales requerdos', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (68, 'Contraseñas no coinciden por favor intente nuevamente',
        'Contraseñas no coinciden por favor intente nuevamente', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (69, 'Contraseña generada correctamente', 'Contraseña generada correctamente', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (70, 'Contraseña no se puedo actualizar', 'Contraseña no se puedo actualizar', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (71, 'No se pudo enviar correo ', 'No se pudo enviar correo ', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (72, 'No se pudo consultar el rol', 'No se pudo consultar el rol', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (73, 'Ha superado la cantidad de conexiones activas permitidas',
        'Ha superado la cantidad de conexiones activas permitidas', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (74, 'No se pudo obtener Fechas permitidas de ingreso', 'No se pudo obtener Fechas permitidas de ingreso', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (75, 'No puede ingresar al sistema en esta fecha. No esta autorizado',
        'No puede ingresar al sistema en esta fecha. No esta autorizado', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (76, 'Ip no autorizada', 'Ip no autorizada', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (77, 'No fue posible consultar informacion de la IP', 'No fue posible consultar informacion de la IP', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (78, 'No se pudo consultar la ultima accion realizada por el usuario',
        'No se pudo consultar la ultima accion realizada por el usuario', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (79, 'No se pudo generar el hash del archivo', 'No se pudo generar el hash del archivo', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (80, 'No existen registros', 'No existe registros', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (81, 'No se pudo consultar pagina del documento', 'No se pudo consultar pagina del documento', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (82, 'No se pudo descifrar el archivo', 'No se pudo descifrar el archivo', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (83, 'El documento ha sido modificado con respecto a la firma',
        'El documento ha sido modificado con respecto a la firma', '1', CAST('2020-08-06 19:05:21.083' AS TIMESTAMP),
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (84, 'No se pudo decifrar el contendio enviado', 'No se pudo decifrar el contenido enviado', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (85, 'El usuario no pertenece a ningun role', 'El usuario no pertenece a ningun role', '1',
        CAST('2020-08-06 19:05:21.083' AS TIMESTAMP), CAST('2020-08-06 19:05:21.083' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (86, 'No está autorizado a realizar esta acción en este recurso',
        'No está autorizado a realizar esta acción en este recurso', '1', CAST('2020-08-12 11:47:08.120' AS TIMESTAMP),
        CAST('2020-08-12 11:47:08.120' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (87, 'No encontró registro para procesar', 'No encontro registro para actualizar', '2',
        CAST('2020-08-12 12:15:36.053' AS TIMESTAMP), CAST('2020-08-12 12:15:36.053' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (88, 'No se pudo parcear el template', 'No se pudo parcear el template', '1',
        CAST('2020-08-12 12:15:36.053' AS TIMESTAMP), CAST('2020-08-12 12:15:36.053' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (89, 'No se pudo crear el nombre del template', 'No se pudo crear el nombre del template', '1',
        CAST('2020-08-12 12:15:36.053' AS TIMESTAMP), CAST('2020-08-12 12:15:36.053' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (90, 'No se pudo crear el template en el repositorio temporal',
        'No se pudo crear el template en el repositorio temporal', '1', CAST('2020-08-12 12:15:36.053' AS TIMESTAMP),
        CAST('2020-08-12 12:15:36.053' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (91, 'No se pudo generar el documento en PDF', 'No se pudo generar el documento en PDF', '1',
        CAST('2020-08-12 12:15:36.053' AS TIMESTAMP), CAST('2020-08-12 12:15:36.053' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (92, 'No se pudo aceptar la data porque la validación de rostros no ha sido realizado',
        'No se pudo aceptar la data porque la validación de rostros no ha sido realizado', '1',
        CAST('2020-08-12 12:15:36.053' AS TIMESTAMP), CAST('2020-08-12 12:15:36.053' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (93, 'Datos validado y aceptados correctamente', 'Datos validado y aceptados correctamente', '1',
        CAST('2020-08-12 12:15:36.053' AS TIMESTAMP), CAST('2020-08-12 12:15:36.053' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (94, 'No se pudo obtener los datos del usuario embebidos en el token',
        'No se pudo obtener los datos del usuario embebidos en el token', '1',
        CAST('2020-08-12 12:15:36.053' AS TIMESTAMP), CAST('2020-08-12 12:15:36.053' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (95, 'El usuario no tiene roles asignados', 'El usuario no tiene roles asignados', '1',
        CAST('2020-08-12 12:15:36.053' AS TIMESTAMP), CAST('2020-08-12 12:15:36.053' AS TIMESTAMP));
INSERT INTO cfg.messages (id, spa, eng, type_message, created_at, updated_at)
VALUES (96, 'El usuario no esta autorizado para cargar la selfie',
        'El usuario no esta autorizado para cargar la selfie', '1',
        CAST('2020-08-12 12:15:36.053' AS TIMESTAMP), CAST('2020-08-12 12:15:36.053' AS TIMESTAMP));
-- +migrate Down
