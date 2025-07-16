from ctypes import *
import os
import subprocess
import sys
import json

class LorawanWrapperError(Exception):
    """Exception raised for LoRaWAN wrapper errors."""
    pass

class LorawanWrapper:
    """
    LoRaWAN Wrapper with CGO fallback support.
    
    This class automatically detects whether the CGO shared library is available
    and falls back to the command-line tool if needed.
    """
    
    def __init__(self):
        self.use_cgo = False
        self.lib = None
        self.cmd_path = None
        self._initialize()
    
    def _initialize(self):
        """Initialize the wrapper, trying CGO first, then falling back to command-line tool."""
        dirname = os.path.dirname(__file__)
        
        # Try to load CGO shared library
        shared_lib_path = os.path.join(dirname, 'utils/lorawanWrapper.so')
        if os.path.exists(shared_lib_path):
            try:
                self.lib = cdll.LoadLibrary(shared_lib_path)
                self.use_cgo = True
                print("LoRaWAN Wrapper: Using CGO shared library")
                return
            except Exception as e:
                print(f"LoRaWAN Wrapper: Failed to load CGO library: {e}")
        
        # Try to use command-line tool
        cmd_path = os.path.join(dirname, 'utils/lorawanWrapper')
        if os.path.exists(cmd_path):
            try:
                # Test if the command-line tool works
                result = subprocess.run([cmd_path], capture_output=True, text=True)
                if result.returncode == 1:  # Expected return code for help
                    self.cmd_path = cmd_path
                    self.use_cgo = False
                    print("LoRaWAN Wrapper: Using command-line tool (no CGO)")
                    return
            except Exception as e:
                print(f"LoRaWAN Wrapper: Failed to use command-line tool: {e}")
        
        raise LorawanWrapperError(
            "Neither CGO shared library nor command-line tool is available. "
            "Please run 'make build' in the utils directory."
        )
    
    def _run_cmd(self, args):
        """Run command-line tool with given arguments."""
        if not self.cmd_path:
            raise LorawanWrapperError("Command-line tool not available")
        
        try:
            result = subprocess.run([self.cmd_path] + args, capture_output=True, text=True)
            if result.returncode != 0:
                raise LorawanWrapperError(f"Command failed: {result.stderr}")
            return result.stdout
        except Exception as e:
            raise LorawanWrapperError(f"Failed to run command: {e}")
    
    def _ensure_bytes(self, value):
        """Ensure value is bytes for CGO functions."""
        if isinstance(value, str):
            return bytes(value, encoding='utf-8')
        return value
    
    def _ensure_string(self, value):
        """Ensure value is string for command-line functions."""
        if isinstance(value, bytes):
            return value.decode('utf-8')
        return value or ""

    def printPHYPayload(self, phyPayload, key=None):
        """Parse and print PHY payload."""
        if self.use_cgo:
            if isinstance(key, str) and len(key) > 0:
                key = self._ensure_bytes(key)
            else:
                key = None
            
            self.lib.printPHYPayload.argtypes = [c_char_p, c_char_p]
            self.lib.printPHYPayload.restype = c_char_p
            
            try:
                result = self.lib.printPHYPayload(self._ensure_bytes(phyPayload), key)
                return result.decode('utf-8')
            except (AttributeError, TypeError):
                result = self.lib.printPHYPayload(phyPayload, key)
                return result.decode('utf-8')
        else:
            args = ["printPHYPayload", self._ensure_string(phyPayload)]
            if key:
                args.append(self._ensure_string(key))
            return self._run_cmd(args)

    def testAppKeysWithJoinAccept(self, keys, data, dontGenerateKeys):
        """Test application keys with Join Accept."""
        if self.use_cgo:
            keysArr = (c_char_p * len(keys))(*keys)
            generateKeys = 0 if dontGenerateKeys else 1
            
            self.lib.testAppKeysWithJoinAccept.argtypes = [type(keysArr), c_int, c_char_p, c_int]
            self.lib.testAppKeysWithJoinAccept.restype = c_char_p
            
            result = self.lib.testAppKeysWithJoinAccept(
                keysArr, len(keysArr), self._ensure_bytes(data), generateKeys
            )
            return result.decode('utf-8')
        else:
            generateKeys = "0" if dontGenerateKeys else "1"
            args = ["testAppKeysWithJoinAccept"] + [self._ensure_string(k) for k in keys] + [
                "--", self._ensure_string(data), generateKeys
            ]
            return self._run_cmd(args)

    def testAppKeysWithJoinRequest(self, keys, data, dontGenerateKeys):
        """Test application keys with Join Request."""
        if self.use_cgo:
            keysArr = (c_char_p * len(keys))(*keys)
            generateKeys = 0 if dontGenerateKeys else 1
            
            self.lib.testAppKeysWithJoinRequest.argtypes = [type(keysArr), c_int, c_char_p, c_int]
            self.lib.testAppKeysWithJoinRequest.restype = c_char_p
            
            result = self.lib.testAppKeysWithJoinRequest(
                keysArr, len(keysArr), self._ensure_bytes(data), generateKeys
            )
            return result.decode('utf-8')
        else:
            generateKeys = "0" if dontGenerateKeys else "1"
            args = ["testAppKeysWithJoinRequest"] + [self._ensure_string(k) for k in keys] + [
                "--", self._ensure_string(data), generateKeys
            ]
            return self._run_cmd(args)

    def getDevAddr(self, key, data):
        """Get device address from Join Accept."""
        if self.use_cgo:
            self.lib.getDevAddr.argtypes = [c_char_p, c_char_p]
            self.lib.getDevAddr.restype = c_char_p
            result = self.lib.getDevAddr(self._ensure_bytes(key), self._ensure_bytes(data))
            return result.decode('utf-8')
        else:
            args = ["getDevAddr", self._ensure_string(key), self._ensure_string(data)]
            return self._run_cmd(args)

    def getDevEUI(self, data):
        """Get device EUI from payload."""
        if self.use_cgo:
            self.lib.getDevEUI.argtypes = [c_char_p]
            self.lib.getDevEUI.restype = c_char_p
            result = self.lib.getDevEUI(self._ensure_bytes(data))
            return result.decode('utf-8')
        else:
            args = ["getDevEUI", self._ensure_string(data)]
            return self._run_cmd(args)

    def getDevAddrFromMACPayload(self, data):
        """Get device address from MAC payload."""
        if self.use_cgo:
            self.lib.getDevAddrFromMACPayload.argtypes = [c_char_p]
            self.lib.getDevAddrFromMACPayload.restype = c_char_p
            result = self.lib.getDevAddrFromMACPayload(self._ensure_bytes(data))
            return result.decode('utf-8')
        else:
            args = ["getDevAddrFromMACPayload", self._ensure_string(data)]
            return self._run_cmd(args)

    def generateSessionKeysFromJoins(self, joinRequest, joinAccept, appKey):
        """Generate session keys from Join Request and Join Accept."""
        if self.use_cgo:
            self.lib.generateSessionKeysFromJoins.argtypes = [c_char_p, c_char_p, c_char_p]
            self.lib.generateSessionKeysFromJoins.restype = c_char_p
            result = self.lib.generateSessionKeysFromJoins(
                self._ensure_bytes(joinRequest),
                self._ensure_bytes(joinAccept),
                self._ensure_bytes(appKey)
            )
            return result.decode('utf-8')
        else:
            args = ["generateSessionKeysFromJoins", 
                   self._ensure_string(joinRequest),
                   self._ensure_string(joinAccept),
                   self._ensure_string(appKey)]
            return self._run_cmd(args)

    def getDevNonce(self, jr):
        """Get device nonce from Join Request."""
        if self.use_cgo:
            self.lib.getDevNonce.argtypes = [c_char_p]
            self.lib.getDevNonce.restype = c_int
            return self.lib.getDevNonce(self._ensure_bytes(jr))
        else:
            args = ["getDevNonce", self._ensure_string(jr)]
            result = self._run_cmd(args)
            return int(result.strip())

    def getCounter(self, datapayload):
        """Get frame counter from data payload."""
        if self.use_cgo:
            self.lib.getCounter.argtypes = [c_char_p]
            self.lib.getCounter.restype = c_int
            return self.lib.getCounter(self._ensure_bytes(datapayload))
        else:
            args = ["getCounter", self._ensure_string(datapayload)]
            result = self._run_cmd(args)
            return int(result.strip())

    def generateValidMIC(self, data, key, jakey=None):
        """Generate valid MIC for payload."""
        if self.use_cgo:
            if jakey is not None:
                jakey = self._ensure_bytes(jakey)
            
            self.lib.generateValidMIC.argtypes = [c_char_p, c_char_p, c_char_p]
            self.lib.generateValidMIC.restype = c_char_p
            
            try:
                result = self.lib.generateValidMIC(
                    self._ensure_bytes(data), self._ensure_bytes(key), jakey
                )
                return result.decode('utf-8')
            except (AttributeError, TypeError):
                return self.lib.generateValidMIC(data, self._ensure_bytes(key), jakey)
        else:
            args = ["generateValidMIC", self._ensure_string(data), self._ensure_string(key)]
            if jakey is not None:
                args.append(self._ensure_string(jakey))
            return self._run_cmd(args)

    def marshalJsonToPHYPayload(self, json_data, key=None, nwkskey=None):
        """Marshal JSON to PHY payload."""
        if self.use_cgo:
            if key:
                key = self._ensure_bytes(key)
            if nwkskey:
                nwkskey = self._ensure_bytes(nwkskey)
            
            self.lib.marshalJsonToPHYPayload.argtypes = [c_char_p, c_char_p, c_char_p]
            self.lib.marshalJsonToPHYPayload.restype = c_char_p
            
            result = self.lib.marshalJsonToPHYPayload(
                self._ensure_bytes(json_data), key, nwkskey
            )
            return result.decode('utf-8')
        else:
            args = ["marshalJsonToPHYPayload", self._ensure_string(json_data)]
            if key:
                args.append(self._ensure_string(key))
            if nwkskey:
                args.append(self._ensure_string(nwkskey))
            return self._run_cmd(args)

    def getMType(self, datapayload):
        """Get message type from payload."""
        if self.use_cgo:
            self.lib.getMType.argtypes = [c_char_p]
            self.lib.getMType.restype = int
            return self.lib.getMType(self._ensure_bytes(datapayload))
        else:
            args = ["getMType", self._ensure_string(datapayload)]
            result = self._run_cmd(args)
            return int(result.strip())

    def getMajor(self, datapayload):
        """Get major version from payload."""
        if self.use_cgo:
            self.lib.getMajor.argtypes = [c_char_p]
            self.lib.getMajor.restype = int
            return self.lib.getMajor(self._ensure_bytes(datapayload))
        else:
            args = ["getMajor", self._ensure_string(datapayload)]
            result = self._run_cmd(args)
            return int(result.strip())

    def getJoinEUI(self, datapayload):
        """Get Join EUI from payload."""
        if self.use_cgo:
            self.lib.getJoinEUI.argtypes = [c_char_p]
            self.lib.getJoinEUI.restype = c_char_p
            result = self.lib.getJoinEUI(self._ensure_bytes(datapayload))
            return result.decode('utf-8')
        else:
            args = ["getJoinEUI", self._ensure_string(datapayload)]
            return self._run_cmd(args)


