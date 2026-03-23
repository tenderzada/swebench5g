FROM golang:1.25-bookworm

LABEL maintainer="swebench5g"
LABEL description="SWE-Bench 5G Pilot - PCF Issue #879: nil pointer panic in policyauthorization"
LABEL task.issue="free5gc/free5gc#879"
LABEL task.difficulty="easy"

# Install system dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    curl \
    jq \
    && rm -rf /var/lib/apt/lists/*

# Use Chinese Go proxy to avoid timeout on proxy.golang.org
ENV GOPROXY=https://goproxy.cn,direct

# Set workspace
WORKDIR /opt/free5gc-pcf

# Clone PCF repo at the BUGGY commit (before PR #65 fix)
# Use ghproxy mirror to avoid GitHub connectivity issues in China
RUN git clone https://ghproxy.com/https://github.com/free5gc/pcf.git . && \
    git remote set-url origin https://github.com/free5gc/pcf.git && \
    git checkout df535f5

# Download Go dependencies
RUN go mod download

# Verify the code compiles
RUN go build ./...

# Copy test suite and scripts
COPY test-suite/ /opt/test-suite/
COPY problem_statement.md /opt/task/problem_statement.md
COPY scripts/ /opt/scripts/

# Make scripts executable
RUN chmod +x /opt/test-suite/*.sh /opt/scripts/*.sh

# Set working directory to the PCF source
WORKDIR /opt/free5gc-pcf

# Default command
CMD ["/bin/bash"]
