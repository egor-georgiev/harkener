import argparse
import logging
import socket

MESSAGE_SIZE = 2

logging.basicConfig(
    level=logging.INFO,
    format="[%(asctime)s.%(msecs)03d] %(message)s",
    datefmt="%Y-%m-%dT%H:%M:%S",
)
logger = logging.getLogger(__name__)


def parse_args() -> tuple[str, int]:
    parser = argparse.ArgumentParser()
    parser.add_argument("host")
    parser.add_argument("port", type=int)

    args = parser.parse_args()
    host, port = args.host, args.port

    return host, port


def main() -> None:
    client = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    addr = parse_args()
    message = b"hi"
    client.sendto(message, addr)
    _, local_port = client.getsockname()
    logger.info(f"opened a socket on local port: {local_port}")
    while True:
        data, _ = client.recvfrom(MESSAGE_SIZE)
        decoded = int.from_bytes(data, byteorder="big")
        logger.info(decoded)


if __name__ == "__main__":
    main()
