#!/bin/bash

# Check if protobuf compiler is installed
if command -v protoc >/dev/null 2>&1; then
    echo "Protobuf compiler is installed."
else
    echo "Protobuf compiler is not installed."
fi

# Check if git is installed
if command -v git >/dev/null 2>&1; then
    echo "Git is installed."
else
    echo "Git is not installed."
fi

# Check if openssl is installed
if command -v openssl >/dev/null 2>&1; then
    echo "OpenSSL is installed."
else
    echo "OpenSSL is not installed."
fi

# Check if jq is installed
if command -v jq >/dev/null 2>&1; then
    echo "jq is installed."
else
    echo "jq is not installed."
fi

# Check if graphviz is installed
if command -v dot >/dev/null 2>&1; then
    echo "Graphviz is installed."
else
    echo "Graphviz is not installed."
fi

# Check if Go is installed
if command -v go >/dev/null 2>&1; then
    echo "Go is installed."
else
    echo "Go is not installed."
fi
