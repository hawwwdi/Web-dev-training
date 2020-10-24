for i in $(seq 1 $1); do
  curl --request GET -sL \
    --url 'http://localhost:8080/req'
  echo " $i"
done
