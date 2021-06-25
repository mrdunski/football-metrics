pipeline {
    environment {
        IMAGE_TAG = "${(GIT_BRANCH == 'master') ? "0.$BUILD_NUMBER.0.$GIT_COMMIT" : "$GIT_COMMIT"}"
        IMAGE_REPO_AND_TAG = "docker.kende.pl/football-metrics:$IMAGE_TAG"
        CHART_VERSION = "${(GIT_BRANCH == 'master') ? "0.$BUILD_NUMBER.0+$GIT_COMMIT" : "0.0.1+$GIT_COMMIT"}"
        DOCKER_BUILDKIT = "1"
        CHART_REPOSITORY_URL="https://charts.kende.pl/api/charts"
        FLEET_URL="ssh://git@ssh.git.kende.pl:8022/kende/fleet.git"
    }
    agent {
        node { label 'docker' }
    }

    stages {
        stage('Test') {
            steps {
                sh "docker build  --pull --no-cache -o build --target artifacts ."
                junit 'build/report.xml'
            }
        }

        stage('Build docker') {
            steps {
                script {
                    currentBuild.description = "image: $IMAGE_REPO_AND_TAG"
                }
                sh "docker build --pull -t $IMAGE_REPO_AND_TAG ."
                sh "docker push $IMAGE_REPO_AND_TAG"
            }
        }
        
        stage('Create chart') {
            steps {
                script {
                    currentBuild.description = """image: $IMAGE_REPO_AND_TAG
chart-version: $CHART_VERSION
app-version: $IMAGE_TAG"""
                }
                sh "apk add --no-cache curl tar gzip"
                sh "curl https://get.helm.sh/helm-v3.1.0-linux-amd64.tar.gz --output helm.tar.gz; tar -zxvf helm.tar.gz; cp linux-amd64/helm /usr/local/bin; rm helm.tar.gz linux-amd64 -r"
                sh "helm package ./charts/app --version $CHART_VERSION --app-version $IMAGE_TAG"
                withCredentials([usernamePassword(credentialsId: 'charts', passwordVariable: 'chartPassword', usernameVariable: 'chartUser')]) {
                    sh "curl --user $chartUser:$chartPassword --data-binary '@football-metrics-${CHART_VERSION}.tgz' ${CHART_REPOSITORY_URL}"
                }
            }
        }

        stage('Deploy') {
            when {
                expression {
                    return GIT_BRANCH == 'master'
                }
            }
            steps {
                git changelog: false, credentialsId: 'jenkins-ssh', url: "${FLEET_URL}"

                withCredentials([usernamePassword(credentialsId: 'charts', passwordVariable: 'chartPassword', usernameVariable: 'chartUser')]) {
                    script {
                        def defaultFleetYaml = """defaultNamespace: football-metrics
helm:
  releaseName: football-metrics
  chart: football-metrics
  version: 0.0.1
  repo: https://$chartUser:$chartPassword@charts.kende.pl
  values:
    ingress:
      enabled: true
      annotations:
        kubernetes.io/ingress.class: nginx
        kubernetes.io/tls-acme: "true"
      hosts:
      - host: football-metrics.dev.kende.pl
        paths:
        - path: /
      tls:
      - secretName: football-metrics.dev.kende.pl
        hosts:
        - football-metrics.dev.kende.pl
    resources:
      limits:
        cpu: 500m
        memory: 128Mi
      requests:
        cpu: 10m
        memory: 64Mi"""

                        def previous = defaultFleetYaml

                        if (fileExists(file: 'football-metrics/fleet.yaml')) {
                            previous = readFile(file: 'football-metrics/fleet.yaml')
                        }
                        def next = previous.replaceAll(/version:.*/, "version: $CHART_VERSION")
                        writeFile(file: 'football-metrics/fleet.yaml', text: next)
                    }
                }

                sh "git config --global user.email \"ci@kende.pl\""
                sh "git config --global user.name \"Jenkins\""
                sh "git add football-metrics/fleet.yaml"
                sh "git commit -m 'Auto-deploy ($JOB_NAME $CHART_VERSION)'"

                sshagent(['jenkins-ssh']) {
                    sh "git push origin master"
                }
            }
        }
    }
}