#!/bin/bash
# =============================================================================
# FIX DOCKER GPU SUPPORT ON WSL2 (WITHOUT DOCKER DESKTOP) - CORRECTED VERSION
# =============================================================================

set -e  # Exit on error

echo "🔧 Fixing Docker GPU support for WSL2..."

# -----------------------------------------------------------------------------
# STEP 1: Verify NVIDIA Driver on Windows Host
# -----------------------------------------------------------------------------
echo ""
echo "📋 STEP 1: Checking NVIDIA driver on Windows host..."
echo "Please run this command in Windows PowerShell:"
echo "  nvidia-smi"
echo ""
echo "If this fails, install NVIDIA GPU Driver for Windows (NOT WSL driver)"
echo "Download from: https://www.nvidia.com/Download/index.aspx"
echo ""
read -p "Press Enter after verifying nvidia-smi works in Windows..."

# -----------------------------------------------------------------------------
# STEP 2: Check WSL CUDA Access
# -----------------------------------------------------------------------------
echo ""
echo "📋 STEP 2: Verifying CUDA access in WSL2..."
if command -v nvidia-smi &> /dev/null; then
    echo "✅ nvidia-smi found in WSL"
    nvidia-smi || echo "⚠️  nvidia-smi failed but might work after setup"
else
    echo "ℹ️  nvidia-smi not in PATH (normal for WSL2, uses Windows driver)"
fi

# -----------------------------------------------------------------------------
# STEP 3: Clean Up Old Installation
# -----------------------------------------------------------------------------
echo ""
echo "📋 STEP 3: Cleaning up old installations..."
sudo rm -f /etc/apt/sources.list.d/nvidia-container-toolkit.list
sudo rm -f /etc/apt/sources.list.d/nvidia*.list
sudo rm -f /usr/share/keyrings/nvidia-container-toolkit-keyring.gpg
sudo apt-get remove --purge -y nvidia-docker2 nvidia-container-toolkit nvidia-container-runtime 2>/dev/null || true
sudo apt-get autoremove -y

# -----------------------------------------------------------------------------
# STEP 4: Install NVIDIA Container Toolkit (CORRECTED METHOD)
# -----------------------------------------------------------------------------
echo ""
echo "📋 STEP 4: Installing NVIDIA Container Toolkit..."

# Detect distribution
distribution=$(. /etc/os-release; echo $ID$VERSION_ID)
echo "Detected distribution: $distribution"

# Add NVIDIA GPG key
curl -fsSL https://nvidia.github.io/libnvidia-container/gpgkey | \
    sudo gpg --dearmor -o /usr/share/keyrings/nvidia-container-toolkit-keyring.gpg

# Add repository with correct format
echo "deb [signed-by=/usr/share/keyrings/nvidia-container-toolkit-keyring.gpg] https://nvidia.github.io/libnvidia-container/stable/deb/\$(ARCH) /" | \
    sudo tee /etc/apt/sources.list.d/nvidia-container-toolkit.list

# Update and install
echo "Updating apt cache..."
sudo apt-get update

echo "Installing nvidia-container-toolkit..."
sudo apt-get install -y nvidia-container-toolkit

# Verify installation
echo ""
echo "Verifying installed packages..."
dpkg -l | grep nvidia-container

# -----------------------------------------------------------------------------
# STEP 5: Configure NVIDIA Container Runtime
# -----------------------------------------------------------------------------
echo ""
echo "📋 STEP 5: Configuring NVIDIA container runtime..."

# Use nvidia-ctk to configure Docker
sudo nvidia-ctk runtime configure --runtime=docker

echo "✅ Docker configured with NVIDIA runtime"

# Show configuration
echo ""
echo "Current /etc/docker/daemon.json:"
sudo cat /etc/docker/daemon.json

# -----------------------------------------------------------------------------
# STEP 6: Restart Docker
# -----------------------------------------------------------------------------
echo ""
echo "📋 STEP 6: Restarting Docker daemon..."
sudo systemctl daemon-reload
sudo systemctl restart docker

# Wait for Docker to start
echo "Waiting for Docker to start..."
sleep 5

# Check Docker status
if sudo systemctl is-active --quiet docker; then
    echo "✅ Docker is running"
else
    echo "❌ Docker failed to start"
    sudo systemctl status docker
    exit 1
fi

# -----------------------------------------------------------------------------
# STEP 7: Verify Configuration
# -----------------------------------------------------------------------------
echo ""
echo "📋 STEP 7: Verifying Docker configuration..."
echo ""
echo "Available runtimes:"
docker info | grep -A 3 "Runtimes:"
echo ""
echo "Default runtime:"
docker info | grep "Default Runtime:"

