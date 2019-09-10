RT_INSTANCE_NAME=rt_testing
RT_VERSION=4.4

cd testdata/

docker rm -f $RT_INSTANCE_NAME
docker run -d --name $RT_INSTANCE_NAME -p 127.0.0.1:8080:80 -e RT_WEB_PORT=8080 netsandbox/request-tracker:$RT_VERSION
docker cp rt_initialdata.pl $RT_INSTANCE_NAME:/root/.
docker exec $RT_INSTANCE_NAME /opt/rt4/sbin/rt-setup-database --action insert --datafile /root/rt_initialdata.pl

curl -F "content=<ticket_formdata" -F user=root -F pass=password localhost:8080/REST/1.0/ticket/new

cd ..

go test