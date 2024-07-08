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
    - (**Optional**) If you prefer using your own TLS certs, you can follow this documentation: https://docs.gitlab.com/omnibus/settings/ssl/#configure-https-manually

4. Reconfigure GitLab
    ```bash
    sudo gitlab-ctl reconfigure
    ```

5. Logging In for the first time
    - Visit the domain name of your GitLab server in your web browser: `https://your_domain`
    - GitLab generates an `initial secure password` for you. It is stored in `/etc/gitlab/initial_root_password`
    - On the login page, enter `Username: root` and `Password: [the password listed on /etc/gitlab/initial_root_password]`

        <p align="center"> <img src="images/login.png"> </p>


## Create Gitlab Repository
Let’s start by creating a GitLab project and adding an HTML file to it. You will later copy the HTML file into an Nginx Docker image, then deploy it to the server

1. Create new project
    - Log in to your GitLab instance and click **New project**
        <p align="center"> <img src="images/new-project.png"> </p>
    - Select **Create blank project**
    - Give it a proper **Project name**
    - Set the **Visibility Level** to **Public**
    - Finally click **Create project**
        <p align="center"> <img src="images/create-project.png"> </p>

2. Add html file
    - On your Project’s overview page, click **+** > **New file**
      <p align="center"> <img src="images/add-file.png"> </p>  
    - Set the File name to `index.html`
    - Add the following HTML to the file body:
        ```html
        <html>
            <body>
                <h1>Hello, World</h1>
            </body>
        </html>
        ```
    - Click **Commit changes** at the bottom of the page

3. Add Dockerfile
    - Go back to the Project’s overview page, click the **+** button and select the **New file**
    - Set the File name to `Dockerfile`
    - Add these instructions to the file body:
        ```Dockerfile
        FROM nginx:stable-alpine3.19-slim
        COPY index.html /usr/share/nginx/html
        ```
    - Click **Commit changes** at the bottom of the page


## Setup Gitlab Runner
Besides GitLab, the `GitLab runner` is yet another system that your private key will enter. For each pipeline, `GitLab` uses `runners` to perform the heavy work, that is, execute the jobs you have specified in the `CI/CD` configuration. That means the deployment job will ultimately be executed on a `GitLab runner`.

1. Login to `gitlab-runner` server
    ```bash
    ssh root @ gitlab_runner_ip
    ```

2. Install depedencies
    ```bash
    sudo apt update
    sudo apt install ca-certificates curl openssh-server postfix tzdata perl
    ```

3. Install gitlab-runner
    ```bash
    # add the official GitLab repository
    curl -L "https://packages.gitlab.com/install/repositories/runner/gitlab-runner/script.deb.sh" | sudo bash

    # install the gitlab-runner package
    sudo apt-get install gitlab-runner

    # verify the installation by checking the service status
    systemctl status gitlab-runner
    ```

4. Obtain project token and GitLab URL
    - In your GitLab project, navigate to **Settings** > **CI/CD** > **Runners** > **Expand**
        <p align="center"> <img src="images/gitlab-runner.png"> </p>
    - In the **Project runners** section, click on **New project runner** and follow the form to create a new runner for your project
    - Fill a proper **Tags** e.g. deployment
    - Once a runner is in place, you’ll find the **registration token** and the **GitLab URL**. Copy both to a text editor
        <p align="center"> <img src="images/gitlab-runner-2.png"> </p>

5. Register the runner
    ```
    # Back to your terminal, register the runner for your project
    sudo gitlab-runner register -n --url https://your_gitlab.com --registration-token project_token --executor docker --description "Deployment Runner" --docker-image "docker:stable" --tag-list deployment --docker-privileged
    
    # Output
    Runner registered successfully. Feel free to start it, but if it's running already the config should be automatically reloaded!
    ```

6. Verify the registration process
    - Go to **Settings** > **CI/CD** > **Runners** in GitLab, where the registered runner will show up
        <p align="center"> <img src="images/gitlab-runner-3.png"> </p>


## Setup Web Server
The CI/CD will run a web server in this VM. We are going to create a user that is dedicated for the deployment task. You will later configure the CI/CD pipeline to log in to the server with that user.

1. Login to `web-server`
    ```bash
    ssh root @ web-server_ip
    ```

2. Create a deployment user
    ```bash
    # create a new user
    sudo adduser deployer

    # add the user to the Docker group
    sudo usermod -aG docker deployer
    ```

3. Setting up an SSH key for the deployment user
    ```bash
    # switch to user deployer
    su deployer

    # generate a SSH key
    ssh-keygen

    Generating public/private rsa key pair.
    Enter file in which to save the key (~/.ssh/id_rsa): ENTER
    Enter passphrase (empty for no passphrase): ENTER
    Enter same passphrase again: ENTER

    Your identification has been saved in ~/.ssh/id_rsa
    Your public key has been saved in ~/.ssh/id_rsa.pub
    ```

4. To authorize the SSH key for the deployer user, you need to append the public key to the authorized_keys file
    ```bash
    cat ~/.ssh/id_rsa.pub >> ~/.ssh/authorized_keys
    ```

