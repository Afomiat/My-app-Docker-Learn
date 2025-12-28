# --- STAGE 1: The Builder (The Factory) ---
    FROM golang:1.23-alpine AS builder

    # Create the working directory
    WORKDIR /app
    
    # Copy your Go files (main.go, go.mod, etc.)
    COPY . .
    
    # Compile the code into a "Static Binary" 
    # CGO_ENABLED=0 ensures it runs perfectly on the tiny Alpine image
    RUN CGO_ENABLED=0 go build -o main .
    
# --- STAGE 2: The Runner ---
    FROM alpine:latest
    WORKDIR /app
    
    # 1. Copy the binary
    COPY --from=builder /app/main .
    
    # 2. ADD THIS LINE: Copy your static folder into the runner stage
    COPY --from=builder /app/static ./static
    
    ENV MONGO_DB_USERNAME=admin \
        MONGO_DB_PWD=password
    
    # Tell the container to run your app
    CMD ["./main"]