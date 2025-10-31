from watchdog.observers.polling import PollingObserver as Observer
from watchdog.events import FileSystemEventHandler
import subprocess, time, threading


class RestartHandler(FileSystemEventHandler):
    def __init__(self, cmd, debounce_seconds=1.5):
        self.cmd = cmd
        self.process = subprocess.Popen(cmd)
        self.debounce_seconds = debounce_seconds
        self._timer = None
        print(f"🧠 Started process PID {self.process.pid}")

    def restart(self):
        print("🔄 File change detected, restarting server...")
        self.process.terminate()
        try:
            self.process.wait(timeout=5)
        except subprocess.TimeoutExpired:
            self.process.kill()

        self.process = subprocess.Popen(self.cmd)
        print(f"♻️ Restarted process PID {self.process.pid}")

    def on_any_event(self, event):
        if event.is_directory:
            return

        # Huỷ timer cũ nếu vẫn đang đợi
        if self._timer and self._timer.is_alive():
            self._timer.cancel()

        # Tạo timer mới → chỉ restart khi không có event mới trong khoảng debounce_seconds
        self._timer = threading.Timer(self.debounce_seconds, self.restart)
        self._timer.start()


if __name__ == "__main__":
    path = "./src"
    cmd = ["python3", "-u", "-m", "src.modelserver.server"]
    event_handler = RestartHandler(cmd)
    observer = Observer()
    observer.schedule(event_handler, path, recursive=True)
    observer.start()

    print(f"👀 Watching for changes in {path} (debounce={event_handler.debounce_seconds}s)...")

    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        observer.stop()
        print("🛑 Stopped watcher.")
    observer.join()
