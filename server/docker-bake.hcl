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
        "service_attendance",
        "service_auth",
        "service_device",
        "service_notify",
        "service_workforce",
        "service_ws_delivery",
        "service_ai_with_cpu",
        "service_analytic",
        "service_profile_update",
    ]
}

# ============================================
# Targets
# ============================================
target "service_profile_update" {
    context = "./service_profile_update"
    dockerfile = "Dockerfile"
    tags = [
        "${docker_hub_info}/${proj_name}_service_profile_update:latest",
    ]
  
}

target "service_analytic" {
    context = "./service_analytic"
    dockerfile = "Dockerfile"
    tags = [
        "${docker_hub_info}/${proj_name}_service_analytic:latest",
    ]
}

target "service_ws_delivery" {
    context = "./service_ws_delivery"
    dockerfile = "Dockerfile"
    tags = [
        "${docker_hub_info}/${proj_name}_service_ws_delivery:latest",
    ]
    
}

target "service_ai_with_cpu" {
    context = "./service_ai"
    dockerfile = "Dockerfile"
    tags = [
        "${docker_hub_info}/${proj_name}_service_ai:cpu",
    ]
}

target "service_workforce" {
    context = "./service_workforce"
    dockerfile = "Dockerfile"
    tags = [
        "${docker_hub_info}/${proj_name}_service_workforce:latest",
    ]
}

target "service_notify" {
    context = "./service_notify"
    dockerfile = "Dockerfile"
    tags = [
        "${docker_hub_info}/${proj_name}_service_notify:latest",
    ]
}

target "service_device" {
    context = "./service_device"
    dockerfile = "Dockerfile"
    tags = [
        "${docker_hub_info}/${proj_name}_service_device:latest",
    ]
}

target "service_auth" {
    context = "./service_auth"
    dockerfile = "Dockerfile"
    tags = [
        "${docker_hub_info}/${proj_name}_service_auth:latest",
    ]
}

target "service_attendance" {
    context = "./service_attendance"
    dockerfile = "Dockerfile"
    tags = [
        "${docker_hub_info}/${proj_name}_service_attendance:latest",
    ]
}
