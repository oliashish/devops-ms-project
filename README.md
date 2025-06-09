# DevOps Microservices Project: GitOps with Argo CD & Kustomize

![Build Status](https://img.shields.io/badge/CI/CD-GitOps%20Automated-brightgreen)
![Kubernetes](https://img.shields.io/badge/Kubernetes-K3s-blue)
![ArgoCD](https://img.shields.io/badge/ArgoCD-Deployed-blueviolet)
![Go](https://img.shields.io/badge/Go-Microservices-00ADD8?logo=go)
![Docker](https://img.shields.io/badge/Docker-Images-2496ED?logo=docker)

---

## Overview

This project demonstrates a modern, automated **GitOps-driven CI/CD pipeline** for deploying Go microservices to a Kubernetes cluster (K3s on Raspberry Pi). It leverages **Argo CD** as the Continuous Delivery tool and **Argo CD Image Updater** to automate image tag updates directly within the Git repository, adhering strictly to GitOps principles. **Kustomize** is used for declarative customization of Kubernetes manifests.

The core idea is to maintain the desired state of the infrastructure and applications in a Git repository, allowing for declarative deployments, easy rollbacks, auditability, and self-healing capabilities.

---

**NOTE:** This is a single repository setup. While this kind of setup is not recommended. This is for learning purposes only.

---

## Technologies Used

- [**Go**](https://golang.org/): Programming language for the backend `user-service` and `product-service` microservices.

- [**Docker**](https://www.docker.com/): Containerization of microservices.

- [**Kubernetes (K3s)**](https://k3s.io/): Lightweight Kubernetes distribution running on a Raspberry Pi cluster, serving as the target deployment environment.

- [**Argo CD**](https://argoproj.github.io/argo-cd/): Declarative, GitOps continuous delivery tool for Kubernetes. It continuously monitors the manifest repository and ensures the cluster's state matches the declared desired state.

- [**Argo CD Image Updater**](https://argocd-image-updater.readthedocs.io/): An Argo CD extension that automatically updates container image tags in Kubernetes manifests (Git) when new images are pushed to a container registry.

- [**Kustomize**](https://kustomize.io/): A template-free configuration management tool for Kubernetes, built into `kubectl`. It allows for declarative customization of raw YAML files using overlays and patches.

- [**GitHub**](https://github.com/): Hosts both the application source code (conceptually) and the Kubernetes manifest repository (this repository), serving as the single source of truth for desired state.

- [**Docker Hub**](https://hub.docker.com/): Public container registry for storing the built Docker images.

---

## Architectural Principles & How It Works (GitOps Flow)

This project strictly follows the **GitOps methodology**, where Git is the single source of truth for defining the desired state of both infrastructure and applications. Changes are **pulled** into the cluster by specialized operators, rather than being **pushed** by external CI pipelines.

Here's the detailed flow:

1.  **Application Development:**

    - Developers write and commit Go code for `user-service` or `product-service` to their respective source code repositories (e.g., `https://github.com/oliashish/devops-ms-project` in a monorepo setup, or separate repos for each service).

2.  **Continuous Integration (CI) Pipeline:**

    - A CI system (in this case - GitHub Actions) is triggered by code commits in the application source repository.
    - The CI pipeline's sole responsibility is to:
      - Build the Go microservice binary.
      - Build a new Docker image with a **unique, immutable tag**.
      - Run tests(if any) on the built image.
      - Push the newly tagged Docker image to Docker Hub(Or any repository of your choice).
    - **Crucially:** The CI pipeline **DOES NOT** directly update Kubernetes manifests or interact with the cluster(intentional for security).

3.  **Kubernetes Manifests as Desired State (This Repository):**

    - This repository (`oliashish/devops-ms-project`) serves as the **GitOps repository**. It contains declarative Kubernetes YAML files (`deployment.yaml`, `service.yaml`) for each microservice.
    - **Kustomize's Role:** Instead of hardcoding image tags in `deployment.yaml`, Kustomize is used to define the image details (name and tag) in a `kustomization.yaml` file for each service. This allows for declarative overrides and enables automation.

4.  **Argo CD Image Updater (Automated Git Write-Back):**

    - Running within the K3s cluster, Argo CD Image Updater periodically polls Docker Hub for new image tags for the configured microservices.
    - When a new, valid tag is detected (e.g., `v1.0.2` for `user-service`), Image Updater automatically performs a `git commit` to _this_ GitOps repository (`oliashish/devops-ms-project`).
    - The commit updates the `newTag` field within the `kustomization.yaml` for the respective microservice (e.g., `k8s-manifests/user-service/kustomization.yaml`).
    - **Why this approach?** This is the core of automated GitOps. The source of truth (Git) is updated by an in-cluster agent, maintaining the declarative nature and audit trail of all deployments.

5.  **Argo CD (Continuous Reconciliation):**

    - The main Argo CD controller continuously monitors this GitOps repository (`oliashish/devops-ms-project`).
    - It detects the new commit made by Argo CD Image Updater.
    - Argo CD then uses its **internal Kustomize engine** to render the full Kubernetes manifests for the application, injecting the new image tag from `kustomization.yaml` into `deployment.yaml`.
    - It compares the "desired state" (the newly rendered manifests from Git) with the "actual state" of your K3s cluster.
    - If a drift is detected (e.g., `user-service` is running `v1.0.1` but Git now declares `v1.0.2`), Argo CD automatically synchronizes the cluster by applying the new manifests.

6.  **Kubernetes Deployment:**
    - The Kubernetes API server receives the updated deployment manifest from Argo CD.
    - Kubernetes initiates a rolling update for your microservice, pulling the newly specified Docker image and spinning up new pods.
    - Your microservice is now running the latest version!

This entire process provides a robust, auditable, and self-healing deployment mechanism, greatly reducing manual intervention and increasing deployment confidence.

---

## Directory Structure

```
user-service
    - minimal code to run user-service
product-service
    - minimal code to run product-service
k8s-manifest
    - user-service # user-service related deployment, service and kustomization configuration
    - product-service
        - product-service related deployment, service and kustomization configuration
.github/workflows
    - build.yaml # contains CI pipeline
```

---

**NOTE** -

> 1. Need to setup permissions on PR and Merges and lock on main branch
> 2. Proper versioning needs to be done for container image tag (currently using latest for testing and local environment but not recommended)
> 3. These services don't talk to each other yet
> 4. Using multi-arch build for docker images for k3s running on rpi
> 5. No tests written
> 6. No multi-environment support i.e no support of flags and optiong for dev, stage or prod.
