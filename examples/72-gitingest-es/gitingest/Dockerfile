FROM python:3.10-slim

RUN <<EOF
apt-get update
apt-get install -y curl git
EOF

WORKDIR /app

COPY requirements.txt .

RUN pip install -r requirements.txt

COPY . .

CMD ["python", "app.py"]