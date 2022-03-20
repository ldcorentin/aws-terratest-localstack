echo "waiting for localstack.."
until $(nc -zv localhost 4566); do
    printf '.'
    sleep 1
done