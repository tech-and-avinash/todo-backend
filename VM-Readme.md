# README - Deploying the `todo` Application

This guide walks you through setting up the backend and frontend for the `todo` application, configuring Nginx, securing the domains with SSL (Let's Encrypt), and setting up a systemd service for the API.

### Permission after copying binary
```bash
chmod +x todo-backend
```
---

## 1. Backend API Setup (`api.todo.nomadule.com`)

### Step 1: Create Nginx Site Config for API
```bash
sudo nano /etc/nginx/sites-available/api.nomadule.com
```
*(

    server {
    listen 80;
    server_name api.nomadule.com;

    location / {
        proxy_pass http://localhost:8080; # Your backend app port
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    access_log /var/log/nginx/api.nomadule.com.access.log;
    error_log  /var/log/nginx/api.nomadule.com.error.log;
}


)*

### Step 2: Enable the Site
```bash
sudo ln -s /etc/nginx/sites-available/api.nomadule.com /etc/nginx/sites-enabled/
```

### Step 3: Test and Reload Nginx
```bash
sudo nginx -t
sudo systemctl reload nginx
```

### Step 4: Install Certbot and Obtain SSL Certificate
```bash
sudo apt update
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d api.nomadule.com
```
Example nginx config file
```nginx
# Redirect all HTTP traffic to HTTPS
server {
    listen 80;
    server_name api.todo.nomadule.com;

    # Redirect to HTTPS
    return 301 https://$host$request_uri;
}

# HTTPS server with reverse proxy to Go backend
server {
    listen 443 ssl;
    server_name api.todo.nomadule.com;

    # SSL config (managed by Certbot)
    ssl_certificate /etc/letsencrypt/live/api.todo.nomadule.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.todo.nomadule.com/privkey.pem;
    include /etc/letsencrypt/options-ssl-nginx.conf;
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem;

    # Reverse proxy to Go API
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### Step 5: Verify and Reload Nginx
```bash
sudo nginx -t
sudo systemctl reload nginx
```

### Step 6: Run Backend API
You can start your Go app (example command):
```bash
./go-migrate-example
```
---

## 2. Set up systemd service for backend

### Step 1: Create systemd service file
```bash
sudo nano /etc/systemd/system/api-nomadule.service
```

Example service file:
```ini
[Unit]
Description=API Nomadule Service
After=network.target

[Service]
ExecStart=/home/azureuser/app/nomadule/nomadule-backend
WorkingDirectory=/home/azureuser/app/nomadule
Restart=always
User=azureuser
EnvironmentFile=/home/azureuser/app/nomadule/.env

[Install]
WantedBy=multi-user.target

```

### Step 2: Reload systemd and Restart Service
```bash
sudo systemctl daemon-reload
sudo systemctl restart api-nomadule.service
```

---

## 3. Testing

- Test API:
```bash
curl --location 'https://api.nomadule.com/users'
```
- Test Local API:
```bash
curl --location 'http://localhost:8080/users'
```
- Check Nginx configuration anytime:
```bash
sudo nginx -t
```
- Reload Nginx if any configuration changes:
```bash
sudo systemctl reload nginx
```

---

# Done! 🚀
Now your API and frontend are served with HTTPS, and the API runs as a service.

journalctl -u api-nomadule.service -n 50 --no-pager
chmod +x /home/azureuser/app/nomadule/nomadule-backend
nomadule@vm-nomadule:/$ sudo systemctl restart api-nomadule.service
nomadule@vm-nomadule:/$ sudo systemctl status api-nomadule.service