for i in {0..50}
do
  curl --request GET -sL \
       --url 'http://localhost:8080/req'
  echo " $i"
done