# Check if nvidia-container-runtime is available
echo ""
echo "Checking nvidia-container-runtime-hook:"
which nvidia-container-runtime-hook || echo "⚠️  nvidia-container-runtime-hook not in PATH"

# -----------------------------------------------------------------------------
# STEP 8: Test GPU Access
# -----------------------------------------------------------------------------
echo ""
echo "📋 STEP 8: Testing GPU access..."
echo ""
echo "Test 1: Using --gpus all flag"
echo "Running: docker run --rm --gpus all nvidia/cuda:12.2.0-base-ubuntu22.04 nvidia-smi"
echo ""

if docker run --rm --gpus all nvidia/cuda:12.2.0-base-ubuntu22.04 nvidia-smi; then
    echo ""
    echo "✅✅✅ SUCCESS! Docker can access GPU with --gpus flag"
    GPU_WORKS=true
else
    echo ""
    echo "❌ Test 1 failed"
    GPU_WORKS=false
fi

# Test without --gpus if configured as default runtime
if grep -q '"default-runtime": "nvidia"' /etc/docker/daemon.json 2>/dev/null; then
    echo ""
    echo "Test 2: Using default runtime (no --gpus flag needed)"
    echo "Running: docker run --rm nvidia/cuda:12.2.0-base-ubuntu22.04 nvidia-smi"
    echo ""

    if docker run --rm nvidia/cuda:12.2.0-base-ubuntu22.04 nvidia-smi; then
        echo ""
        echo "✅✅✅ SUCCESS! Docker can access GPU with default runtime"
        GPU_WORKS=true
    else
        echo ""
        echo "❌ Test 2 failed"
    fi
fi

# -----------------------------------------------------------------------------
# STEP 9: Additional Diagnostics
# -----------------------------------------------------------------------------
echo ""
echo "📋 STEP 9: Additional diagnostics..."
echo ""

echo "Checking NVIDIA devices in WSL:"
ls -la /dev/nvidia* 2>/dev/null || echo "⚠️  No /dev/nvidia* devices found"

echo ""
echo "Checking for dxg device (WSL2 GPU support):"
ls -la /dev/dxg 2>/dev/null || echo "⚠️  No /dev/dxg device found"

echo ""
echo "WSL version info:"
cat /proc/version

# -----------------------------------------------------------------------------
# FINAL STATUS
# -----------------------------------------------------------------------------
echo ""
echo "=========================================="
if [ "$GPU_WORKS" = true ]; then
    echo "🎉🎉🎉 SETUP SUCCESSFUL! 🎉🎉🎉"
    echo "=========================================="
    echo ""
    echo "Your Docker can now access the GPU!"
    echo ""
    echo "Usage examples:"
    echo "  docker run --rm --gpus all nvidia/cuda:12.2.0-base-ubuntu22.04 nvidia-smi"
    echo "  docker run --rm --gpus all pytorch/pytorch:latest python -c 'import torch; print(torch.cuda.is_available())'"
else
    echo "❌ SETUP INCOMPLETE - TROUBLESHOOTING NEEDED"
    echo "=========================================="
    echo ""
    echo "Common issues and solutions:"
    echo ""
    echo "1. ⚠️  MISSING NVIDIA DEVICES (/dev/nvidia*):"
    echo "   Solution:"
    echo "   - Open PowerShell as Administrator"
    echo "   - Run: wsl --shutdown"
    echo "   - Wait 10 seconds"
    echo "   - Start WSL again"
    echo "   - Check: ls -la /dev/nvidia*"
    echo ""
    echo "2. ⚠️  OUTDATED WSL KERNEL:"
    echo "   Solution:"
    echo "   - Open PowerShell as Administrator"
    echo "   - Run: wsl --update"
    echo "   - Run: wsl --shutdown"
    echo "   - Restart WSL"
    echo ""
    echo "3. ⚠️  NVIDIA DRIVER NOT INSTALLED ON WINDOWS:"
    echo "   Solution:"
    echo "   - Download from: https://www.nvidia.com/Download/index.aspx"
    echo "   - Install Windows NVIDIA driver (NOT WSL-specific driver)"
    echo "   - Restart Windows"
    echo ""
    echo "4. ⚠️  WINDOWS VERSION TOO OLD:"
    echo "   - Windows 10 Build 21362+ or Windows 11 required"
    echo "   - Update Windows to latest version"
    echo ""
    echo "5. 🔧 MANUAL VERIFICATION STEPS:"
    echo "   - In Windows PowerShell: nvidia-smi (should work)"
    echo "   - In WSL: ls -la /dev/nvidia* (should show devices)"
    echo "   - In WSL: cat /proc/version (kernel should be 5.10.43.3+)"
    echo ""
    echo "After fixing issues, rerun this script."
fi
echo "=========================================="