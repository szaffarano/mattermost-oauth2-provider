# Oauth2 Provider para integrar authenticación de mattermost

1. Crear estructura de directorios para estado de los contenedores

        ./setup.sh
2. Modificar el archivo `nginx/static/login.html` para que la URL apunte a la IP local
3. Modificar configuración en `main.go` para adaptar a las URLs locales
4. Ejecutar `docker-compose`

        docker-compose build
        docker-compose up -d
5. En `http://<IP LOCAL>` estará escuchando mattermost
6. En `http://<IP LOCAL>:8080` estará ejecutando un simulador de login externo
7. En `http://<IP LOCAL>:9096` estará ejecutando el oauth provider
8. Acceder a `http://<IP LOCAL>` y autenticarse con contraseña por única vez para obtener un usuario admin
9. Ingresar a la consola de administración de mattermost: `http://<IP LOCAL>/admin_console/system_analytics` y configurar:
   1. En GENERAL -> Configuration opción "Site URL"
   2. En AUTHENTICATION -> GitLab, opciones
      - Application ID: ver en la configuración, por default dejé 222222
      - Application Secret Key: ver en la configuración
      - GitLab Site URL: `http://<IP LOCAL>:9096`
10. Salvar.  Hacer logout e ingresar con la opción de GitLab o desde `http://<IP LOCAL>:8080`