FROM docker.io/library/golang:1.25 AS BUILDER

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY main.go .
RUN go build -v -o /usr/local/bin/odfvalidatorparser main.go

# Set the base image to use for subsequent instructions
FROM docker.io/library/eclipse-temurin:25

# Set the working directory inside the container
WORKDIR /usr/src

ADD https://repo1.maven.org/maven2/org/odftoolkit/odfvalidator/0.12.0/odfvalidator-0.12.0-jar-with-dependencies.jar ./odfvalidator.jar

COPY --from=BUILDER /usr/local/bin/odfvalidatorparser /usr/local/bin/odfvalidatorparser

# Copy any source file(s) required for the action
COPY entrypoint.sh .

# Configure the container to be run as an executable
ENTRYPOINT ["/usr/src/entrypoint.sh"]
