go build -o gitdoc-cli

if [ $? -ne 0 ]; then
    echo "build failed"
    exit 1
fi

echo "build success"
