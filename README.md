# Wisteria 💮

[![GitHub license](https://img.shields.io/github/license/vickyphang/wisteria)](https://github.com/vickyphang/wisteria/blob/main/LICENSE)
![GitHub stars](https://img.shields.io/github/stars/vickyphang/wisteria)

### Gitlab CI/CD Home Project

The purpose of this project is to build a `ci/cd pipeline` with `Gitlab`. I will be using my own `Gitlab server`, register a `gitlab-runner`, configure the pipeline to build a `Docker image`, push it to the `Docker hub`, and deploy it to a server using `SSH`. For demonstration purpose, I will deploy a small, `static web page` using Nginx

### Prerequisite
- 3 VMs using Ubuntu 22.04 (gitlab-server, gitlab-runner, web-server)
- A registered domain-name
- Docker installed in each VM
- Docker hub account

## Setup Gitlab Server
1. Install depedencies
```bash
sudo apt update
sudo apt install ca-certificates curl openssh-server postfix tzdata perl
```

2. Install Gitlab
```bash
# move into the /tmp directory
cd /tmp

# download the installation script
curl -LO https://packages.gitlab.com/install/repositories/gitlab/gitlab-ce/script.deb.sh

# run the installer
sudo bash /tmp/script.deb.sh

# install the actual GitLab application with apt
sudo apt install gitlab-ce
```

3. Edit the Gitlab configuration file
    - Open config file: `sudo nano /etc/gitlab/gitlab.rb`
    - Match `external_url` with your `domain` and make sure to change `http` to `https` to automatically redirect users to the site protected by the Let’s Encrypt certificate
    - Fill `letsencrypt['contact_emails']` with your email. In case there are problems with your domain, Let’s Encrypt project can use to contact you
    - (**Optional**) If you prefer using your own TLS certs, you can follow this documentation: https://docs.gitlab.com/omnibus/settings/ssl/

4. Reconfigure GitLab
```bash
sudo gitlab-ctl reconfigure
```

5. Logging In for the first time
    - Visit the domain name of your GitLab server in your web browser: `https://your_domain`
    - GitLab generates an `initial secure password` for you. It is stored in `/etc/gitlab/initial_root_password`
    - On the login page, enter `Username: root` and `Password: [the password listed on /etc/gitlab/initial_root_password]`

<p align="center"> <img src="images/login.png"> </p>