5.  Store the SSH private key to a notepad. Later we will store it as variable in the CI/CD pipeline
    ```bash
    cat ~/.ssh/id_rsa

    # output
    -----BEGIN RSA PRIVATE KEY-----
    ...
    -----END RSA PRIVATE KEY-----
    ```

## Configuring the .gitlab-ci.yml File
We are going to configure the GitLab CI/CD pipeline. The pipeline will build a Docker image and push it to the Docker hub. The final step in your pipeline is to log in to your server, pull the latest Docker image, remove the old container, and start a new container.

1. Navigate to **Settings** > **CI/CD** > **Variables** in your GitLab project and click **Add Variable**. Fill out the form as follows:
    - `ID_RSA`
        - Key : **ID_RSA**
        - Value : Paste your SSH private key 
        - Type  : **File**
        - Environment Scope: **All (default)**
        - Protect variable: **Checked**
        - Mask variable: **Unchecked**
    - `SERVER_IP`
        - Key: **SERVER_IP**
        - Value: `your_server_IP`
        - Type: **Variable**
        - Environment scope: **All (default)**
        - Protect variable: **Checked**
        - Mask variable: **Checked**
    - `SERVER_USER`
        - Key: **SERVER_USER**
        - Value: **deployer**
        - Type: **Variable**
        - Environment scope: **All (default)**
        - Protect variable: **Checked**
        - Mask variable: **Checked**
    - `CI_REGISTRY`
        - Key: **CI_REGISTRY**
        - Value: **docker.io**
        - Type: **Variable**
        - Environment scope: **All (default)**
        - Protect variable: **Unchecked**
        - Mask variable: **Unchecked**
    - `CI_REGISTRY_USER`
        - Key: **CI_REGISTRY_USER**
        - Value: `your_docker_username`
        - Type: **Variable**
        - Environment scope: **All (default)**
        - Protect variable: **Unchecked**
        - Mask variable: **Unchecked**
    - `CI_REGISTRY_PASSWORD`
        - Key: **CI_REGISTRY_PASSWORD**
        - Value: `your_docker_pat`
        - Type: **Variable**
        - Environment scope: **All (default)**
        - Protect variable: **Checked**
        - Mask variable: **Checked**

2. Go to the **Project overview** page, click the **+** button and select **New file**. Then set the File name to `.gitlab-ci.yml`.

3. Add these contents to the file body:
    ```yaml
    stages:
      - publish
      - deploy

    variables:
      DOCKER_DRIVER: overlay2
      DOCKER_TLS_CERTDIR: ""
      DOCKER_IMAGE: $CI_REGISTRY_USER/<image-name>:$CI_COMMIT_SHORT_SHA

    publish:
      image: docker:latest
      stage: publish
      services:
        - docker:dind
      script:
        - apk add --no-cache docker
        - echo $CI_REGISTRY_PASSWORD | docker login -u $CI_REGISTRY_USER --password-stdin $CI_REGISTRY
        - docker build -t $DOCKER_IMAGE .
        - docker push $DOCKER_IMAGE

    deploy:
      image: alpine:latest
      stage: deploy
      tags:
        - deployment
      script:
        - chmod og= $ID_RSA
        - apk update && apk add openssh-client
        - ssh -i $ID_RSA -o StrictHostKeyChecking=no $SERVER_USER@$SERVER_IP "echo $CI_REGISTRY_PASSWORD | docker login -u $CI_REGISTRY_USER --password-stdin $CI_REGISTRY"
        - ssh -i $ID_RSA -o StrictHostKeyChecking=no $SERVER_USER@$SERVER_IP "docker pull $DOCKER_IMAGE"
        - ssh -i $ID_RSA -o StrictHostKeyChecking=no $SERVER_USER@$SERVER_IP "docker container rm -f my-app || true"
        - ssh -i $ID_RSA -o StrictHostKeyChecking=no $SERVER_USER@$SERVER_IP "docker run -d -p 80:80 --name my-app $DOCKER_IMAGE"
      environment:
        name: production
        url: http://<server-ip>
    ```

4. Finally click **Commit changes** at the bottom of the page in GitLab to create the `.gitlab-ci.yml` file

5. Validating the CI/CD pipeline
    - When a `.gitlab-ci.yml` file is pushed to the repository, GitLab will automatically detect it and start a CI/CD pipeline. At the time you created the `.gitlab-ci.yml` file, GitLab started the first pipeline.
    - Go to **Build** > **Pipelines** in your GitLab project to see the pipeline’s status.
        <p align="center"> <img src="images/ci-cd.png"> </p>

6. Validating the deployment on the web-server
    - We want to check the deployed container on your web-server. Head over to your terminal and make sure to log in again
        ```bash
        ssh root @ web_server_ip
        ```
    - List the running containers
        ```bash
        docker ps

        # output
        CONTAINER ID   IMAGE                        COMMAND                  CREATED        STATUS        PORTS                               NAMES
        fb86cae159e2   docker-user/my-app:9f4a2e61  "/docker-entrypoint.…"   23 hours ago   Up 23 hours   0.0.0.0:80->80/tcp, :::80->80/tcp   my-app
        ```