server {
    server_name api.todo.nomadule.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    listen 443 ssl; # managed by Certbot
    ssl_certificate /etc/letsencrypt/live/api.todo.nomadule.com/fullchain.pem; # managed by Certbot
    ssl_certificate_key /etc/letsencrypt/live/api.todo.nomadule.com/privkey.pem; # managed by Certbot
    include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot

}
server {
    if ($host = api.todo.nomadule.com) {
        return 301 https://$host$request_uri;
    } # managed by Certbot


    listen 80;
    server_name api.todo.nomadule.com;
    return 404; # managed by Certbot


}

nomadule@vm-nomadule:~$ sudo ln -s /etc/nginx/sites-available/api.todo.nomadule.com /etc/nginx/sites-enabled/

nomadule@vm-nomadule:~$ sudo nginx -t
nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
nginx: configuration file /etc/nginx/nginx.conf test is successful

nomadule@vm-nomadule:~$ sudo systemctl reload nginx