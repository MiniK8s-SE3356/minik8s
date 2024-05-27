FROM python:3.10-slim

# Set the working directory to /app
WORKDIR /app

COPY ./ /app/

RUN pip install --no-cache-dir -r requirements.txt

# Make port 5000 available to the world outside this container
EXPOSE 5000

ENTRYPOINT ["python"]

CMD ["base.py"]

# To build the image, run the following command:
# docker build -t serverless-base-image .