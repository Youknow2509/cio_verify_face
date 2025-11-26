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
        "${docker_hub_info}/${proj_name}_service_ai:cpu_test",
    ]
}


target "service_analytic" {
    context = "."
    dockerfile = "service_analytic/Dockerfile"
    tags = [
        "${docker_hub_info}/${proj_name}_service_analytic:latest_test",
    ]
}

target "service_attendance" {
    context = "."
    dockerfile = "service_attendance/Dockerfile"
    tags = [
        "${docker_hub_info}/${proj_name}_service_attendance:latest_test",
    ]
}

target "service_auth" {
    context = "."
    dockerfile = "service_auth/Dockerfile"
    tags = [
        "${docker_hub_info}/${proj_name}_service_auth:latest_test",
    ]
}

target "service_device" {
    context = "./service_device"
    dockerfile = "Dockerfile"
    tags = [
        "${docker_hub_info}/${proj_name}_service_device:latest_test",
    ]
}

target "service_identity" {
    context = "./service_identity"
    dockerfile = "Dockerfile"
    tags = [
        "${docker_hub_info}/${proj_name}_service_identity:latest_test",
    ]
}

target "service_notify" {
    context = "./service_notify"
    dockerfile = "Dockerfile"
    tags = [
        "${docker_hub_info}/${proj_name}_service_notify:latest_test",
    ]
}

target "service_profile_update" {
    context = "./service_profile_update"
    dockerfile = "Dockerfile"
    tags = [
        "${docker_hub_info}/${proj_name}_service_profile_update:latest_test",
    ]
}

target "service_workforce" {
    context = "./service_workforce"
    dockerfile = "Dockerfile"
    tags = [
        "${docker_hub_info}/${proj_name}_service_workforce:latest_test",
    ]
}





