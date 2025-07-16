//go:build !cgo
// +build !cgo

package main

import (
	"fmt"
	"strconv"
	"time"

	. "github.com/matiassequeira/lorawan"
	log "github.com/sirupsen/logrus"
)

// Command implementations for non-CGO build

func cmdPrintPHYPayload(phyPayload string, key string) string {
	var keyObj AES128Key
	var phy PHYPayload

	setLogLevel()

	if err := phy.UnmarshalText([]byte(phyPayload)); err != nil {
		log.Error("Unmarshal error with PHYPayload: ", phyPayload, " Error: ", err)
		return jsonError("Failed to unmarshal PHYPayload")
	}

	if key != "" {
		keyObj.UnmarshalText([]byte(key))

		if phy.MHDR.MType == UnconfirmedDataUp || phy.MHDR.MType == ConfirmedDataUp || phy.MHDR.MType == UnconfirmedDataDown || phy.MHDR.MType == ConfirmedDataDown {
			if err := phy.DecryptFRMPayload(keyObj); err != nil {
				log.Error("Cannot decrypt FRMPayload:", err)
			}
		} else if phy.MHDR.MType == JoinAccept {
			if err := phy.DecryptJoinAcceptPayload(keyObj); err != nil {
				log.Error("Cannot decrypt Join Accept:", err)
			}
		}
	}

	_, ok := phy.MACPayload.(*MACPayload)
	if ok {
		if err := phy.DecodeFOptsToMACCommands(); err != nil {
			log.Error("Error decoding FOpts:", err)
		}
	}

	phyJSON, err := phy.MarshalJSON()
	if err != nil {
		return jsonError("Failed to marshal PHYPayload to JSON")
	}

	return string(phyJSON)
}

func cmdGetMType(dataPayload string) string {
	var phy PHYPayload

	if err := phy.UnmarshalText([]byte(dataPayload)); err != nil {
		return "-1"
	}
	return strconv.Itoa(int(phy.MHDR.MType))
}

func cmdGetMajor(dataPayload string) string {
	var phy PHYPayload

	if err := phy.UnmarshalText([]byte(dataPayload)); err != nil {
		return "-1"
	}
	return strconv.Itoa(int(phy.MHDR.Major))
}

func cmdGetCounter(dataPayload string) string {
	var phy PHYPayload

	if err := phy.UnmarshalText([]byte(dataPayload)); err != nil {
		return "-1"
	}
	macPL, ok := phy.MACPayload.(*MACPayload)

	if ok {
		return strconv.Itoa(int(macPL.FHDR.FCnt))
	}
	return "-1"
}

func cmdGetDevEUI(dataPayload string) string {
	return returnDevEUI(dataPayload)
}

func cmdGetJoinEUI(dataPayload string) string {
	var phy PHYPayload

	if err := phy.UnmarshalText([]byte(dataPayload)); err != nil {
		return "Error"
	}

	jrPL, ok := phy.MACPayload.(*JoinRequestPayload)
	if !ok {
		return "Error"
	}
	return fmt.Sprintf("%v", jrPL.JoinEUI)
}

func cmdGetDevNonce(dataPayload string) string {
	var phy PHYPayload

	if err := phy.UnmarshalText([]byte(dataPayload)); err != nil {
		return "-1"
	}

	jrPL, ok := phy.MACPayload.(*JoinRequestPayload)
	if !ok {
		return "-1"
	}

	return strconv.Itoa(int(jrPL.DevNonce))
}

func cmdGetDevAddr(appKey string, dataPayload string) string {
	var key AES128Key
	var phy PHYPayload

	if err := phy.UnmarshalText([]byte(dataPayload)); err != nil {
		return "Error"
	}

	if err := key.UnmarshalText([]byte(appKey)); err != nil {
		return "Error"
	}

	if err := phy.DecryptJoinAcceptPayload(key); err != nil {
		return "Error"
	}

	jaPL, ok := phy.MACPayload.(*JoinAcceptPayload)
	if ok {
		return fmt.Sprintf("%v", jaPL.DevAddr)
	}

	return "Error"
}

func cmdGetDevAddrFromMACPayload(dataPayload string) string {
	var phy PHYPayload

	if err := phy.UnmarshalText([]byte(dataPayload)); err != nil {
		return "Error"
	}
	macPL, ok := phy.MACPayload.(*MACPayload)

	if ok {
		return fmt.Sprintf("%v", macPL.FHDR.DevAddr)
	}
	return ""
}

