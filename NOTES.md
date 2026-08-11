
i just went though my command history and pulled a bunch of stuff that might be useful
```
git clone git@github.com:RED-M0CKING-LINE/RED-M0CKING-LINE.github.io.git

firewall-cmd --permanent --add-service=http;
firewall-cmd --permanent --add-service=https;
firewall-cmd --reload

journalctl -t personal-website-main-go-deploy -f
crontab -e

systemctl --user enable podman-restart
loginctl enable-linger 1000

export DOCKER_HOST="unix://$XDG_RUNTIME_DIR/podman/podman.sock"

sudo ausearch -m AVC,USER_AVC -ts recent -i
podman image inspect localhost/compose_nginx:latest --format 'configured user: {{.Config.User}}'
podman unshare chown -R 65532:65532 deploy/nginx/logs

```

