# ============================================
# Variables
# ============================================
variable "proj_name" {
    default = "cio_verify_face"
    description = "Project Name"
}

variable "docker_hub_info" {
    default = "someone2509"
    description = "Docker Hub Push Image"
}

# ============================================
# Groups
# ============================================
group "default" {
    targets = [
        "service_ai_with_cpu",
        "service_analytic",
        "service_attendance",
        "service_auth",
        "service_device",
        "service_identity",
        "service_notify",
        "service_profile_update",
        "service_workforce",
    ]
}

# ============================================
# Targets
# ============================================
target "service_ai_with_cpu" {
    context = "./service_ai"
    dockerfile = "Dockerfile"
    tags = [
        "${docker_hub_info}/${proj_name}_service_ai:cpu",
    ]
}


target "service_analytic" {
    context = "."
    dockerfile = "service_analytic/Dockerfile"
    tags = [
        "${docker_hub_info}/${proj_name}_service_analytic:latest",
    ]
}

target "service_attendance" {
    context = "."
    dockerfile = "service_attendance/Dockerfile"
    tags = [
        "${docker_hub_info}/${proj_name}_service_attendance:latest",
    ]
}

target "service_auth" {
    context = "."
    dockerfile = "service_auth/Dockerfile"
    tags = [
        "${docker_hub_info}/${proj_name}_service_auth:latest",
    ]
}

target "service_device" {
    context = "."
    dockerfile = "service_device/Dockerfile"
    tags = [
        "${docker_hub_info}/${proj_name}_service_device:latest",
    ]
}

target "service_identity" {
    context = "."
    dockerfile = "service_identity/Dockerfile"
    tags = [
        "${docker_hub_info}/${proj_name}_service_identity:latest",
    ]
}

target "service_notify" {
    context = "."
    dockerfile = "service_notify/Dockerfile"
    tags = [
        "${docker_hub_info}/${proj_name}_service_notify:latest",
    ]
}

target "service_profile_update" {
    context = "."
    dockerfile = "service_profile_update/Dockerfile"
    tags = [
        "${docker_hub_info}/${proj_name}_service_profile_update:latest",
    ]
}

target "service_workforce" {
    context = "."
    dockerfile = "service_workforce/Dockerfile"
    tags = [
        "${docker_hub_info}/${proj_name}_service_workforce:latest",
    ]
}





