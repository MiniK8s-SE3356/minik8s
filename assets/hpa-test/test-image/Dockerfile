FROM python:3.10-slim

# Set the working directory to /app
WORKDIR /app

COPY ./ /app

RUN pip install --no-cache-dir -r requirements.txt


ENTRYPOINT ["python"]

CMD ["main.py"]

# To build the image, run the following command:
# docker build -t hpa-test-image .