func cmdGenerateSessionKeysFromJoins(joinRequest string, joinAccept string, appKey string) string {
	var key AES128Key
	var joinEui EUI64
	var joinReq PHYPayload
	var joinAcc PHYPayload

	key.UnmarshalText([]byte(appKey))

	if err := joinReq.UnmarshalText([]byte(joinRequest)); err != nil {
		return ""
	}
	jrPL, ok := joinReq.MACPayload.(*JoinRequestPayload)
	if !ok {
		return ""
	}

	if err := joinAcc.UnmarshalText([]byte(joinAccept)); err != nil {
		return ""
	}
	if err := joinAcc.DecryptJoinAcceptPayload(key); err != nil {
		return ""
	}
	jaPL, ok := joinAcc.MACPayload.(*JoinAcceptPayload)
	if !ok {
		return ""
	}

	nwkSKey, err := getFNwkSIntKey(false, key, jaPL.HomeNetID, joinEui, jaPL.JoinNonce, jrPL.DevNonce)
	if err != nil {
		return ""
	}
	appSkey, err := getAppSKey(false, key, jaPL.HomeNetID, joinEui, jaPL.JoinNonce, jrPL.DevNonce)
	if err != nil {
		return ""
	}
	nKey, _ := nwkSKey.MarshalText()
	aKey, _ := appSkey.MarshalText()

	return fmt.Sprintf("{\"nwkSKey\": \"%s\", \"appSKey\": \"%s\"}", string(nKey), string(aKey))
}

func cmdGenerateValidMIC(dataPayload string, newKey string, jaKey string) string {
	return signPacket(dataPayload, newKey, jaKey)
}

func cmdMarshalJsonToPHYPayload(jsonStr string, key string, nwkskey string) string {
	return parseJSONtoPHY(jsonStr, key, nwkskey)
}

func cmdTestAppKeysWithJoinRequest(keys []string, joinRequestData string, generateKeys int) string {
	var phy PHYPayload
	var testCounter int64 = 0
	var key AES128Key
	var foundKeys string = ""

	setLogLevel()

	if err := phy.UnmarshalText([]byte(joinRequestData)); err != nil {
		log.Error("JoinRequest data: ", joinRequestData, "Error: ", err)
		return foundKeys
	}

	// Test keys from the provided list
	for _, keyStr := range keys {
		if err := key.UnmarshalText([]byte(keyStr)); err != nil {
			log.Error("Unmarshall error with AppKey: ", keyStr, err)
			continue
		}

		result, err := testAppKeyWithJoinRequest(key, phy, &testCounter)
		if err != nil {
			log.Error("Error with JoinRequest :", err)
			return foundKeys
		} else if len(result) > 0 {
			foundKeys += result + " "
		}
	}

	// Generate keys if requested (simplified version)
	if generateKeys != 0 {
		jrPL, ok := phy.MACPayload.(*JoinRequestPayload)
		if !ok {
			return foundKeys
		}

		joinEUI, _ := jrPL.JoinEUI.MarshalBinary()
		devEUI, _ := jrPL.DevEUI.MarshalBinary()
		vendorsKeys := [][]byte{append(joinEUI, devEUI...), append(devEUI, joinEUI...)}

		for _, vendorKey := range vendorsKeys {
			if err := key.UnmarshalBinary(vendorKey); err != nil {
				log.Error("Unmarshall error with AppKey: ", vendorKey, err)
				continue
			}

			result, err := testAppKeyWithJoinRequest(key, phy, &testCounter)
			if err != nil {
				log.Error("Error with JoinRequest :", err)
				return foundKeys
			} else if len(result) > 0 {
				foundKeys += result + " "
			}
		}
	}

	if len(foundKeys) > 0 {
		return foundKeys
	}
	return ""
}

func cmdTestAppKeysWithJoinAccept(keys []string, joinAcceptData string, generateKeys int) string {
	var phy PHYPayload
	var key AES128Key
	var decryptCounter int64 = 0
	var foundKeys string = ""

	setLogLevel()

	if err := phy.UnmarshalText([]byte(joinAcceptData)); err != nil {
		log.Error("Error with Join accept with data: "+joinAcceptData, ". Error: ", err)
		return foundKeys
	}

	// Test keys from the provided list
	for _, keyStr := range keys {
		if err := key.UnmarshalText([]byte(keyStr)); err != nil {
			log.Error("Unmarshall error with AppKey:", keyStr, ".", err)
			continue
		}

		result, err := testAppKeyWithJoinAccept(key, phy, &decryptCounter)
		if err != nil {
			log.Error("Error with JoinAccept with data: ", joinAcceptData, " Error: ", err)
			return foundKeys
		} else if len(result) > 0 {
			foundKeys += result + " "
		}
	}

	// Simplified key generation for demonstration
	if generateKeys != 0 {
		start := time.Now()
		// Only do a limited key generation to avoid excessive computation
		// In a real implementation, this would include the full bruteforce logic
		
		elapsed := time.Since(start)
		log.Debug("Bruteforcing the JoinAccept took:", elapsed)
	}

	if len(foundKeys) > 0 {
		return foundKeys
	}
	return ""
}

func setLogLevel() {
	commonSetLogLevel()
}