# Create global instance for backward compatibility
_wrapper = None

def _get_wrapper():
    global _wrapper
    if _wrapper is None:
        _wrapper = LorawanWrapper()
    return _wrapper

# Backward compatibility functions
def printPHYPayload(phyPayload, key=None):
    return _get_wrapper().printPHYPayload(phyPayload, key)

def testAppKeysWithJoinAccept(keys, data, dontGenerateKeys):
    return _get_wrapper().testAppKeysWithJoinAccept(keys, data, dontGenerateKeys)

def testAppKeysWithJoinRequest(keys, data, dontGenerateKeys):
    return _get_wrapper().testAppKeysWithJoinRequest(keys, data, dontGenerateKeys)

def getDevAddr(key, data):
    return _get_wrapper().getDevAddr(key, data)

def getDevEUI(data):
    return _get_wrapper().getDevEUI(data)

def getDevAddrFromMACPayload(data):
    return _get_wrapper().getDevAddrFromMACPayload(data)

def generateSessionKeysFromJoins(joinRequest, joinAccept, appKey):
    return _get_wrapper().generateSessionKeysFromJoins(joinRequest, joinAccept, appKey)

def getDevNonce(jr):
    return _get_wrapper().getDevNonce(jr)

def getCounter(datapayload):
    return _get_wrapper().getCounter(datapayload)

def generateValidMIC(data, key, jakey=None):
    return _get_wrapper().generateValidMIC(data, key, jakey)

def marshalJsonToPHYPayload(json_data, key=None, nwkskey=None):
    return _get_wrapper().marshalJsonToPHYPayload(json_data, key, nwkskey)

def getMType(datapayload):
    return _get_wrapper().getMType(datapayload)

def getMajor(datapayload):
    return _get_wrapper().getMajor(datapayload)

def getJoinEUI(datapayload):
    return _get_wrapper().getJoinEUI(datapayload)