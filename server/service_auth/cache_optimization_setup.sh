#!/bin/bash

# Cache Optimization Setup Script
# This script helps setup and test the cache optimization features

echo "🚀 Cache Optimization Setup & Test Script"
echo "========================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go first."
        exit 1
    fi
    print_status "Go is installed: $(go version)"
}

# Check dependencies
check_dependencies() {
    print_status "Checking project dependencies..."
    
    if [ ! -f "go.mod" ]; then
        print_error "go.mod not found. Please run this script from the project root."
        exit 1
    fi
    
    # Check for required dependencies
    go mod tidy
    print_status "Dependencies checked and updated."
}

# Build the project
build_project() {
    print_status "Building the project..."
    
    if go build -o bin/service_auth cmd/server/main.go; then
        print_status "Build successful!"
    else
        print_error "Build failed. Please check the code for errors."
        exit 1
    fi
}

# Run tests
run_tests() {
    print_status "Running tests..."
    
    # Run unit tests
    go test -v ./internal/application/service/impl/... -count=1
    
    # Run cache strategy tests
    if [ $? -eq 0 ]; then
        print_status "All tests passed!"
    else
        print_warning "Some tests failed. Please check the test output."
    fi
}

# Setup cache configuration
setup_cache_config() {
    print_status "Setting up cache configuration..."
    
    # Create cache config if not exists
    if [ ! -f "config/cache.yaml" ]; then
        cat > config/cache.yaml << EOF
cache:
  local:
    type: "ristretto"
    max_cost: 1000000  # 1MB
    num_counters: 100000
    buffer_items: 64
  distributed:
    type: "redis"
    host: "localhost:6379"
    password: ""
    db: 0
    pool_size: 10
    min_idle_conns: 5
  ttl:
    user_info: 600        # 10 minutes
    access_token: 7200    # 2 hours
    permission: 600       # 10 minutes
    device_check: 300     # 5 minutes
    local_user_info: 120  # 2 minutes
    local_access_token: 300  # 5 minutes
    local_permission: 60  # 1 minute
    local_device_check: 60   # 1 minute
EOF
        print_status "Created cache configuration file: config/cache.yaml"
    else
        print_warning "Cache configuration already exists: config/cache.yaml"
    fi
}

# Check Redis connection
check_redis() {
    print_status "Checking Redis connection..."
    
    if command -v redis-cli &> /dev/null; then
        if redis-cli ping > /dev/null 2>&1; then
            print_status "Redis is running and accessible."
        else
            print_warning "Redis is not running. Please start Redis server."
            print_status "To start Redis: redis-server"
        fi
    else
        print_warning "redis-cli not found. Please install Redis."
        print_status "Installation:"
        print_status "  macOS: brew install redis"
        print_status "  Ubuntu: sudo apt-get install redis-server"
        print_status "  CentOS: sudo yum install redis"
    fi
}

# Performance benchmark
run_benchmark() {
    print_status "Running cache performance benchmark..."
    
    cat > benchmark_test.go << 'EOF'
package main

import (
    "context"
    "fmt"
    "sync"
    "time"
    
    "github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/service/impl"
)

func main() {
    // Initialize cache service
    cacheService, err := impl.NewAuthCacheService()
    if err != nil {
        fmt.Printf("Error initializing cache service: %v\n", err)
        return
    }
    
    ctx := context.Background()
    
    // Benchmark parameters
    numRequests := 1000
    numWorkers := 10
    
    fmt.Printf("🏃 Starting benchmark: %d requests with %d workers\n", numRequests, numWorkers)
    
    // Warmup
    fmt.Println("⏳ Warming up cache...")
    userIDs := make([]string, 10)
    for i := 0; i < 10; i++ {
        userIDs[i] = fmt.Sprintf("test-user-%d", i)
    }
    cacheService.PreloadUserData(ctx, userIDs)
    
    // Benchmark
    start := time.Now()
    var wg sync.WaitGroup
    requestsPerWorker := numRequests / numWorkers
    
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for j := 0; j < requestsPerWorker; j++ {
                userID := fmt.Sprintf("test-user-%d", j%10)
                _, err := cacheService.GetUserInfoCached(ctx, userID)
                if err != nil {
                    fmt.Printf("Error in worker %d: %v\n", workerID, err)
                }
            }
        }(i)
    }
    
    wg.Wait()
    duration := time.Since(start)
    
    // Results
    fmt.Printf("\n📊 Benchmark Results:\n")
    fmt.Printf("Total Requests: %d\n", numRequests)
    fmt.Printf("Total Time: %v\n", duration)
    fmt.Printf("Requests/Second: %.2f\n", float64(numRequests)/duration.Seconds())
    fmt.Printf("Average Response Time: %v\n", duration/time.Duration(numRequests))
    
    // Get cache stats
    if stats, err := cacheService.GetCacheStats(ctx); err == nil {
        fmt.Printf("\n📈 Cache Statistics:\n")
        fmt.Printf("Local Cache Hits: %d\n", stats.LocalCacheHits)
        fmt.Printf("Distributed Cache Hits: %d\n", stats.DistributedCacheHits)
        fmt.Printf("Cache Misses: %d\n", stats.CacheMisses)
        fmt.Printf("Hit Ratio: %.2f%%\n", stats.HitRatio*100)
    }
}
EOF

    if go run benchmark_test.go; then
        print_status "Benchmark completed successfully!"
    else
        print_warning "Benchmark failed or cache service not properly initialized."
    fi
    
    # Clean up
    rm -f benchmark_test.go
}

