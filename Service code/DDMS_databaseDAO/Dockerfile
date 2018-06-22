#Parent image
FROM java:8
#Maintainer 
MAINTAINER padraigh.jarvis
#Working dir 
RUN mkdir /databaseAccess
WORKDIR /databaseAccess
#CP contents to workdir
ADD . /databaseAccess
#Install maven
RUN apt-get update 
RUN apt-get install -y maven
#Download dependencies 
RUN ["mvn", "dependency:resolve"]
RUN ["mvn", "verify"]
#Package everything into a uberjar 
RUN ["mvn", "compile", "assembly:single"]
#Run app when container launches 
CMD ["/usr/lib/jvm/java-8-openjdk-amd64/bin/java", "-jar", "target/databaseAccess-1-jar-with-dependencies.jar"]
