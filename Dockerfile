# Set the base image to use for subsequent instructions
FROM eclipse-temurin:21

# Set the working directory inside the container
WORKDIR /usr/src

ADD https://repo1.maven.org/maven2/org/odftoolkit/odfvalidator/0.12.0/odfvalidator-0.12.0-jar-with-dependencies.jar ./odfvalidator.jar

# Copy any source file(s) required for the action
COPY entrypoint.sh .

# Configure the container to be run as an executable
ENTRYPOINT ["/usr/src/entrypoint.sh"]
