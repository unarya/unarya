import torch
import socket

def main():
    # Check CUDA
    if not torch.cuda.is_available():
        print("❌ CUDA not available")
        exit(1)

    # Check port binding
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    result = sock.connect_ex(("localhost", 6000))
    sock.close()
    if result != 0:
        print("❌ Model server not responding on port 6000")
        exit(1)

    print("✅ GPU and Model Server healthy")
    exit(0)

if __name__ == "__main__":
    main()
