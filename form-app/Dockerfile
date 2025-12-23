# 1. Install Go (The Environment)
FROM golang:1.23-alpine

# 2. Set environment variables for Mongo
ENV MONGO_DB_USERNAME=admin \
    MONGO_DB_PWD=password

# 3. Create and set the folder (Replaces 'mkdir')
WORKDIR /app

# 4. Copy your current folder files into /app
COPY . .

# 5. NEW STEP FOR GO: Compile your code into a binary named 'main'
RUN go build -o main .

# 6. Start the app by running the binary
CMD ["./main"]