# Generate usage examples
generate_examples() {
    print_status "Generating usage examples..."
    
    mkdir -p examples
    
    # Example 1: Basic cache usage
    cat > examples/basic_cache_usage.go << 'EOF'
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/service/impl"
)

func main() {
    // Initialize cache service
    cacheService, err := impl.NewAuthCacheService()
    if err != nil {
        log.Fatal("Failed to initialize cache service:", err)
    }
    
    ctx := context.Background()
    
    // Example: Get user info with caching
    userID := "example-user-id"
    userInfo, err := cacheService.GetUserInfoCached(ctx, userID)
    if err != nil {
        log.Printf("Error getting user info: %v", err)
        return
    }
    
    if userInfo != nil {
        fmt.Printf("User Info: %+v\n", userInfo)
    } else {
        fmt.Println("User not found")
    }
    
    // Example: Check permissions with caching
    companyID := "example-company-id"
    hasPermission, err := cacheService.CheckUserPermissionCached(ctx, companyID, userID)
    if err != nil {
        log.Printf("Error checking permission: %v", err)
        return
    }
    
    fmt.Printf("User has permission: %t\n", hasPermission)
}
EOF

    # Example 2: Middleware usage
    cat > examples/middleware_usage.go << 'EOF'
package main

import (
    "log"
    
    "github.com/gin-gonic/gin"
    "github.com/youknow2509/cio_verify_face/server/service_auth/internal/infrastructure/middleware"
)

func main() {
    // Initialize optimized middleware
    optimizedAuth, err := middleware.NewOptimizedAuthMiddleware()
    if err != nil {
        log.Fatal("Failed to initialize optimized middleware:", err)
    }
    
    router := gin.Default()
    
    // Basic authentication middleware
    router.Use(optimizedAuth.Apply())
    
    // Routes with role checking
    adminRoutes := router.Group("/admin")
    adminRoutes.Use(optimizedAuth.ApplyWithRoleCheck(0)) // Admin role only
    adminRoutes.GET("/dashboard", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Admin dashboard"})
    })
    
    // Routes with company permission checking
    companyRoutes := router.Group("/company/:company_id")
    companyRoutes.Use(optimizedAuth.ApplyWithCompanyPermission())
    companyRoutes.GET("/info", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Company info"})
    })
    
    router.Run(":8080")
}
EOF

    print_status "Examples generated in ./examples/ directory"
}

# Main menu
show_menu() {
    echo ""
    echo "Select an option:"
    echo "1) Check dependencies"
    echo "2) Setup cache configuration"
    echo "3) Check Redis connection"
    echo "4) Build project"
    echo "5) Run tests"
    echo "6) Run performance benchmark"
    echo "7) Generate usage examples"
    echo "8) Full setup (all above)"
    echo "9) Exit"
    echo ""
}

# Main script execution
main() {
    check_go
    
    while true; do
        show_menu
        read -p "Enter your choice [1-9]: " choice
        
        case $choice in
            1)
                check_dependencies
                ;;
            2)
                setup_cache_config
                ;;
            3)
                check_redis
                ;;
            4)
                build_project
                ;;
            5)
                run_tests
                ;;
            6)
                run_benchmark
                ;;
            7)
                generate_examples
                ;;
            8)
                print_status "Running full setup..."
                check_dependencies
                setup_cache_config
                check_redis
                build_project
                run_tests
                generate_examples
                print_status "Full setup completed!"
                ;;
            9)
                print_status "Goodbye!"
                exit 0
                ;;
            *)
                print_error "Invalid option. Please choose 1-9."
                ;;
        esac
    done
}

# Run main function
main