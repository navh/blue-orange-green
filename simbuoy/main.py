import socket
import time
from proto.buoy_pb2 import BuoyStatus

def create_sample_message():
    msg = BuoyStatus()
    msg.buoy_id = 1234
    msg.report_id = 1333
    msg.timestamp = int(time.time())
    
    # Position data
    msg.latitude = 45.5231
    msg.longitude = -122.6765
    msg.depth_meters = 0.0
    
    # Environmental data
    msg.temp_celsius = 15.5
    
    # Accelerometer readings
    msg.accel_x = 0.0
    msg.accel_y = 0.0
    msg.accel_z = 0.0
    
    # System status
    msg.battery_level_percent = 95
    
    return msg

def send_message(host='localhost', port=8080):
    # Create a TCP/IP socket
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    
    try:
        # Connect to the server
        sock.connect((host, port))
        
        # Create and serialize message
        message = create_sample_message()
        serialized_message = message.SerializeToString()
        
        # Send message length first (as 4 bytes, big endian)
        message_length = len(serialized_message)
        sock.send(message_length.to_bytes(4, byteorder='big'))
        
        # Send the actual message
        sock.send(serialized_message)
        print(f"Sent message: {message}")
        
    finally:
        sock.close()

if __name__ == "__main__":
    try:
        send_message()
        print("Message sent successfully")
    except Exception as e:
        print(f"Error sending message: {e}")
