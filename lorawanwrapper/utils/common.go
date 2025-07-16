package main

import (
	"fmt"
	"os"

	. "github.com/matiassequeira/lorawan"
	log "github.com/sirupsen/logrus"
)

// Common functions used by both CGO and non-CGO builds

func returnDevEUI(dataBytes string) string {
	var phy PHYPayload

	if err := phy.UnmarshalText([]byte(dataBytes)); err != nil {
		fmt.Println(err)
		fmt.Println("Join request data: " + dataBytes)
		return ""
	}

	jrPL, ok := phy.MACPayload.(*JoinRequestPayload)
	if !ok {
		fmt.Println("MACPayload must be a *JoinRequestPayload")
		return "Error"
	}
	return fmt.Sprintf("%v", jrPL.DevEUI)
}

func testAppKeyWithJoinRequest(key AES128Key, phy PHYPayload, counter *int64) (string, error) {
	*counter++

	result, err := phy.ValidateUplinkJoinMIC(key)

	if err != nil {
		log.Error("Error validating JoinRequest MIC: ", err)
		return "", nil
	}

	if result {
		foundKey, _ := key.MarshalText()
		return string(foundKey), nil
	} else {
		return "", nil
	}
}

func testAppKeyWithJoinAccept(key AES128Key, phy PHYPayload, counter *int64) (string, error) {
	*counter++

	if err := phy.DecryptJoinAcceptPayload(key); err != nil {
		// Here we return nil instead of the error since this error is caused by a wrong key
		return "", nil
	}
	joinEUI := EUI64{8, 7, 6, 5, 4, 3, 2, 1}
	devNonce := DevNonce(258)

	result, err := phy.ValidateDownlinkJoinMIC(JoinRequestType, joinEUI, devNonce, key)

	if err != nil {
		// Here we return nil instead of the error since this error is caused by a wrong key
		return "", nil
	}

	if result == true {
		foundKey, _ := key.MarshalText()
		return string(foundKey), nil
	} else {
		return "", nil
	}
}

func commonSetLogLevel() {
	if env := os.Getenv("ENVIRONMENT"); env == "PROD" {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}
}